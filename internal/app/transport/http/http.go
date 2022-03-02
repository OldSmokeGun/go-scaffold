package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go-scaffold/internal/app/config"
	"go-scaffold/internal/app/pkg/responsex"
	"go-scaffold/internal/app/service/v1/greet"
	"go-scaffold/internal/app/service/v1/user"
	"go-scaffold/internal/app/transport/http/handler"
	"go-scaffold/internal/app/transport/http/router"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"time"
)

var ProviderSet = wire.NewSet(
	handler.ProviderSet,
	router.ProviderSet,
	NewServer,
)

// NewServer 创建 HTTP 服务器
func NewServer(
	loggerWriter *rotatelogs.RotateLogs,
	logger log.Logger,
	zLogger *zap.Logger,
	cm *config.Config,
	router *gin.Engine,
	greetService *greet.Service,
	userService *user.Service,
) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(
				recovery.WithLogger(logger),
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					fmt.Println(33333333333333333)
					return errors.New(responsex.ServerErrorCode.String())
				}),
			),
			logging.Server(logger),
			tracing.Server(),
			validate.Validator(),
		),
		khttp.Logger(logger),
		khttp.ErrorEncoder(func(writer http.ResponseWriter, request *http.Request, err error) {
			fmt.Println(1111111111)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Header().Set("Content-Type", "application/json;charset=utf-8")
			body := responsex.NewServerErrorBody(responsex.WithMsg(err.Error()))
			jsonEncode := encoding.GetCodec("json")
			b, err := jsonEncode.Marshal(body)
			if err != nil {
				log.NewHelper(logger).Error(err)
				return
			}
			if _, err := writer.Write(b); err != nil {
				log.NewHelper(logger).Error(err)
				return
			}
		}),
		khttp.ResponseEncoder(func(writer http.ResponseWriter, request *http.Request, i interface{}) error {
			fmt.Println(222222222222)
			writer.Header().Set("Content-Type", "application/json;charset=utf-8")
			body := responsex.NewSuccessBody(responsex.WithData(nilToObject(i)))
			jsonEncode := encoding.GetCodec("json")
			b, err := jsonEncode.Marshal(body)
			if err != nil {
				return err
			}
			if _, err := writer.Write(b); err != nil {
				return err
			}
			return nil
		}),
	}

	if cm.App.Http.Network != "" {
		opts = append(opts, khttp.Network(cm.App.Http.Network))
	}

	if cm.App.Http.Addr != "" {
		opts = append(opts, khttp.Address(cm.App.Http.Addr))
	}

	if cm.App.Http.Timeout != 0 {
		opts = append(opts, khttp.Timeout(time.Duration(cm.App.Http.Timeout)*time.Second))
	}

	srv := khttp.NewServer(opts...)
	srv.HandlePrefix("/", router)
	// srv.HandlePrefix("/", router.New(loggerWriter, logger, zLogger, cm))
	// if cm.App.Env == "local" {
	// 	srv.HandlePrefix("/q/", openapiv2.NewHandler(openapiv2.WithGeneratorOptions()))
	// }
	//
	// greetpb.RegisterGreetHTTPServer(srv, greetService)
	// userpb.RegisterUserHTTPServer(srv, userService)

	return srv
}

func nilToObject(i interface{}) interface{} {
	if i != nil {
		val := reflect.ValueOf(i)
		switch val.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
			if val.IsNil() {
				i = struct{}{}
			}
		case reflect.Slice:
			if val.IsNil() {
				i = make([]interface{}, 0)
			}
		}
	} else {
		i = struct{}{}
	}

	return i
}
