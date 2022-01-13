package bindx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/pkg/validatorx"
	"go-scaffold/internal/app/rest/pkg/responsex"
)

var (
	ErrValidateMessageTransformFailed = errors.New("校验信息转换失败")
)

// BindModel 模型绑定需要实现的接口
type BindModel interface {
	ErrorMessage() map[string]string
}

func shouldBind(ctx *gin.Context, m BindModel, b interface{}, bindBody bool) bool {
	var err error

	if bindBody {
		err = ctx.ShouldBindBodyWith(m, b.(binding.BindingBody))
	} else {
		switch b {
		case nil:
			err = ctx.ShouldBind(m)
		case binding.JSON:
			err = ctx.ShouldBindJSON(m)
		case binding.XML:
			err = ctx.ShouldBindXML(m)
		case binding.Query:
			err = ctx.ShouldBindQuery(m)
		case binding.YAML:
			err = ctx.ShouldBindYAML(m)
		case binding.Header:
			err = ctx.ShouldBindHeader(m)
		case binding.Uri:
			err = ctx.ShouldBindUri(m)
		default:
			err = ctx.ShouldBindWith(m, b.(binding.Binding))
		}
	}

	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			errsMap := validatorx.Translate(errs, m.ErrorMessage())

			if len(errsMap) == 0 {
				responsex.ServerError(ctx, responsex.WithMsg(ErrValidateMessageTransformFailed.Error()))
				return false
			}

			for _, e := range errsMap {
				responsex.ValidateError(ctx, responsex.WithMsg(e))
				return false
			}

			return true
		}

		global.Logger().Error(err.Error())

		responsex.ServerError(ctx)

		return false
	}

	return true
}

// ShouldBindDefault *gin.Context.ShouldBind 方法的扩展
func ShouldBindDefault(ctx *gin.Context, m BindModel) bool {
	return shouldBind(ctx, m, nil, false)
}

// ShouldBindJSON *gin.Context.ShouldBindJSON 方法的扩展
func ShouldBindJSON(ctx *gin.Context, m BindModel) bool {
	return shouldBind(ctx, m, binding.JSON, false)
}

// ShouldBindXML *gin.Context.ShouldBindXML 方法的扩展
func ShouldBindXML(ctx *gin.Context, m BindModel) bool {
	return shouldBind(ctx, m, binding.XML, false)
}

// ShouldBindQuery *gin.Context.ShouldBindQuery 方法的扩展
func ShouldBindQuery(ctx *gin.Context, m BindModel) bool {
	return shouldBind(ctx, m, binding.Query, false)
}

// ShouldBindYAML *gin.Context.ShouldBindYAML 方法的扩展
func ShouldBindYAML(ctx *gin.Context, m BindModel) bool {
	return shouldBind(ctx, m, binding.YAML, false)
}

// ShouldBindHeader *gin.Context.ShouldBindHeader 方法的扩展
func ShouldBindHeader(ctx *gin.Context, m BindModel) bool {
	return shouldBind(ctx, m, binding.Header, false)
}

// ShouldBindUri *gin.Context.ShouldBindUri 方法的扩展
func ShouldBindUri(ctx *gin.Context, m BindModel) bool {
	return shouldBind(ctx, m, binding.Uri, false)
}

// ShouldBindWith *gin.Context.ShouldBindWith 方法的扩展
func ShouldBindWith(ctx *gin.Context, b binding.Binding, m BindModel) bool {
	return shouldBind(ctx, m, b, false)
}

// ShouldBindBodyWith *gin.Context.ShouldBindBodyWith 方法的扩展
func ShouldBindBodyWith(ctx *gin.Context, b binding.Binding, m BindModel) bool {
	return shouldBind(ctx, m, b, true)
}
