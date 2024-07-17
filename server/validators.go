package server

import (
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type CustomValidator struct {
	trans     ut.Translator
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()

	en_translations.RegisterDefaultTranslations(validate, trans)

	return &CustomValidator{validator: validate, trans: trans}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)

	if err != nil {
		object, _ := err.(validator.ValidationErrors)

		for _, key := range object {
			return errors.New(key.Translate(cv.trans))
		}
	}

	return nil
}
