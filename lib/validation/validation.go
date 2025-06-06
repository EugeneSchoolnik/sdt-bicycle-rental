package validation

import (
	"fmt"
	"strings"

	"errors"

	"github.com/go-playground/validator/v10"
)

func PrettyError(validationErrs ...validator.ValidationErrors) error {
	var errMsgs []string

	for _, vErr := range validationErrs {
		for _, err := range vErr {
			switch err.ActualTag() {
			case "required":
				errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
			case "email":
				errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid email", err.Field()))
			default:
				errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
			}
		}
	}

	return errors.New(strings.Join(errMsgs, ", "))
}
