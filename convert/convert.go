package convert

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	return strconv.Atoi(s.String())
}

func (s StrTo) Int8() (int8, error) {
	n, err := strconv.Atoi(s.String())
	return int8(n), err
}

func (s StrTo) Int32() (int32, error) {
	n, err := strconv.Atoi(s.String())
	return int32(n), err
}

func (s StrTo) UInt32() (uint32, error) {
	n, err := strconv.Atoi(s.String())
	return uint32(n), err
}

func (s StrTo) Int64() (int64, error) {
	n, err := strconv.Atoi(s.String())
	return int64(n), err
}

func (s StrTo) UInt64() (uint64, error) {
	n, err := strconv.Atoi(s.String())
	return uint64(n), err
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) MustInt8() int8 {
	v, _ := s.Int8()
	return v
}

func (s StrTo) MustInt32() int32 {
	v, _ := s.Int32()
	return v
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

func (s StrTo) MustInt64() int64 {
	v, _ := s.Int64()
	return v
}

func (s StrTo) MustUInt64() uint64 {
	v, _ := s.UInt64()
	return v
}

func (s StrTo) ToIntArr() []int32 {
	arr := make([]int32, 0, 1)
	strArr := strings.Split(s.String(), ",")
	for _, ID := range strArr {
		if ID == "" {
			continue
		}
		n, _ := strconv.Atoi(ID)
		arr = append(arr, int32(n))
	}
	return arr
}

// Div 把十进制数字转成 ABCD
func Div(Num int) string {
	var (
		Str  string = ""
		k    int
		temp []int //保存转化后每一位数据的值，然后通过索引的方式匹配A-Z
	)
	//用来匹配的字符A-Z
	Slice := []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	if Num > 26 { //数据大于26需要进行拆分
		for {
			k = Num % 26 //从个位开始拆分，如果求余为0，说明末尾为26，也就是Z，如果是转化为26进制数，则末尾是可以为0的，这里必须为A-Z中的一个
			if k == 0 {
				temp = append(temp, 26)
				k = 26
			} else {
				temp = append(temp, k)
			}
			Num = (Num - k) / 26 //减去Num最后一位数的值，因为已经记录在temp中
			if Num <= 26 {       //小于等于26直接进行匹配，不需要进行数据拆分
				temp = append(temp, Num)
				break
			}
		}
	} else {
		return Slice[Num]
	}

	for _, value := range temp {
		Str = Slice[value] + Str //因为数据切分后存储顺序是反的，所以Str要放在后面
	}
	return Str
}

// ToMap 结构体转为Map[string]interface{}
func ToMap(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}
