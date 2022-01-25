package user

import (
	"bou.ke/monkey"
	"context"
	"errors"
	"fmt"
	"github.com/alicebob/miniredis/v2"

	// "github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/test"
	"gorm.io/gorm"
	"testing"
	"time"
)

func Test_repository_FindByKeyword(t *testing.T) {

	t.Run("keyword_is_empty", func(t *testing.T) {
		mockDB, err := test.NewMockDB()
		if err != nil {
			t.Fatal(err)
		}
		defer mockDB.MDB.Close()

		repo := New()
		repo.db = mockDB.GDB
		now := time.Now()

		exceptedUsers := []*model.User{
			{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test1", Age: 18, Phone: "13000000000"},
			{BaseModel: model.BaseModel{ID: 2, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test2", Age: 28, Phone: "13800000000"},
		}

		rows := mockDB.Mock.NewRows([]string{"id", "name", "age", "phone", "created_at", "updated_at", "deleted_at"})
		for _, exceptedUser := range exceptedUsers {
			rows.AddRow(
				exceptedUser.ID,
				exceptedUser.Name,
				exceptedUser.Age,
				exceptedUser.Phone,
				exceptedUser.CreatedAt,
				exceptedUser.UpdatedAt,
				exceptedUser.DeletedAt,
			)
		}

		mockDB.Mock.ExpectQuery("SELECT \\* FROM (.+) WHERE (.+)\\.`deleted_at` = \\? ORDER BY updated_at DESC").
			WillReturnRows(rows)

		users, err := repo.FindByKeyword([]string{"*"}, "", "updated_at DESC")

		assert.Equal(t, exceptedUsers, users)
		assert.NoError(t, err)

		if err = mockDB.Mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("keyword_is_not_empty", func(t *testing.T) {
		mockDB, err := test.NewMockDB()
		if err != nil {
			t.Fatal(err)
		}
		defer mockDB.MDB.Close()

		repo := New()
		repo.db = mockDB.GDB
		now := time.Now()

		exceptedUsers := []*model.User{
			{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test1", Age: 18, Phone: "13000000000"},
			{BaseModel: model.BaseModel{ID: 2, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test2", Age: 28, Phone: "13800000000"},
		}

		rows := mockDB.Mock.NewRows([]string{"id", "name", "age", "phone", "created_at", "updated_at", "deleted_at"})
		for _, exceptedUser := range exceptedUsers {
			rows.AddRow(
				exceptedUser.ID,
				exceptedUser.Name,
				exceptedUser.Age,
				exceptedUser.Phone,
				exceptedUser.CreatedAt,
				exceptedUser.UpdatedAt,
				exceptedUser.DeletedAt,
			)
		}

		mockDB.Mock.ExpectQuery("SELECT \\* FROM (.+) WHERE \\(name LIKE (.+) OR phone LIKE (.+)\\) AND (.+)\\.`deleted_at` = \\?  ORDER BY updated_at DESC").
			WillReturnRows(rows)

		users, err := repo.FindByKeyword([]string{"*"}, "test", "updated_at DESC")

		assert.Equal(t, exceptedUsers, users)
		assert.NoError(t, err)

		if err = mockDB.Mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("query_error", func(t *testing.T) {
		mockDB, err := test.NewMockDB()
		if err != nil {
			t.Fatal(err)
		}
		defer mockDB.MDB.Close()

		repo := New()
		repo.db = mockDB.GDB

		mockDB.Mock.ExpectQuery("SELECT \\* FROM (.+) WHERE (.+)\\.`deleted_at` = \\? ORDER BY updated_at DESC").
			WillReturnError(gorm.ErrInvalidValue)

		users, err := repo.FindByKeyword([]string{"*"}, "", "updated_at DESC")

		assert.Nil(t, users)
		assert.ErrorIs(t, err, gorm.ErrInvalidValue)

		if err = mockDB.Mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}

func Test_repository_FindOneByID(t *testing.T) {

	t.Run("cache_is_valid", func(t *testing.T) {
		now := time.Now()
		exceptedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}
		exceptedUserJson, err := jsoniter.Marshal(exceptedUser)
		if err != nil {
			t.Fatal(err)
		}

		rdb := miniredis.RunT(t)

		if err := rdb.Set(fmt.Sprintf(cacheKeyFormat, exceptedUser.ID), string(exceptedUserJson)); err != nil {
			t.Fatal(err)
		}

		repo := New()
		repo.rdb = test.NewMockRedisClient(rdb.Addr())

		user, err := repo.FindOneByID(exceptedUser.ID, []string{"*"})

		assert.Equal(t, exceptedUser, user)
		assert.NoError(t, err)
	})

	t.Run("get_cache_error", func(t *testing.T) {
		rdb := miniredis.RunT(t)

		if err := rdb.Set(fmt.Sprintf(cacheKeyFormat, 1), ""); err != nil {
			t.Fatal(err)
		}

		client := test.NewMockRedisClient(rdb.Addr())
		client.Close()

		repo := New()
		repo.rdb = client

		user, err := repo.FindOneByID(1, []string{"*"})

		assert.Nil(t, user)
		assert.ErrorIs(t, err, redis.ErrClosed)
	})

	t.Run("cache_unmarshal_error", func(t *testing.T) {
		rdb := miniredis.RunT(t)

		if err := rdb.Set(fmt.Sprintf(cacheKeyFormat, 1), ""); err != nil {
			t.Fatal(err)
		}

		repo := New()
		repo.rdb = test.NewMockRedisClient(rdb.Addr())

		user, err := repo.FindOneByID(1, []string{"*"})

		assert.Nil(t, user)
		assert.EqualError(t, err, "readObjectStart: expect { or n, but found \x00, error found in #0 byte of ...||..., bigger context ...||...")
	})

	t.Run("cache_is_invalid", func(t *testing.T) {

		t.Run("query_ok", func(t *testing.T) {
			mockDB, err := test.NewMockDB()
			if err != nil {
				t.Fatal(err)
			}
			defer mockDB.MDB.Close()

			rdb := miniredis.RunT(t)

			redisClient := test.NewMockRedisClient(rdb.Addr())

			repo := New()
			repo.db = mockDB.GDB
			repo.rdb = redisClient

			now := time.Now()

			exceptedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}
			exceptedUserJson, err := jsoniter.Marshal(exceptedUser)
			if err != nil {
				t.Fatal(err)
			}

			rows := mockDB.Mock.NewRows([]string{"id", "name", "age", "phone", "created_at", "updated_at", "deleted_at"}).
				AddRow(exceptedUser.ID, exceptedUser.Name, exceptedUser.Age, exceptedUser.Phone, exceptedUser.CreatedAt, exceptedUser.UpdatedAt, exceptedUser.DeletedAt)

			mockDB.Mock.ExpectQuery("SELECT \\* FROM (.+) WHERE id = (.+) AND (.+)\\.`deleted_at` = \\? LIMIT 1").
				WillReturnRows(rows)

			user, err := repo.FindOneByID(exceptedUser.ID, []string{"*"})

			assert.Equal(t, exceptedUser, user)
			assert.NoError(t, err)
			assert.JSONEq(t, string(exceptedUserJson), redisClient.Get(context.Background(), fmt.Sprintf(cacheKeyFormat, exceptedUser.ID)).Val())

			if err = mockDB.Mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("query_not_found", func(t *testing.T) {
			mockDB, err := test.NewMockDB()
			if err != nil {
				t.Fatal(err)
			}
			defer mockDB.MDB.Close()

			rdb := miniredis.RunT(t)

			repo := New()
			repo.db = mockDB.GDB
			repo.rdb = test.NewMockRedisClient(rdb.Addr())

			mockDB.Mock.ExpectQuery("SELECT \\* FROM (.+) WHERE id = (.+) AND (.+)\\.`deleted_at` = \\? LIMIT 1").
				WillReturnRows(mockDB.Mock.NewRows([]string{}))

			user, err := repo.FindOneByID(0, []string{"*"})

			assert.Nil(t, user)
			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

			if err = mockDB.Mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("query_record_marshal_error", func(t *testing.T) {
			mockDB, err := test.NewMockDB()
			if err != nil {
				t.Fatal(err)
			}
			defer mockDB.MDB.Close()

			rdb := miniredis.RunT(t)

			repo := New()
			repo.db = mockDB.GDB
			repo.rdb = test.NewMockRedisClient(rdb.Addr())

			now := time.Now()

			exceptedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}

			rows := mockDB.Mock.NewRows([]string{"id", "name", "age", "phone", "created_at", "updated_at", "deleted_at"}).
				AddRow(exceptedUser.ID, exceptedUser.Name, exceptedUser.Age, exceptedUser.Phone, exceptedUser.CreatedAt, exceptedUser.UpdatedAt, exceptedUser.DeletedAt)

			mockDB.Mock.ExpectQuery("SELECT \\* FROM (.+) WHERE id = (.+) AND (.+)\\.`deleted_at` = \\? LIMIT 1").
				WillReturnRows(rows)

			monkey.Patch(jsoniter.Marshal, func(v interface{}) ([]byte, error) {
				return nil, errors.New("test error")
			})
			defer monkey.Unpatch(jsoniter.Marshal)

			user, err := repo.FindOneByID(0, []string{"*"})

			assert.Nil(t, user)
			assert.EqualError(t, err, "test error")

			if err = mockDB.Mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("cache_set_error", func(t *testing.T) {
		mockDB, err := test.NewMockDB()
		if err != nil {
			t.Fatal(err)
		}
		defer mockDB.MDB.Close()

		rdb, rmock := redismock.NewClientMock()

		repo := New()
		repo.db = mockDB.GDB
		repo.rdb = rdb

		now := time.Now()

		exceptedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}
		exceptedUserJson, err := jsoniter.Marshal(exceptedUser)
		if err != nil {
			t.Fatal(err)
		}

		rows := mockDB.Mock.NewRows([]string{"id", "name", "age", "phone", "created_at", "updated_at", "deleted_at"}).
			AddRow(exceptedUser.ID, exceptedUser.Name, exceptedUser.Age, exceptedUser.Phone, exceptedUser.CreatedAt, exceptedUser.UpdatedAt, exceptedUser.DeletedAt)

		mockDB.Mock.ExpectQuery("SELECT \\* FROM (.+) WHERE id = (.+) AND (.+)\\.`deleted_at` = \\? LIMIT 1").
			WillReturnRows(rows)

		rmock.ExpectGet(fmt.Sprintf(cacheKeyFormat, exceptedUser.ID)).SetVal("")
		rmock.ExpectSet(
			fmt.Sprintf(cacheKeyFormat, exceptedUser.ID),
			string(exceptedUserJson),
			time.Duration(cacheExpire)*time.Second,
		).SetErr(errors.New("test error"))

		user, err := repo.FindOneByID(exceptedUser.ID, []string{"*"})

		assert.Nil(t, user)
		assert.EqualError(t, err, "test error")

		if err = mockDB.Mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func Test_repository_Create(t *testing.T) {
}

func Test_repository_Save(t *testing.T) {
}

func Test_repository_Delete(t *testing.T) {
}
