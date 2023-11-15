package i18n

import (
	"github.com/hanakogo/i18n/internal/errors"
	"github.com/hanakogo/i18n/internal/status"
	"github.com/hanakogo/i18n/internal/structs"
)

var i18nFS *structs.I18nFS
var (
	DefaultLang  string
	FallbackLang string
)

func Init(opts Opts) (err error) {
	if status.Initialized {
		return errors.ErrorAlreadyInitialized
	}

	err = initI18nFS(opts)
	if err != nil {
		return err
	}
	status.Initialized = true

	err = initLoadLanguages(opts.Languages)
	if err != nil {
		return err
	}

	err = initSetDefault(opts.DefaultLang, opts.FallbackLang)
	if err != nil {
		return err
	}

	return nil
}

func initI18nFS(opts Opts) (err error) {
	i18nFS, err = structs.NewI18nFS(
		opts.FSOpts.FSMode,
		opts.FSOpts.Prefix,
		opts.FSOpts.EmbedFS,
	)
	return
}

func initLoadLanguages(languages []string) error {
	for _, language := range languages {
		err := Load(language)
		if err != nil {
			return err
		}
	}
	return nil
}

func initSetDefault(defLang string, fbLang string) (err error) {
	if Has(fbLang) {
		FallbackLang = fbLang
	} else {
		return errors.GetSpecificTypeLangNotFound("fallback", fbLang)
	}

	if Has(defLang) {
		DefaultLang = defLang
	} else {
		return errors.GetSpecificTypeLangNotFound("default", fbLang)
	}
	return
}

func Initialized() bool {
	return status.Initialized
}

func Reset() {
	if !status.Initialized {
		return
	}
	status.Initialized = false
	i18nFS = nil
}
