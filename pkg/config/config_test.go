package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
	"time"
)

type testConfigModel struct {
	Test  string
	Test1 struct {
		Key1 int
		Key2 string
	}
	Test2 []struct {
		Key1 int
		Key2 string
	}
}

func TestLoad(t *testing.T) {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	t.Run("config_file_normal", func(t *testing.T) {
		configFilePath := fmt.Sprintf("%s/test_load_%d%d.yaml", t.TempDir(), time.Now().UnixNano(), newRand.Intn(99999))
		configContent := `Test: "value"
Test1:
  Key1: 123
  Key2: "value2"
Test2:
  - Key1: 123
    Key2: "value2"`
		testConfigCompareModel := &testConfigModel{
			Test: "value",
			Test1: struct {
				Key1 int
				Key2 string
			}{
				Key1: 123,
				Key2: "value2",
			},
			Test2: []struct {
				Key1 int
				Key2 string
			}{
				{
					Key1: 123,
					Key2: "value2",
				},
			},
		}

		f, err := os.Create(configFilePath)
		if err != nil {
			t.Fatal(err)
		}

		_, err = f.WriteString(configContent)
		if err != nil {
			t.Fatal(err)
		}
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err = os.Remove(configFilePath); err != nil {
				t.Fatal(err)
			}
		})

		m := new(testConfigModel)
		assert.NoError(t, Load(configFilePath, m))
		assert.Equal(t, testConfigCompareModel, m)
	})
	t.Run("config_file_not_exist", func(t *testing.T) {
		assert.ErrorAs(t, Load("./6620aa46-4ff3-4943-b227-af3a334b81a3.yaml", new(testConfigModel)), &viper.ConfigFileNotFoundError{})
	})
	t.Run("config_file_is_empty", func(t *testing.T) {
		m := new(testConfigModel)
		assert.ErrorIs(t, Load("", m), ErrFileNotSpecified)
		assert.Equal(t, new(testConfigModel), m)
	})
	t.Run("config_file_read_error", func(t *testing.T) {
		// 设置不支持的配置文件类型
		configFilePath := fmt.Sprintf("%s/test_load_%d%d.test", t.TempDir(), time.Now().UnixNano(), newRand.Intn(99999))

		f, err := os.Create(configFilePath)
		if err != nil {
			t.Fatal(err)
		}
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err = os.Remove(configFilePath); err != nil {
				t.Fatal(err)
			}
		})

		assert.ErrorAs(t, Load(configFilePath, new(testConfigModel)), &viper.ConfigFileNotFoundError{})
	})
	t.Run("config_file_unmarshal_error", func(t *testing.T) {
		// 设置不支持的配置文件类型
		configFilePath := fmt.Sprintf("%s/test_load_%d%d.yaml", t.TempDir(), time.Now().UnixNano(), newRand.Intn(99999))
		configContent := `Test: "value"`

		f, err := os.Create(configFilePath)
		if err != nil {
			t.Fatal(err)
		}

		_, err = f.WriteString(configContent)
		if err != nil {
			t.Fatal(err)
		}

		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err = os.Remove(configFilePath); err != nil {
				t.Fatal(err)
			}
		})

		exceptedError := "'' expected type 'string', got unconvertible type 'map[string]interface {}', value: 'map[test:value]'"
		assert.EqualError(t, Load(configFilePath, new(string)), exceptedError)
	})
	t.Run("config_file_live_reload", func(t *testing.T) {
		configFilePath := fmt.Sprintf("%s/test_load_%d%d.yaml", t.TempDir(), time.Now().UnixNano(), newRand.Intn(99999))
		configContent := `Test: "value"`
		testConfigCompareModel := &testConfigModel{
			Test: "value",
		}

		f, err := os.Create(configFilePath)
		if err != nil {
			t.Fatal(err)
		}

		_, err = f.WriteString(configContent)
		if err != nil {
			t.Fatal(err)
		}
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err = os.Remove(configFilePath); err != nil {
				t.Fatal(err)
			}
		})

		var ch = make(chan struct{}, 1)

		m := new(testConfigModel)
		assert.NoError(t, Load(configFilePath, m, func(v *viper.Viper, model interface{}, e fsnotify.Event) {
			if err := v.MergeInConfig(); err != nil {
				panic(err)
			}
			if err := v.Unmarshal(model); err != nil {
				panic(err)
			}

			ch <- struct{}{}
		}))
		assert.Equal(t, testConfigCompareModel, m)

		f, err = os.OpenFile(configFilePath, os.O_RDWR, 0666)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = f.WriteAt([]byte("VALUE"), 7); err != nil {
			t.Fatal(err)
		}
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}

		<-ch

		assert.Equal(t, "VALUE", m.Test)
	})
}

func TestMustLoad(t *testing.T) {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	t.Run("no_panic", func(t *testing.T) {
		configFilePath := fmt.Sprintf("%s/test_load_%d%d.yaml", t.TempDir(), time.Now().UnixNano(), newRand.Intn(99999))
		f, err := os.Create(configFilePath)
		if err != nil {
			t.Fatal(err)
		}
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err = os.Remove(configFilePath); err != nil {
				t.Fatal(err)
			}
		})

		m := new(testConfigModel)
		assert.NotPanics(t, func() {
			MustLoad(configFilePath, m)
		})
	})
	t.Run("panic", func(t *testing.T) {
		m := new(testConfigModel)
		assert.Panics(t, func() {
			MustLoad("./2106fae0-33ca-464a-b73a-52c5d88c8938.yaml", m)
		})
	})
}
