package app

import "gin-scaffold/internal/utils/validator"

func Constructor() error {
	err := validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		return err
	}

	return nil
}
