package status

import (
	"github.com/hanakogo/i18n/internal/errors"
)

var (
	Initialized = false
)

func checkInitialized() error {
	if !Initialized {
		return errors.ErrorNotInitialized
	}
	return nil
}

func MustInitialized() {
	if err := checkInitialized(); err != nil {
		panic(err)
	}
}
