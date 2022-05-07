package trace

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/api/scaffold/v1/greet"
	errorsx "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/transport/http/pkg/response"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	"net/http"
)

// Example 示例方法
// @Router       /v1/trace [get]
// @Summary      示例接口
// @Description  示例接口
// @Tags         示例
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Success      200  {object}  example.Success           "成功响应"
// @Failure      500  {object}  example.ServerError       "服务器出错"
// @Failure      400  {object}  example.ClientError       "客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
// @Failure      401  {object}  example.Unauthorized      "登陆失效"
// @Failure      403  {object}  example.PermissionDenied  "没有权限"
// @Failure      404  {object}  example.ResourceNotFound  "资源不存在"
// @Failure      429  {object}  example.TooManyRequest    "请求过于频繁"
func (h *Handler) Example(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	h.example(reqCtx)

	conn, err := h.grpcClient.DialInsecure(reqCtx, h.conf.Services.Self)
	if err != nil {
		h.logger.Error(err)
		response.Error(ctx, errorsx.ServerError())
		return
	}

	client := greet.NewGreetClient(conn)
	resp, err := client.Hello(reqCtx, &greet.HelloRequest{Name: "Example"})
	if err != nil {
		e := errorsx.FromGRPCError(err)
		response.Error(ctx, fmt.Errorf("GRPC 调用错误：%s", e.Message))
		return
	}
	h.logger.Infof("请求结果：%s", resp.Msg)

	// 获取当前请求 span
	span := trace.SpanFromContext(otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Request.Header)))
	defer span.End()

	requestUrl := fmt.Sprintf("http://%s/api/v1/greet?name=tracer", h.conf.App.Http.Addr)
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		response.Error(ctx, errorsx.ServerError())
		return
	}

	// 携带 baggage 信息
	mem, err := baggage.NewMember("exampleKey", "exampleValue")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		response.Error(ctx, errorsx.ServerError())
		return
	}

	bag, err := baggage.New(mem)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		response.Error(ctx, errorsx.ServerError())
		return
	}
	reqCtx = baggage.ContextWithBaggage(reqCtx, bag)

	// 注入上下文信息
	otel.GetTextMapPropagator().Inject(reqCtx, propagation.HeaderCarrier(request.Header))

	if _, err = http.DefaultClient.Do(request); err != nil {
		err = fmt.Errorf("请求 %s 失败: %w", requestUrl, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		response.Error(ctx, fmt.Errorf("请求 %s 失败", requestUrl))
		return
	}

	response.Success(ctx)
	return
}

// example 示例方法
func (h *Handler) example(ctx context.Context) {
	_, span := h.trace.Tracer("").Start(
		ctx,
		"Handler.example",
		trace.WithAttributes(
			attribute.String("exampleKey1", "exampleValue1"),
			attribute.String("exampleKey2", "exampleValue2"),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	span.AddEvent(
		"exampleEvent",
		trace.WithAttributes(
			attribute.String("exampleKey1", "exampleValue1"),
			attribute.String("exampleKey2", "exampleValue2"),
		),
		trace.WithStackTrace(true),
	)

	if rand.Intn(10) > 5 {
		span.RecordError(errors.New("example rand error"))
		span.SetStatus(codes.Error, "example rand error")
	}

	defer span.End()
}
