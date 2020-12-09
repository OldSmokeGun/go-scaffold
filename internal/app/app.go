package app

import (
	"gin-scaffold/internal/utils/validator"
	"github.com/gin-gonic/gin"
)

func Constructor(r *gin.Engine) error {
	err := validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		return err
	}

	return nil
}
