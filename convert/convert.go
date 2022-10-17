package convert

import (
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
	var arr []int32
	strArr := strings.Split(s.String(), ",")
	for _, ID := range strArr {
		n, _ := strconv.Atoi(ID)
		arr = append(arr, int32(n))
	}
	return arr
}
