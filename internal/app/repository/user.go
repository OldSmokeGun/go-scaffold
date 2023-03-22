package repository

import (
	"context"

	"go-scaffold/internal/app/domain"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ domain.UserRepository = (*UserRepository)(nil)

// UserRepository 用户仓储
type UserRepository struct {
	db *gorm.DB
	// rdb *redis.Client
}

// NewUserRepository 构造用户仓储
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
		// rdb: rdb,
	}
}

// var (
// 	cacheKeyFormat = model.User{}.TableName() + "_%d"
// 	cacheExpire    = 3600
// )

// FindListParam 列表查询参数
type FindListParam struct {
	Keyword string
}

// FindList 列表查询
func (r *UserRepository) FindList(ctx context.Context, param domain.FindUserListParam) ([]*domain.User, error) {
	var models []*userModel

	query := r.db.WithContext(ctx).Select("*")

	if param.Keyword != "" {
		query.Where(
			r.db.Where("name LIKE ?", "%"+param.Keyword+"%").
				Or("phone LIKE ?", "%"+param.Keyword+"%"),
		)
	}

	err := query.Order("updated_at DESC").Find(&models).Error
	if err != nil {
		return nil, convertError(err)
	}

	entities := make([]*domain.User, 0, len(models))
	for _, m := range models {
		entities = append(entities, m.toEntity())
	}

	return entities, nil
}

// FindOneByID 根据 id 查询详情
func (r *UserRepository) FindOneByID(ctx context.Context, id domain.ID) (*domain.User, error) {
	m := new(userModel)

	err := r.db.WithContext(ctx).Select("*").
		Where("id = ?", id).
		Take(m).Error
	if err != nil {
		return nil, convertError(err)
	}

	return m.toEntity(), nil
}

// Exist 数据是否存在
func (r *UserRepository) Exist(ctx context.Context, id domain.ID) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(userModel{}).Where("id = ?", id).Count(&count).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, convertError(err)
	}

	return count == 1, nil
}

// Create 创建数据
func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	m := &userModel{
		Name:  user.Name,
		Age:   user.Age,
		Phone: user.Phone,
	}

	return errors.WithStack(r.db.WithContext(ctx).Create(m).Error)
}

// Update 保存数据
func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
	m := &userModel{
		BaseModel: BaseModel{
			ID: user.ID.Int64(),
		},
		Name:  user.Name,
		Age:   user.Age,
		Phone: user.Phone,
	}

	return errors.WithStack(r.db.WithContext(ctx).Save(m).Error)
}

// Delete 删除数据
func (r *UserRepository) Delete(ctx context.Context, user domain.User) error {
	return errors.WithStack(r.db.WithContext(ctx).Delete(&userModel{}, user.ID).Error)
}

// userModel 用户模型
type userModel struct {
	BaseModel
	Name  string `gorm:"column:name;type:varchar(64);not null;default:'';comment:名称"`
	Age   int8   `gorm:"column:age;type:tinyint(3);not null;default:0;comment:年龄"`
	Phone string `gorm:"column:phone;type:varchar(11);not null;default:'';comment:手机号码"`
}

// TableName 表名
func (u userModel) TableName() string {
	return "users"
}

// Migrate 迁移
func (u userModel) Migrate(db *gorm.DB) error {
	if err := db.Set(
		"gorm:table_options",
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表'",
	).AutoMigrate(u); err != nil {
		return err
	}

	return nil
}

func (u userModel) toEntity() *domain.User {
	return &domain.User{
		ID:    domain.ID(u.ID),
		Name:  u.Name,
		Age:   u.Age,
		Phone: u.Phone,
	}
}
