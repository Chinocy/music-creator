package util

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type CustomValidation struct {
	tagName        string
	errorMsg       string
	validationFunc validator.Func
}

func ValidateStruct(body interface{}) (clientErrs []string) {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	enTranslations.RegisterDefaultTranslations(validate, trans)
	customValidations := []CustomValidation{
		{
			tagName:        "is-language",
			errorMsg:       fmt.Sprintf("{0} must be a valid language (%s)", strings.Join(Languages, ", ")),
			validationFunc: validateLanguage,
		},
		{
			tagName:        "is-genre",
			errorMsg:       fmt.Sprintf("{0} must be a valid genre (%s)", strings.Join(Genres, ", ")),
			validationFunc: validateGenre,
		},
		{
			tagName:        "is-emotion",
			errorMsg:       fmt.Sprintf("{0} must be a valid emotion (%s)", strings.Join(Emotions, ", ")),
			validationFunc: validateEmotion,
		}}
	for _, customValidation := range customValidations {
		registerCustomValidation(validate, trans, customValidation)
	}
	err := validate.Struct(body)
	clientErrs = translateError(err, trans)
	return
}

func translateError(err error, trans ut.Translator) (clientErrs []string) {
	if err == nil {
		return
	}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		clientErrs = append(clientErrs, "Invalid payload")
		return
	} else if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
		clientErrs = append(
			clientErrs,
			fmt.Sprintf(
				"validation failed on field '%s'. Invalid type: %s, expected: %s",
				jsonError.Field,
				jsonError.Value,
				jsonError.Type,
			),
		)
		return
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := e.Translate(trans)
		clientErrs = append(clientErrs, translatedErr)
	}
	return
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(fe.Tag(), fe.Field())
	return t
}

func registerCustomValidation(
	validate *validator.Validate,
	trans ut.Translator,
	customValidation CustomValidation,
) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(customValidation.tagName, customValidation.errorMsg, true)
	}
	validate.RegisterTranslation(customValidation.tagName, trans, registerFn, translateFunc)
	validate.RegisterValidation(customValidation.tagName, customValidation.validationFunc)
}

func validateLanguage(fl validator.FieldLevel) bool {
	language := fl.Field().String()
	_, ok := LanguagesMap[language]
	return ok
}

func validateGenre(fl validator.FieldLevel) bool {
	genre := fl.Field().String()
	_, ok := GenresMap[genre]
	return ok
}

func validateEmotion(fl validator.FieldLevel) bool {
	emotion := fl.Field().String()
	_, ok := EmotionsMap[emotion]
	return ok
}
