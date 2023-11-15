package errors

import "fmt"

var (
	ErrorAlreadyInitialized = fmt.Errorf("i18n has already been initialized")
	ErrorNotInitialized     = fmt.Errorf("i18n not Initialized")
)

func GetLangNotFound(lang string) error {
	return fmt.Errorf("language [%s] is not found", lang)
}

func GetLangNotExists(lang string) error {
	return fmt.Errorf("language [%s] is not exists", lang)
}

func GetSpecificTypeLangNotFound(typ string, lang string) error {
	return fmt.Errorf("%s language [%s] is not found", typ, lang)
}
