package structs

import (
	"github.com/hanakogo/i18n/i18nfs"
	"github.com/hanakogo/i18n/internal/utils"
	"io/fs"
	"os"
	"path/filepath"
)

// WalkLangYAML walk all yaml file of language, and process content of them with walkFunc
func (i *I18nFS) WalkLangYAML(lang string, walkFunc func(content string)) error {
	langDir := i.filePathJoin(i.FSPrefix, lang)
	switch i.FsMode {
	case i18nfs.ModeEmbed:
		dirEntries, err := i.LangFSEmbed.ReadDir(langDir)
		if err != nil {
			return err
		}
		for _, entry := range dirEntries {
			if !utils.CheckYaml(entry.Name()) {
				continue
			}
			langFile := i.filePathJoin(langDir, entry.Name())
			bytes, err := i.LangFSEmbed.ReadFile(langFile)
			if err != nil {
				return err
			}
			walkFunc(string(bytes))
		}
	case i18nfs.ModeFileSystem:
		err := filepath.WalkDir(langDir, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !utils.CheckYaml(entry.Name()) {
				return nil
			}
			bytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			walkFunc(string(bytes))
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
