package utils

import (
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
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
