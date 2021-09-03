package handlers

import (
	"gin-scaffold/global"
	"gin-scaffold/internal/pkg/http/response"
	"gin-scaffold/internal/pkg/validator"
	"github.com/gin-gonic/gin"
	validator2 "github.com/go-playground/validator/v10"
	"net/http"
)

func responseValidateError(ctx *gin.Context, err error, errTrans map[string]string) {
	validationErrors, ok := err.(validator2.ValidationErrors)
	if ok {
		errs := validator.Translate(validationErrors, errTrans)
		for _, err := range errs {
			ctx.JSON(http.StatusOK, response.FailedFormat(err))
			return
		}
	} else {
		ctx.JSON(http.StatusOK, response.Format(response.ArgumentsInvalidCode, response.ArgumentsInvalidCodeMessage, nil))
		return
	}
}

func ValidateQueryError(ctx *gin.Context, obj interface{}, errTrans map[string]string) bool {
	if err := ctx.ShouldBindQuery(obj); err != nil {
		global.Logger().Error(err)
		responseValidateError(ctx, err, errTrans)
		return false
	}
	return true
}

func ValidateFormError(ctx *gin.Context, obj interface{}, errTrans map[string]string) bool {
	if err := ctx.ShouldBind(obj); err != nil {
		global.Logger().Error(err)
		responseValidateError(ctx, err, errTrans)
		return false
	}
	return true
}
