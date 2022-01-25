package user

import (
	"bou.ke/monkey"
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

func Test_service_Save(t *testing.T) {

	t.Run("save_success", func(t *testing.T) {
		saveParam := &SaveParam{ID: 1, Name: "test", Age: 18, Phone: "13000000000"}
		columns := []string{"*"}

		userModel := new(model.User)
		if err := copier.Copy(userModel, saveParam); err != nil {
			t.Fatal(err)
		}
		userModel.CreatedAt = time.Now().Unix()
		userModel.UpdatedAt = time.Now().Unix()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		gomock.InOrder(
			mockRepository.EXPECT().
				FindOneByID(saveParam.ID, columns).
				Return(userModel, nil),

			mockRepository.EXPECT().
				Save(userModel).
				Return(userModel, nil),
		)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Save(saveParam)

		assert.NoError(t, err)
	})

	t.Run("save_find_one_not_found", func(t *testing.T) {
		saveParam := &SaveParam{ID: 0, Name: "test", Age: 18, Phone: "13000000000"}
		columns := []string{"*"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindOneByID(saveParam.ID, columns).
			Return(nil, gorm.ErrRecordNotFound)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Save(saveParam)

		assert.ErrorIs(t, err, ErrUserNotExist)
	})

	t.Run("save_find_one_error", func(t *testing.T) {
		saveParam := &SaveParam{ID: 1, Name: "test", Age: 18, Phone: "13000000000"}
		columns := []string{"*"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindOneByID(saveParam.ID, columns).
			Return(nil, errors.New("test error"))

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Save(saveParam)

		assert.EqualError(t, err, ErrDataQueryFailed.Error())
	})

	t.Run("param_copy_error", func(t *testing.T) {
		saveParam := &SaveParam{ID: 1, Name: "test", Age: 18, Phone: "13000000000"}
		columns := []string{"*"}

		userModel := new(model.User)
		if err := copier.Copy(userModel, saveParam); err != nil {
			t.Fatal(err)
		}
		userModel.CreatedAt = time.Now().Unix()
		userModel.UpdatedAt = time.Now().Unix()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindOneByID(saveParam.ID, columns).
			Return(userModel, nil)

		monkey.Patch(copier.Copy, func(toValue interface{}, fromValue interface{}) error {
			return errors.New("test error")
		})
		defer monkey.Unpatch(copier.Copy)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Save(saveParam)

		assert.EqualError(t, err, responsex.ServerErrorCode.String())
	})

	t.Run("save_failed", func(t *testing.T) {
		saveParam := &SaveParam{ID: 1, Name: "test", Age: 18, Phone: "13000000000XXX"}
		columns := []string{"*"}

		userModel := new(model.User)
		if err := copier.Copy(userModel, saveParam); err != nil {
			t.Fatal(err)
		}
		userModel.CreatedAt = time.Now().Unix()
		userModel.UpdatedAt = time.Now().Unix()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		gomock.InOrder(
			mockRepository.EXPECT().
				FindOneByID(saveParam.ID, columns).
				Return(userModel, nil),

			mockRepository.EXPECT().
				Save(userModel).
				Return(nil, errors.New("test error")),
		)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Save(saveParam)

		assert.ErrorIs(t, err, ErrDataStoreFailed)
	})
}
