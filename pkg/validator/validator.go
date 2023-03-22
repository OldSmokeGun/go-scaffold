package validator

import (
	"errors"

	"github.com/nyaruka/phonenumbers"
)

var (
	// ErrAssertTypeToStringFailed assertion type is string failed
	ErrAssertTypeToStringFailed = errors.New("assert type to string failed")

	// ErrInvalidPhoneNumber invalid phone number
	ErrInvalidPhoneNumber = errors.New("phone number format is invalid")
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
