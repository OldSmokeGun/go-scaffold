package user

import (
	"bou.ke/monkey"
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
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
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		repo := New()
		repo.db = gdb

		now := time.Now()

		expectedUsers := []*model.User{
			{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test1", Age: 18, Phone: "13000000000"},
			{BaseModel: model.BaseModel{ID: 2, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test2", Age: 28, Phone: "13800000000"},
		}

		rows := dmock.NewRows([]string{"id", "name", "age", "phone", "created_at", "updated_at", "deleted_at"})
		for _, expectedUser := range expectedUsers {
			rows.AddRow(
				expectedUser.ID,
				expectedUser.Name,
				expectedUser.Age,
				expectedUser.Phone,
				expectedUser.CreatedAt,
				expectedUser.UpdatedAt,
				expectedUser.DeletedAt,
			)
		}

		dmock.ExpectQuery("SELECT \\* FROM (.+) WHERE (.+)\\.`deleted_at` = \\? ORDER BY updated_at DESC").
			WillReturnRows(rows)

		users, err := repo.FindByKeyword([]string{"*"}, "", "updated_at DESC")

		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, err)

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("keyword_is_not_empty", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		repo := New()
		repo.db = gdb

		now := time.Now()

		expectedUsers := []*model.User{
			{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test1", Age: 18, Phone: "13000000000"},
			{BaseModel: model.BaseModel{ID: 2, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test2", Age: 28, Phone: "13800000000"},
		}

		rows := dmock.NewRows([]string{"id", "name", "age", "phone", "created_at", "updated_at", "deleted_at"})
		for _, expectedUser := range expectedUsers {
			rows.AddRow(
				expectedUser.ID,
				expectedUser.Name,
				expectedUser.Age,
				expectedUser.Phone,
				expectedUser.CreatedAt,
				expectedUser.UpdatedAt,
				expectedUser.DeletedAt,
			)
		}

		dmock.ExpectQuery("SELECT \\* FROM (.+) WHERE \\(name LIKE (.+) OR phone LIKE (.+)\\) AND (.+)\\.`deleted_at` = \\?  ORDER BY updated_at DESC").
			WillReturnRows(rows)

		users, err := repo.FindByKeyword([]string{"*"}, "test", "updated_at DESC")

		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, err)

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("query_failed", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		repo := New()
		repo.db = gdb

		dmock.ExpectQuery("SELECT \\* FROM (.+) WHERE (.+)\\.`deleted_at` = \\? ORDER BY updated_at DESC").
			WillReturnError(errors.New("test error"))

		users, err := repo.FindByKeyword([]string{"*"}, "", "updated_at DESC")

		assert.Nil(t, users)
		assert.EqualError(t, err, "test error")

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}

func Test_repository_FindOneByID(t *testing.T) {

	t.Run("cache_is_valid", func(t *testing.T) {
		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}
		expectedUserJson, err := jsoniter.Marshal(expectedUser)
		if err != nil {
			t.Fatal(err)
		}

		rdb, rmock := redismock.NewClientMock()
		defer rdb.Close()

		repo := New()
		repo.rdb = rdb

		rmock.ExpectGet(fmt.Sprintf(cacheKeyFormat, expectedUser.ID)).SetVal(string(expectedUserJson))

		user, err := repo.FindOneByID(expectedUser.ID, []string{"*"})

		assert.Equal(t, expectedUser, user)
		assert.NoError(t, err)

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("get_cache_failed", func(t *testing.T) {
		rdb, rmock := redismock.NewClientMock()
		defer rdb.Close()

		repo := New()
		repo.rdb = rdb

		rmock.ExpectGet(fmt.Sprintf(cacheKeyFormat, 1)).SetErr(errors.New("test error"))

		user, err := repo.FindOneByID(1, []string{"*"})

		assert.Nil(t, user)
		assert.EqualError(t, err, "test error")

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("cache_unmarshal_error", func(t *testing.T) {
		rdb, rmock := redismock.NewClientMock()
		defer rdb.Close()

		repo := New()
		repo.rdb = rdb

		rmock.ExpectGet(fmt.Sprintf(cacheKeyFormat, 1)).SetVal("test")

		user, err := repo.FindOneByID(1, []string{"*"})

		assert.Nil(t, user)
		assert.EqualError(t, err, "readObjectStart: expect { or n, but found t, error found in #1 byte of ...|test|..., bigger context ...|test|...")

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("cache_is_invalid", func(t *testing.T) {

		t.Run("query_ok", func(t *testing.T) {
			mdb, dmock, gdb, err := test.NewDBMock()
			if err != nil {
				t.Fatal(err)
			}
			defer mdb.Close()

			rdb, rmock := redismock.NewClientMock()
			defer rdb.Close()

			repo := New()
			repo.db = gdb
			repo.rdb = rdb

			now := time.Now()

			expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}
			expectedUserJson, err := jsoniter.Marshal(expectedUser)
			if err != nil {
				t.Fatal(err)
			}

			rows := dmock.NewRows([]string{"id", "name", "age", "phone", "created_at", "updated_at", "deleted_at"}).
				AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Age, expectedUser.Phone, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.DeletedAt)

			dmock.ExpectQuery("SELECT \\* FROM (.+) WHERE id = (.+) AND (.+)\\.`deleted_at` = \\? LIMIT 1").
				WillReturnRows(rows)

			rmock.ExpectGet(fmt.Sprintf(cacheKeyFormat, expectedUser.ID)).SetVal("")
			rmock.ExpectSet(
				fmt.Sprintf(cacheKeyFormat, expectedUser.ID),
				string(expectedUserJson),
				time.Duration(cacheExpire)*time.Second,
			).SetVal(string(expectedUserJson))
			rmock.ExpectGet(fmt.Sprintf(cacheKeyFormat, expectedUser.ID)).SetVal(string(expectedUserJson))

			user, err := repo.FindOneByID(expectedUser.ID, []string{"*"})

			assert.Equal(t, expectedUser, user)
			assert.NoError(t, err)
			assert.JSONEq(t, string(expectedUserJson), rdb.Get(context.Background(), fmt.Sprintf(cacheKeyFormat, expectedUser.ID)).Val())

			if err = dmock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}

			if err = rmock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})

		t.Run("query_not_found", func(t *testing.T) {
			mdb, dmock, gdb, err := test.NewDBMock()
			if err != nil {
				t.Fatal(err)
			}
			defer mdb.Close()

			rdb, rmock := redismock.NewClientMock()
			defer rdb.Close()

			repo := New()
			repo.db = gdb
			repo.rdb = rdb

			dmock.ExpectQuery("SELECT \\* FROM (.+) WHERE id = (.+) AND (.+)\\.`deleted_at` = \\? LIMIT 1").
				WillReturnRows(dmock.NewRows([]string{}))

			rmock.ExpectGet(fmt.Sprintf(cacheKeyFormat, 0)).SetVal("")

			user, err := repo.FindOneByID(0, []string{"*"})

			assert.Nil(t, user)
			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

			if err = dmock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}

			if err = rmock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})

		t.Run("query_record_marshal_error", func(t *testing.T) {
			mdb, dmock, gdb, err := test.NewDBMock()
			if err != nil {
				t.Fatal(err)
			}
			defer mdb.Close()

			rdb, rmock := redismock.NewClientMock()
			defer rdb.Close()

			repo := New()
			repo.db = gdb
			repo.rdb = rdb

			now := time.Now()

			expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}

			rows := dmock.NewRows([]string{"id", "name", "age", "phone", "created_at", "updated_at", "deleted_at"}).
				AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Age, expectedUser.Phone, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.DeletedAt)

			dmock.ExpectQuery("SELECT \\* FROM (.+) WHERE id = (.+) AND (.+)\\.`deleted_at` = \\? LIMIT 1").
				WillReturnRows(rows)

			rmock.ExpectGet(fmt.Sprintf(cacheKeyFormat, 0)).SetVal("")

			monkey.Patch(jsoniter.Marshal, func(v interface{}) ([]byte, error) {
				return nil, errors.New("test error")
			})
			defer monkey.Unpatch(jsoniter.Marshal)

			user, err := repo.FindOneByID(0, []string{"*"})

			assert.Nil(t, user)
			assert.EqualError(t, err, "test error")

			if err = dmock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}

			if err = rmock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("cache_set_failed", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		rdb, rmock := redismock.NewClientMock()
		defer rdb.Close()

		repo := New()
		repo.db = gdb
		repo.rdb = rdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}
		expectedUserJson, err := jsoniter.Marshal(expectedUser)
		if err != nil {
			t.Fatal(err)
		}

		rows := dmock.NewRows([]string{"id", "name", "age", "phone", "created_at", "updated_at", "deleted_at"}).
			AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Age, expectedUser.Phone, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.DeletedAt)

		dmock.ExpectQuery("SELECT \\* FROM (.+) WHERE id = (.+) AND (.+)\\.`deleted_at` = \\? LIMIT 1").
			WillReturnRows(rows)

		rmock.ExpectGet(fmt.Sprintf(cacheKeyFormat, expectedUser.ID)).SetVal("")
		rmock.ExpectSet(
			fmt.Sprintf(cacheKeyFormat, expectedUser.ID),
			string(expectedUserJson),
			time.Duration(cacheExpire)*time.Second,
		).SetErr(errors.New("test error"))

		user, err := repo.FindOneByID(expectedUser.ID, []string{"*"})

		assert.Nil(t, user)
		assert.EqualError(t, err, "test error")

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func Test_repository_Create(t *testing.T) {

	t.Run("create_success", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		rdb, rmock := redismock.NewClientMock()
		defer rdb.Close()

		repo := New()
		repo.db = gdb
		repo.rdb = rdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}
		expectedUserJson, err := jsoniter.Marshal(expectedUser)
		if err != nil {
			t.Fatal(err)
		}

		dmock.ExpectExec("INSERT INTO (.+)").
			WillReturnResult(sqlmock.NewResult(0, 1))

		rmock.ExpectSet(
			fmt.Sprintf(cacheKeyFormat, expectedUser.ID),
			string(expectedUserJson),
			time.Duration(cacheExpire)*time.Second,
		).SetVal(string(expectedUserJson))

		user, err := repo.Create(expectedUser)

		assert.Equal(t, expectedUser, user)
		assert.NoError(t, err)

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("create_failed", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		repo := New()
		repo.db = gdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}

		dmock.ExpectExec("INSERT INTO (.+)").
			WillReturnError(errors.New("test error"))

		user, err := repo.Create(expectedUser)

		assert.Nil(t, user)
		assert.EqualError(t, err, "test error")

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("model_marshal_error", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		repo := New()
		repo.db = gdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}

		dmock.ExpectExec("INSERT INTO (.+)").
			WillReturnResult(sqlmock.NewResult(0, 1))

		monkey.Patch(jsoniter.Marshal, func(v interface{}) ([]byte, error) {
			return nil, errors.New("test error")
		})
		defer monkey.Unpatch(jsoniter.Marshal)

		user, err := repo.Create(expectedUser)

		assert.Nil(t, user)
		assert.EqualError(t, err, "test error")

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("cache_set_failed", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		rdb, rmock := redismock.NewClientMock()
		defer rdb.Close()

		repo := New()
		repo.db = gdb
		repo.rdb = rdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}
		expectedUserJson, err := jsoniter.Marshal(expectedUser)
		if err != nil {
			t.Fatal(err)
		}

		dmock.ExpectExec("INSERT INTO (.+)").
			WillReturnResult(sqlmock.NewResult(0, 1))

		rmock.ExpectSet(
			fmt.Sprintf(cacheKeyFormat, expectedUser.ID),
			string(expectedUserJson),
			time.Duration(cacheExpire)*time.Second,
		).SetErr(errors.New("test error"))

		user, err := repo.Create(expectedUser)

		assert.Nil(t, user)
		assert.EqualError(t, err, "test error")

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func Test_repository_Save(t *testing.T) {

	t.Run("save_success", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		rdb, rmock := redismock.NewClientMock()
		defer rdb.Close()

		repo := New()
		repo.db = gdb
		repo.rdb = rdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}
		expectedUserJson, err := jsoniter.Marshal(expectedUser)
		if err != nil {
			t.Fatal(err)
		}

		dmock.ExpectExec("UPDATE (.+) SET (.+)").
			WillReturnResult(sqlmock.NewResult(0, 1))

		rmock.ExpectSet(
			fmt.Sprintf(cacheKeyFormat, expectedUser.ID),
			string(expectedUserJson),
			time.Duration(cacheExpire)*time.Second,
		).SetVal(string(expectedUserJson))

		user, err := repo.Save(expectedUser)

		assert.Equal(t, expectedUser, user)
		assert.NoError(t, err)

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("save_failed", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		repo := New()
		repo.db = gdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}

		dmock.ExpectExec("UPDATE (.+) SET (.+)").
			WillReturnError(errors.New("test error"))

		user, err := repo.Save(expectedUser)

		assert.Nil(t, user)
		assert.EqualError(t, err, "test error")

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("model_marshal_error", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		repo := New()
		repo.db = gdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}

		dmock.ExpectExec("UPDATE (.+) SET (.+)").
			WillReturnResult(sqlmock.NewResult(0, 1))

		monkey.Patch(jsoniter.Marshal, func(v interface{}) ([]byte, error) {
			return nil, errors.New("test error")
		})
		defer monkey.Unpatch(jsoniter.Marshal)

		user, err := repo.Save(expectedUser)

		assert.Nil(t, user)
		assert.EqualError(t, err, "test error")

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("cache_set_failed", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		rdb, rmock := redismock.NewClientMock()
		defer rdb.Close()

		repo := New()
		repo.db = gdb
		repo.rdb = rdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}
		expectedUserJson, err := jsoniter.Marshal(expectedUser)
		if err != nil {
			t.Fatal(err)
		}

		dmock.ExpectExec("UPDATE (.+) SET (.+)").
			WillReturnResult(sqlmock.NewResult(0, 1))

		rmock.ExpectSet(
			fmt.Sprintf(cacheKeyFormat, expectedUser.ID),
			string(expectedUserJson),
			time.Duration(cacheExpire)*time.Second,
		).SetErr(errors.New("test error"))

		user, err := repo.Save(expectedUser)

		assert.Nil(t, user)
		assert.EqualError(t, err, "test error")

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func Test_repository_Delete(t *testing.T) {

	t.Run("delete_success", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		rdb, rmock := redismock.NewClientMock()
		defer rdb.Close()

		repo := New()
		repo.db = gdb
		repo.rdb = rdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}

		dmock.ExpectExec("UPDATE (.+) SET (.+)").
			WillReturnResult(sqlmock.NewResult(0, 1))

		rmock.ExpectDel(fmt.Sprintf(cacheKeyFormat, expectedUser.ID)).SetVal(1)

		err = repo.Delete(expectedUser)

		assert.NoError(t, err)

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("delete_failed", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		repo := New()
		repo.db = gdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}

		dmock.ExpectExec("UPDATE (.+) SET (.+)").
			WillReturnError(errors.New("test error"))

		err = repo.Delete(expectedUser)

		assert.EqualError(t, err, "test error")

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("cache_del_failed", func(t *testing.T) {
		mdb, dmock, gdb, err := test.NewDBMock()
		if err != nil {
			t.Fatal(err)
		}
		defer mdb.Close()

		rdb, rmock := redismock.NewClientMock()
		defer rdb.Close()

		repo := New()
		repo.db = gdb
		repo.rdb = rdb

		now := time.Now()

		expectedUser := &model.User{BaseModel: model.BaseModel{ID: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix(), DeletedAt: 0}, Name: "test", Age: 18, Phone: "13000000000"}

		dmock.ExpectExec("UPDATE (.+) SET (.+)").
			WillReturnResult(sqlmock.NewResult(0, 1))

		rmock.ExpectDel(fmt.Sprintf(cacheKeyFormat, expectedUser.ID)).SetErr(errors.New("test error"))

		err = repo.Delete(expectedUser)

		assert.EqualError(t, err, "test error")

		if err = dmock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}

		if err = rmock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}
