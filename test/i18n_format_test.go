package test

import (
	"github.com/gookit/goutil/testutil/assert"
	"github.com/hanakogo/i18n"
	"testing"
)

func TestFormat(t *testing.T) {
	TestLoadEmbed(t)

	as := assert.New(t)

	as.Eq(
		"参数测试 参数1 参数1.1 参数[abc def]",
		i18n.GetStringF("test.format_test", "测试", 1, 1.1, []string{"abc", "def"}),
	)

	// with translation mode
	as.Eq(
		"参数测试 参数1 参数1.1 参数[abc def]",
		i18n.GetStringTrF("zh-CN", "test.format_test", "测试", 1, 1.1, []string{"abc", "def"}),
	)
}

func TestTemplate(t *testing.T) {
	TestLoadEmbed(t)

	as := assert.New(t)

	as.Eq(
		"测试字符串1:测试",
		i18n.Get[string]("test.temp_test1", i18n.ConvertString, "default"),
	)

	as.Eq(
		"测试字符串1:测试",
		// also support default value
		i18n.GetString("test.temp_test1", "default"),
	)

	as.Eq(
		"测试数字2:123.456",
		i18n.GetString("test.temp_test2"),
	)

	as.Eq(
		"测试按语言引用:测试 eng <NotFound>",
		// explanation of "<NotFound>": default placeholder for template which is not found
		// because of "common" language isn't loaded at initial, ${common:common.str1} isn't found
		i18n.GetString("test.temp_test3"),
	)

	err := i18n.Load("common")
	as.Eq(nil, err)

	as.Eq(
		"测试按语言引用:测试 eng common_str", // now, ${common:common.str1} can be found
		i18n.GetString("test.temp_test3"),
	)

	// with translation mode
	as.Eq(
		"测试按语言引用:测试 eng common_str", // now, ${common:common.str1} can be found
		i18n.GetStringTr("zh-CN", "test.temp_test3"),
	)
}
