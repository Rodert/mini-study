package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// ContentCategoryRepository 内容分类仓储。
type ContentCategoryRepository struct {
	db *gorm.DB
}

// NewContentCategoryRepository 创建分类仓库实例。
func NewContentCategoryRepository(db *gorm.DB) *ContentCategoryRepository {
	return &ContentCategoryRepository{db: db}
}

// ListByRole 根据角色查询可用分类，role 为空时返回全部。
func (r *ContentCategoryRepository) ListByRole(role string) ([]model.ContentCategory, error) {
	var categories []model.ContentCategory
	query := r.db.Where("status = ?", true).Order("sort_order asc, id asc")
	if role != "" && role != "both" {
		query = query.Where("role_scope = ? OR role_scope = ?", role, "both")
	}
	if err := query.Find(&categories).Error; err != nil {
		return nil, errors.Wrap(err, "list categories")
	}
	return categories, nil
}

// FindByID 通过 ID 获取分类。
func (r *ContentCategoryRepository) FindByID(id uint) (*model.ContentCategory, error) {
	var category model.ContentCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, errors.Wrap(err, "find category")
	}
	return &category, nil
}

// EnsureDefaults 确保默认分类存在，缺失则插入。
func (r *ContentCategoryRepository) EnsureDefaults(defaults []model.ContentCategory) error {
	for _, item := range defaults {
		var count int64
		if err := r.db.Model(&model.ContentCategory{}).
			Where("name = ?", item.Name).
			Count(&count).Error; err != nil {
			return errors.Wrap(err, "check category")
		}
		if count > 0 {
			continue
		}
		if err := r.db.Create(&item).Error; err != nil {
			return errors.Wrap(err, "seed category")
		}
	}
	return nil
}
