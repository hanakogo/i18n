package i18n

import (
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/mathutil"
)

type ConvertFunc[T any] func(value any) T

var (
	ConvertString ConvertFunc[string]  = ConvertAnyToString
	ConvertInt64  ConvertFunc[int64]   = ConvertAnyToInt64
	ConvertFloat  ConvertFunc[float64] = ConvertAnyToFloat
)

func ConvertAnyToString(value any) string {
	valString, err := goutil.ToString(value)
	if err != nil {
		return DefaultString
	}
	return valString
}

func ConvertAnyToInt64(value any) int64 {
	valInt, err := goutil.ToInt64(value)
	if err != nil {
		return DefaultInt
	}
	return valInt
}

func ConvertAnyToFloat(value any) float64 {
	valFloat, err := mathutil.ToFloat(value)
	if err != nil {
		return DefaultFloat
	}
	return valFloat
}
