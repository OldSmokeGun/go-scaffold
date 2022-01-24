package user

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/repository/user"
	"go-scaffold/internal/app/test"
	"gorm.io/gorm"
	"testing"
	"time"
)

func init() {
	test.Init()
}

func Test_service_Delete(t *testing.T) {

	t.Run("delete_success", func(t *testing.T) {
		deleteParam := &DeleteParam{ID: 1}
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
		gomock.InOrder(
			mockRepository.EXPECT().
				FindOneByID(deleteParam.ID, columns).
				Return(userModel, nil),

			mockRepository.EXPECT().
				Delete(userModel).
				Return(nil),
		)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Delete(deleteParam)

		assert.NoError(t, err)
	})

	t.Run("delete_find_one_not_found", func(t *testing.T) {
		deleteParam := &DeleteParam{ID: 0}
		columns := []string{"*"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindOneByID(deleteParam.ID, columns).
			Return(nil, gorm.ErrRecordNotFound)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Delete(deleteParam)

		assert.ErrorIs(t, err, ErrUserNotExist)
	})

	t.Run("delete_find_one_error", func(t *testing.T) {
		deleteParam := &DeleteParam{ID: 1}
		columns := []string{"*"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindOneByID(deleteParam.ID, columns).
			Return(nil, gorm.ErrInvalidValue)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Delete(deleteParam)

		assert.EqualError(t, err, ErrDataQueryFailed.Error())
	})

	t.Run("delete_failed", func(t *testing.T) {
		deleteParam := &DeleteParam{ID: 1}
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
		gomock.InOrder(
			mockRepository.EXPECT().
				FindOneByID(deleteParam.ID, columns).
				Return(userModel, nil),

			mockRepository.EXPECT().
				Delete(userModel).
				Return(gorm.ErrInvalidValue),
		)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Delete(deleteParam)

		assert.ErrorIs(t, err, ErrDataDeleteFailed)
	})
}
