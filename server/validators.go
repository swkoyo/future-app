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
	validate.RegisterStructValidation(AppointmentFromToValidation, GetTrainerAppointmentsReq{})
	validate.RegisterStructValidation(AvailabilityFromToValidation, GetTrainerAvailabilityReq{})
	validate.RegisterValidation("is-future-date", ValidateFutureDate)

	en_translations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTranslation("is-future-date", trans, func(ut ut.Translator) error {
		return ut.Add("is-future-date", "{0} must be a future date", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("is-future-date", fe.Field())
		return t
	})

	validate.RegisterTranslation("timeframe-invalid", trans, func(ut ut.Translator) error {
		return ut.Add("timeframe-invalid", "Invalid timeframe", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("timeframe-invalid", fe.Field())
		return t
	})

	validate.RegisterTranslation("timeframe-max", trans, func(ut ut.Translator) error {
		return ut.Add("timeframe-max", "Timeframe must be 90 days or lower", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("timeframe-max", fe.Field())
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

type GetTrainerAppointmentsReq struct {
	TrainerID int    `param:"trainer_id" validate:"required,min=1"`
	From      string `query:"from" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	To        string `query:"to" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

func AppointmentFromToValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(GetTrainerAppointmentsReq)

	parsedFrom, err := time.Parse(time.RFC3339, req.From)
	if err != nil {
		sl.ReportError(parsedFrom, "from", "From", "datetime", "")
	}

	parsedTo, err := time.Parse(time.RFC3339, req.To)
	if err != nil {
		sl.ReportError(parsedTo, "to", "To", "datetime", "")
	}

	if parsedFrom.After(parsedTo) {
		sl.ReportError(parsedFrom, "from", "From", "timeframe-invalid", "")
	}

	if parsedTo.Sub(parsedFrom) > 90*24*time.Hour {
		sl.ReportError(parsedTo, "to", "To", "timeframe-max", "")
	}
}

type GetTrainerAvailabilityReq struct {
	TrainerID int    `param:"trainer_id" validate:"required,min=1"`
	From      string `query:"from" validate:"required,datetime=2006-01-02T15:04:05Z07:00,is-future-date"`
	To        string `query:"to" validate:"required,datetime=2006-01-02T15:04:05Z07:00,is-future-date"`
}

func AvailabilityFromToValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(GetTrainerAvailabilityReq)

	parsedFrom, err := time.Parse(time.RFC3339, req.From)
	if err != nil {
		sl.ReportError(parsedFrom, "from", "From", "datetime", "")
	}

	parsedTo, err := time.Parse(time.RFC3339, req.To)
	if err != nil {
		sl.ReportError(parsedTo, "to", "To", "datetime", "")
	}

	if parsedFrom.After(parsedTo) {
		sl.ReportError(parsedFrom, "from", "From", "timeframe-invalid", "")
	}

	if parsedTo.Sub(parsedFrom) > 90*24*time.Hour {
		sl.ReportError(parsedTo, "to", "To", "timeframe-max", "")
	}
}
