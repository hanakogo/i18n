package i18n

import (
	"embed"
	"github.com/hanakogo/i18n/i18nfs"
)

type Opts struct {
	FSOpts       FSOpts
	DefaultLang  string
	FallbackLang string
	Languages    []string
}

type FSOpts struct {
	FSMode  i18nfs.FSMode
	Prefix  string
	EmbedFS *embed.FS
}
