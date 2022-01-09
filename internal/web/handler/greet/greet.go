package greet

import "github.com/gin-gonic/gin"

type Handler interface {
	Hello(ctx *gin.Context)
}

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}
