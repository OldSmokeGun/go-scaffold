package validator

import (
	validator2 "github.com/go-playground/validator/v10"
	"regexp"
)

var IsPhone validator2.Func = func(fl validator2.FieldLevel) bool {
	phone := fl.Field().String()

	ok, _ := regexp.MatchString("^[1](([3][0-9])|([4][5-9])|([5][0-3,5-9])|([6][5,6])|([7][0-8])|([8][0-9])|([9][1,8,9]))[0-9]{8}$", phone)

	if !ok {
		return false
	}

	return true
}