package greet

import (
	"github.com/gin-gonic/gin"
)

type Service interface {
	Hello(param *HelloParam) (string, error)
}

type service struct {
	ctx *gin.Context
}

func NewService(ctx *gin.Context) *service {
	return &service{
		ctx: ctx,
	}
}
