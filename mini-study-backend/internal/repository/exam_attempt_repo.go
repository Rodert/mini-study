package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// ExamAttemptRepository handles exam attempt persistence.
type ExamAttemptRepository struct {
	db *gorm.DB
}

// NewExamAttemptRepository creates a repository.
func NewExamAttemptRepository(db *gorm.DB) *ExamAttemptRepository {
	return &ExamAttemptRepository{db: db}
}

// Create saves an attempt record.
func (r *ExamAttemptRepository) Create(attempt *model.ExamAttempt) error {
	if err := r.db.Create(attempt).Error; err != nil {
		return errors.Wrap(err, "create exam attempt")
	}
	return nil
}

// FindLatestByUserAndExam returns latest attempt for user & exam.
func (r *ExamAttemptRepository) FindLatestByUserAndExam(userID, examID uint) (*model.ExamAttempt, error) {
	var attempt model.ExamAttempt
	if err := r.db.
		Where("user_id = ? AND exam_id = ?", userID, examID).
		Order("created_at DESC").
		Preload("Exam").
		First(&attempt).Error; err != nil {
		return nil, errors.Wrap(err, "find latest exam attempt")
	}
	return &attempt, nil
}

// ListByUser returns all attempts for a user ordered by submission time desc.
func (r *ExamAttemptRepository) ListByUser(userID uint) ([]model.ExamAttempt, error) {
	var attempts []model.ExamAttempt
	if err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Preload("Exam").
		Find(&attempts).Error; err != nil {
		return nil, errors.Wrap(err, "list attempts by user")
	}
	return attempts, nil
}

// ListLatestByUsers returns latest attempt per user.
func (r *ExamAttemptRepository) ListLatestByUsers(userIDs []uint) (map[uint]model.ExamAttempt, error) {
	if len(userIDs) == 0 {
		return map[uint]model.ExamAttempt{}, nil
	}

	var attempts []model.ExamAttempt
	if err := r.db.
		Where("user_id IN ?", userIDs).
		Order("created_at DESC").
		Preload("Exam").
		Find(&attempts).Error; err != nil {
		return nil, errors.Wrap(err, "list attempts by user ids")
	}

	result := make(map[uint]model.ExamAttempt, len(userIDs))
	for _, attempt := range attempts {
		if _, exists := result[attempt.UserID]; exists {
			continue
		}
		result[attempt.UserID] = attempt
	}
	return result, nil
}

// ExamAggregateRow stores aggregated stats.
type ExamAggregateRow struct {
	ExamID       uint
	AttemptCount int64
	PassCount    int64
	AvgScore     float64
}

// AggregateByExamForUsers aggregates attempt stats for given users.
func (r *ExamAttemptRepository) AggregateByExamForUsers(userIDs []uint) ([]ExamAggregateRow, error) {
	if len(userIDs) == 0 {
		return []ExamAggregateRow{}, nil
	}

	var rows []ExamAggregateRow
	if err := r.db.Table("exam_attempts").
		Select("exam_id, COUNT(*) as attempt_count, SUM(CASE WHEN pass = 1 THEN 1 ELSE 0 END) as pass_count, AVG(score) as avg_score").
		Where("user_id IN ?", userIDs).
		Group("exam_id").
		Scan(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "aggregate attempts by exam")
	}
	return rows, nil
}
