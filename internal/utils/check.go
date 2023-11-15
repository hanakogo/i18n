package utils

import (
	"github.com/gookit/goutil/strutil"
)

// CheckYaml check extension of filename
func CheckYaml(filename string) bool {
	return strutil.IsEndOf(filename, ".yaml") || strutil.IsEndOf(filename, ".yml")
}
