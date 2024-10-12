package server

import (
	"context"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	gserv "go-scaffold/internal/app/facade/server/grpc"
	hserv "go-scaffold/internal/app/facade/server/http"
	"go-scaffold/internal/config"
)

var ProviderSet = wire.NewSet(
	gserv.ProviderSet,
	hserv.ProviderSet,
	New,
)

// Server kratos server
type Server struct {
	server *kratos.App
}

// New build kratos server
func New(
	ctx context.Context,
	appName config.AppName,
	hs *http.Server,
	gs *grpc.Server,
	// discovery discovery.Discovery, // optional service registered
) *Server {
	hostname, _ := os.Hostname()

	options := []kratos.Option{
		kratos.Context(ctx),
		kratos.ID(hostname),
		kratos.Name(appName.String()),
	}

	var servers []transport.Server
	if hs != nil {
		servers = append(servers, hs)
	}
	if gs != nil {
		servers = append(servers, gs)
	}

	if len(servers) > 0 {
		options = append(options, kratos.Server(servers...))
	}

	// options = append(options, kratos.Registrar(discovery)) // optional service registered

	server := kratos.New(options...)

	return &Server{server}
}

// Start kratos server
func (t *Server) Start() error {
	return t.server.Run()
}

// Stop kratos server
func (t *Server) Stop() error {
	return t.server.Stop()
}
