package structs

import (
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/strutil"
	"github.com/hanakogo/i18n/internal/utils"
)

const DefaultTemplatePlaceholder = "<NotFound>"

// ParseTemplateString parse path in template to refer others string
// format of template like "${path}"
func (i *I18nFS) ParseTemplateString(stringVal string, lang string) any {
	return utils.ParseTemplates(stringVal, func(template string) string {
		// specific language in template
		if strutil.IContains(template, ":") {
			lang = strutil.Split(template, ":")[0]
			template = strutil.Split(template, ":")[1]
		}
		if lang == "" || template == "" {
			return DefaultTemplatePlaceholder
		}
		parsedTempVal, err := i.GetValByPath(lang, template)
		if err != nil {
			return DefaultTemplatePlaceholder
		}
		parsedTempString, err := goutil.ToString(parsedTempVal)
		if err != nil {
			return DefaultTemplatePlaceholder
		}
		return parsedTempString
	})
}
