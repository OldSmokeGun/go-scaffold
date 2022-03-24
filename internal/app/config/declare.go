package config

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
