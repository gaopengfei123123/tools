package logs

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"time"
)

// LoggerByID 带id的logger
type LoggerByID struct {
	Logger
	ID string
	IP string
}

// NewLoggerByID 生成一个新的
func NewLoggerByID(channelLens ...int64) *LoggerByID {
	lg := new(LoggerByID)
	lg.level = LevelDebug
	lg.loggerFuncCallDepth = 2
	lg.msgChanLen = append(channelLens, 0)[0]
	if lg.msgChanLen <= 0 {
		lg.msgChanLen = defaultAsyncMsgLen
	}
	lg.signalChan = make(chan string, 1)
	lg.setLogger(AdapterConsole)
	return lg
}

func (lg *LoggerByID) writeMsg(logLevel int, msg string, v ...interface{}) error {
	if !lg.init {
		lg.lock.Lock()
		lg.setLogger(AdapterConsole)
		lg.lock.Unlock()
	}
	// 格式: time [id][ip][time][event][category] body
	format := "[%s][%s][%v]%s[%s]"

	if lg.ID == "" {
		lg.ID = "--id--"
	}

	if lg.IP == "" {
		lg.IP = "--ip--"
	}

	category := "--category--"
	eventLevel := "--level--"

	if len(v) > 0 {
		if len(v) > 1 {
			category = v[0].(string)
			v = v[1:]
			msg = fmt.Sprintf(msg, v...)
		} else {
			category = v[0].(string)
		}
	}
	when := time.Now()

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
		eventLevel = levelPrefix[logLevel]
	}

	head := fmt.Sprintf(format, lg.ID, lg.IP, when.Unix(), eventLevel, category)
	msg = head + " " + msg

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

// Info Log INFO level message.
// compatibility alias for Informational()
func (lg *LoggerByID) Info(f interface{}, v ...interface{}) {
	format, v := interfaceToStr(f, v...)
	fmt.Println(format, v)
	if LevelInfo > lg.level {
		return
	}
	lg.writeMsg(LevelInfo, format, v...)
}

// Debug Log DEBUG level message.
func (lg *LoggerByID) Debug(f interface{}, v ...interface{}) {
	format, v := interfaceToStr(f, v...)
	if LevelDebug > lg.level {
		return
	}
	lg.writeMsg(LevelDebug, format, v...)
}

// Error Log ERROR level message.
func (lg *LoggerByID) Error(f interface{}, v ...interface{}) {
	format, v := interfaceToStr(f, v...)
	if LevelError > lg.level {
		return
	}
	lg.writeMsg(LevelError, format, v...)
}
