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

func Test_service_List(t *testing.T) {

	t.Run("get_list_success", func(t *testing.T) {
		listParam := &ListParam{Keyword: ""}
		columns := []string{"*"}

		users := []*model.User{
			{
				BaseModel: model.BaseModel{ID: 1, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(), DeletedAt: 0},
				Name:      "test1",
				Age:       18,
				Phone:     "13000000000",
			},
			{
				BaseModel: model.BaseModel{ID: 2, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(), DeletedAt: 0},
				Name:      "test2",
				Age:       28,
				Phone:     "13800000000",
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindByKeyword(columns, listParam.Keyword, "updated_at DESC").
			Return(users, nil)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		var listResult ListResult
		if err := copier.Copy(&listResult, users); err != nil {
			t.Fatal(err)
		}

		ret, err := newService.List(listParam)

		assert.NoError(t, err)
		assert.Equal(t, listResult, ret)
	})

	t.Run("get_list_failed", func(t *testing.T) {
		listParam := &ListParam{Keyword: ""}
		columns := []string{"*"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindByKeyword(columns, listParam.Keyword, "updated_at DESC").
			Return(nil, gorm.ErrInvalidField)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		ret, err := newService.List(listParam)

		assert.ErrorIs(t, err, ErrDataQueryFailed)
		assert.Nil(t, ret)
	})

	t.Run("result_copy_error", func(t *testing.T) {
		listParam := &ListParam{Keyword: ""}
		columns := []string{"*"}

		users := []*model.User{
			{
				BaseModel: model.BaseModel{ID: 1, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(), DeletedAt: 0},
				Name:      "test1",
				Age:       18,
				Phone:     "13000000000",
			},
			{
				BaseModel: model.BaseModel{ID: 2, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(), DeletedAt: 0},
				Name:      "test2",
				Age:       28,
				Phone:     "13800000000",
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepository := user.NewMockRepository(ctrl)
		mockRepository.EXPECT().
			FindByKeyword(columns, listParam.Keyword, "updated_at DESC").
			Return(users, nil)

		newService := New()
		newService.Logger = test.Logger()
		newService.Repository = mockRepository

		monkey.Patch(copier.Copy, func(toValue interface{}, fromValue interface{}) error {
			return copier.ErrInvalidCopyDestination
		})
		defer monkey.Unpatch(copier.Copy)

		ret, err := newService.List(listParam)

		assert.EqualError(t, err, responsex.ServerErrorCode.String())
		assert.Nil(t, ret)
	})
}
