package test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"io"
)

var (
	// logger 测试用 logger
	// 注意：日志内容将会被丢弃，可用于替换实际的 logger
	logger *zap.Logger
)

func Init() {
	var err error

	// logger
	logger, err = zap.NewDevelopment(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(io.Discard),
			// zapcore.AddSync(os.Stdout),
			zapcore.ErrorLevel,
		)
	}))
	if err != nil {
		panic(err)
	}
}

// Logger 测试用 logger
func Logger() *zap.Logger {
	return logger
}

// MockDB 测试用 db
type MockDB struct {
	MDB  *sql.DB
	Mock sqlmock.Sqlmock
	GDB  *gorm.DB
}

// NewMockDB 返回用于测试的 db
func NewMockDB() (*MockDB, error) {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	gdb, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn:                      mdb,
			SkipInitializeWithVersion: true,
		}),
		&gorm.Config{
			// Logger: glogger.Discard,
			Logger: glogger.Default,
		},
	)
	if err != nil {
		return nil, err
	}

	return &MockDB{MDB: mdb, Mock: mock, GDB: gdb}, nil
}

// NewMockRedisClient 返回用于测试的 redis 客户端
func NewMockRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
