package server

import (
	"errors"
	"future-app/models"
	"time"

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
	validate.RegisterValidation("is-future-date", ValidateFutureDate)

	en_translations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTranslation("is-future-date", trans, func(ut ut.Translator) error {
		return ut.Add("is-future-date", "{0} must be a future date", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("is-future-date", fe.Field())
		return t
	})

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

type PostAppointmentReq struct {
	UserID    int    `json:"user_id" validate:"required,min=1"`
	TrainerID int    `json:"trainer_id" validate:"required,min=1"`
	StartedAt string `json:"started_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00,is-future-date"`
	EndedAt   string `json:"ended_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00,is-future-date"`
}

func ValidateFutureDate(fl validator.FieldLevel) bool {
	parsedDate, err := models.ParseDateStr(fl.Field().String())
	if err != nil {
		return false
	}
	return parsedDate.After(time.Now())
}
