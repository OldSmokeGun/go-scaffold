package user

import (
	"bou.ke/monkey"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"go-scaffold/internal/app/rest/pkg/responsex"
	"go-scaffold/internal/app/service/user"
	"go-scaffold/internal/app/test"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	test.Init()
}

// mockDetailRequest 模拟请求
func mockDetailRequest(t *testing.T, h Interface, req *DetailReq) *httptest.ResponseRecorder {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Handle("GET", "/api/v1/user/:id", h.Detail)

	r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/user/%d", req.ID), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, r)

	return w
}

func TestDetailReq_ErrorMessage(t *testing.T) {
	cr := new(DetailReq)
	assert.NotNil(t, cr.ErrorMessage())
	assert.Greater(t, len(cr.ErrorMessage()), 0)
}

func Test_handler_Detail(t *testing.T) {

	t.Run("get_detail_success", func(t *testing.T) {
		detailReq := &DetailReq{1}

		detailParam := new(user.DetailParam)
		if err := copier.Copy(detailParam, detailReq); err != nil {
			t.Fatal(err)
		}

		detailResult := &user.DetailResult{ID: 1, Name: "Tom", Age: 18, Phone: "13500000000"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Detail(detailParam).
			Return(detailResult, nil)

		detailResp := new(DetailResp)
		if err := copier.Copy(detailResp, detailResult); err != nil {
			t.Fatal(err)
		}

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockDetailRequest(t, newHandler, detailReq)

		respJson, err := jsoniter.Marshal(responsex.NewSuccessBody(responsex.WithData(detailResp)))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("req_id_required", func(t *testing.T) {
		detailReq := &DetailReq{}
		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = nil

		w := mockDetailRequest(t, newHandler, detailReq)

		respJson, err := jsoniter.Marshal(responsex.NewValidateErrorBody(responsex.WithMsg(detailReq.ErrorMessage()["ID.required"])))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("req_copy_error", func(t *testing.T) {
		detailReq := &DetailReq{1}
		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = nil

		monkey.Patch(copier.Copy, func(toValue interface{}, fromValue interface{}) error {
			return errors.New("test error")
		})
		defer monkey.Unpatch(copier.Copy)

		w := mockDetailRequest(t, newHandler, detailReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("result_copy_error", func(t *testing.T) {
		detailReq := &DetailReq{1}

		detailParam := new(user.DetailParam)
		if err := copier.Copy(detailParam, detailReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Detail(detailParam).
			Return(nil, nil)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockDetailRequest(t, newHandler, detailReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("get_detail_failed", func(t *testing.T) {
		detailReq := &DetailReq{1}

		detailParam := new(user.DetailParam)
		if err := copier.Copy(detailParam, detailReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Detail(detailParam).
			Return(nil, user.ErrUserNotExist)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockDetailRequest(t, newHandler, detailReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody(responsex.WithMsg(user.ErrUserNotExist.Error())))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})
}
