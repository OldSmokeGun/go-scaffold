package transport

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/service"
	gtr "go-scaffold/internal/app/transport/grpc"
	htr "go-scaffold/internal/app/transport/http"
	"os"
)

var hostname, _ = os.Hostname()

var ProviderSet = wire.NewSet(
	repository.ProviderSet,
	service.ProviderSet,
	gtr.ProviderSet,
	htr.ProviderSet,
	New,
)

type Transport struct {
	logger *log.Helper
	server *kratos.App
}

func New(
	logger log.Logger,
	cm *config.Config,
	hs *http.Server,
	gs *grpc.Server,
	discovery *etcd.Registry,
) *Transport {
	options := []kratos.Option{
		kratos.ID(hostname),
		kratos.Name(cm.App.Name),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(hs, gs),
	}

	if discovery != nil {
		options = append(options, kratos.Registrar(discovery))
	}

	server := kratos.New(options...)

	return &Transport{
		logger: log.NewHelper(logger),
		server: server,
	}
}

func (t *Transport) Start() error {
	t.logger.Info("transport server starting ...")

	if err := t.server.Run(); err != nil {
		return err
	}
	return nil
}

func (t *Transport) Stop() error {
	if err := t.server.Stop(); err != nil {
		return err
	}

	t.logger.Info("transport server stopping ...")
	return nil
}
