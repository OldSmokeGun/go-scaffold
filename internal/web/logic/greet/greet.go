package greet

import (
	"github.com/gin-gonic/gin"
)

type Logic interface {
	Hello(param *HelloParam) (string, error)
}

type logic struct {
	ctx *gin.Context
}

func NewLogic(ctx *gin.Context) *logic {
	return &logic{
		ctx: ctx,
	}
}
