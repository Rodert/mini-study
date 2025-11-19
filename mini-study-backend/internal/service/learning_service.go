package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
)

// LearningService handles learning progress.
type LearningService struct {
	records  *repository.LearningRecordRepository
	contents *repository.ContentRepository
	users    *repository.UserRepository
}

// NewLearningService builds learning service.
func NewLearningService(
	recordRepo *repository.LearningRecordRepository,
	contentRepo *repository.ContentRepository,
	userRepo *repository.UserRepository,
) *LearningService {
	return &LearningService{
		records:  recordRepo,
		contents: contentRepo,
		users:    userRepo,
	}
}

// UpdateProgress updates the learning record for a user/content.
func (s *LearningService) UpdateProgress(userID uint, req dto.LearningProgressRequest) (*dto.LearningProgressResponse, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}

	content, err := s.contents.FindByID(req.ContentID)
	if err != nil {
		return nil, err
	}

	if err := s.ensureContentAccessible(user, content); err != nil {
		return nil, err
	}
	if content.Status != "published" && user.Role != model.RoleAdmin {
		return nil, errors.New("内容未发布")
	}

	record, err := s.records.FirstOrCreate(user.ID, content.ID)
	if err != nil {
		return nil, err
	}

	newPos := req.VideoPosition
	if newPos < record.VideoPosition {
		newPos = record.VideoPosition
	}
	record.VideoPosition = newPos

	duration := content.DurationSeconds
	if duration <= 0 {
		duration = 1
	}
	progress := int((record.VideoPosition * 100) / duration)
	if progress > 100 {
		progress = 100
	}
	record.Progress = progress

	if progress >= 99 && record.Status != "completed" {
		record.Status = "completed"
		now := time.Now()
		record.CompletedAt = &now
	} else if record.Status == "" {
		record.Status = "in_progress"
	}

	if err := s.records.Upsert(record); err != nil {
		return nil, err
	}
	return s.buildProgressResponse(record, content), nil
}

// GetProgress returns progress for user/content.
func (s *LearningService) GetProgress(userID uint, contentID uint) (*dto.LearningProgressResponse, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}

	content, err := s.contents.FindByID(contentID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureContentAccessible(user, content); err != nil {
		return nil, err
	}

	record, err := s.records.FindByUserAndContent(userID, contentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 返回默认进度
			return &dto.LearningProgressResponse{
				ContentID:       contentID,
				VideoPosition:   0,
				DurationSeconds: content.DurationSeconds,
				Progress:        0,
				Status:          "not_started",
			}, nil
		}
		return nil, err
	}
	return s.buildProgressResponse(record, content), nil
}

// ListProgress returns all learning progress for user.
func (s *LearningService) ListProgress(userID uint) ([]dto.LearningProgressResponse, error) {
	records, err := s.records.ListByUser(userID)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return []dto.LearningProgressResponse{}, nil
	}

	responses := make([]dto.LearningProgressResponse, 0, len(records))
	for idx := range records {
		content, err := s.contents.FindByID(records[idx].ContentID)
		if err != nil {
			return nil, err
		}
		resp := s.buildProgressResponse(&records[idx], content)
		responses = append(responses, *resp)
	}
	return responses, nil
}

func (s *LearningService) buildProgressResponse(record *model.LearningRecord, content *model.Content) *dto.LearningProgressResponse {
	return &dto.LearningProgressResponse{
		ContentID:       record.ContentID,
		VideoPosition:   record.VideoPosition,
		DurationSeconds: content.DurationSeconds,
		Progress:        record.Progress,
		Status:          record.Status,
	}
}

func (s *LearningService) ensureContentAccessible(user *model.User, content *model.Content) error {
	if user.Role == model.RoleAdmin {
		return nil
	}
	if content.VisibleRoles != "both" && content.VisibleRoles != user.Role {
		return errors.New("无权访问该内容")
	}
	return nil
}
