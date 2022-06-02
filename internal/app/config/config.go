package config

import (
	"github.com/google/wire"
	"go-scaffold/internal/app/component/casbin"
	"go-scaffold/internal/app/component/orm"
	"go-scaffold/internal/app/component/redis"
	"go-scaffold/internal/app/component/trace"
	"go-scaffold/internal/app/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	wire.FieldsOf(
		new(*Config),
		"App",
		"HTTP",
		"DB",
		"Redis",
		"Trace",
		"JWT",
		"Casbin",
	),
	wire.FieldsOf(new(*casbin.Config), "Model", "Adapter"),
)

// SupportedEnvs 支持的环境
var SupportedEnvs = []string{Local.String(), Test.String(), Prod.String()}

// Env 当前运行环境
type Env string

func (e Env) String() string {
	return string(e)
}

const (
	Local Env = "local"
	Test  Env = "test"
	Prod  Env = "prod"
)

type Config struct {
	App    *App           `json:"app"`
	HTTP   *HTTP          `json:"http"`
	DB     *orm.Config    `json:"db"`
	Redis  *redis.Config  `json:"redis"`
	Trace  *trace.Config  `json:"trace"`
	JWT    *JWT           `json:"jwt"`
	Casbin *casbin.Config `json:"casbin"`
}

type App struct {
	Name    string `json:"name"`
	Env     Env    `json:"env"`
	Timeout int64  `json:"timeout"`
}

type HTTP struct {
	Network      string `json:"network"`
	Addr         string `json:"addr"`
	Timeout      int64  `json:"timeout"`
	ExternalAddr string `json:"externalAddr"`
}

type JWT struct {
	Key string `json:"key"`
}

// Loaded 配置加载后调用的钩子函数
func Loaded(logger *zap.Logger, configModel *Config) error {
	if configModel.Trace != nil {
		configModel.Trace.ServiceName = configModel.App.Name
		configModel.Trace.Env = configModel.App.Env.String()
		configModel.Trace.Timeout = configModel.App.Timeout
	}

	if configModel.Casbin != nil {
		if configModel.Casbin.Adapter != nil {
			if configModel.Casbin.Adapter.Gorm != nil {
				configModel.Casbin.Adapter.Gorm.SetMigration(func(db *gorm.DB) error {
					return (&model.CasbinRule{}).Migrate(db)
				})
			}
		}
	}

	return nil
}
