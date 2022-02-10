package trace

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-scaffold/internal/app/rest/pkg/responsex"
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
func (h *handler) Example(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	h.example(reqCtx)

	// 获取当前请求 span
	span := trace.SpanFromContext(otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Request.Header)))
	defer span.End()

	requestUrl := fmt.Sprintf("http://localhost:%d/api/v1/greet?name=tracer", h.Config.REST.Port)
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		responsex.ServerError(ctx)
		return
	}

	// 携带 baggage 信息
	mem, err := baggage.NewMember("exampleKey", "exampleValue")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		responsex.ServerError(ctx)
		return
	}

	bag, err := baggage.New(mem)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		responsex.ServerError(ctx)
		return
	}
	reqCtx = baggage.ContextWithBaggage(reqCtx, bag)

	// 注入上下文信息
	otel.GetTextMapPropagator().Inject(reqCtx, propagation.HeaderCarrier(request.Header))

	if _, err = http.DefaultClient.Do(request); err != nil {
		err = fmt.Errorf("请求 %s 失败: %w", requestUrl, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		responsex.ServerError(ctx, responsex.WithMsg(fmt.Sprintf("请求 %s 失败", requestUrl)))
		return
	}

	responsex.Success(ctx)
	return
}

// example 示例方法
func (h *handler) example(ctx context.Context) {
	_, span := h.Tracer.Start(
		ctx,
		"handler.example",
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
