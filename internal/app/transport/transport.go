package transport

import (
	"go-scaffold/internal/app/component/discovery"
	"go-scaffold/internal/app/config"
	gtr "go-scaffold/internal/app/transport/grpc"
	htr "go-scaffold/internal/app/transport/http"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
)

var hostname, _ = os.Hostname()

var ProviderSet = wire.NewSet(
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
	appConf *config.App,
	hs *http.Server,
	gs *grpc.Server,
	discovery discovery.Discovery,
) *Transport {
	var servers []transport.Server
	if hs != nil {
		servers = append(servers, hs)
	}
	if gs != nil {
		servers = append(servers, gs)
	}

	options := []kratos.Option{
		kratos.ID(hostname),
		kratos.Name(appConf.Name),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(servers...),
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
