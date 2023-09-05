package validator

import (
	"errors"
	"unicode"

	"github.com/nyaruka/phonenumbers"
)

var (
	// ErrAssertTypeToStringFailed assertion type is string failed
	ErrAssertTypeToStringFailed = errors.New("assert type to string failed")

	// ErrInvalidPhoneNumber invalid phone number
	ErrInvalidPhoneNumber = errors.New("phone number format is invalid")

	// ErrPasswordComplexityTooLow password complexity too low
	ErrPasswordComplexityTooLow = errors.New("password complexity too low")
)

// IsPhoneNumber check whether it is a phone number
//
// based on Google's libphonenumber library
func IsPhoneNumber(value any) error {
	phone, ok := value.(string)
	if !ok {
		return ErrAssertTypeToStringFailed
	}

	phoneNumber, err := phonenumbers.Parse(phone, "CN")
	if err != nil {
		return err
	}

	if !phonenumbers.IsValidNumber(phoneNumber) {
		return ErrInvalidPhoneNumber
	}

	return nil
}

// PasswordComplexity validate password complexity
func PasswordComplexity(value any) error {
	password, ok := value.(string)
	if !ok {
		return ErrAssertTypeToStringFailed
	}

	var (
		existNumber bool
		existUpper  bool
		existLower  bool
		existSymbol bool
	)

	for _, v := range password {
		if unicode.IsNumber(v) && !existNumber {
			existNumber = true
		}
		if unicode.IsUpper(v) && !existUpper {
			existUpper = true
		}
		if unicode.IsLower(v) && !existLower {
			existLower = true
		}
		if unicode.IsPunct(v) && !existSymbol {
			existSymbol = true
		}
	}

	if !existNumber || !existUpper || !existLower || !existSymbol {
		return ErrPasswordComplexityTooLow
	}

	return nil
}
