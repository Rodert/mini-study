package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
)

// NoticeService handles notice business logic.
type NoticeService struct {
	repo     *repository.NoticeRepository
	userRepo *repository.UserRepository
	audit    *AuditService
}

// NewNoticeService creates notice service.
func NewNoticeService(repo *repository.NoticeRepository, userRepo *repository.UserRepository, audit *AuditService) *NoticeService {
	return &NoticeService{repo: repo, userRepo: userRepo, audit: audit}
}

// GetLatestNotice returns the latest active notice for current user.
func (s *NoticeService) GetLatestNotice(userID uint) (*model.Notice, error) {
	if userID == 0 {
		return nil, errors.New("未登录")
	}
	// ensure user exists
	if _, err := s.userRepo.FindByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return s.repo.FindLatestActive(time.Now())
}

// AdminListNotices lists notices for admin.
func (s *NoticeService) AdminListNotices(adminID uint, status *bool) ([]model.Notice, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}
	return s.repo.ListAdmin(status)
}

// AdminCreateNotice creates notice.
func (s *NoticeService) AdminCreateNotice(adminID uint, req dto.AdminCreateNoticeRequest) (*model.Notice, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	notice := &model.Notice{
		Title:   req.Title,
		Content: req.Content,
		ImageURL: req.ImageURL,
		StartAt:  req.StartAt,
		EndAt:    req.EndAt,
		Status:   true,
	}
	if req.Status != nil {
		notice.Status = *req.Status
	}

	if err := s.repo.Create(notice); err != nil {
		return nil, err
	}
	_ = s.audit.Record(adminID, "create_notice", "notices", notice.Title, "success")
	return notice, nil
}

// AdminUpdateNotice updates notice fields.
func (s *NoticeService) AdminUpdateNotice(adminID, noticeID uint, req dto.AdminUpdateNoticeRequest) (*model.Notice, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	notice, err := s.repo.FindByID(noticeID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		notice.Title = req.Title
	}
	if req.Content != "" {
		notice.Content = req.Content
	}
	if req.ImageURL != "" {
		notice.ImageURL = req.ImageURL
	}
	if req.Status != nil {
		notice.Status = *req.Status
	}
	if req.StartAt != nil || req.EndAt != nil {
		notice.StartAt = req.StartAt
		notice.EndAt = req.EndAt
	}

	if err := s.repo.Update(notice); err != nil {
		return nil, err
	}
	_ = s.audit.Record(adminID, "update_notice", "notices", notice.Title, "success")
	return notice, nil
}

func (s *NoticeService) ensureAdmin(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("管理员不存在")
		}
		return err
	}
	if user.Role != model.RoleAdmin {
		return errors.New("仅管理员可操作")
	}
	return nil
}
