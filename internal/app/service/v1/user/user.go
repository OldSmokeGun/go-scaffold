package user

import (
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	pb "go-scaffold/internal/app/api/v1/user"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/repository/user"
)

var (
	ErrDataStoreFailed  = errors.New("数据保存失败")
	ErrDataQueryFailed  = errors.New("数据查询失败")
	ErrDataDeleteFailed = errors.New("数据删除失败")
	ErrUserNotExist     = errors.New("用户不存在")
)

type Service struct {
	pb.UnimplementedUserServer
	logger *log.Helper
	conf   *config.Config
	repo   user.Interface
}

func NewService(
	logger log.Logger,
	conf *config.Config,
	repo user.Interface,
) *Service {
	return &Service{
		logger: log.NewHelper(logger),
		conf:   conf,
		repo:   repo,
	}
}
