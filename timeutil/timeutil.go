package timeutil

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	cst *time.Location
)

// CSTLayout China standard timezone layout
const CSTLayout = "2006-01-02 15:04:05"
const CSTLayoutDate = "2006-01-02"
const CSTLayoutDate2 = "2006-1-02"
const CSTLayoutDate3 = "2006-1-2"
const CSTLayoutDate4 = "20060102"

func init() {
	var err error
	if cst, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}

	// the default value is China timezone
	time.Local = cst
}

// RFC3339ToCST convert rfc3339 value to China standard timezone layout
// convert 2020-11-08T08:18:46+08:00 to 2020-11-08 08:18:46
func RFC3339ToCST(value string) (string, error) {
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", err
	}

	return ts.In(cst).Format(CSTLayout), nil
}

// StrToTime 将日期转换成time 对象, 尝试解析多个日期模板
func StrToTime(value string, loc *time.Location) (tt time.Time, err error) {
	tplList := []string{CSTLayout, CSTLayoutDate, CSTLayoutDate2, CSTLayoutDate3, CSTLayoutDate4}
	for _, curTpl := range tplList {
		tt, err = time.ParseInLocation(curTpl, value, loc)
		if err == nil {
			return
		}
	}
	return
}

// RFC3339ToUnix convert rfc3339 value to unix
// convert 2020-11-08T08:18:46+08:00 to 1579871471
func RFC3339ToUnix(value string) (int64, error) {
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return 0, err
	}

	return ts.In(cst).Unix(), nil
}

// CSTCurrentString formatting time
// returns a style value like "2006-01-02 15:04:05"
func CSTCurrentString() string {
	ts := time.Now()
	return ts.In(cst).Format(CSTLayout)
}

// ParseCSTInLocation formatting time
// 2022-03-24 12:33:33 to time.Time类型
func ParseCSTInLocation(datetime string) (time.Time, error) {
	return time.ParseInLocation(CSTLayout, datetime, cst)
}

// ParseCSTDateInLocation formatting date
// 2022-03-24 to time.Time类型
func ParseCSTDateInLocation(date string) (time.Time, error) {
	return time.ParseInLocation(CSTLayoutDate, date, cst)
}

// CSTStringToUnix return timestamp
// convert 2020-01-24 21:11:11 to 1579871471
func CSTStringToUnix(cstLayoutString string) (int64, error) {
	stamp, err := time.ParseInLocation(CSTLayout, cstLayoutString, cst)
	if err != nil {
		return 0, err
	}
	return stamp.Unix(), nil
}

// GMTString formatting time
// returns a style value like "Mon, 02 Jan 2006 15:04:05 GMT"
func GMTString() string {
	return time.Now().In(cst).Format(http.TimeFormat)
}

// ParseGMTInLocation formatting time
func ParseGMTInLocation(date string) (time.Time, error) {
	return time.ParseInLocation(http.TimeFormat, date, cst)
}

// SubInLocation calculate time diff ormd on location
func SubInLocation(ts time.Time) float64 {
	return math.Abs(time.Now().In(cst).Sub(ts).Seconds())
}

func DurationToSecond(str string) (int64, error) {
	var durationStr string
	if str == "" {
		durationStr = "0s"
	}
	tokens := strings.Split(str, ":")
	tokenCount := len(tokens)
	if tokenCount <= 0 || tokenCount > 3 {
		durationStr = "0s"
	} else if tokenCount == 1 {
		durationStr = fmt.Sprintf("%ss", tokens[0])
	} else if tokenCount == 2 {
		durationStr = fmt.Sprintf("%sm%ss", tokens[0], tokens[1])
	} else if tokenCount == 3 {
		durationStr = fmt.Sprintf("%sh%sm%ss", tokens[0], tokens[1], tokens[2])
	}
	parseDuration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, err
	}
	return int64(parseDuration.Seconds()), nil
}

func GToTime(value interface{}) (time.Time, error) {
	return ToTime(value, "")
}

func ToTime(value interface{}, format string) (time.Time, error) {
	var sTime string
	switch value.(type) {
	case string:
		if format == "" {
			return guessTime(value.(string))
		}
		return time.Parse(format, value.(string))
	case int8:
		return time.Unix(int64(value.(int8)), 0), nil
	case int:
		sTime = strconv.Itoa(value.(int))
		break
	case int32:
		sTime = string(value.(int32))
		break
	case int64:
		sTime = strconv.FormatInt(value.(int64), 10)
		break
	case time.Time:
		return value.(time.Time), nil
	}
	if sTime != "" {
		sSec := beego.Substr(sTime, 0, 10)
		sNsec := (sTime + strings.Repeat("0", 19))[10:19]
		sec, err := strconv.ParseInt(sSec, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		nSec, err := strconv.ParseInt(sNsec, 10, 64)
		if err != nil {
			nSec = 0
		}
		return time.Unix(sec, nSec), nil
	}
	return time.Time{}, errors.New("无法解析时间")
}

type byLen struct {
	bv   string
	strA []string
}

func (bl byLen) Len() int {
	return len(bl.strA)
}

func (bl byLen) Swap(i, j int) {
	bl.strA[i], bl.strA[j] = bl.strA[j], bl.strA[i]
}

func (bl byLen) Less(i, j int) bool {
	return math.Abs(float64(len(bl.bv)-len(bl.strA[i]))) < math.Abs(float64(len(bl.bv)-len(bl.strA[j])))
}

func guessTime(dateString string) (time.Time, error) {
	config := []string{
		"2006-01-02 15:04:05", "2006-01-02", "2006-01-02 15:04", "2006-01-02 15", "2006-01-02 15.04.05", "2006-01-02 15.04",
		"2006/01/02 15:04:05", "2006/01/02", "2006/01/02 15:04", "2006/01/02 15", "2006/01/02 15.04.05", "2006/01/02 15.04",
		"01/02/06", "15:04:05", "15:04", "04:05", "15.04", "04.05", "01-02", "01/02", "02-01", "02/01", time.Layout, time.ANSIC,
		time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339,
		time.RFC3339Nano, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano,
	}
	b := byLen{
		bv:   dateString,
		strA: config,
	}
	sort.Sort(b)
	for _, format := range b.strA {
		t, e := time.Parse(format, dateString)
		if e == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("不支持的时间格式类型")
}

// TimeCnt 时间统计函数
func TimeCnt(msgTag string) func() (msg string, cnt time.Duration, formatMsg string) {
	start := time.Now()
	return func() (msg string, cnt time.Duration, formatMsg string) {
		elapsed := time.Now().Sub(start)
		return msgTag, elapsed, fmt.Sprintf("%s: 该函数执行完成耗时：%s", msgTag, elapsed)
	}
}
