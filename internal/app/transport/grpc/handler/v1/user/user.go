package user

import (
	"github.com/go-kratos/kratos/v2/log"
	pb "go-scaffold/internal/app/api/scaffold/v1/user"
	userrepo "go-scaffold/internal/app/repository/user"
	usersvc "go-scaffold/internal/app/service/user"
)

var _ pb.UserServer = (*Handler)(nil)

type Handler struct {
	pb.UnimplementedUserServer
	logger  *log.Helper
	service *usersvc.Service
	repo    userrepo.RepositoryInterface
}

func NewHandler(
	logger log.Logger,
	service *usersvc.Service,
	repo userrepo.RepositoryInterface,
) *Handler {
	return &Handler{
		logger:  log.NewHelper(logger),
		service: service,
		repo:    repo,
	}
}
