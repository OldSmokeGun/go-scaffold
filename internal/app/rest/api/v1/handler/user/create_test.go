package user

import (
	"bou.ke/monkey"
	"bytes"
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

// mockCreateRequest 模拟请求
func mockCreateRequest(t *testing.T, h Interface, req *CreateReq) *httptest.ResponseRecorder {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Handle("POST", "/api/v1/user", h.Create)

	jsonReq, err := jsoniter.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("POST", "/api/v1/user", bytes.NewReader(jsonReq))
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, r)

	return w
}

func TestCreateReq_ErrorMessage(t *testing.T) {
	cr := new(CreateReq)
	assert.NotNil(t, cr.ErrorMessage())
	assert.Greater(t, len(cr.ErrorMessage()), 0)
}

func Test_handler_Create(t *testing.T) {

	t.Run("create_success", func(t *testing.T) {
		createReq := &CreateReq{Name: "test", Age: 18, Phone: "13000000000"}

		createParam := new(user.CreateParam)
		if err := copier.Copy(createParam, createReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Create(createParam).
			Return(nil)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockCreateRequest(t, newHandler, createReq)

		respJson, err := jsoniter.Marshal(responsex.NewSuccessBody())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("req_validate_failed", func(t *testing.T) {

		t.Run("req_name_required", func(t *testing.T) {
			createReq := &CreateReq{"", 1, "13000000000"}
			newHandler := New()
			newHandler.Logger = test.Logger()
			newHandler.Service = nil

			w := mockCreateRequest(t, newHandler, createReq)

			respJson, err := jsoniter.Marshal(responsex.NewValidateErrorBody(responsex.WithMsg(createReq.ErrorMessage()["Name.required"])))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, string(respJson), w.Body.String())
		})

		t.Run("req_age_min", func(t *testing.T) {
			createReq := &CreateReq{"test", 0, "13000000000"}
			newHandler := New()
			newHandler.Logger = test.Logger()
			newHandler.Service = nil

			w := mockCreateRequest(t, newHandler, createReq)

			respJson, err := jsoniter.Marshal(responsex.NewValidateErrorBody(responsex.WithMsg(createReq.ErrorMessage()["Age.min"])))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, string(respJson), w.Body.String())
		})

		t.Run("req_phone_phone", func(t *testing.T) {
			createReq := &CreateReq{"test", 1, "100"}
			newHandler := New()
			newHandler.Logger = test.Logger()
			newHandler.Service = nil

			w := mockCreateRequest(t, newHandler, createReq)

			respJson, err := jsoniter.Marshal(responsex.NewValidateErrorBody(responsex.WithMsg(createReq.ErrorMessage()["Phone.phone"])))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, string(respJson), w.Body.String())
		})
	})

	t.Run("req_copy_error", func(t *testing.T) {
		createReq := &CreateReq{Name: "test", Age: 18, Phone: "13000000000"}
		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = nil

		monkey.Patch(copier.Copy, func(toValue interface{}, fromValue interface{}) error {
			return copier.ErrInvalidCopyDestination
		})
		defer monkey.Unpatch(copier.Copy)

		w := mockCreateRequest(t, newHandler, createReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("create_failed", func(t *testing.T) {
		createReq := &CreateReq{Name: "test", Age: 18, Phone: "13000000000"}

		createParam := new(user.CreateParam)
		if err := copier.Copy(createParam, createReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Create(createParam).
			Return(user.ErrDataStoreFailed)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockCreateRequest(t, newHandler, createReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody(responsex.WithMsg(user.ErrDataStoreFailed.Error())))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})
}
