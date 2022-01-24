package user

import (
	"bou.ke/monkey"
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

func Test_service_Create(t *testing.T) {

	t.Run("create_success", func(t *testing.T) {
		createParam := &CreateParam{Name: "test", Age: 18, Phone: "13000000000"}

		userModel := new(model.User)
		if err := copier.Copy(userModel, createParam); err != nil {
			t.Fatal(err)
		}

		createdUserModel := *userModel
		createdUserModel.ID = 1
		createdUserModel.CreatedAt = time.Now().Unix()
		createdUserModel.UpdatedAt = time.Now().Unix()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			Create(userModel).
			Return(&createdUserModel, nil)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Create(createParam)

		assert.NoError(t, err)
	})

	t.Run("param_copy_error", func(t *testing.T) {
		createParam := &CreateParam{Name: "test", Age: 18, Phone: "13000000000"}

		monkey.Patch(copier.Copy, func(toValue interface{}, fromValue interface{}) error {
			return copier.ErrInvalidCopyDestination
		})
		defer monkey.Unpatch(copier.Copy)

		newService := New()
		newService.Logger = test.Logger()

		err := newService.Create(createParam)

		assert.EqualError(t, err, responsex.ServerErrorCode.String())
	})

	t.Run("create_failed", func(t *testing.T) {
		createParam := &CreateParam{Name: "test", Age: 18, Phone: "13000000000XXX"}

		userModel := new(model.User)
		if err := copier.Copy(userModel, createParam); err != nil {
			t.Fatal(err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			Create(userModel).
			Return(nil, gorm.ErrInvalidValueOfLength)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		err := newService.Create(createParam)

		assert.ErrorIs(t, err, ErrDataStoreFailed)
	})
}
