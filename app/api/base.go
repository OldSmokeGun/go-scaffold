package api

import (
	"gin-scaffold/app/utils/http/response"
	"gin-scaffold/app/utils/validator"
	"gin-scaffold/internal/global"
	"github.com/gin-gonic/gin"
	validator2 "github.com/go-playground/validator/v10"
	"net/http"
)

func responseValidateError(c *gin.Context, err error, errTrans map[string]string) {
	validationErrors, ok := err.(validator2.ValidationErrors)
	if ok {
		errs := validator.Translate(validationErrors, errTrans)
		for _, err := range errs {
			c.JSON(http.StatusOK, response.FailedFormat(err))
			return
		}
	} else {
		c.JSON(http.StatusOK, response.Format(response.ArgumentsInvalidCode, response.ArgumentsInvalidCodeMessage, nil))
		return
	}
}

func ValidateQueryError(c *gin.Context, obj interface{}, errTrans map[string]string) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		global.GetLogger().Error(err)
		responseValidateError(c, err, errTrans)
		return false
	}
	return true
}

func ValidateFormError(c *gin.Context, obj interface{}, errTrans map[string]string) bool {
	if err := c.ShouldBind(obj); err != nil {
		global.GetLogger().Error(err)
		responseValidateError(c, err, errTrans)
		return false
	}
	return true
}
