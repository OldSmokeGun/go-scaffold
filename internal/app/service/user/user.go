package user

//go:generate mockgen -source=user.go -destination=user_mock.go -package=user -mock_names=Interface=MockService

import (
	"context"
	"errors"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/repository/user"
	"go.uber.org/zap"
)

type Interface interface {
	Create(ctx context.Context, param *CreateParam) error
	List(ctx context.Context, param *ListParam) (ListResult, error)
	Detail(ctx context.Context, param *DetailParam) (*DetailResult, error)
	Save(ctx context.Context, param *SaveParam) error
	Delete(ctx context.Context, param *DeleteParam) error
}

var (
	ErrDataStoreFailed  = errors.New("数据保存失败")
	ErrDataQueryFailed  = errors.New("数据查询失败")
	ErrDataDeleteFailed = errors.New("数据删除失败")
	ErrUserNotExist     = errors.New("用户不存在")
)

type service struct {
	Logger     *zap.Logger
	Repository user.Interface
}

func New() *service {
	return &service{
		Logger:     global.Logger(),
		Repository: user.New(),
	}
}
