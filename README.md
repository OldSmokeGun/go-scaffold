# `Go` 开发基础脚手架

## 架构图

![image](./docs/images/脚手架架构图.png)

## 目录结构

- `bin`：二进制文件输入目录
- `cmd`：编译入口
  - `app`：主程序编译入口
  - 辅助程序编译入口目录（如 `ctl` 等命令行生成工具）
  - ...
- `deploy`：环境和部署相关目录
  - `docker-compose`：`docker-compose` 容器编排目录
  - `kubernetes` 编排配置目录
  - ...
- `docs`：文档目录
- `etc`：配置文件目录
- `internal`：
  - `app`：主程序逻辑代码
    - `cli`：命令行功能
      - `command`：命令行功能入口
        - 模块目录
        - ...
      - `pkg`：功能增强包目录
        - ...
      - `script`：临时脚本
    - `config`：主程序配置模型
    - `global`：全局对象
    - `model`：数据库模型
    - `pkg`：功能增强包目录
      - ...
    - `rest`：`HTTP` 功能
      - `api`：`api` 文档相关定义
        - `docs`：`swagger` 生成的定义文件
      - `config`：`rest` 配置模型
      - `handler`：`HTTP` 请求处理入口
        - 模块目录
        - ...
      - `middleware`：中间件目录
      - `pkg`：功能增强包目录
        - ...
      - `router`：路由定义目录
    - `rpc`：`RPC` 功能
    - `service`：逻辑代码
      - 模块目录
      - ...
    - ...
- `logs`：日志文件生成目录
- `pkg`：功能类库
  - ...
  
# 如何运行

## `go build` 或 `go run`

1. `go build` 方式

```shell
$ go build -o bin/app cmd/app/main.go 
$ ./bin/app
```

2. `go run` 方式

```shell
$ go run cmd/app/main.go 
```

## `make`

```shell
 # 下载依赖
$ make download
$ make build

# 或依据平台编译
$ make linux-build
$ make windows-build
$ make mac-build

# 运行
$ ./bin/app
```

## `docker-compose`

`docker-compose` 的启动方式有两种，一种是基于 `air` 镜像，一种是基于 `Dockerfile` 来构建镜像

- 基于 `air` 镜像的方式只适用于开发阶段，请勿用于生产环境
- 基于 `Dockerfile` 的方式如果用于开发阶段，修改的代码将不会更新，除非在 `docker-compose` 启动时指定 `--build` 参数，但是这将会导致每次启动时都重新构建镜像，可能需要等待很长时间

> 注意：基于 `air` 镜像启动时，在 `Windows` 系统环境下，热更新可能不会生效，这是因为 `fsnotify` 无法收到 `wsl` 文件系统的变更通知

```shell
# 基于 air 
$ docker-compose -f deploy/docker-compose/docker-compose-dev.yaml up

# 基于 Dockerfile
$ docker-compose -f deploy/docker-compose/docker-compose.yaml up
```

## 热重启

热重启功能基于 `air` 包，[文档地址](https://github.com/cosmtrek/air)

```shell
$ air
```

## 运行命令行程序或脚本

命令行程序功能基于 `cobra` 包，[文档地址](https://github.com/cosmtrek/air)

```shell
$ ./bin/app [标志] <子命令> [标志] [参数]

# 帮助信息
$ ./bin/app -h
$ ./bin/app <子命令> -h
```

# 全局对象

在程序启动时，会根据配置文件，将数据库、`Redis`、日志等实例初始化到全局对象上

全局对象实例的定义位于 `internal/app/global/global.go` 中

```go
var (
	loggerOutput *rotatelogs.RotateLogs // 日志输出实例
	conf         *config.Config         // 配置实例
	logger       *zap.Logger            // 日志实例
	db           *gorm.DB               // 数据库实例
	redisClient  *redis.Client          // redis 实例
)
```

# 配置

配置文件位于 `etc` 目录中，如果位于其它位置，需要在运行程序时通过 `-f` 参数指定其位置。默认配置文件路径：`etc/app.yaml`

配置文件的内容在程序启动的最初就会被加载到 `app` 全局配置实例中，为了方便管理和维护各个子服务的配置，在 `app` 配置实例中，对各个子服务的配置进行了拆分，分散到其独立的配置实例中

`app` 全局配置定义

```go
type (
	Config struct {
		App  `mapstructure:",squash"`
		REST restconfig.Config `mapstructure:"REST"` // rest 服务配置
	}

	App struct {
		Env              Env
		ShutdownWaitTime int
		Log              *logger.Config
		DB               *orm.Config
		Redis            *redisclient.Config
	}
)
```

`rest` 服务配置定义

```go
type (
	Config struct {
		Host        string
		Port        int
		ExternalUrl string
		Jwt         Jwt
	}

	Jwt struct {
		Key    string
		Expire time.Duration
	}
)
```

在配置文件更改后，无需重启服务，更改内容会自动同步到配置实例中

- `app` 全局配置实例的定义位于 `internal/app/config/config.go` 中
- `rest` 服务配置实例的定义位于 `internal/app/rest/config/config.go` 中

不提供直接使用 `viper` 包直接通过键名来获取配置值的方法，而统一使用 `global.Config()` 来获取

# 日志

日志实例获取：`global.Logger()`

日志基于 `zap` 包，[文档地址](https://github.com/uber-go/zap)

# `rest` 服务

`rest` 服务基于 `gin` 提供 `HTTP` 的相关服务，[文档地址](https://github.com/gin-gonic/gin)

## 数据校验和绑定

在 `internal/app/rest/pkg/bindx` 包中，对 `gin` 的 `ShouldBind` 系列方法，以及参数校验，`validator` 包错误信息的翻译进行了封装

在参数校验失败后，会对错误信息进行翻译，然后进行 `JSON` 响应

`ShouldBind` 系列方法需要绑定的结构体必须实现 `bindx` 包中的 `BindModel` 接口，`ErrorMessage` 方法返回一个 `map`，此 `map` 是“字段.校验规则”与错误信息的键值对

接口定义：

```go
type BindModel interface {
    ErrorMessage() map[string]string
}
```

例：

```go
type HelloReq struct {
    Name string `form:"name" binding:"required"`
}

func (HelloReq) ErrorMessage() map[string]string {
	return map[string]string{
		"Name.required": "名称不能为空",
	}
}
```

绑定方法会返回一个布尔值表示参数是否校验成功，如果失败，则应立即 `return` 中止函数的执行

函数定义（和 `gin` 的`ShouldBind` 系列方法对应）：

```go
func ShouldBindDefault(ctx *gin.Context, m BindModel)
func ShouldBindJSON(ctx *gin.Context, m BindModel)
func ShouldBindXML(ctx *gin.Context, m BindModel)
func ShouldBindQuery(ctx *gin.Context, m BindModel)
func ShouldBindYAML(ctx *gin.Context, m BindModel)
func ShouldBindHeader(ctx *gin.Context, m BindModel)
func ShouldBindUri(ctx *gin.Context, m BindModel)
func ShouldBindWith(ctx *gin.Context, b binding.Binding, m BindModel)
func ShouldBindBodyWith(ctx *gin.Context, b binding.Binding, m BindModel)
```

例：

```go
// 自动根据请求头的 Content-Type 来进行对应的绑定
if !bindx.ShouldBindDefault(ctx, req) {
    return
}

// 绑定 JSON 类型的请求数据
if !bindx.ShouldBindJSON(ctx, req) {
    return
}

// 绑定 form-data 和 x-www-form-urlencoded 类型的请求数据
if !bindx.ShouldBindQuery(ctx, req) {
    return
}
```

## 响应

在 `internal/app/rest/pkg/responsex` 包中，对 `JSON` 数据的响应进行了封装，统一了 `HTTP` 响应码，错误码的返回

函数定义：

```go
// Success 成功响应
func Success(ctx *gin.Context, ops ...OptionFunc)
func ServerError(ctx *gin.Context, ops ...OptionFunc)
func ClientError(ctx *gin.Context, ops ...OptionFunc)
func ValidateError(ctx *gin.Context, ops ...OptionFunc)
func Unauthorized(ctx *gin.Context, ops ...OptionFunc)
func PermissionDenied(ctx *gin.Context, ops ...OptionFunc)
func ResourceNotFound(ctx *gin.Context, ops ...OptionFunc)
func TooManyRequest(ctx *gin.Context, ops ...OptionFunc)
```

例：

> 注意：调用方法后，需 `return` 结束方法继续执行

```go
bindx.Success(ctx) // 成功响应
bindx.Success(ctx, bindx.WithData(data)) // 返回数据

bindx.ServerError(ctx) // 服务器错误响应
bindx.ServerError(ctx, bindx.WithMsg(msg)) // 返回错误信息

bindx.ClientError(ctx) // 客户端错误响应
bindx.ClientError(ctx, bindx.WithMsg(msg)) // 返回错误信息

bindx.ValidateError(ctx) // 参数校验错误响应
bindx.ValidateError(ctx, bindx.WithMsg(msg)) // 返回错误信息

// ...
```

## `swagger` 文档生成

`swagger` 文档的生成基于 `swag` 包，[文档地址](https://github.com/swaggo/swag)

`swagger` 文档统一生成到 `internal/app/rest/api/docs` 目录下，否则无法访问

生成 `swagger` 文档的方式有三种

1. `swag` 命令方式

```shell
$ swag fmt -d internal/app -g app.go
$ swag init -d internal/app -g app.go -o internal/app/rest/api/docs
```

2. `make` 方式

```shell
$ make doc
```

3. `go generate` 方式

```shell
$ go generate
```

## 如何访问 `swagger` 文档

浏览器打开 `<yourAddress>:/docs`