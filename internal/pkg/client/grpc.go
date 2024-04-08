package client

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"

	"go-scaffold/internal/pkg/discovery"
)

// ErrDiscoveryIsNotSet discovery is not be set
var ErrDiscoveryIsNotSet = errors.New("the discovery is not set")

const (
	defaultTimeout          = 5 * time.Second
	discoveryProtocolPrefix = "discovery://"
)

// GRPC gRPC client
type GRPC struct {
	discovery discovery.Discovery
}

// NewGRPC build gRPC client
func NewGRPC() *GRPC {
	return &GRPC{}
}

func (c *GRPC) WithDiscovery(disc discovery.Discovery) *GRPC {
	c.discovery = disc
	return c
}

// Dial returns a GRPC connection
func (c *GRPC) Dial(ctx context.Context, endpoint string, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	clientOptions, err := c.defaultClientConfig(endpoint)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		clientOptions = append(clientOptions, opt)
	}

	return dial(ctx, false, clientOptions...)
}

// DialInsecure returns a insecure GRPC connection
func (c *GRPC) DialInsecure(ctx context.Context, endpoint string, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	clientOptions, err := c.defaultClientConfig(endpoint)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		clientOptions = append(clientOptions, opt)
	}

	return dial(ctx, true, clientOptions...)
}

func (c *GRPC) defaultClientConfig(endpoint string) ([]kgrpc.ClientOption, error) {
	clientOptions := []kgrpc.ClientOption{
		kgrpc.WithEndpoint(endpoint),
		kgrpc.WithTimeout(defaultTimeout),
		kgrpc.WithMiddleware(
			logging.Client(log.GetLogger()),
			tracing.Client(),
			metadata.Client(),
		),
	}

	if isDiscoveryProtocol(endpoint) {
		if c.discovery == nil {
			return nil, ErrDiscoveryIsNotSet
		}
		clientOptions = append(clientOptions, kgrpc.WithDiscovery(c.discovery))
	}

	return clientOptions, nil
}

func dial(ctx context.Context, insecure bool, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	if insecure {
		return kgrpc.DialInsecure(ctx, opts...)
	} else {
		return kgrpc.Dial(ctx, opts...)
	}
}

func isDiscoveryProtocol(endpoint string) bool {
	return strings.HasPrefix(endpoint, discoveryProtocolPrefix)
}
