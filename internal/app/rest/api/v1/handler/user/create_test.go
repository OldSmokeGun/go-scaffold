package user

import (
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

// mockRequest 模拟请求
func mockRequest(t *testing.T, h Interface, req *CreateReq) *httptest.ResponseRecorder {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Handle("POST", "/api/v1/user", h.Create)

	jsonCreateReq, err := jsoniter.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("POST", "/api/v1/user", bytes.NewReader(jsonCreateReq))
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
		createReq := &CreateReq{
			Name:  "test",
			Age:   18,
			Phone: "13000000000",
		}

		createParam := new(user.CreateParam)
		if err := copier.Copy(createParam, createReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Create(createParam).
			Return(nil)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockRequest(t, newHandler, createReq)

		respBody := &responsex.Body{}
		if err := jsoniter.Unmarshal(w.Body.Bytes(), respBody); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, responsex.SuccessCode, respBody.Code)
		assert.Equal(t, responsex.SuccessCode.String(), respBody.Msg)
	})

	t.Run("req_validate_failed", func(t *testing.T) {

		t.Run("req_name_required", func(t *testing.T) {
			createReq := &CreateReq{"", 1, "13000000000"}
			newHandler := New()
			newHandler.Logger = test.Logger()
			newHandler.Service = nil

			w := mockRequest(t, newHandler, createReq)

			respBody := &responsex.Body{}
			if err := jsoniter.Unmarshal(w.Body.Bytes(), respBody); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, responsex.ValidateErrorCode, respBody.Code)
			assert.Equal(t, createReq.ErrorMessage()["Name.required"], respBody.Msg)
		})

		t.Run("req_age_min", func(t *testing.T) {
			createReq := &CreateReq{"test", 0, "13000000000"}
			newHandler := New()
			newHandler.Logger = test.Logger()
			newHandler.Service = nil

			w := mockRequest(t, newHandler, createReq)

			respBody := &responsex.Body{}
			if err := jsoniter.Unmarshal(w.Body.Bytes(), respBody); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, responsex.ValidateErrorCode, respBody.Code)
			assert.Equal(t, createReq.ErrorMessage()["Age.min"], respBody.Msg)
		})

		t.Run("req_phone_phone", func(t *testing.T) {
			createReq := &CreateReq{"test", 1, "100"}
			newHandler := New()
			newHandler.Logger = test.Logger()
			newHandler.Service = nil

			w := mockRequest(t, newHandler, createReq)

			respBody := &responsex.Body{}
			if err := jsoniter.Unmarshal(w.Body.Bytes(), respBody); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, responsex.ValidateErrorCode, respBody.Code)
			assert.Equal(t, createReq.ErrorMessage()["Phone.phone"], respBody.Msg)
		})
	})

	t.Run("req_create_error", func(t *testing.T) {
		createReq := &CreateReq{
			Name:  "test",
			Age:   18,
			Phone: "13000000000",
		}

		createParam := new(user.CreateParam)
		if err := copier.Copy(createParam, createReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Create(createParam).
			Return(user.ErrDataStoreFailed)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockRequest(t, newHandler, createReq)

		respBody := &responsex.Body{}
		if err := jsoniter.Unmarshal(w.Body.Bytes(), respBody); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, responsex.ServerErrorCode, respBody.Code)
		assert.Equal(t, user.ErrDataStoreFailed.Error(), respBody.Msg)
	})
}
