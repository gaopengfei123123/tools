// 参照 github.com/astaxie/beego/logs 写的日志包
// 日志输出 console, file 模式

package logs

import (
	"fmt"
	// "log"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const levelLoggerImpl = -1

// RFC5424 log message levels.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// 日志输出格式
const (
	AdapterConsole = "console"
	AdapterFile    = "file"
)

// 别名
const (
	LevelInfo  = LevelInformational
	LevelTrace = LevelDebug
	LevelWarn  = LevelWarning
)

type newLoggerFunc func() LoggerInterface

// FormaterFunc 自定义文件格式构造方法
type FormaterFunc func(logLevel int, msg string, v ...interface{}) string

// LoggerInterface 各模式处理类需要实现的接口
type LoggerInterface interface {
	Init(config string) error
	WriteMsg(when time.Time, msg string, level int) error
	Destroy()
	Flush()
}

var adapters = make(map[string]newLoggerFunc)
var levelPrefix = [LevelDebug + 1]string{"[M] ", "[A] ", "[C] ", "[E] ", "[W] ", "[N] ", "[I] ", "[D] "}

// Register 注册处理日志的实体, 只是注册,并没有放在logger里面的outputs里面
func Register(name string, log newLoggerFunc) {
	if log == nil {
		panic("logs: 需要一个日志处理实体")
	}
	if _, dup := adapters[name]; dup {
		panic("logs: 重复注册 " + name)
	}
	adapters[name] = log
}

// Logger 日志结构
type Logger struct {
	lock                sync.Mutex
	level               int
	init                bool
	enableFuncCallDepth bool
	loggerFuncCallDepth int
	asynchronous        bool
	prefix              string
	msgChanLen          int64
	msgChan             chan *logMsg
	signalChan          chan string
	wg                  sync.WaitGroup
	outputs             []*nameLogger
	formatter           FormaterFunc // 可以自定义日志格式化方法
}

const defaultAsyncMsgLen = 1e3

type nameLogger struct {
	LoggerInterface
	Name string
}

type logMsg struct {
	level int
	msg   string
	when  time.Time
}

var logMsgPool *sync.Pool

// NewLogger 实例化一个日志处理
func NewLogger(channelLen ...int64) *Logger {
	lg := new(Logger)
	lg.level = LevelDebug
	lg.loggerFuncCallDepth = 2
	// 默认缓冲channel长度1k
	lg.msgChanLen = append(channelLen, 0)[0]
	if lg.msgChanLen <= 0 {
		lg.msgChanLen = defaultAsyncMsgLen
	}
	lg.signalChan = make(chan string, 1)
	return lg
}

// Async 设置是否开启异步模式
func (lg *Logger) Async(msgLen ...int64) *Logger {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	if lg.asynchronous {
		return lg
	}

	lg.asynchronous = true
	if len(msgLen) > 0 && msgLen[0] > 0 {
		lg.msgChanLen = msgLen[0]
	}
	lg.msgChan = make(chan *logMsg, lg.msgChanLen)
	logMsgPool = &sync.Pool{
		New: func() interface{} {
			return &logMsg{}
		},
	}
	lg.wg.Add(1)
	go lg.startLogger()
	return lg
}

func (lg *Logger) startLogger() {
	shutdown := false
LOOP:
	for {
		select {
		// 接收日志类消息
		case bm := <-lg.msgChan:
			lg.writeToLoggers(bm.when, bm.msg, bm.level)
			logMsgPool.Put(bm)

			// 接受指令类消息  flush 清空消息, 或者 close 关闭channel
		case sg := <-lg.signalChan:
			lg.flush()
			if sg == "close" {
				for _, l := range lg.outputs {
					l.Destroy()
				}
				lg.outputs = nil
				shutdown = true
			}
			// 这条是放在 close 中比较合适吗?
			lg.wg.Done()

			if shutdown {
				break LOOP
			}
		}
	}
}

// 清空堆积的消息
func (lg *Logger) flush() {
	if lg.asynchronous {
		for {
			if len(lg.msgChan) > 0 {
				bm := <-lg.msgChan
				lg.writeToLoggers(bm.when, bm.msg, bm.level)
				logMsgPool.Put(bm)
				continue
			}
			break
		}
	}
	for _, l := range lg.outputs {
		l.Flush()
	}
}

// SetLevel 设置日志等级
func (lg *Logger) SetLevel(l int) {
	lg.level = l
}

// SetFormatter 设置自定义的格式化方法
func (lg *Logger) SetFormatter(formatter FormaterFunc, force ...bool) error {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	isForce := append(force, false)[0]
	if isForce {
		lg.formatter = formatter
		return nil
	}

	if lg.formatter != nil {
		return fmt.Errorf("logs: 已经存在formatter, 不再设置")
	}
	lg.formatter = formatter
	return nil
}

// 配置日志适配器
func (lg *Logger) setLogger(adapterName string, configs ...string) error {
	// fmt.Printf("设置适配器: %s, config: %v \n", adapterName, configs)
	// 设置默认值
	config := append(configs, "{}")[0]
	for _, l := range lg.outputs {
		if l.Name == adapterName {
			return fmt.Errorf("logs: 重复设置适配器 %q , 只能配置一次", adapterName)
		}
	}

	adapterFunc, ok := adapters[adapterName]
	if !ok {
		return fmt.Errorf("logs: 找不到适配器的名字 %q, 是不是还没注册?", adapterName)
	}

	adapter := adapterFunc()
	err := adapter.Init(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "logs: Logger.SetLogger: "+err.Error())
		return err
	}
	lg.outputs = append(lg.outputs, &nameLogger{Name: adapterName, LoggerInterface: adapter})
	return nil
}

// SetLogger 对外暴露的设置接口
func (lg *Logger) SetLogger(adapterName string, config ...string) error {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	// fmt.Printf("SetLogger: %s, config: %v \n", adapterName, config)

	if !lg.init {
		lg.outputs = []*nameLogger{}
		lg.init = true
	}
	return lg.setLogger(adapterName, config...)
}

// DelLogger 删除适配器
func (lg *Logger) DelLogger(adapterName string) error {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	outputs := []*nameLogger{}
	for _, logger := range lg.outputs {
		if logger.Name == adapterName {
			logger.Destroy()
		} else {
			outputs = append(outputs, logger)
		}
	}

	if len(outputs) == len(lg.outputs) {
		return fmt.Errorf("logs: 删除了一个不存在的适配器 %q (注册了没?)", adapterName)
	}
	lg.outputs = outputs
	return nil
}

func (lg *Logger) writeToLoggers(when time.Time, msg string, level int) {
	// fmt.Printf("writeToLoggers: %s, outputs: %v \n", msg, lg.outputs)
	for _, l := range lg.outputs {
		// fmt.Printf("in outputs: %s \n", l.Name)
		err := l.WriteMsg(when, msg, level)
		if err != nil {
			fmt.Fprintf(os.Stdout, "logs: 写入 %v 出错, error: %v \n", l.Name, err)
		}
	}
}

func (lg *Logger) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	// 存在 \n 字符时截取它之前的
	if p[len(p)-1] == '\n' {
		p = p[0 : len(p)-1]
	}

	err = lg.writeMsg(levelLoggerImpl, string(p))
	if err != nil {
		return len(p), err
	}
	return 0, err
}

func (lg *Logger) writeMsg(logLevel int, msg string, v ...interface{}) error {
	if !lg.init {
		lg.lock.Lock()
		lg.setLogger(AdapterConsole)
		lg.lock.Unlock()
	}

	// 如果存在自定义格式的内容, 将直接
	if lg.formatter != nil {
		msg = lg.formatter(logLevel, msg, v...)
	} else {
		msg = lg.prefix + " " + msg
		if lg.enableFuncCallDepth {
			_, file, line, ok := runtime.Caller(lg.loggerFuncCallDepth)
			if !ok {
				file = "???"
				line = 0
			}
			_, filename := path.Split(file)
			msg = "[" + filename + ":" + strconv.Itoa(line) + "] " + msg
		}

		//set level info in front of filename info
		if logLevel == levelLoggerImpl {
			// set to emergency to ensure all log will be print out correctly
			logLevel = LevelEmergency
		} else {
			msg = levelPrefix[logLevel] + msg
		}
	}

	when := time.Now()

	// fmt.Printf("writeMsg: %s \n", msg)

	if lg.asynchronous {
		lm := logMsgPool.Get().(*logMsg)
		lm.level = logLevel
		lm.msg = msg
		lm.when = when
		lg.msgChan <- lm
	} else {
		lg.writeToLoggers(when, msg, logLevel)
	}

	return nil
}

// Reset 重启
func (lg *Logger) Reset() {
	lg.Flush()
	for _, l := range lg.outputs {
		l.Destroy()
	}
	lg.outputs = nil
}

// Flush 清空数据
func (lg *Logger) Flush() {
	if lg.asynchronous {
		lg.signalChan <- "flush"
		lg.wg.Wait()
		lg.wg.Add(1)
		return
	}
	lg.flush()
}

// Close 关闭
func (lg *Logger) Close() {
	if lg.asynchronous {
		lg.signalChan <- "close"
		lg.wg.Wait()
		close(lg.msgChan)
	} else {
		lg.flush()
		for _, l := range lg.outputs {
			l.Destroy()
		}
		lg.outputs = nil
	}
	close(lg.signalChan)
}

// Emergency 紧急
func (lg *Logger) Emergency(format string, v ...interface{}) {
	if LevelEmergency > lg.level {
		return
	}
	lg.writeMsg(LevelEmergency, format, v...)
}

// Alert 警告
func (lg *Logger) Alert(format string, v ...interface{}) {
	if LevelAlert > lg.level {
		return
	}
	lg.writeMsg(LevelAlert, format, v...)
}

// Critical 重要
func (lg *Logger) Critical(format string, v ...interface{}) {
	if LevelCritical > lg.level {
		return
	}
	lg.writeMsg(LevelCritical, format, v...)
}

// Error 错误
func (lg *Logger) Error(format string, v ...interface{}) {
	if LevelError > lg.level {
		return
	}
	lg.writeMsg(LevelError, format, v...)
}

// Warning 警告
func (lg *Logger) Warning(format string, v ...interface{}) {
	if LevelWarning > lg.level {
		return
	}
	lg.writeMsg(LevelWarning, format, v...)
}

// Warn 警告
func (lg *Logger) Warn(format string, v ...interface{}) {
	if LevelWarn > lg.level {
		return
	}
	lg.writeMsg(LevelWarn, format, v...)
}

// Notice 提醒
func (lg *Logger) Notice(format string, v ...interface{}) {
	if LevelNotice > lg.level {
		return
	}
	lg.writeMsg(LevelNotice, format, v...)
}

// Info 消息
func (lg *Logger) Info(format string, v ...interface{}) {
	if LevelInfo > lg.level {
		return
	}
	lg.writeMsg(LevelInfo, format, v...)
}

// Debug debug(这里将自动开启文件路径追踪)
func (lg *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug > lg.level {
		return
	}
	needReset := false
	if lg.enableFuncCallDepth == false {
		lg.enableFuncCallDepth = true
		needReset = true
	}

	lg.writeMsg(LevelDebug, format, v...)

	if needReset {
		lg.enableFuncCallDepth = false
	}
}
