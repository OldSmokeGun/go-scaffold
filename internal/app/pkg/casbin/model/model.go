package model

import (
	"go-scaffold/internal/config"

	"github.com/casbin/casbin/v2/model"
)

// New creates model
func New(conf config.CasbinModel) (model.Model, error) {
	return model.NewModelFromFile(conf.Path)
}
