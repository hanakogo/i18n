# hanakogo-i18n

[**hanakogo-i18n**](https://github.com/hanakogo/i18n) is a simple, fast, and user-friendly tool for managing multiple
languages. It is based on `YAML` files.

## Features

- **Easy to Use**: Provides a simple and clear usage interface. Supports multiple languages and configuration files in
  any language.
- **Clear Key Structure**: Allows reading configuration items using a dot-separated path, e.g., `main.businessA.str1`.
- **Multi-Type Support**: Reads `String`, `Int64`, and `Float64` types from language configurations. Supports reading
  any type of `Slice`.
- **Formatting Support**: Supports reading configuration items with formatted values using regular format specifiers.
- **Template String**: Supports using template strings with placeholders, e.g., `${refer}`. The `refer` is a full path
  to the target item.
- **Flexible Configuration Sources**: Supports reading language configurations from both `embed.FS` in Golang and
  traditional file systems (directory mode).
- **Language Settings**: Allows setting a default language and a fallback language.
- **Singleton Design**: Designed with the singleton pattern, eliminating the need to create any struct instances.

## Installation

```shell
go get github.com/hanakogo/i18n
```

## Usage
#### Simple File Structure

The following is an example of a simple file structure for language configurations:

```plaintext
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

In this structure, language configurations are organized in separate directories based on language codes. Each language directory contains YAML files for different parts of the application, such as `main.yaml` and `database.yaml`.

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
var languageFiles embed.FS

func main() {
  // Initialize i18n
  err := i18n.Init(i18n.Opts{
    FSOpts: i18n.FSOpts{
      FSMode:   i18nfs.ModeEmbed,
      Prefix:   "lang",
      EmbedFS:  &languageFiles,
    },
    DefaultLang:  "zh-CN",
    FallbackLang: "en",
    Languages: []string{
      "en", "zh-CN",
    },
  })
  if err != nil {
    fmt.Println(err)
    return
  }

  // example data from i18n:
  // zh-CN (Default)
  // main:
  //   title: 测试
  // en (Fallback)
  // main:
  //   title: test
  //   str1: string one

  // Example usage
  // Retrieves the translated string for the key "main.title" in the current language.
  // If the key doesn't exist, fallback to the default value "def".
  title := i18n.Get[string]("main.title", i18n.ConvertString, "def")
  fmt.Println(title) // Output: "测试"

  // Retrieve the translated string for the key "main.title.not.exists" in the current language.
  // Since the key doesn't exist, fallback to the default value "def".
  notExists := i18n.Get[string]("main.title.not.exists", i18n.ConvertString, "def")
  fmt.Println(notExists) // Output: "def"

  // Retrieve the translated string for the key "main.str1" in the current language.
  // If the key doesn't exist in the current language, fallback to the "en" language.
  str1 := i18n.Get[string]("main.str1", i18n.ConvertString, "def")
  fmt.Println(str1) // Output: "string one"

  // Advanced usage
  // Retrieve the value for the key "main.str1" as an interface{}.
  // This allows retrieving values of any type for custom usage.
  val, err := i18n.GetValue("main.str1")
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(val)
}
```

#### Load language

```go
// manual load language after initialized
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
// check i18n is initialized
i18n.Initialized()

// reset all status (include loaded languages)
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
i18n.GetStringF("main.format", "abc") // "test abc" 
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
i18n.Get[string]("main.test", func (value any) string {
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
//   - a
//   - b
i18n.GetSlice[string]("main.list", i18n.ConvertString) // []string{"a", "b"}

// test:
//   strList:
//     - abc
//     - def
//   numList:
//     - 1
//     - 2
//   objList:
//     - sublist:
//       - sublist_element1
//     - substr: substr
i18n.GetInt64("test.numList[0]") // 1
i18n.GetString("test.strList[0]") // "abc"
i18n.GetString("test.objList[0].sublist[0]") // "sublist_element1"
i18n.GetString("test.objList[1].substr") // "substr"
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