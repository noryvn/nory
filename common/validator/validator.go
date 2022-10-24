package validator

import (
	"fmt"
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

func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	vErr, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	message := ""
	for _, fe := range vErr {
		message = fmt.Sprintf("failed to validate %q, because %q rule", fe.Namespace(), fe.Tag())
	}
	return response.NewBadRequest(message)
}
