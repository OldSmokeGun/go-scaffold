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

// mockSaveRequest 模拟请求
func mockSaveRequest(t *testing.T, h Interface, req *SaveReq) *httptest.ResponseRecorder {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Handle("PUT", "/api/v1/user", h.Save)

	jsonReq, err := jsoniter.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("PUT", "/api/v1/user", bytes.NewReader(jsonReq))
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, r)

	return w
}

func TestSaveReq_ErrorMessage(t *testing.T) {
	cr := new(SaveReq)
	assert.NotNil(t, cr.ErrorMessage())
	assert.Greater(t, len(cr.ErrorMessage()), 0)
}

func Test_handler_Save(t *testing.T) {

	t.Run("save_success", func(t *testing.T) {
		saveReq := &SaveReq{ID: 1, Name: "test", Age: 18, Phone: "13000000000"}

		saveParam := new(user.SaveParam)
		if err := copier.Copy(saveParam, saveReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Save(saveParam).
			Return(nil)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockSaveRequest(t, newHandler, saveReq)

		respJson, err := jsoniter.Marshal(responsex.NewSuccessBody())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("req_validate_failed", func(t *testing.T) {

		t.Run("req_id_required", func(t *testing.T) {
			saveReq := &SaveReq{0, "test", 18, "13000000000"}
			newHandler := New()
			newHandler.Logger = test.Logger()
			newHandler.Service = nil

			w := mockSaveRequest(t, newHandler, saveReq)

			respJson, err := jsoniter.Marshal(responsex.NewValidateErrorBody(responsex.WithMsg(saveReq.ErrorMessage()["ID.required"])))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, string(respJson), w.Body.String())
		})

		t.Run("req_name_required", func(t *testing.T) {
			saveReq := &SaveReq{1, "", 1, "13000000000"}
			newHandler := New()
			newHandler.Logger = test.Logger()
			newHandler.Service = nil

			w := mockSaveRequest(t, newHandler, saveReq)

			respJson, err := jsoniter.Marshal(responsex.NewValidateErrorBody(responsex.WithMsg(saveReq.ErrorMessage()["Name.required"])))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, string(respJson), w.Body.String())
		})

		t.Run("req_age_min", func(t *testing.T) {
			saveReq := &SaveReq{1, "test", 0, "13000000000"}
			newHandler := New()
			newHandler.Logger = test.Logger()
			newHandler.Service = nil

			w := mockSaveRequest(t, newHandler, saveReq)

			respJson, err := jsoniter.Marshal(responsex.NewValidateErrorBody(responsex.WithMsg(saveReq.ErrorMessage()["Age.min"])))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, string(respJson), w.Body.String())
		})

		t.Run("req_phone_phone", func(t *testing.T) {
			saveReq := &SaveReq{1, "test", 1, "100"}
			newHandler := New()
			newHandler.Logger = test.Logger()
			newHandler.Service = nil

			w := mockSaveRequest(t, newHandler, saveReq)

			respJson, err := jsoniter.Marshal(responsex.NewValidateErrorBody(responsex.WithMsg(saveReq.ErrorMessage()["Phone.phone"])))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, string(respJson), w.Body.String())
		})
	})

	t.Run("req_copy_error", func(t *testing.T) {
		saveReq := &SaveReq{ID: 1, Name: "test", Age: 18, Phone: "13000000000"}
		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = nil

		monkey.Patch(copier.Copy, func(toValue interface{}, fromValue interface{}) error {
			return copier.ErrInvalidCopyDestination
		})
		defer monkey.Unpatch(copier.Copy)

		w := mockSaveRequest(t, newHandler, saveReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})

	t.Run("save_failed", func(t *testing.T) {
		saveReq := &SaveReq{ID: 1, Name: "test", Age: 18, Phone: "13000000000"}

		saveParam := new(user.SaveParam)
		if err := copier.Copy(saveParam, saveReq); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := user.NewMockService(ctrl)
		mockService.EXPECT().
			Save(saveParam).
			Return(user.ErrDataStoreFailed)

		newHandler := New()
		newHandler.Logger = test.Logger()
		newHandler.Service = mockService

		w := mockSaveRequest(t, newHandler, saveReq)

		respJson, err := jsoniter.Marshal(responsex.NewServerErrorBody(responsex.WithMsg(user.ErrDataStoreFailed.Error())))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, string(respJson), w.Body.String())
	})
}
