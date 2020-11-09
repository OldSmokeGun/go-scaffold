package migrate

import (
	"gin-scaffold/internal/app/models"
	"gin-scaffold/internal/global"
)

func Run() error {
	var list = []interface{}{
		&models.User{},
	}

	if err := global.DB.Set("gorm:table_options", "ENGINE=InnoDB charset=utf8mb4").AutoMigrate(list...); err != nil {
		return err
	}

	return nil
}
