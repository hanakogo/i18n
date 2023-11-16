# hanakogo-i18n
[**hanakogo-i18n**](https://github.com/hanakogo/i18n) is a simple, fast and easy-to-use tool for multi-languages management. Based on `YAML`.

## Features
- easy to use: simple and clear usage, support multi languages and multi configuration files in any language.
- clear key: read config item by path which is split by dot, e.g.`main.businessA.str1`.
- multi type: read `String`,`Int64`,`Float64` from config of language, even **any** type of `Slice`.
- formatting: support reading config item with format, use regular format specifier.
- template string: support uses template being like `${refer}`. `refer` is a full path of target item.
- read mode: support read configs of language from `embed.FS` of golang and traditional filesystem(directory mode).
- language settings: support default language and fallback language settings.
- designed with the singleton pattern, without the creation of any struct

## Installation
```shell
go get github.com/hanakogo/i18n
```

## Usage
```shell
lang/
    en/
      main.yaml
      database.yaml
      ...
    zh-CN/
      main.yaml
      database.yaml
      ...
```

#### Basic usage

```go
package main

import (
	"embed"
	"fmt"
	"github.com/hanakogo/i18n"
	"github.com/hanakogo/i18n/i18nfs"
)

//go:embed lang
var testdata embed.FS

func main() {
	// initialize
	err := i18n.Init(i18n.Opts{
		FSOpts: i18n.FSOpts{
			// use i18nfs.ModeFileSystem to read from directory
			FSMode: i18nfs.ModeEmbed,
			// directory prefix
			Prefix: "lang",
			// set EmbedFS if you want to use i18nfs.ModeEmbed
			EmbedFS: &testdata,
		},
		DefaultLang:  "zh-CN",
		FallbackLang: "en",
		// languages list to load, DefaultLang and FallbackLang must be contained
		Languages: []string{
			"en", "zh-CN",
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// configuration template of current language (DefaultLang):
	// zh-CN
	// main:
	//   title: 测试
	// en
	// main:
	//   title: test
	//   str1: string one
	
	// basic Get[T]() function must specify ConvertFunc used to convert any to the target type
	// you can use ConvertFunc from package or define it by yourself
	// default value of Get() and GetTr() is must specified
	i18n.Get[string]("main.title", i18n.ConvertString, "def") // "test"
	i18n.Get[string]("main.title.not.exists", i18n.ConvertString, "def") // "def"
	
	// this path doesn't exist in language "zh-CN", fallback to language "en"
	i18n.Get[string]("main.str1", i18n.ConvertString, "def") // "string one"
	
	// advanced usage
	val, err := i18n.GetValue("main.str1") // get any type, for custom usage
}
```

#### Load language

```go
// manual load language
// must after initialize
err := i18n.Load("common")
if err != nil {
	fmt.Println(err)
	return
}
```

#### Change global exposed value

```go
// you can change default global return value
i18n.DefaultString = "abc"
i18n.DefaultInt = 123
i18n.DefaultFloat = 1.1

// after initialized, you set below settings manually
i18n.DefaultLang = "en"
i18n.FallbackLang = "zh-CN"
```

#### Status check

```go
// check is initialized
i18n.Initialized()

// reset all status (include loaded language)
i18n.Reset()

// check language is existing or not
i18n.Has("zh-CN")

// check a target path is existing or not
i18n.HasPath("main.dst") // default check all languages
i18n.HasPath("main.dst", "en", ...) // you can specify some language to check
```

#### Format string

```go
// main:
//   format: test %s
i18n.GetStringF("main.test", "abc") // "test abc" 
```

#### Get a value of specific type

```go
// main:
//   test: 12.3
i18n.GetString("main.test") // "12.3"
i18n.GetInt64("main.test") // 12
i18n.GetFloat("main.test") // 12.30000
```

#### Custom ConvertFunc

```go
// convert any to string
// main:
//   test: 12.3
i18n.Get[string]("main.test", func(value any) string {
	res, err := mathutil.ToString(value)
	if err != nil {
		return ""
	}
	return res
}, "def") // "12.3"
```

#### Get a Slice

```go
// main:
//   list:
//     - a
//     - b
i18n.GetSlice[string]("main.list", i18n.ConvertString) // []string{"a", "b"}
```

#### Set default value

```go
// all `getXX` functions are support default value (include translation func)
// optional, if not specify, will use default global value
i18n.GetString("main.test", "def")
i18n.GetInt64("main.test", 0)
i18n.GetFloat("main.test", 1.0)
// default global value of slice is empty slice
i18n.GetSlice("main.test", i18n.ConvertString, []string{})
```

#### Translation mode

```go
// it's just a nice alias meaning that "get value from specified language"

// zh-CN (Default language)
// main:
//   title: 测试
// en (Fallback language)
// main:
//   title: test
i18n.GetTr[string]("en", "main.test", i18n.ConvertString, "def") // "test"
i18n.GetStringTr("en", "main.test") // "test"
i18n.GetStringTr("zh-CN", "main.test") // "测试"
// same usage on GetInt64Tr(), GetFloatTr(), GetSliceTr(), GetValueTr()
```

#### Template string

```go
// zh-CN (Default language)
// main:
//   refer: 我去，初音未来
//   template1: 引用: ${main.refer}
//   template2: 引用: ${main.engOnlyStr}
//   template3: 引用: ${en:main.refer}

// en (Fallback language)
// main:
//   refer: meow
//   engOnlyStr: eng

// parse simple template
i18n.GetString("main.template1") // "引用: 我去，初音未来"
// auto fallback if template reference doesn't exist in this language
i18n.GetString("main.template2") // "引用: eng"
// refer item from other language
i18n.GetString("main.template3") // "引用: meow"
```

## Dependencies

- [go-yaml/yaml](https://github.com/go-yaml/yaml)
- [gookit/goutil](https://github.com/gookit/goutil)

## License

**[LGPL-3.0](LICENSE)**