package i18n

import (
	"fmt"
	"github.com/hanakogo/i18n/internal/errors"
	"github.com/hanakogo/i18n/internal/status"
	"reflect"
	"slices"
)

func Has(lang string) bool {
	status.MustInitialized()

	return i18nFS.HasLang(lang)
}

func Load(lang string) (err error) {
	if !status.Initialized {
		return errors.ErrorNotInitialized
	}
	if !i18nFS.IsLangExists(lang) {
		return errors.GetLangNotExists(lang)
	}

	err = i18nFS.Read(lang)

	return
}

func HasPath(path string, languages ...string) (ok bool, contains []string) {
	status.MustInitialized()

	if len(languages) == 0 {
		languages = i18nFS.GetLanguages()
	}

	for _, language := range languages {
		value := getAnyOfLang(language, path)
		//goland:noinspection GoTypeAssertionOnErrors
		if _, ok := value.(error); ok {
			continue
		}
		ok = true
		contains = append(contains, language)
	}

	return
}

func getAnyOfLang(lang string, path string) (value any) {
	status.MustInitialized()

	value, err := i18nFS.GetValByPath(lang, path)
	if err != nil {
		return err
	}

	return
}

func GetValueTr(lang string, path string, def ...any) (val any, err error) {
	if len(def) == 0 {
		err = fmt.Errorf("must provide default value")
		return
	}

	defRes := def[0]
	value := getAnyOfLang(lang, path)
	//goland:noinspection GoTypeAssertionOnErrors
	if _, ok := value.(error); ok {
		return defRes, nil
	}

	return value, nil
}

func GetTr[T any](lang string, path string, convertFunc ConvertFunc[T], defRes T) (val T) {
	value := getAnyOfLang(lang, path)
	//goland:noinspection GoTypeAssertionOnErrors
	if _, ok := value.(error); ok {
		return defRes
	}
	if value == nil {
		return defRes
	}

	return convertFunc(value)
}

func GetStringTr(lang string, path string, def ...string) string {
	defRes := requireDefault[string](DefaultString, def...)
	value := getAnyOfLang(lang, path)
	//goland:noinspection GoTypeAssertionOnErrors
	if _, ok := value.(error); ok {
		return defRes
	}

	if stringVal := ConvertAnyToString(value); stringVal != DefaultString {
		return stringVal
	}

	return defRes
}

func GetStringTrF(lang string, path string, args ...any) (stringVal string) {
	stringVal = GetStringTr(lang, path, DefaultString)
	if stringVal == DefaultString {
		return
	}

	return fmt.Sprintf(stringVal, args...)
}

func GetInt64Tr(lang string, path string, def ...int64) int64 {
	defRes := requireDefault[int64](DefaultInt, def...)
	value := getAnyOfLang(lang, path)
	//goland:noinspection GoTypeAssertionOnErrors
	if _, ok := value.(error); ok {
		return defRes
	}

	if intVal := ConvertAnyToInt64(value); intVal != DefaultInt {
		return intVal
	}

	return defRes
}

func GetFloatTr(lang string, path string, def ...float64) float64 {
	defRes := requireDefault[float64](DefaultFloat, def...)
	value := getAnyOfLang(lang, path)
	//goland:noinspection GoTypeAssertionOnErrors
	if _, ok := value.(error); ok {
		return defRes
	}

	if floatVal := ConvertAnyToFloat(value); floatVal != DefaultFloat {
		return floatVal
	}

	return defRes
}

func GetSliceTr[T comparable](lang string, path string, convertFunc ConvertFunc[T], def ...[]T) []T {
	defRes := requireDefault[[]T]([]T{}, def...)
	value := getAnyOfLang(lang, path)
	//goland:noinspection GoTypeAssertionOnErrors
	if _, ok := value.(error); ok {
		return defRes
	}

	if valueType := reflect.TypeOf(value); valueType.Kind() == reflect.Slice {
		var resSlice []T
		for _, elem := range value.([]any) {
			resSlice = append(resSlice, convertFunc(elem))
		}
		return resSlice
	}

	return defRes
}

func GetValue(path string, def ...any) (val any, err error) {
	val, err = GetValueTr(DefaultLang, path, def...)
	if err != nil {
		return
	}
	defRes := def[0]

	// fallback
	if defRes == val {
		val, err = GetValueTr(FallbackLang, path, def...)
		if err != nil {
			return defRes, nil
		}
	}

	return
}

func Get[T comparable](path string, convertFunc ConvertFunc[T], def T) (val T) {
	value := GetTr[T](DefaultLang, path, convertFunc, def)
	if value == def {
		value = GetTr[T](FallbackLang, path, convertFunc, def)
		if value == def {
			return def
		}
	}

	return value
}

func GetString(path string, def ...string) (stringVal string) {
	stringVal = GetStringTr(DefaultLang, path, def...)

	// fallback
	if stringVal == requireDefault[string](DefaultString, def...) {
		stringVal = GetStringTr(FallbackLang, path, def...)
	}

	return
}

func GetStringF(path string, args ...any) (stringVal string) {
	stringVal = GetString(path, DefaultString)
	if stringVal == DefaultString {
		return
	}

	return fmt.Sprintf(stringVal, args...)
}

func GetInt64(path string, def ...int64) (intVal int64) {
	intVal = GetInt64Tr(DefaultLang, path, def...)

	// fallback
	if intVal == requireDefault[int64](DefaultInt, def...) {
		intVal = GetInt64Tr(FallbackLang, path, def...)
	}

	return
}

func GetFloat(path string, def ...float64) (floatVal float64) {
	floatVal = GetFloatTr(DefaultLang, path, def...)

	// fallback
	if floatVal == requireDefault[float64](DefaultFloat, def...) {
		floatVal = GetFloatTr(FallbackLang, path, def...)
	}

	return
}

func GetSlice[T comparable](path string, convertFunc ConvertFunc[T], def ...[]T) (valueList []T) {
	valueList = GetSliceTr[T](DefaultLang, path, convertFunc, def...)

	// fallback
	if slices.Equal(valueList, requireDefault[[]T]([]T{}, def...)) {
		valueList = GetSliceTr[T](FallbackLang, path, convertFunc, def...)
	}

	return valueList
}
