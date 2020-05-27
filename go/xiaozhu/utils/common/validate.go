package common

import (
	"fmt"
	zhongwen "github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	. "gopkg.in/go-playground/validator.v9/translations/zh"
)

var validate *validator.Validate
var trans ut.Translator

func init() {
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")

	validate = validator.New()
	err := RegisterDefaultTranslations(validate, trans)

	if err != nil {
		fmt.Println(err)
	}
}

func Check(mystruct interface{}) (bool, string) {
	err := validate.Struct(mystruct)

	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		//errs, ok := err.(validator.ValidationErrors)
		if errs, ok := err.(validator.ValidationErrors); ok {
			for k, v := range errs.Translate(trans) {
				return false, k + v
			}
		}
	}

	return true, ""
}
