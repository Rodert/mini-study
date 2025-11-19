package repository

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// BannerRepository handles banner persistence.
type BannerRepository struct {
	db *gorm.DB
}

// NewBannerRepository creates banner repo.
func NewBannerRepository(db *gorm.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

// Create inserts a banner.
func (r *BannerRepository) Create(banner *model.Banner) error {
	if err := r.db.Create(banner).Error; err != nil {
		return errors.Wrap(err, "create banner")
	}
	return nil
}

// Update updates banner in DB.
func (r *BannerRepository) Update(banner *model.Banner) error {
	if err := r.db.Save(banner).Error; err != nil {
		return errors.Wrap(err, "update banner")
	}
	return nil
}

// FindByID finds banner by id.
func (r *BannerRepository) FindByID(id uint) (*model.Banner, error) {
	var banner model.Banner
	if err := r.db.First(&banner, id).Error; err != nil {
		return nil, errors.Wrap(err, "find banner")
	}
	return &banner, nil
}

// ListAdmin returns all banners optionally filtered by status.
func (r *BannerRepository) ListAdmin(status *bool) ([]model.Banner, error) {
	query := r.db.Order("sort_order asc, id desc")
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	var banners []model.Banner
	if err := query.Find(&banners).Error; err != nil {
		return nil, errors.Wrap(err, "list banners")
	}
	return banners, nil
}

// ListVisible returns banners visible to role now.
func (r *BannerRepository) ListVisible(role string, now time.Time) ([]model.Banner, error) {
	query := r.db.
		Where("status = ?", true).
		Where("start_at IS NULL OR start_at <= ?", now).
		Where("end_at IS NULL OR end_at >= ?", now).
		Order("sort_order asc, id desc")
	if role != "" && role != "both" {
		query = query.Where("visible_roles = ? OR visible_roles = ?", role, "both")
	}

	var banners []model.Banner
	if err := query.Find(&banners).Error; err != nil {
		return nil, errors.Wrap(err, "list visible banners")
	}
	return banners, nil
}

