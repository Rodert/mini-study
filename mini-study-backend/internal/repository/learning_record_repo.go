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
