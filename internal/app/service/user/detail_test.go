package user

import (
	"bou.ke/monkey"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/repository/user"
	"go-scaffold/internal/app/rest/pkg/responsex"
	"go-scaffold/internal/app/test"
	"gorm.io/gorm"
	"testing"
	"time"
)

func init() {
	test.Init()
}

func Test_service_Detail(t *testing.T) {

	t.Run("get_detail_success", func(t *testing.T) {
		detailParam := &DetailParam{ID: 1}
		columns := []string{"*"}

		userModel := &model.User{
			BaseModel: model.BaseModel{ID: 1, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(), DeletedAt: 0},
			Name:      "test",
			Age:       18,
			Phone:     "13000000000",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindOneByID(context.TODO(), detailParam.ID, columns).
			Return(userModel, nil)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		detailResult := new(DetailResult)
		if err := copier.Copy(detailResult, userModel); err != nil {
			t.Fatal(err)
		}

		ret, err := newService.Detail(context.TODO(), detailParam)

		assert.NoError(t, err)
		assert.Equal(t, detailResult, ret)
	})

	t.Run("get_detail_not_found", func(t *testing.T) {
		detailParam := &DetailParam{ID: 0}
		columns := []string{"*"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindOneByID(context.TODO(), detailParam.ID, columns).
			Return(nil, gorm.ErrRecordNotFound)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		ret, err := newService.Detail(context.TODO(), detailParam)

		assert.ErrorIs(t, err, ErrUserNotExist)
		assert.Nil(t, ret)
	})

	t.Run("get_detail_failed", func(t *testing.T) {
		detailParam := &DetailParam{ID: 1}
		columns := []string{"*"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindOneByID(context.TODO(), detailParam.ID, columns).
			Return(nil, errors.New("test error"))

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		ret, err := newService.Detail(context.TODO(), detailParam)

		assert.ErrorIs(t, err, ErrDataQueryFailed)
		assert.Nil(t, ret)
	})

	t.Run("result_copy_error", func(t *testing.T) {
		detailParam := &DetailParam{ID: 1}
		columns := []string{"*"}

		userModel := &model.User{
			BaseModel: model.BaseModel{ID: 1, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(), DeletedAt: 0},
			Name:      "test",
			Age:       18,
			Phone:     "13000000000",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindOneByID(context.TODO(), detailParam.ID, columns).
			Return(userModel, nil)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		monkey.Patch(copier.Copy, func(toValue interface{}, fromValue interface{}) error {
			return errors.New("test error")
		})
		defer monkey.Unpatch(copier.Copy)

		ret, err := newService.Detail(context.TODO(), detailParam)

		assert.EqualError(t, err, responsex.ServerErrorCode.String())
		assert.Nil(t, ret)
	})
}
