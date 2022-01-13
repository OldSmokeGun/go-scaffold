# `Go` 开发基础脚手架

## 架构图

![image](./docs/images/脚手架架构图.png)

## 目录结构

- `bin`：二进制文件输入目录
- `cmd`：编译入口
  - `app`：主程序编译入口
  - 辅助程序编译入口目录（如 `ctl` 等命令行生成工具）
  - ...
- `config`：配置文件目录
- `deploy`：环境和部署相关目录
  - `docker-compose`：`docker-compose` 容器编排目录
  - `kubernetes` 编排配置目录
  - ...
- `docs`：文档目录
- `internal`：
  - `app`：主程序逻辑代码
    - `config`：主程序配置模型
    - `global`：全局对象
    - `logic`：逻辑代码
      - 模块目录
      - ...
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
    - ...
- `logs`：日志文件生成目录
- `pkg`：功能类库
  - ...

# 开发规范

## 命名规范

### 包定义

- `Go` 语言中的包是以目录为最小粒度来隔离代码的，所以，目录的定义必须简短、有意义
- 不要使用 `common`、`util` 或者 `lib` 这类宽泛的、无意义的包名，应该从包的名字就能直接理解这个包的作用，例如：如果要编写操作字符串的公共包，那么这个包的名称应为 `stringx`，表明是对 `string` 操作的增强
- 包名和目录名称保持一致
- 包名应为小写单词，不要使用下划线或者混合大小写
- 包名不要使用复数

### 文件命名

- 文件的名称应简短有意义
- 文件名应该为小写，并且单词之间使用下划线分隔

### 变量命名

- 变量的命名采用驼峰的风格
- 全局变量需要特别注意首字母的大小写，在编码时，首字母应小写，只有在其它包需要引用此变量时，才修改其首字母为大写，保持包的封闭
- 若变量类型为布尔类型，则名称应以 `Has`、`Is`、`Can` 或 `Allow` 等开头
- 如果是变量是定义的 `error`，变量名称应以 `Err` 开头，例如：`ErrFileIsNotExist`

### 结构体命名

- 结构体的命名采用驼峰的风格
- 在编码时，结构体名称首字母应小写，只有在需要时才修改为大写
- 结构体的属性命名也是同样的规则

### 函数和方法命名

- 驼峰风格
- 只有在需要时首字母才为大写
- 方法的接收者名称应为结构体的首字母小写，例如：`User` 的 `Hello` 方法，其方法接收者名称应为 `u`

### 接口命名

- 驼峰风格，只有在需要时首字母才为大写
- 接口的名称尽量以 `er` 为后缀，例如：`Writer`、`Reader`

### 常量命名

- 常量名称均为大写，并使用下划线分隔
- 定义枚举时，重新定义其基础类型

例如：

```go
type Color string

const (
	Red Color = "Red"
	Green Color = "Green"
)
```

## 错误处理

对于错误处理，采用防御式编程的原则

- 有 `error` 必须处理，不允许使用下划线 `_` 丢弃变量
- 如果不需要处理 `error`，必须向上层抛出
- 函数或方式如果要返回 `error`，必须是所有返回值的最后一个
- 如果函数正常返回值，`error` 必须为 `nil`，如果返回 `error`，则正常值为空

## 文档注释

书写注释时，双下划线 `//` 应与后面的注释内容保持一个空格

例如：

```go
// 示例注释
```

导出的变量、结构体、函数等，需要为其写注释时，根据规范，双下划线后面的注释内容应以变量、结构体或函数的名称开头，并且，名称与后面的注释内容保持一个空格

例如：

```go
// Foo 示例结构体
type Foo struct {}
```

## 业务分层

不采用传统的 `MVC` 架构模式，而是一个业务模块在 `handler` 和 `logic` 中就是一个单独的包

这样做的原因是，通过目录来强制在编写业务模块的逻辑时达到高内聚，低耦合的目标

如果两个包互相引用造成循环依赖，那么应该考虑这部分逻辑分层的合理性，并且应该将其抽离为一个新的包

`handler` 和 `logic` 中，一个业务模块的方法就是一个文件，文件名为方法的名称，这是为了减少单个文件中的代码量，能够一眼找到某个业务的方法定义在什么地方

`handler` 中业务的方法需要书写 `swagger` 注释，通常会很长，将方法拆分到独立的文件中也能排除维护代码时的干扰

## `goroutine`

由于 `Go` 中创建协程非常容易，所以经常会导致滥用

- 在使用 `goroutine` 时，必须要考虑 `goroutine` 的退出和回收
- 使用 `goroutine` 时尽量搭配 `sync` 等同步原语或者 `context` 包，以防止 `goroutine` 的泄露

## 单元测试

为了写出能够测试和方便测试的代码，需遵循以下原则

- 不允许在函数和方法中直接使用外部依赖，而是应该在构造函数中进行注入
- 在编写代码时应采用面向接口编程的原则，例如：所有 `handler` 和 `logic` 中的业务模块都必须定义描述此业务功能的接口，然后定义此业务模块的实体，此实体必须实现接口

必须定义接口和实现接口的原因是，需要利用接口的多态性，在编写单元测试时，单元测试框架能够根据接口 `mock` 假对象进行测试

- **单元测试不允许发生真实的 `HTTP` 和数据库请求**

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

```shell
$ air
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

配置文件位于 `config` 目录中，如果位于其它位置，需要在运行程序时通过 `-f` 参数指定其位置。默认配置文件路径：`config/app.yaml`

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

`rest` 服务提供 `HTTP` 的相关服务

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
