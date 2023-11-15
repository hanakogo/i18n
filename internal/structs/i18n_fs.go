package structs

import (
	"embed"
	"fmt"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/hanakogo/i18n/i18nfs"
	"github.com/hanakogo/i18n/internal/errors"
	"github.com/hanakogo/i18n/internal/utils"
	"gopkg.in/yaml.v3"
	"reflect"
)

type I18nFS struct {
	LangFSEmbed *embed.FS
	FSPrefix    string
	FsMode      i18nfs.FSMode

	langStringMaps map[string]map[string]any
}

func NewI18nFS(fsMode i18nfs.FSMode, prefix string, embedFS *embed.FS) (*I18nFS, error) {
	switch fsMode {
	case i18nfs.ModeEmbed:
		if embedFS == nil {
			return nil, fmt.Errorf("must provide a vaild embedFS if use ModeEmbed")
		}
		_, err := embedFS.ReadDir(prefix)
		if err != nil {
			return nil, err
		}
	case i18nfs.ModeFileSystem:
		isDir := fsutil.IsDir(prefix)
		if !isDir {
			return nil, fmt.Errorf("directory not found: %s", prefix)
		}
	}

	return &I18nFS{
		LangFSEmbed:    embedFS,
		FSPrefix:       prefix,
		FsMode:         fsMode,
		langStringMaps: make(map[string]map[string]any),
	}, nil
}

// GetLanguages get list of languages
func (i *I18nFS) GetLanguages() []string {
	return maputil.Keys(i.langStringMaps)
}

// IsLangExists check language exists or not on filesystem
func (i *I18nFS) IsLangExists(lang string) bool {
	langDir := i.filePathJoin(i.FSPrefix, lang)
	switch i.FsMode {
	case i18nfs.ModeEmbed:
		_, err := i.LangFSEmbed.ReadDir(langDir)
		if err != nil {
			return false
		}
	case i18nfs.ModeFileSystem:
		return fsutil.IsDir(langDir)
	}
	return true
}

// HasLang check language in langStringMaps
func (i *I18nFS) HasLang(lang string) bool {
	return maputil.HasKey(i.langStringMaps, lang)
}

// Read all yaml of language on filesystem into langStringMaps
func (i *I18nFS) Read(lang string) error {
	if !i.IsLangExists(lang) {
		return errors.GetLangNotExists(lang)
	}
	var langMaps []map[string]any
	var err error
	err = i.WalkLangYAML(lang, func(content string) {
		dst := make(map[string]any)
		err = yaml.Unmarshal([]byte(content), &dst)
		langMaps = append(langMaps, dst)
	})
	if err != nil {
		return err
	}
	mergedLangMap := make(map[string]any)
	for _, langMap := range langMaps {
		err = utils.MergeStringMap(&mergedLangMap, langMap, lang)
		if err != nil {
			return err
		}
	}
	i.langStringMaps[lang] = mergedLangMap
	return nil
}

// GetValByPath get value by paths which are split by dot
func (i *I18nFS) GetValByPath(lang string, path string) (any, error) {
	if !i.HasLang(lang) {
		return "", errors.GetLangNotFound(lang)
	}

	// validate path
	paths, err := utils.ParsePath(path)
	if err != nil {
		return "", err
	}

	// walk all nodes of path, unless last one
	langMap := i.langStringMaps[lang]
	for i := range paths[:len(paths)-1] {
		value := langMap[paths[i]]

		// take out a Map, then continue
		if value, ok := value.(map[string]any); ok {
			langMap = value
		} else {
			// can't take out a Map, so we can't continue to walk deeper structure
			return "", fmt.Errorf(
				"destination path[%s.<%s>.%s] isn't point to a object, can't continue to get value",
				strutil.Join(".", paths[:i]...),
				paths[i],
				strutil.Join(".", paths[i+1:]...),
			)
		}
	}

	// need some special handling for last one node
	lastPath := paths[len(paths)-1]
	value := langMap[lastPath]
	// if value is nil or Map
	if value == nil || reflect.TypeOf(value).Kind() == reflect.Map {
		return "", fmt.Errorf(
			"destination path[%s] isn't point to a value",
			strutil.Join(".", paths...),
		)
	}

	// do template formatting
	if valString, ok := value.(string); ok {
		if valStringParsed := i.ParseTemplateString(valString, lang); valStringParsed != nil {
			value = valStringParsed
		}
	}

	return value, nil
}
