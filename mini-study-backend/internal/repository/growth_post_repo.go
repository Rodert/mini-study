package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// GrowthPostRepository handles growth circle post persistence.
type GrowthPostRepository struct {
	db *gorm.DB
}

// NewGrowthPostRepository creates growth post repo.
func NewGrowthPostRepository(db *gorm.DB) *GrowthPostRepository {
	return &GrowthPostRepository{db: db}
}

// Create inserts a growth post.
func (r *GrowthPostRepository) Create(post *model.GrowthPost) error {
	if err := r.db.Create(post).Error; err != nil {
		return errors.Wrap(err, "create growth post")
	}
	return nil
}

// Update updates a growth post.
func (r *GrowthPostRepository) Update(post *model.GrowthPost) error {
	if err := r.db.Save(post).Error; err != nil {
		return errors.Wrap(err, "update growth post")
	}
	return nil
}

// Delete soft-deletes a growth post.
func (r *GrowthPostRepository) Delete(post *model.GrowthPost) error {
	if err := r.db.Delete(post).Error; err != nil {
		return errors.Wrap(err, "delete growth post")
	}
	return nil
}

// FindByID finds a growth post by id.
func (r *GrowthPostRepository) FindByID(id uint) (*model.GrowthPost, error) {
	var post model.GrowthPost
	if err := r.db.Preload("Creator").First(&post, id).Error; err != nil {
		return nil, errors.Wrap(err, "find growth post")
	}
	return &post, nil
}

// ListPublic returns approved posts visible to everyone, optionally filtered by keyword.
func (r *GrowthPostRepository) ListPublic(keyword string) ([]model.GrowthPost, error) {
	query := r.db.Preload("Creator").
		Where("status = ?", "approved").
		Order("created_at desc, id desc")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("content LIKE ?", like)
	}

	var posts []model.GrowthPost
	if err := query.Find(&posts).Error; err != nil {
		return nil, errors.Wrap(err, "list public growth posts")
	}
	return posts, nil
}

// ListByCreator returns posts created by a specific user.
func (r *GrowthPostRepository) ListByCreator(creatorID uint, keyword, status string) ([]model.GrowthPost, error) {
	query := r.db.Preload("Creator").
		Where("creator_id = ?", creatorID).
		Order("created_at desc, id desc")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("content LIKE ?", like)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var posts []model.GrowthPost
	if err := query.Find(&posts).Error; err != nil {
		return nil, errors.Wrap(err, "list my growth posts")
	}
	return posts, nil
}

// AdminList returns posts for admin with optional status and keyword filters.
func (r *GrowthPostRepository) AdminList(keyword, status string) ([]model.GrowthPost, error) {
	query := r.db.Preload("Creator").Order("created_at desc, id desc")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("content LIKE ?", like)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var posts []model.GrowthPost
	if err := query.Find(&posts).Error; err != nil {
		return nil, errors.Wrap(err, "admin list growth posts")
	}
	return posts, nil
}
