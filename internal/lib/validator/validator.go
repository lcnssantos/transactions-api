package validator

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"schneider.vip/problem"
)

func findComma(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			return i
		}
	}
	return -1
}

func getJSONFieldName(s interface{}, fieldName string) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if field, found := t.FieldByName(fieldName); found {
		tags := []string{"json", "param", "query", "header", "form", "xml"}

		for _, tag := range tags {
			tagValue := field.Tag.Get(tag)

			if tagValue != "" {
				if comma := findComma(tagValue); comma != -1 {
					return tagValue[:comma]
				}
				return tagValue
			}
		}
	}

	return fieldName
}

type Locale string

const (
	Locale_en Locale = "en"
)

type CustomValidator struct {
	validate    *validator.Validate
	translators map[Locale]ut.Translator
	locale      Locale
}

func englishTranslator() (ut.Translator, error) {
	en := en.New()
	uni := ut.New(en, en)
	translator, ok := uni.GetTranslator(en.Locale())
	if !ok {
		return nil, errors.New("en translator not found")
	}
	return translator, nil
}

func NewCustomValidator(locale Locale) (*CustomValidator, error) {
	validate := validator.New()

	en, err := englishTranslator()
	if err != nil {
		return nil, err
	}

	var translators = map[Locale]ut.Translator{
		Locale_en: en,
	}

	switch locale {
	case Locale_en:
		err = en_translations.RegisterDefaultTranslations(validate, translators[locale])
	default:
		err = errors.New("locale not found")
	}

	if err != nil {
		return nil, err
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tags := []string{"json", "param", "query", "header", "form", "xml"}

		for _, tagName := range tags {
			name := fld.Tag.Get(tagName)
			if name != "" {
				if comma := findComma(name); comma != -1 {
					return name[:comma]
				}
				return name
			}
		}

		return fld.Name
	})

	return &CustomValidator{
		validate:    validate,
		translators: translators,
		locale:      locale,
	}, nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validate.Struct(i); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)

		if !ok {
			return err
		}
		errorMessages := make(map[string]string)

		for _, e := range validationErrors {
			errorMessages[getJSONFieldName(i, e.Field())] = e.Translate(cv.translators[cv.locale])
			fmt.Println(errorMessages[getJSONFieldName(i, e.Field())])
		}

		problemDetails := problem.Of(http.StatusBadRequest).Append(
			problem.Status(http.StatusBadRequest),
			problem.Detail("Request validation error"),
			problem.Custom("errors", errorMessages),
		)

		return echo.NewHTTPError(http.StatusBadRequest, problemDetails)
	}
	return nil
}
