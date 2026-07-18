package controller

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// Trans 全局翻译器，供 controller 层翻译校验错误信息
var Trans ut.Translator

// InitTrans 初始化 validator 翻译器，locale 一般为 "zh"
func InitTrans(locale string) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取 json tag 作为字段名的函数，让错误信息更友好
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// 创建中文翻译器
		zhLocale := zh.New()
		uni := ut.New(zhLocale, zhLocale)
		Trans, _ = uni.GetTranslator(locale)

		// 注册默认的中文翻译
		return zhTranslations.RegisterDefaultTranslations(v, Trans)
	}
	return nil
}
