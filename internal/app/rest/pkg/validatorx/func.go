package validatorx

import (
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

// ValidatePhone 校验手机号
// 基于 google 的 libphonenumber 库的 go 版本
func ValidatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	phoneNumber, err := phonenumbers.Parse(phone, "CN")
	if err != nil {
		return false
	}

	return phonenumbers.IsValidNumber(phoneNumber)
}
