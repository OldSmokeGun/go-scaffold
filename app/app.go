package app

import (
	"gin-scaffold/app/util/validator"
	"github.com/gin-gonic/gin"
)

func Run(r *gin.Engine) error {
	err := validator.RegisterValidator([]validator.CustomValidator{
		{"phone", validator.IsPhone},
	})
	if err != nil {
		return err
	}

	return nil
}
