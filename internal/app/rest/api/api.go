package api

//go:generate swag fmt -g api.go
//go:generate swag init -g api.go -o docs --parseInternal

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
