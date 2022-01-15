package migratorx

import "gorm.io/gorm"

// Migrator 数据迁移
type Migrator struct {
	// db gorm 数据库实例
	db *gorm.DB
}

// New 构造函数
func New(db *gorm.DB) *Migrator {
	return &Migrator{
		db: db,
	}
}

// Exec 执行迁移任务
func (m Migrator) Exec(tasks []*Task) (err error) {
	for _, task := range tasks {
		if task.Execute == nil {
			if err = m.db.Set(
				"gorm:table_options",
				"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='"+task.Comment+"'",
			).AutoMigrate(task.Model); err != nil {
				return err
			}

		} else {

			if err = task.Execute(m.db); err != nil {
				return err
			}
		}
	}

	return nil
}

// Task 迁移任务
type Task struct {
	// Comment 数据表的注释
	Comment string

	// Model 要迁移的模型
	Model interface{}

	// Execute 是迁移任务的执行函数
	// 如果设置此函数，Comment 和 Model 属性的设置将无效
	Execute func(db *gorm.DB) error
}
