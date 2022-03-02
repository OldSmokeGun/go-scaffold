package handler

import (
	"github.com/google/wire"
	"go-scaffold/internal/app/transport/http/handler/v1/greet"
	"go-scaffold/internal/app/transport/http/handler/v1/trace"
	"go-scaffold/internal/app/transport/http/handler/v1/user"
)

//go:generate swag fmt -g handler.go
//go:generate swag init -g handler.go -o docs --parseInternal

// @title                       API 接口文档
// @description                 API 接口文档
// @version                     0.0.0
// @host                        localhost
// @BasePath                    /api
// @schemes                     http https
// @accept                      json
// @accept                      x-www-form-urlencoded
// @produce                     json
// @securityDefinitions.apikey  LoginAuth
// @in                          header
// @name                        Token

var ProviderSet = wire.NewSet(
	greet.New,
	trace.New,
	user.New,
)
