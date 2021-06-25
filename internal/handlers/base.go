package handlers

import (
	"gin-scaffold/internal/ctx"
	"gin-scaffold/internal/utils/http/response"
	"gin-scaffold/internal/utils/validator"
	"github.com/gin-gonic/gin"
	validator2 "github.com/go-playground/validator/v10"
	"net/http"
)

func responseValidateError(ctx *gin.Context, appCtx *ctx.Context, err error, errTrans map[string]string) {
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

func ValidateQueryError(ctx *gin.Context, appCtx *ctx.Context, obj interface{}, errTrans map[string]string) bool {
	if err := ctx.ShouldBindQuery(obj); err != nil {
		appCtx.GetLogger().Error(err)
		responseValidateError(ctx, appCtx, err, errTrans)
		return false
	}
	return true
}

func ValidateFormError(ctx *gin.Context, appCtx *ctx.Context, obj interface{}, errTrans map[string]string) bool {
	if err := ctx.ShouldBind(obj); err != nil {
		appCtx.GetLogger().Error(err)
		responseValidateError(ctx, appCtx, err, errTrans)
		return false
	}
	return true
}
