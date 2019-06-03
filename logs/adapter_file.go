package logs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 文件类型写入类
type fileLogWriter struct {
	sync.RWMutex
	Filename string
	// fileWriter *os.File

	// Rotate at line
	MaxLines         int `json:"maxlines"`
	maxLinesCurLines int

	MaxFiles         int `json:"maxfiles"`
	MaxFilesCurFiles int

	// Rotate at size
	MaxSize        int `json:"maxsize"`
	maxSizeCurSize int

	// Rotate daily
	Daily         bool  `json:"daily"`
	MaxDays       int64 `json:"maxdays"`
	dailyOpenDate int
	dailyOpenTime time.Time

	// Rotate hourly
	Hourly         bool  `json:"hourly"`
	MaxHours       int64 `json:"maxhours"`
	hourlyOpenDate int
	hourlyOpenTime time.Time

	Rotate bool `json:"rotate"`

	Level int `json:"level"`

	Perm string `json:"perm"`

	RotatePerm string `json:"rotateperm"`

	fileNameOnly, suffix string // like "project.log", project is fileNameOnly and .log is suffix
}

// 文件句柄池, 如果出现多个 Newlogger 初始化出来的日志对象, 将共享一个文件写入
// 这样做是避免在类似for循环这种情景下
// TODO: 修复一下读写不安全的问题
type filePool struct {
	pool map[string]*fileStruct
}

type fileStruct struct {
	sync.Mutex
	Fd *os.File
}

var fp *filePool

func init() {
	// 初始化文件池
	fp = &filePool{
		pool: make(map[string]*fileStruct, 1),
	}

	Register(AdapterFile, NewFileWriter)
}

func (fp *filePool) Set(filename string, fd *os.File) error {
	_, ok := fp.pool[filename]
	if ok {
		fd.Close()
		return nil
	}
	fp.pool[filename] = &fileStruct{
		Fd: fd,
	}
	return nil
}

func (fp *filePool) Get(filename string) (*fileStruct, error) {
	fb, ok := fp.pool[filename]
	if !ok {
		return fb, errors.New("不存在的句柄id")
	}
	return fb, nil
}
func (fp *filePool) Del(filename string) error {
	file, ok := fp.pool[filename]
	if !ok {
		return errors.New("不存在的句柄id")
	}
	file.Lock()
	defer file.Unlock()
	err := file.Fd.Close()
	if err != nil {
		return err
	}
	delete(fp.pool, filename)
	return nil
}

// NewFileWriter 初始化一个日志类
func NewFileWriter() LoggerInterface {
	w := &fileLogWriter{
		Daily:      true,
		MaxDays:    60,
		Hourly:     false,
		MaxHours:   168,
		Rotate:     true,
		RotatePerm: "0444",
		Level:      LevelTrace,
		Perm:       "0744",
		MaxLines:   10000000,
		MaxFiles:   999,
		MaxSize:    1 << 28,
	}
	return w
}

func (w *fileLogWriter) Init(jsonConfig string) error {
	err := json.Unmarshal([]byte(jsonConfig), w)

	if err != nil {
		return err
	}

	if len(w.Filename) == 0 {
		return errors.New("请输入 filename 的值")
	}
	w.suffix = filepath.Ext(w.Filename)
	w.fileNameOnly = strings.TrimSuffix(w.Filename, w.suffix)

	return w.startLogger()
}

// logger 初始化的前期准备
func (w *fileLogWriter) startLogger() error {
	return w.createLogFile()
}

func (w *fileLogWriter) createLogFile() error {
	perm, err := strconv.ParseInt(w.Perm, 8, 64)
	if err != nil {
		return err
	}

	dirpath := path.Dir(w.Filename)
	err = os.MkdirAll(dirpath, os.FileMode(perm))
	if err != nil {
		return err
	}

	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	return fp.Set(w.Filename, fd)
}

func (w *fileLogWriter) WriteMsg(when time.Time, msg string, level int) error {
	if level > w.Level {
		return nil
	}
	hd, _, _ := formatTimeHeader(when)
	msg = string(hd) + msg + "\n"

	lf, err := fp.Get(w.Filename)
	if err != nil {
		return err
	}

	lf.Lock()
	lf.Unlock()

	_, err = lf.Fd.Write([]byte(msg))
	if err == nil {
		w.maxLinesCurLines++
		w.maxSizeCurSize += len(msg)
	}
	return err
}

func (w *fileLogWriter) Destroy() {
	fp.Del(w.Filename)
}
func (w *fileLogWriter) Flush() {

}
