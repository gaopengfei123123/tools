package logs

import (
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

type logWriter struct {
	sync.Mutex
	writer io.Writer
}

const (
	y1  = `0123456789`
	y2  = `0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789`
	y3  = `0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999`
	y4  = `0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789`
	mo1 = `000000000111`
	mo2 = `123456789012`
	d1  = `0000000001111111111222222222233`
	d2  = `1234567890123456789012345678901`
	h1  = `000000000011111111112222`
	h2  = `012345678901234567890123`
	mi1 = `000000000011111111112222222222333333333344444444445555555555`
	mi2 = `012345678901234567890123456789012345678901234567890123456789`
	s1  = `000000000011111111112222222222333333333344444444445555555555`
	s2  = `012345678901234567890123456789012345678901234567890123456789`
	ns1 = `0123456789`
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})

	w32Green   = string([]byte{27, 91, 52, 50, 109})
	w32White   = string([]byte{27, 91, 52, 55, 109})
	w32Yellow  = string([]byte{27, 91, 52, 51, 109})
	w32Red     = string([]byte{27, 91, 52, 49, 109})
	w32Blue    = string([]byte{27, 91, 52, 52, 109})
	w32Magenta = string([]byte{27, 91, 52, 53, 109})
	w32Cyan    = string([]byte{27, 91, 52, 54, 109})

	reset = string([]byte{27, 91, 48, 109})
)

func newLogWriter(wr io.Writer) *logWriter {
	return &logWriter{writer: wr}
}

func (lg *logWriter) println(when time.Time, msg string) {
	lg.Lock()
	defer lg.Unlock()

	h, _, _ := formatTimeHeader(when)
	lg.writer.Write(append(append(h, msg...), '\n'))
}

// yyyy-mm-dd hh:ii:ss
func formatTimeHeader(when time.Time) ([]byte, int, int) {
	y, mo, d := when.Date()
	h, mi, s := when.Clock()
	// ns := when.Nanosecond() / 1000000

	var buf [20]byte
	buf[0] = y1[y/1000%10]
	buf[1] = y2[y/100]
	buf[2] = y3[y-y/100*100]
	buf[3] = y4[y-y/100*100]
	buf[4] = '-'
	buf[5] = mo1[mo-1]
	buf[6] = mo2[mo-1]
	buf[7] = '-'
	buf[8] = d1[d-1]
	buf[9] = d2[d-1]
	buf[10] = ' '
	buf[11] = h1[h]
	buf[12] = h2[h]
	buf[13] = ':'
	buf[14] = mi1[mi]
	buf[15] = mi2[mi]
	buf[16] = ':'
	buf[17] = s1[s]
	buf[18] = s2[s]
	buf[19] = ' '

	// buf[19] = '.'
	// buf[20] = ns1[ns/100]
	// buf[21] = ns1[ns%100/10]
	// buf[22] = ns1[ns%10]

	// buf[23] = ' '

	return buf[0:], d, h
}

// yyyy-mm-dd
func formatDateHeader(when time.Time) (res []byte, year, month, day string) {
	y, mo, d := when.Date()

	var buf [10]byte
	buf[0] = y1[y/1000%10]
	buf[1] = y2[y/100]
	buf[2] = y3[y-y/100*100]
	buf[3] = y4[y-y/100*100]
	buf[4] = '-'
	buf[5] = mo1[mo-1]
	buf[6] = mo2[mo-1]
	buf[7] = '-'
	buf[8] = d1[d-1]
	buf[9] = d2[d-1]
	return buf[0:], string(buf[0:4]), string(buf[5:7]), string(buf[8:])
}

// 获取指定目录下的文件列表
func listFile(myfolder string) []string {
	fileList := []string{}
	var fileName string

	files, _ := ioutil.ReadDir(myfolder)
	for _, file := range files {
		if file.IsDir() {
			dirFile := listFile(myfolder + "/" + file.Name())
			fileList = append(fileList, dirFile...)
		} else {
			fileName = fmt.Sprintln(myfolder + "/" + file.Name())
			fileList = append(fileList, fileName)
		}
	}
	return fileList
}
