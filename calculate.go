package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"strconv"
)

func CombineFieldInt32(old interface{}, ne interface{}) (int32, error) {
	res := old.(int32) + ne.(int32)
	return res, nil
}

func CombineFieldFloat32(old interface{}, ne interface{}) (float32, error) {
	tt := decimal.NewFromFloat32(old.(float32)).Add(decimal.NewFromFloat32(ne.(float32)))
	re, exact := tt.Float64()
	if exact {
		return float32(re), nil
	}
	return 0, errors.New("CombineFieldFloat32 exec error")
}

func DivideInt32(dividend int32, divisor int32, mul ...int32) (float64, error) {
	if divisor == 0 {
		return 0, nil
	}

	multi := int32(1)
	if len(mul) != 0 {
		multi = mul[0]
	}

	return strconv.ParseFloat(fmt.Sprintf("%.2f", float64(dividend)/float64(divisor)*float64(multi)), 64)
}
