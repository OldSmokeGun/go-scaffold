package grpc

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"go-scaffold/internal/app/component/discovery"
	"google.golang.org/grpc"
	"strings"
	"time"
)

var ErrDiscoveryIsNotSet = errors.New("the discovery is not set")

const (
	defaultTimeout          = 5 * time.Second
	discoveryProtocolPrefix = "discovery://"
)

type Client struct {
	logger    log.Logger
	discovery discovery.Discovery
}

func New(logger log.Logger, disc discovery.Discovery) *Client {
	return &Client{
		logger:    logger,
		discovery: disc,
	}
}

func dial(ctx context.Context, insecure bool, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	if insecure {
		return kgrpc.DialInsecure(ctx, opts...)
	} else {
		return kgrpc.Dial(ctx, opts...)
	}
}

func (c *Client) defaultClientConfig(endpoint string) ([]kgrpc.ClientOption, error) {
	clientOptions := []kgrpc.ClientOption{
		kgrpc.WithEndpoint(endpoint),
		kgrpc.WithLogger(c.logger),
		kgrpc.WithTimeout(defaultTimeout),
		kgrpc.WithMiddleware(
			logging.Client(c.logger),
			tracing.Client(),
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

func (c *Client) Dial(ctx context.Context, endpoint string, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	clientOptions, err := c.defaultClientConfig(endpoint)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		clientOptions = append(clientOptions, opt)
	}

	return dial(ctx, false, clientOptions...)
}

func (c *Client) DialInsecure(ctx context.Context, endpoint string, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	clientOptions, err := c.defaultClientConfig(endpoint)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		clientOptions = append(clientOptions, opt)
	}

	return dial(ctx, true, clientOptions...)
}

func isDiscoveryProtocol(endpoint string) bool {
	return strings.HasPrefix(endpoint, discoveryProtocolPrefix)
}
