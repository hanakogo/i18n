package utils

import (
	"fmt"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
	"reflect"
	"strings"
)

// WalkStringMap deep walk map[string]any
func WalkStringMap(stringMap map[string]any, walkFun func(value any, path []string), path ...string) {
	for key, value := range stringMap {
		if value, ok := value.(map[string]any); ok {
			WalkStringMap(value, walkFun, append(path, key)...)
			continue
		}
		walkFun(value, append(path, key))
	}
}

// MergeStringMap deep merge map[string]any
func MergeStringMap(dst *map[string]any, src map[string]any, lang string) (err error) {
	WalkStringMap(src, func(value any, path []string) {
		pathStr := strutil.Join(".", path...)

		// check override of data
		//if value, ok := maputil.GetByPath(pathStr, *dst); ok {
		//	fmt.Printf("WARN: overriding existed data of language<%s>: {path: %s, value: %v}\n", lang, pathStr, value)
		//}

		err = maputil.SetByPath(dst, pathStr, value)
		if err != nil {
			return
		}
	})
	return
}

// TakeStringMap take out the value from map[string]any
func TakeStringMap(src *map[string]any, key string) (value any) {
	var (
		sliceIdx      int
		sliceIdxValid bool
	)
	if lastOpenBracketIdx := strings.LastIndex(key, "["); strutil.IsEndOf(key, "]") && lastOpenBracketIdx > 0 {
		idx := key[lastOpenBracketIdx+1 : len(key)-1]
		key = key[:lastOpenBracketIdx]

		err := fmt.Errorf("invalid index of key <%s[%s]>", key, idx)
		if idx == "" {
			return err
		}

		sliceIdxInt, errOfToInt := mathutil.ToInt(idx)
		if errOfToInt != nil {
			return err
		}

		sliceIdx = sliceIdxInt
		sliceIdxValid = true
	}

	value = (*src)[key]
	if value == nil {
		return nil
	}

	valueType := reflect.TypeOf(value)
	if sliceIdxValid && valueType.Kind() == reflect.Slice {
		if sliceVal, ok := value.([]any); ok {
			if sliceIdx >= len(sliceVal) || sliceIdx < 0 {
				panic(fmt.Errorf(`index %d is outbound of slice "%s"`, sliceIdx, key))
			}
			return sliceVal[sliceIdx]
		}
		return nil
	}

	return value
}
