package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"go-scaffold/internal/app/component/casbin/adapter"
)

type Config struct {
	Model   *Model
	Adapter *adapter.Config
}

type Model struct {
	Path string
}

func New(m *Model, adp adapter.Adapter) (*casbin.Enforcer, error) {
	if m == nil && adp == nil {
		return nil, nil
	}

	cm, err := model.NewModelFromFile(m.Path)
	if err != nil {
		return nil, err
	}

	ef, err := casbin.NewEnforcer(cm, adp)
	if err != nil {
		return nil, err
	}

	return ef, nil
}
