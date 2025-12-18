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
	points   *PointService
}

// NewLearningService builds learning service.
func NewLearningService(
	recordRepo *repository.LearningRecordRepository,
	contentRepo *repository.ContentRepository,
	userRepo *repository.UserRepository,
	pointSvc *PointService,
) *LearningService {
	return &LearningService{
		records:  recordRepo,
		contents: contentRepo,
		users:    userRepo,
		points:   pointSvc,
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
	prevRecord := *record
	wasCompleted := record.Status == "completed"

	// 如果是文档或图文类型，打开即视为完成
	if content.Type == "doc" || content.Type == "article" {
		if record.Status != "completed" {
			record.Status = "completed"
			record.Progress = 100
			now := time.Now()
			record.CompletedAt = &now
		}
	} else {
		// 视频类型：记录播放位置和进度
		newPos := req.VideoPosition
		if newPos < record.VideoPosition {
			newPos = record.VideoPosition
		}
		record.VideoPosition = newPos

		duration := content.DurationSeconds
		if duration <= 0 {
			duration = 1
		}

		// 计算进度百分比
		progress := int((record.VideoPosition * 100) / duration)
		if progress > 100 {
			progress = 100
		}
		record.Progress = progress

		// 判断是否完成：
		// 1. 进度 >= 95%（考虑误差）
		// 2. 或者播放位置 >= 视频时长的 95%（更准确）
		isCompleted := false
		if record.Status != "completed" {
			// 方式1：进度百分比判断（>= 95%）
			if progress >= 95 {
				isCompleted = true
			}
			// 方式2：播放位置判断（>= 视频时长的 95%）
			if duration > 0 && record.VideoPosition >= duration*95/100 {
				isCompleted = true
			}
		}

		if isCompleted {
			record.Status = "completed"
			record.Progress = 100 // 确保完成时进度为100%
			now := time.Now()
			record.CompletedAt = &now
		} else if record.Status == "" {
			record.Status = "in_progress"
		}
	}

	nowCompleted := record.Status == "completed"

	if err := s.records.Upsert(record); err != nil {
		return nil, err
	}

	if !wasCompleted && nowCompleted && s.points != nil {
		if err := s.points.AwardContentCompletion(user.ID, content); err != nil {
			rollback := prevRecord
			// best-effort rollback to previous state
			_ = s.records.Upsert(&rollback)
			return nil, err
		}
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

// AggregateByUsers aggregates learning progress for given users.
func (s *LearningService) AggregateByUsers(userIDs []uint) (map[uint]repository.LearningProgressAggregate, error) {
	return s.records.AggregateByUsers(userIDs)
}

// GetUserLearningStats returns learning statistics for a specific user.
func (s *LearningService) GetUserLearningStats(userID uint) (*dto.UserLearningStatsResponse, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}

	completed, err := s.records.CountCompletedByUser(userID)
	if err != nil {
		return nil, err
	}

	total, err := s.records.CountTotalByUser(userID)
	if err != nil {
		return nil, err
	}

	// 获取该角色可见的已发布内容总数
	totalContents, err := s.contents.CountPublishedForRole(string(user.Role))
	if err != nil {
		return nil, err
	}

	completionRate := float64(0)
	if totalContents > 0 {
		completionRate = float64(completed) / float64(totalContents) * 100
	}

	return &dto.UserLearningStatsResponse{
		UserID:         userID,
		CompletedCount: completed,
		TotalCount:     total,
		TotalContents:  totalContents,
		CompletionRate: completionRate,
	}, nil
}

// GetContentCompletionStats returns completion statistics for a specific content.
func (s *LearningService) GetContentCompletionStats(contentID uint) (*dto.ContentCompletionStatsResponse, error) {
	content, err := s.contents.FindByID(contentID)
	if err != nil {
		return nil, err
	}

	completed, err := s.records.CountCompletedByContent(contentID)
	if err != nil {
		return nil, err
	}

	total, err := s.records.CountTotalByContent(contentID)
	if err != nil {
		return nil, err
	}

	completionRate := float64(0)
	if total > 0 {
		completionRate = float64(completed) / float64(total) * 100
	}

	return &dto.ContentCompletionStatsResponse{
		ContentID:      contentID,
		ContentTitle:   content.Title,
		CompletedCount: completed,
		TotalCount:     total,
		CompletionRate: completionRate,
	}, nil
}
