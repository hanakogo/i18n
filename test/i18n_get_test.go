package test

import (
	"github.com/gookit/goutil/testutil/assert"
	"github.com/hanakogo/i18n"
	"testing"
)

func TestGet(t *testing.T) {
	TestLoadEmbed(t)

	as := assert.New(t)

	// test basic get
	as.Eq(
		"测试",
		i18n.Get[string]("test.str1", i18n.ConvertString, "default"),
	)
	as.Eq(
		"default",
		i18n.Get[string]("test.str1.not.exist", i18n.ConvertString, "default"),
	)
	value, err := i18n.GetValue("test.str1", i18n.ConvertString)
	as.Eq(nil, err)
	as.Eq(
		"测试",
		value,
	)

	// test convert
	as.Eq(
		int64(123),
		i18n.Get[int64]("test.num2", i18n.ConvertInt64, 0),
	)
	as.Eq(
		int64(0),
		i18n.Get[int64]("test.num2.not.exist", i18n.ConvertInt64, 0),
	)

	as.Eq(
		123.456,
		i18n.Get[float64]("test.num2", i18n.ConvertFloat, 0.1),
	)
	as.Eq(
		0.1,
		i18n.Get[float64]("test.num2.not.exist", i18n.ConvertFloat, 0.1),
	)

	as.Eq(
		"123.456",
		i18n.Get[string]("test.num2", i18n.ConvertString, "default"),
	)
	as.Eq(
		"default",
		i18n.Get[string]("test.num2.not.exist", i18n.ConvertString, "default"),
	)

	// test get with translation mode
	as.Eq(
		"654.321",
		i18n.GetTr[string]("en", "test.num2", i18n.ConvertString, "default"),
	)
	as.Eq(
		"default",
		i18n.GetTr[string]("en", "test.num2.not.exist", i18n.ConvertString, "default"),
	)

	value, err = i18n.GetValueTr("en", "fruits.banana", i18n.ConvertString)
	as.Eq(nil, err)
	as.Eq(
		"banana",
		value,
	)

	// test get slice
	as.Eq(
		[]string{"abc", "def"},
		i18n.GetSlice[string]("test.strList", i18n.ConvertString),
	)
	as.Eq(
		[]string{"1", "2"},
		i18n.GetSlice[string]("test.numList", i18n.ConvertString),
	)
	as.Eq(
		[]string{},
		i18n.GetSlice[string]("test.strList.not.exist", i18n.ConvertString),
	)

	as.Eq(
		[]int64{1, 2},
		i18n.GetSlice[int64]("test.numList", i18n.ConvertInt64),
	)
	as.Eq(
		[]int64{-1, -1},
		i18n.GetSlice[int64]("test.strList", i18n.ConvertInt64),
	)
	as.Eq(
		[]int64{},
		i18n.GetSlice[int64]("test.numList.not.exist", i18n.ConvertInt64),
	)

	// by language
	as.Eq(
		[]string{"a", "b", "c"},
		i18n.GetSliceTr[string]("en", "test.strList", i18n.ConvertString),
	)
}

func TestMultiLevel(t *testing.T) {
	TestLoadEmbed(t)

	as := assert.New(t)

	ok, contains := i18n.HasPath("test.first.second.third")
	as.Eq(true, ok)
	as.Eq([]string{"zh-CN"}, contains)

	as.Eq(
		"3level",
		i18n.GetString("test.first.second.third"),
	)
}

func TestGetString(t *testing.T) {
	TestLoadEmbed(t)

	as := assert.New(t)

	as.Eq(
		"香蕉",
		i18n.GetString("fruits.banana"),
	)

	as.Eq(
		"banana",
		i18n.GetStringTr("en", "fruits.banana"),
	)
}

func TestGetInt64(t *testing.T) {
	TestLoadEmbed(t)

	as := assert.New(t)

	as.Eq(
		int64(123),
		i18n.GetInt64("test.num1"),
	)

	as.Eq(
		int64(654),
		i18n.GetInt64Tr("en", "test.num1"),
	)

	as.Eq(
		int64(-1),
		i18n.GetInt64("test.str1"),
	)
}

func TestGetFloat(t *testing.T) {
	TestLoadEmbed(t)

	as := assert.New(t)

	as.Eq(
		123.456,
		i18n.GetFloat("test.num2"),
	)

	as.Eq(
		654.321,
		i18n.GetFloatTr("en", "test.num2"),
	)

	as.Eq(
		float64(-1),
		i18n.GetFloat("test.str1"),
	)
}
