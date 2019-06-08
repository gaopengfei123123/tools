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

	// // Rotate hourly 不对小时的进行转换了, 费劲
	// Hourly         bool  `json:"hourly"`
	// MaxHours       int64 `json:"maxhours"`
	// hourlyOpenDate int
	// hourlyOpenTime time.Time

	Rotate bool `json:"rotate"`

	Level int `json:"level"`

	Perm string `json:"perm"`

	RotatePerm string `json:"rotateperm"` // 备份文件的权限

	fileNameOnly, filePath, suffix string // like "project.log", project is fileNameOnly and .log is suffix
}

// 文件句柄池, 如果出现多个 Newlogger 初始化出来的日志对象, 将共享一个文件写入
// 这样做是避免在类似for循环这种情景下
// TODO: 修复一下读写不安全的问题
type filePool struct {
	pool map[string]*fileStruct
}

type fileStruct struct {
	sync.Mutex
	Fd        *os.File
	FileSize  int
	FileLine  int
	FileIndex int
	Date      time.Time

	MaxSize int
	MaxLine int
	MaxDays int64

	Perm                           string
	RotatePerm                     string // 备份文件的权限
	Filename                       string
	fileNameOnly, filePath, suffix string // like "project.log", project is fileNameOnly and .log is suffix
}

var fp *filePool

func init() {
	// 初始化文件池
	fp = &filePool{
		pool: make(map[string]*fileStruct, 1),
	}

	Register(AdapterFile, NewFileWriter)
}

// Update 更新文件信息
func (fs *fileStruct) Update(when time.Time, size int) {
	fs.FileSize += size
	fs.FileLine++
	if fs.Date.Nanosecond() < when.Nanosecond() {
		fs.Date = when
	}
}

func (fs *fileStruct) CreateFile() error {
	perm, err := strconv.ParseInt(fs.Perm, 8, 64)
	if err != nil {
		return err
	}

	dirpath := path.Dir(fs.Filename)
	err = os.MkdirAll(dirpath, os.FileMode(perm))
	if err != nil {
		return err
	}

	fd, err := os.OpenFile(fs.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err != nil {
		return err
	}
	fs.Fd = fd
	return nil
}

// FlashFile 根据文件现在状态备份,并重新生成个文件句柄
func (fs *fileStruct) FlashFile() error {
	oldName := fs.Fd.Name()
	_, y2, m2, d2 := formatDateHeader(fs.Date)
	filename := fmt.Sprintf("%s.%s%s%s", oldName, y2, m2, d2)
	// fmt.Printf("fileName: %s  %v \n", filename, fs.FileIndex)
	if fs.FileIndex > 0 {
		filename = fmt.Sprintf("%s.%d", filename, fs.FileIndex)
	}
	// fmt.Printf("flashFile: %v \n", filename)
	fs.Fd.Close()
	err := os.Rename(oldName, filename)
	if err != nil {
		return err
	}
	return fs.CreateFile()
}

// DoRotate 日志转换的判断
func (fs *fileStruct) DoRotate(when time.Time) error {
	currentDate, _, _, _ := formatDateHeader(when)
	fileDate, _, _, _ := formatDateHeader(fs.Date)
	// fmt.Printf("currentDate: %s \n", currentDate)
	// fmt.Printf("fileDate: %s \n", fileDate)

	// 日期不对的情况, 现有文件变成备份, 重新生成一个新文件
	if string(currentDate) != string(fileDate) {
		err := fs.FlashFile()
		fs.FileIndex = 0
		fs.FileSize = 0
		fs.FileLine = 0
		return err
	}

	// fmt.Printf("fileMaxSize: %v, %v \n", fs.MaxSize, fs.FileSize)

	if fs.MaxSize <= fs.FileSize {
		fs.FileIndex++
		fs.FileSize = 0
		return fs.FlashFile()
	}
	// fmt.Printf("fileMaxLine: %v, %v \n", fs.MaxLine, fs.FileLine)

	if fs.MaxLine <= fs.FileLine {
		fs.FileIndex++
		fs.FileLine = 0
		return fs.FlashFile()
	}
	return nil
}

func (fp *filePool) Set(fs *fileStruct) {
	filename := fs.Filename
	_, ok := fp.pool[filename]
	if ok {
		return
	}
	fp.pool[filename] = fs
	return
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
		Daily:   true,
		MaxDays: 60,

		Rotate:     true,
		RotatePerm: "0544",
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
	w.filePath = ""
	strIndex := strings.LastIndex(w.fileNameOnly, "/")
	// 文件名中包含路径就拆分一下
	if strIndex != -1 {
		filename := w.fileNameOnly
		// 只获取文件名, 除去后缀
		w.fileNameOnly = filename[strIndex+1:]
		w.filePath = filename[0 : strIndex+1]
	}

	return w.startLogger()
}

// logger 初始化的前期准备
func (w *fileLogWriter) startLogger() error {
	fs := &fileStruct{
		Filename:     w.Filename,
		FileSize:     0,
		FileLine:     0,
		FileIndex:    0,
		Date:         time.Now(),
		MaxSize:      w.MaxSize,
		MaxLine:      w.MaxLines,
		MaxDays:      w.MaxDays,
		Perm:         w.Perm,
		RotatePerm:   w.RotatePerm,
		fileNameOnly: w.fileNameOnly,
		filePath:     w.filePath,
		suffix:       w.suffix,
	}

	err := fs.CreateFile()
	if err != nil {
		return err
	}
	fp.Set(fs)
	return nil
}

func (w *fileLogWriter) WriteMsg(when time.Time, msg string, level int) error {
	if level > w.Level {
		return nil
	}
	hd, _, _ := formatTimeHeader(when)
	msg = string(hd) + msg + "\n"

	fs, err := fp.Get(w.Filename)
	if err != nil {
		return err
	}

	fs.Lock()
	defer fs.Unlock()

	err = fs.DoRotate(when)
	if err != nil {
		return err
	}

	_, err = fs.Fd.Write([]byte(msg))
	if err == nil {
		w.maxLinesCurLines++
		w.maxSizeCurSize += len(msg)
		fs.Update(when, len(msg))
	}
	return err
}

func (w *fileLogWriter) Destroy() {
	fp.Del(w.Filename)
}
func (w *fileLogWriter) Flush() {

}
