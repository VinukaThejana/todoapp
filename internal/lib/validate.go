package lib

import (
	"net/mail"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidateUsername is a custom validator for validating username
func ValidateUsername(fl validator.FieldLevel) bool {
	username := fl.Parent().FieldByName("Username").String()
	if username == "" {
		return false
	}

	regex, err := regexp.Compile(`^[a-zA-Z0-9_.#]{1,20}$`)
	if err != nil {
		return false
	}

	return regex.MatchString(username)
}

// ValiateEmailOrUsername is a custom validator for validating email or username
func ValiateEmailOrUsername(fl validator.FieldLevel) bool {
	username := fl.Parent().FieldByName("Username").String()
	email := fl.Parent().FieldByName("Email").String()

	if username == "" && email == "" {
		return false
	}

	if email != "" {
		_, err := mail.ParseAddress(email)
		return err == nil
	}

	regex, err := regexp.Compile(`^[a-zA-Z0-9_.#]{1,20}$`)
	if err != nil {
		return false
	}

	return regex.MatchString(username)
}
