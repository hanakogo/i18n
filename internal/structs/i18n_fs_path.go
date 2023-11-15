package structs

import (
	"github.com/gookit/goutil/strutil"
	"github.com/hanakogo/i18n/i18nfs"
	"path/filepath"
)

// filePathJoin special implementation of filepath.Join, add support for embed.FS
func (i *I18nFS) filePathJoin(paths ...string) string {
	path := filepath.Join(paths...)
	if i.FsMode == i18nfs.ModeEmbed {
		path = strutil.Replaces(path, map[string]string{
			"\\": "/",
		})
	}
	return path
}
