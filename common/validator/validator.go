package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	"nory/common/response"
)

var (
	validate = validator.New()

	alphanumericSpaceRegex = regexp.MustCompile("^[a-zA-Z0-9\\s]*$")
)

func init() {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			name = strings.ToLower(field.Name)
		}
		if name == "-" {
			return ""
		}
		return name
	})

	validate.RegisterValidation("alphanumeric_space", func(field validator.FieldLevel) bool {
		str := field.Field().String()
		return alphanumericSpaceRegex.MatchString(str)
	})
}

func ValidateStruct(s any, message string) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	vErr, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	errs := map[string][]string{}
	for _, fe := range vErr {
		path := fe.Namespace()
		path = strings.SplitN(path, ".", 2)[1]
		errs[path] = append(errs[path], fe.Error())
	}
	return response.NewBadRequest(message, errs)
}
