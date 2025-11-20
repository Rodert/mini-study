package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// LearningRecordRepository 学习记录仓储。
type LearningRecordRepository struct {
	db *gorm.DB
}

// NewLearningRecordRepository 创建仓库实例。
func NewLearningRecordRepository(db *gorm.DB) *LearningRecordRepository {
	return &LearningRecordRepository{db: db}
}

// FindByUserAndContent 根据用户与内容查询记录。
func (r *LearningRecordRepository) FindByUserAndContent(userID, contentID uint) (*model.LearningRecord, error) {
	var record model.LearningRecord
	if err := r.db.Where("user_id = ? AND content_id = ?", userID, contentID).First(&record).Error; err != nil {
		return nil, errors.Wrap(err, "find learning record")
	}
	return &record, nil
}

// Upsert 保存学习记录（存在则更新，不存在则创建）。
func (r *LearningRecordRepository) Upsert(record *model.LearningRecord) error {
	if record.ID == 0 {
		if err := r.db.Create(record).Error; err != nil {
			return errors.Wrap(err, "create learning record")
		}
		return nil
	}
	if err := r.db.Save(record).Error; err != nil {
		return errors.Wrap(err, "update learning record")
	}
	return nil
}

// FirstOrCreate 确保记录存在并返回。
func (r *LearningRecordRepository) FirstOrCreate(userID, contentID uint) (*model.LearningRecord, error) {
	record := &model.LearningRecord{
		UserID:    userID,
		ContentID: contentID,
		Status:    "in_progress",
	}
	if err := r.db.Where("user_id = ? AND content_id = ?", userID, contentID).FirstOrCreate(record).Error; err != nil {
		return nil, errors.Wrap(err, "first or create learning record")
	}
	return record, nil
}

// ListByUser 列出用户的全部学习记录。
func (r *LearningRecordRepository) ListByUser(userID uint) ([]model.LearningRecord, error) {
	var records []model.LearningRecord
	if err := r.db.Where("user_id = ?", userID).Find(&records).Error; err != nil {
		return nil, errors.Wrap(err, "list learning records by user")
	}
	return records, nil
}

// LearningProgressAggregate holds aggregated learning stats per user.
type LearningProgressAggregate struct {
	UserID    uint
	Completed int64
	Total     int64
}

// AggregateByUsers aggregates learning statuses for the given users.
func (r *LearningRecordRepository) AggregateByUsers(userIDs []uint) (map[uint]LearningProgressAggregate, error) {
	if len(userIDs) == 0 {
		return map[uint]LearningProgressAggregate{}, nil
	}

	var rows []LearningProgressAggregate
	if err := r.db.Table("learning_records").
		Select("user_id, SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) AS completed, COUNT(*) AS total").
		Where("user_id IN ?", userIDs).
		Group("user_id").
		Scan(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "aggregate learning records")
	}

	result := make(map[uint]LearningProgressAggregate, len(rows))
	for _, row := range rows {
		result[row.UserID] = row
	}
	return result, nil
}

// ContentCompletionAggregate holds aggregated completion stats per content.
type ContentCompletionAggregate struct {
	ContentID uint
	Completed int64
	Total     int64
}

// AggregateByContents aggregates learning completion statuses for the given contents.
func (r *LearningRecordRepository) AggregateByContents(contentIDs []uint) (map[uint]ContentCompletionAggregate, error) {
	if len(contentIDs) == 0 {
		return map[uint]ContentCompletionAggregate{}, nil
	}

	var rows []ContentCompletionAggregate
	if err := r.db.Table("learning_records").
		Select("content_id, SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) AS completed, COUNT(*) AS total").
		Where("content_id IN ?", contentIDs).
		Group("content_id").
		Scan(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "aggregate learning records by content")
	}

	result := make(map[uint]ContentCompletionAggregate, len(rows))
	for _, row := range rows {
		result[row.ContentID] = row
	}
	return result, nil
}

// CountCompletedByUser counts completed learning records for a user.
func (r *LearningRecordRepository) CountCompletedByUser(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.LearningRecord{}).
		Where("user_id = ? AND status = ?", userID, "completed").
		Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, "count completed records by user")
	}
	return count, nil
}

// CountTotalByUser counts total learning records for a user.
func (r *LearningRecordRepository) CountTotalByUser(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.LearningRecord{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, "count total records by user")
	}
	return count, nil
}

// CountCompletedByContent counts completed learning records for a content.
func (r *LearningRecordRepository) CountCompletedByContent(contentID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.LearningRecord{}).
		Where("content_id = ? AND status = ?", contentID, "completed").
		Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, "count completed records by content")
	}
	return count, nil
}

// CountTotalByContent counts total learning records for a content.
func (r *LearningRecordRepository) CountTotalByContent(contentID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.LearningRecord{}).
		Where("content_id = ?", contentID).
		Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, "count total records by content")
	}
	return count, nil
}
