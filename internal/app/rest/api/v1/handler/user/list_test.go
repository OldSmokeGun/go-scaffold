package user

import (
	"bou.ke/monkey"
	"bytes"
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

// mockListRequest 模拟请求
func mockListRequest(t *testing.T, h Interface, req *ListReq) *httptest.ResponseRecorder {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Handle("GET", "/api/v1/users", h.List)

	jsonReq, err := jsoniter.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users"), bytes.NewReader(jsonReq))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	engine.ServeHTTP(w, r)

	return w
}

func TestListReq_ErrorMessage(t *testing.T) {
	cr := new(ListReq)
	assert.Nil(t, cr.ErrorMessage())
}

func Test_handler_List(t *testing.T) {

	t.Run("get_list_success", func(t *testing.T) {
		listReq := &ListReq{""}

		listParam := new(user.ListParam)

		listResult := user.ListResult{
			{ID: 1, Name: "Tom", Age: 18, Phone: "13500000000"},
			{ID: 2, Name: "Jack", Age: 19, Phone: "13800000000"},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			List(listParam).
			Return(listResult, nil)

		listResp := new(ListResp)
		if err := copier.Copy(listResp, listResult); err != nil {
			t.Fatal(err)
		}

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockListRequest(t, newHandler, listReq)

		respJson, err := jsoniter.Marshal(responsex.NewSuccessBody(responsex.WithData(listResp)))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("req_copy_req_error", func(t *testing.T) {
		listReq := &ListReq{""}
		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = nil

		monkey.Patch(copier.Copy, func(toValue interface{}, fromValue interface{}) error {
			return copier.ErrInvalidCopyDestination
		})
		defer monkey.Unpatch(copier.Copy)

		w := mockListRequest(t, newHandler, listReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("get_list_error", func(t *testing.T) {
		listReq := &ListReq{""}

		listParam := new(user.ListParam)
		if err := copier.Copy(listParam, listReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			List(listParam).
			Return(nil, user.ErrDataQueryFailed)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockListRequest(t, newHandler, listReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody(responsex.WithMsg(user.ErrDataQueryFailed.Error())))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})
}
