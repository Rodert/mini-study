package repository

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// NoticeRepository handles notice persistence.
type NoticeRepository struct {
	db *gorm.DB
}

// NewNoticeRepository creates notice repo.
func NewNoticeRepository(db *gorm.DB) *NoticeRepository {
	return &NoticeRepository{db: db}
}

// Create inserts a notice.
func (r *NoticeRepository) Create(notice *model.Notice) error {
	if err := r.db.Create(notice).Error; err != nil {
		return errors.Wrap(err, "create notice")
	}
	return nil
}

// Update updates notice in DB.
func (r *NoticeRepository) Update(notice *model.Notice) error {
	if err := r.db.Save(notice).Error; err != nil {
		return errors.Wrap(err, "update notice")
	}
	return nil
}

// FindByID finds notice by id.
func (r *NoticeRepository) FindByID(id uint) (*model.Notice, error) {
	var notice model.Notice
	if err := r.db.First(&notice, id).Error; err != nil {
		return nil, errors.Wrap(err, "find notice")
	}
	return &notice, nil
}

// ListAdmin returns all notices optionally filtered by status.
func (r *NoticeRepository) ListAdmin(status *bool) ([]model.Notice, error) {
	query := r.db.Order("start_at desc, id desc")
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	var notices []model.Notice
	if err := query.Find(&notices).Error; err != nil {
		return nil, errors.Wrap(err, "list notices")
	}
	return notices, nil
}

// FindLatestActive returns the latest active notice by time window.
func (r *NoticeRepository) FindLatestActive(now time.Time) (*model.Notice, error) {
	var notice model.Notice
	query := r.db.
		Where("status = ?", true).
		Where("start_at IS NULL OR start_at <= ?", now).
		Where("end_at IS NULL OR end_at >= ?", now).
		Order("start_at desc, id desc")

	if err := query.First(&notice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "find latest active notice")
	}
	return &notice, nil
}
