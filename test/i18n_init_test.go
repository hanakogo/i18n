package test

import (
	"embed"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/hanakogo/i18n"
	"github.com/hanakogo/i18n/i18nfs"
	"testing"
)

//go:embed lang
var testdata embed.FS

func TestLoadEmbed(t *testing.T) {
	as := assert.New(t)

	i18n.Reset()
	err := i18n.Init(i18n.Opts{
		FSOpts: i18n.FSOpts{
			FSMode:  i18nfs.ModeEmbed,
			Prefix:  "lang",
			EmbedFS: &testdata,
		},
		DefaultLang:  "zh-CN",
		FallbackLang: "en",
		Languages: []string{
			"en", "zh-CN",
		},
	})

	as.Eq(nil, err)
}

func TestLoadFileSystem(t *testing.T) {
	as := assert.New(t)

	i18n.Reset()
	err := i18n.Init(i18n.Opts{
		FSOpts: i18n.FSOpts{
			FSMode: i18nfs.ModeFileSystem,
			Prefix: "./lang",
		},
		DefaultLang:  "zh-CN",
		FallbackLang: "en",
		Languages: []string{
			"en", "zh-CN",
		},
	})

	as.Eq(nil, err)
}

func TestInitialized(t *testing.T) {
	TestLoadEmbed(t)

	as := assert.New(t)

	as.Eq(true, i18n.Initialized())
	as.Eq(true, i18n.Has("zh-CN"))

	as.Eq(false, i18n.Has("common"))
	err := i18n.Load("common")
	as.Eq(nil, err)
	as.Eq(true, i18n.Has("common"))
}
