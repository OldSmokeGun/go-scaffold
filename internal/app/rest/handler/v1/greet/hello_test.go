package greet

import (
	"bou.ke/monkey"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"go-scaffold/internal/app/rest/pkg/responsex"
	"go-scaffold/internal/app/service/greet"
	"go-scaffold/internal/app/test"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	test.Init()
}

// mockHelloRequest 模拟请求
func mockHelloRequest(t *testing.T, h Interface, req *HelloReq) *httptest.ResponseRecorder {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Handle("GET", "/api/v1/greet", h.Hello)

	r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/greet?name=%s", req.Name), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, r)

	return w
}

func TestHelloReq_ErrorMessage(t *testing.T) {
	cr := new(HelloReq)
	assert.NotNil(t, cr.ErrorMessage())
	assert.Greater(t, len(cr.ErrorMessage()), 0)
}

func Test_handler_Hello(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		helloReq := &HelloReq{"Tom"}

		helloParam := new(greet.HelloParam)
		if err := copier.Copy(helloParam, helloReq); err != nil {
			t.Fatal(err)
		}

		helloResult := fmt.Sprintf("Hello, %s!", helloParam.Name)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := greet.NewMockService(ctrl)
		mockService.EXPECT().
			Hello(context.TODO(), helloParam).
			Return(helloResult, nil)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockHelloRequest(t, newHandler, helloReq)

		helloResp := HelloResp{helloResult}

		respJson, err := jsoniter.Marshal(responsex.NewSuccessBody(responsex.WithData(helloResp)))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("req_name_required", func(t *testing.T) {
		helloReq := &HelloReq{}

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = nil

		w := mockHelloRequest(t, newHandler, helloReq)

		respJson, err := jsoniter.Marshal(responsex.NewValidateErrorBody(responsex.WithMsg(helloReq.ErrorMessage()["Name.required"])))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("req_copy_error", func(t *testing.T) {
		helloReq := &HelloReq{"Tom"}

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = nil

		monkey.Patch(copier.Copy, func(toValue interface{}, fromValue interface{}) error {
			return errors.New("test error")
		})
		defer monkey.Unpatch(copier.Copy)

		w := mockHelloRequest(t, newHandler, helloReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("call_hello_failed", func(t *testing.T) {
		helloReq := &HelloReq{"Tom"}

		helloParam := new(greet.HelloParam)
		if err := copier.Copy(helloParam, helloReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := greet.NewMockService(ctrl)
		mockService.EXPECT().
			Hello(context.TODO(), helloParam).
			Return("", errors.New("test error"))

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockHelloRequest(t, newHandler, helloReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody(responsex.WithMsg("test error")))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})
}
