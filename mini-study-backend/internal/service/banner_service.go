package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
)

// BannerService handles banner business logic.
type BannerService struct {
	repo     *repository.BannerRepository
	userRepo *repository.UserRepository
	audit    *AuditService
}

// NewBannerService creates banner service.
func NewBannerService(repo *repository.BannerRepository, userRepo *repository.UserRepository, audit *AuditService) *BannerService {
	return &BannerService{repo: repo, userRepo: userRepo, audit: audit}
}

// ListVisible returns banners available to current user role.
func (s *BannerService) ListVisible(userID uint) ([]model.Banner, error) {
	role, _, err := s.resolveUserRole(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListVisible(role, time.Now())
}

// AdminListBanners lists banners for admin.
func (s *BannerService) AdminListBanners(adminID uint, status *bool) ([]model.Banner, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}
	return s.repo.ListAdmin(status)
}

// AdminCreateBanner creates banner.
func (s *BannerService) AdminCreateBanner(adminID uint, req dto.AdminCreateBannerRequest) (*model.Banner, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	banner := &model.Banner{
		Title:        req.Title,
		ImageURL:     req.ImageURL,
		LinkURL:      req.LinkURL,
		VisibleRoles: req.VisibleRoles,
		SortOrder:    req.SortOrder,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		Status:       true,
	}
	if banner.VisibleRoles == "" {
		banner.VisibleRoles = "both"
	}
	if req.Status != nil {
		banner.Status = *req.Status
	}

	if err := s.repo.Create(banner); err != nil {
		return nil, err
	}
	_ = s.audit.Record(adminID, "create_banner", "banners", banner.Title, "success")
	return banner, nil
}

// AdminUpdateBanner updates banner fields.
func (s *BannerService) AdminUpdateBanner(adminID, bannerID uint, req dto.AdminUpdateBannerRequest) (*model.Banner, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	banner, err := s.repo.FindByID(bannerID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		banner.Title = req.Title
	}
	if req.ImageURL != "" {
		banner.ImageURL = req.ImageURL
	}
	if req.LinkURL != "" {
		banner.LinkURL = req.LinkURL
	}
	if req.VisibleRoles != "" {
		banner.VisibleRoles = req.VisibleRoles
	}
	if req.SortOrder != nil {
		banner.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		banner.Status = *req.Status
	}
	if req.StartAt != nil || req.EndAt != nil {
		banner.StartAt = req.StartAt
		banner.EndAt = req.EndAt
	}

	if err := s.repo.Update(banner); err != nil {
		return nil, err
	}
	_ = s.audit.Record(adminID, "update_banner", "banners", banner.Title, "success")
	return banner, nil
}

func (s *BannerService) ensureAdmin(userID uint) error {
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

func (s *BannerService) resolveUserRole(userID uint) (string, *model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", nil, err
	}
	if user.Role == model.RoleAdmin {
		return "", user, nil
	}
	return user.Role, user, nil
}
