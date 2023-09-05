package model

import (
	"github.com/casbin/casbin/v2/model"

	"go-scaffold/internal/config"
)

// New creates model
func New(conf config.CasbinModel) (model.Model, error) {
	return model.NewModelFromFile(conf.Path)
}
