package v1

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	kerr "github.com/go-kratos/kratos/v2/errors"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/trace"

	v1 "go-scaffold/internal/app/facade/server/grpc/api/v1"
	"go-scaffold/internal/config"
	berr "go-scaffold/internal/errors"
	"go-scaffold/internal/pkg/client"
	"go-scaffold/pkg/trace"
)

type TraceHandler struct {
	logger       *slog.Logger
	servicesConf config.Services
	hsConf       config.HTTPServer
	trace        *trace.Trace
	grpcClient   *client.GRPC
}

func NewTraceHandler(
	logger *slog.Logger,
	servicesConf config.Services,
	hsConf config.HTTPServer,
	trace *trace.Trace,
	grpcClient *client.GRPC,
) *TraceHandler {
	return &TraceHandler{
		logger:       logger,
		servicesConf: servicesConf,
		hsConf:       hsConf,
		trace:        trace,
		grpcClient:   grpcClient,
	}
}

// Example 示例方法
//
//	@Router			/v1/trace/example [post]
//	@Summary		示例接口
//	@Description	示例接口
//	@Tags			示例
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	example.Success				"成功响应"
//	@Failure		500	{object}	example.ServerError			"服务器出错"
//	@Failure		400	{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized		"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied	"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *TraceHandler) Example(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()

	h.example(reqCtx)

	conn, err := h.grpcClient.DialInsecure(reqCtx, h.servicesConf.Self)
	if err != nil {
		return errors.Wrap(err, "init grpc client connect error")
	}

	client := v1.NewGreetClient(conn)
	resp, err := client.Hello(reqCtx, &v1.GreetHelloRequest{Name: "Example"})
	if err != nil {
		e := kerr.FromError(err)
		return errors.Errorf("GRPC 调用错误: %s", e.Message)
	}
	h.logger.Info("请求结果", slog.String("msg", resp.Msg))

	// 获取当前请求 span
	span := sdktrace.SpanFromContext(otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Request().Header)))
	defer span.End()

	requestUrl := fmt.Sprintf("http://%s/api/v1/greet?name=tracer", h.hsConf.Addr)
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return berr.ErrInternalError.WithError(errors.WithStack(err))
	}

	// 携带 baggage 信息
	mem, err := baggage.NewMember("exampleKey", "exampleValue")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return berr.ErrInternalError.WithError(errors.WithStack(err))
	}

	bag, err := baggage.New(mem)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return berr.ErrInternalError.WithError(errors.WithStack(err))
	}
	reqCtx = baggage.ContextWithBaggage(reqCtx, bag)

	// 注入上下文信息
	otel.GetTextMapPropagator().Inject(reqCtx, propagation.HeaderCarrier(request.Header))

	if _, err = http.DefaultClient.Do(request); err != nil {
		err = fmt.Errorf("请求 %s 失败: %w", requestUrl, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.Errorf("请求 %s 失败", requestUrl)
	}

	return ctx.JSON(http.StatusOK, nil)
}

// example 示例方法
func (h *TraceHandler) example(ctx context.Context) {
	_, span := h.trace.Tracer("").Start(
		ctx,
		"TraceHandler.example",
		sdktrace.WithAttributes(
			attribute.String("exampleKey1", "exampleValue1"),
			attribute.String("exampleKey2", "exampleValue2"),
		),
		sdktrace.WithSpanKind(sdktrace.SpanKindInternal),
	)
	span.AddEvent(
		"exampleEvent",
		sdktrace.WithAttributes(
			attribute.String("exampleKey1", "exampleValue1"),
			attribute.String("exampleKey2", "exampleValue2"),
		),
		sdktrace.WithStackTrace(true),
	)

	if rand.Intn(10) > 5 {
		span.RecordError(errors.New("example rand error"))
		span.SetStatus(codes.Error, "example rand error")
	}

	defer span.End()
}
