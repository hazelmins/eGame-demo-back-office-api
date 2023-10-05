/*
 * @Description:初始化翻译器（Translator）Gin框架中进行表单验证的国际化处理
 */
package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

func InitTrans(locale string) (ut.Translator, error) {
	var trans ut.Translator
	var err error
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		if locale == "zh" {
			v.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]
				if name == "-" {
					return ""
				}
				return name
			})
		}

		zhT := zh.New()
		enT := en.New()

		uni := ut.New(enT, zhT, enT)

		var ok bool
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			err = fmt.Errorf("uni.GetTranslator(%s) failed", locale)
			return nil, err
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return trans, err
	}
	return trans, nil
}

/*
这段代码的主要作用是根据传入的locale参数初始化一个Translator，
并将其用于Gin框架的表单验证错误信息的国际化处理。根据不同的locale，
会注册不同的标签翻译函数和翻译器。如果传入的locale不是"en"或"zh"，则默认使用英文翻译器。
注意：这段代码依赖于Gin框架和其内部的Validator引擎以及一些第三方包来实现翻译功能。
*/
