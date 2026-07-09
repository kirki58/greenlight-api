package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

/*	--- Validator with Translator Implementation ---	*/
type UniversalValidator struct {
	validate   *validator.Validate
	translator ut.Translator
}

func NewUniversalValidator() (*UniversalValidator, error) {
	enLocale := en.New() // en locale contains the localization engine rules, satisfying ut.Translator interface

	uni := ut.New(enLocale, enLocale) /* ut.UniversalTranslator type is a registry for different locales (ut.Translator types)
	Arguements:
		1. fallback locale: used as a fallback when an unregistered locale (ut.Translator) is requested from the universal translator
		2. a list of supported locales (ut.Translator), in this case it is only the english locale we instantiated just before
	*/

	enTranslator, _ := uni.GetTranslator("en") // Request the english translator we just registered
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := en_translations.RegisterDefaultTranslations(validate, enTranslator)
	if err != nil{
		return nil, fmt.Errorf("Failed to register default translations while constructing the universal validator")
	}

	return &UniversalValidator{
		validate: validate,
		translator: enTranslator,
	}, nil
}

func (uv *UniversalValidator) UseJSONTagNames(){
	uv.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
}

// Returns a structured map of validation errors
// field --> error reason
func (uv *UniversalValidator) ValidateBody(payload any)map[string]string{
	structuredErrors := make(map[string]string)

	err := uv.validate.Struct(payload)
	if err != nil{
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors){
			for _, fieldErr := range validationErrors{
				structuredErrors[fieldErr.Field()] = fieldErr.Translate(uv.translator)
			}

			return structuredErrors
		}
		
		// Any other type of error means that we passed in something to ValidateBody which we should not, so the code using ValidateBody is broken
		panic("Broken validation framework")
	}
	return nil // No validation errors
}

/*	--- Custom Validation tags ---	*/

func yearFrom(fl validator.FieldLevel) bool {
	yearFrom, err := strconv.ParseInt(fl.Param(), 10, 32)
	if err != nil {
		return false
	}

	fieldVal := fl.Field().Int()
	return yearFrom <= fieldVal && fieldVal <= int64(time.Now().Year())
}

func (app *application) RegisterCustomValidations() {
	app.uniValidator.validate.RegisterValidation("yearfrom", yearFrom)
	app.uniValidator.validate.RegisterTranslation("yearfrom", app.uniValidator.translator, 
		// Register the message template to be compiled directly into the given translator's cache at boot
		func(ut ut.Translator) error {
			return ut.Add("yearfrom","{0} must be a valid year between {1} and {2}", false)
		},
		// Register the parameter injection logic to be executed at runtime
		func(ut ut.Translator, fe validator.FieldError) string {
			currentYear := strconv.Itoa(time.Now().Year())
			t, _ := ut.T("yearfrom", fe.Field(), fe.Param(), currentYear)
			return t
		},
	)
}
