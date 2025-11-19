package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// UserRepository 用户数据存取仓库。
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建 UserRepository 实例。
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 新增用户记录。
func (r *UserRepository) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.Wrap(err, "create user")
	}
	return nil
}

// Update 更新已有用户数据。
func (r *UserRepository) Update(user *model.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return errors.Wrap(err, "update user")
	}
	return nil
}

// FindByWorkNo 通过工号查询用户。
func (r *UserRepository) FindByWorkNo(workNo string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("work_no = ?", workNo).First(&user).Error; err != nil {
		return nil, errors.Wrap(err, "find user by work no")
	}
	return &user, nil
}

// ListManagers 获取所有启用中的店长。
func (r *UserRepository) ListManagers() ([]model.User, error) {
	var users []model.User
	if err := r.db.Where("role = ?", model.RoleManager).Where("status = ?", true).Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "list managers")
	}
	return users, nil
}

// FindByID 通过主键 ID 查询用户。
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, errors.Wrap(err, "find user by id")
	}
	return &user, nil
}
