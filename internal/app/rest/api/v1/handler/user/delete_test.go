package user

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

// mockDeleteRequest 模拟请求
func mockDeleteRequest(t *testing.T, h Interface, req *DeleteReq) *httptest.ResponseRecorder {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Handle("DELETE", "/api/v1/user/:id", h.Delete)

	r := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/user/%d", req.ID), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, r)

	return w
}

func TestDeleteReq_ErrorMessage(t *testing.T) {
	cr := new(DeleteReq)
	assert.NotNil(t, cr.ErrorMessage())
	assert.Greater(t, len(cr.ErrorMessage()), 0)
}

func Test_handler_Delete(t *testing.T) {

	t.Run("delete_success", func(t *testing.T) {
		deleteReq := &DeleteReq{1}

		deleteParam := new(user.DeleteParam)
		if err := copier.Copy(deleteParam, deleteReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Delete(context.TODO(), deleteParam).
			Return(nil)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockDeleteRequest(t, newHandler, deleteReq)

		respJson, err := jsoniter.Marshal(responsex.NewSuccessBody())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("req_id_required", func(t *testing.T) {
		deleteReq := &DeleteReq{}

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = nil

		w := mockDeleteRequest(t, newHandler, deleteReq)

		respJson, err := jsoniter.Marshal(responsex.NewValidateErrorBody(responsex.WithMsg(deleteReq.ErrorMessage()["ID.required"])))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("req_copy_error", func(t *testing.T) {
		deleteReq := &DeleteReq{1}

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = nil

		monkey.Patch(copier.Copy, func(toValue interface{}, fromValue interface{}) error {
			return errors.New("test error")
		})
		defer monkey.Unpatch(copier.Copy)

		w := mockDeleteRequest(t, newHandler, deleteReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("delete_failed", func(t *testing.T) {
		deleteReq := &DeleteReq{100}

		deleteParam := new(user.DeleteParam)
		if err := copier.Copy(deleteParam, deleteReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Delete(context.TODO(), deleteParam).
			Return(user.ErrDataDeleteFailed)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockDeleteRequest(t, newHandler, deleteReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody(responsex.WithMsg(user.ErrDataDeleteFailed.Error())))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})
}
