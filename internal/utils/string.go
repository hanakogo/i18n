package utils

import (
	"fmt"
	"github.com/gookit/goutil/strutil"
	"strings"
)

// Find an alias of string.Index()
func Find(s string, keyword string) int {
	return strings.Index(s, keyword)
}

// ParseTemplates parse all templates like ${xx} and replace it by replaceFunc
func ParseTemplates(s string, replaceFunc func(template string) string) string {
	result := s
	for {
		template, exists, end := FindFirstTemplate(s)
		if !exists {
			break
		}
		// do replace template
		fullTemplate := fmt.Sprintf("${%s}", template)
		replaceSt := replaceFunc(template)
		replaceCount := strings.Count(result, fullTemplate)
		result = strings.Replace(result, fullTemplate, replaceSt, replaceCount)

		// prepare for next template
		s = s[end:]
	}
	return result
}

// FindFirstTemplate find first template like ${xx}
func FindFirstTemplate(s string) (template string, exists bool, end int) {
	// eg. "abc${e.f.g}"
	// find first index of char "$"
	startIdx := Find(s, "$")
	// if there has any failed, return result of not found
	if startIdx == -1 {
		return "", false, -1
	}

	// find first index of char "{" after first "$"
	openBraceStartIdx := Find(
		// index of first "$" + 1
		s[startIdx+1:],
		"{",
	)
	// if there has any failed, return result of not found
	if openBraceStartIdx == -1 {
		return "", false, -1
	}

	// find first index of char "}" after first "{"
	closeBraceEndIdx := Find(
		// index of first "{" + 1
		s[startIdx+openBraceStartIdx+1:],
		"}",
	)
	// if there has any failed, return result of not found
	if closeBraceEndIdx == -1 {
		return "", false, -1
	}

	var (
		//dollarIdx = startIdx
		openBraceIdx  = startIdx + openBraceStartIdx + 1
		closeBraceIdx = startIdx + openBraceStartIdx + closeBraceEndIdx + 1
	)

	exists = true
	template = s[openBraceIdx+1 : closeBraceIdx]
	end = closeBraceIdx + 1

	return
}

// ParsePath parse string type path as slice
func ParsePath(path string) (paths []string, err error) {
	if strutil.IContains(path, " ") {
		err = fmt.Errorf(
			"destination path[%s] contains invaild space",
			strutil.Join(".", paths...),
		)
		return
	}
	paths = strutil.Split(path, ".")
	for _, p := range paths {
		if p == "" {
			err = fmt.Errorf(
				"destination path[%s] contains empty node",
				strutil.Join(".", paths...),
			)
			return
		}
	}
	return paths, nil
}
