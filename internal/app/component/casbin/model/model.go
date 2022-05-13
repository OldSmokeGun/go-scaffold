package model

import (
	"github.com/casbin/casbin/v2/model"
)

type Config struct {
	Path string
}

func New(config *Config) (model.Model, error) {
	return model.NewModelFromFile(config.Path)
}
