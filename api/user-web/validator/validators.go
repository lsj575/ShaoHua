package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateLoginEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if email == "" {
		return true
	}
	// use regular expression validate email
	pattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	if ok, _ := regexp.MatchString(pattern, email); !ok {
		return false
	}
	return true
}

func ValidateLoginUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if username == "" {
		return true
	}
	lenOfUsername := len(username)
	if lenOfUsername > 18 {
		return false
	}
	return true
}
