package config

import (
	"time"

	"github.com/samber/lo"
)

var supportedEnvs = []Env{Dev, Test, Prod}

// App application base config
type App struct {
	Timeout time.Duration `json:"timeout"`
}

func (App) GetName() string {
	return "app"
}

// AppName application name
type AppName string

func (e AppName) String() string {
	return string(e)
}

// Env application running environment
type Env string

func (e Env) String() string {
	return string(e)
}

// Check  the environment is set correctly
func (e Env) Check() bool {
	return lo.Contains(supportedEnvs, e)
}

// IsDebug report whether it is a debug environment
func (e Env) IsDebug() bool {
	return e == Dev || e == Test
}

const (
	Dev  Env = "dev"
	Test Env = "test"
	Prod Env = "prod"
)
