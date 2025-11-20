package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// ContentRepository 内容数据仓储。
type ContentRepository struct {
	db *gorm.DB
}

// NewContentRepository 创建内容仓库实例。
func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

// Create 新增内容。
func (r *ContentRepository) Create(content *model.Content) error {
	if err := r.db.Create(content).Error; err != nil {
		return errors.Wrap(err, "create content")
	}
	return nil
}

// Update 更新内容。
func (r *ContentRepository) Update(content *model.Content) error {
	if err := r.db.Save(content).Error; err != nil {
		return errors.Wrap(err, "update content")
	}
	return nil
}

// FindByID 通过 ID 查询内容。
func (r *ContentRepository) FindByID(id uint) (*model.Content, error) {
	var content model.Content
	if err := r.db.Preload("Category").First(&content, id).Error; err != nil {
		return nil, errors.Wrap(err, "find content")
	}
	return &content, nil
}

// ListAdmin 管理端内容列表，支持分类/类型/状态筛选。
func (r *ContentRepository) ListAdmin(categoryID uint, contentType, status string) ([]model.Content, error) {
	query := r.db.Preload("Category").Order("id desc")
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if contentType != "" {
		query = query.Where("type = ?", contentType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var contents []model.Content
	if err := query.Find(&contents).Error; err != nil {
		return nil, errors.Wrap(err, "list admin contents")
	}
	return contents, nil
}

// ListPublishedByRole 根据角色、分类、类型获取已发布内容。
func (r *ContentRepository) ListPublishedByRole(role string, categoryID uint, contentType string) ([]model.Content, error) {
	query := r.db.Preload("Category").
		Where("status = ?", "published").
		Order("publish_at desc, id desc")

	if role != "" && role != "both" {
		query = query.Where("visible_roles = ? OR visible_roles = ?", role, "both")
	}

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if contentType != "" {
		query = query.Where("type = ?", contentType)
	}

	var contents []model.Content
	if err := query.Find(&contents).Error; err != nil {
		return nil, errors.Wrap(err, "list published contents")
	}
	return contents, nil
}

// CountPublishedForRole counts how many contents are visible to a given role.
func (r *ContentRepository) CountPublishedForRole(role string) (int64, error) {
	query := r.db.Model(&model.Content{}).Where("status = ?", "published")
	if role != "" && role != "both" {
		query = query.Where("visible_roles = ? OR visible_roles = ?", role, "both")
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, "count contents for role")
	}
	return count, nil
}

// CountPublishedByCategoryAndRole counts how many published contents are in a category and visible to a given role.
func (r *ContentRepository) CountPublishedByCategoryAndRole(categoryID uint, role string) (int64, error) {
	query := r.db.Model(&model.Content{}).
		Where("status = ?", "published").
		Where("category_id = ?", categoryID)
	
	if role != "" && role != "both" {
		query = query.Where("visible_roles = ? OR visible_roles = ?", role, "both")
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, "count contents by category and role")
	}
	return count, nil
}
