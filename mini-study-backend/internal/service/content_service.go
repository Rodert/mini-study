package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
)

// ContentService handles business logic around categories and contents.
type ContentService struct {
	categories *repository.ContentCategoryRepository
	contents   *repository.ContentRepository
	users      *repository.UserRepository
}

// NewContentService builds a content service.
func NewContentService(
	categoryRepo *repository.ContentCategoryRepository,
	contentRepo *repository.ContentRepository,
	userRepo *repository.UserRepository,
) *ContentService {
	return &ContentService{
		categories: categoryRepo,
		contents:   contentRepo,
		users:      userRepo,
	}
}

// ListCategories returns categories visible to current user role.
func (s *ContentService) ListCategories(userID uint) ([]model.ContentCategory, error) {
	roleFilter, _, err := s.resolveUserRole(userID)
	if err != nil {
		return nil, err
	}
	return s.categories.ListByRole(roleFilter)
}

// ListCategoriesWithCount returns categories visible to current user role with content count.
func (s *ContentService) ListCategoriesWithCount(userID uint) ([]model.ContentCategory, []int64, error) {
	roleFilter, _, err := s.resolveUserRole(userID)
	if err != nil {
		return nil, nil, err
	}
	categories, err := s.categories.ListByRole(roleFilter)
	if err != nil {
		return nil, nil, err
	}
	
	counts := make([]int64, len(categories))
	for i, category := range categories {
		count, err := s.contents.CountPublishedByCategoryAndRole(category.ID, roleFilter)
		if err != nil {
			return nil, nil, err
		}
		counts[i] = count
	}
	
	return categories, counts, nil
}

// AdminCreateContent creates a new content entry.
func (s *ContentService) AdminCreateContent(adminID uint, req dto.AdminCreateContentRequest) (*model.Content, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	category, err := s.categories.FindByID(req.CategoryID)
	if err != nil {
		return nil, err
	}
	if !category.Status {
		return nil, errors.New("分类已禁用")
	}

	if req.Type == "video" && req.DurationSeconds <= 0 {
		return nil, errors.New("视频必须提供 duration_seconds")
	}
	if req.VisibleRoles == "" {
		req.VisibleRoles = category.RoleScope
	}

	content := &model.Content{
		Title:           req.Title,
		Type:            req.Type,
		CategoryID:      req.CategoryID,
		VisibleRoles:    req.VisibleRoles,
		FilePath:        req.FilePath,
		CoverURL:        req.CoverURL,
		Summary:         req.Summary,
		Status:          req.Status,
		CreatorID:       adminID,
		DurationSeconds: req.DurationSeconds,
	}

	if content.Status == "" {
		content.Status = "draft"
	}
	if content.Status == "published" {
		now := time.Now()
		content.PublishAt = &now
	}

	if err := s.contents.Create(content); err != nil {
		return nil, err
	}
	return content, nil
}

// AdminUpdateContent updates content basic info/status.
func (s *ContentService) AdminUpdateContent(adminID, contentID uint, req dto.AdminUpdateContentRequest) (*model.Content, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	content, err := s.contents.FindByID(contentID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		content.Title = req.Title
	}
	if req.Type != "" {
		content.Type = req.Type
	}
	if req.CategoryID > 0 {
		category, err := s.categories.FindByID(req.CategoryID)
		if err != nil {
			return nil, err
		}
		if !category.Status {
			return nil, errors.New("分类已禁用")
		}
		content.CategoryID = req.CategoryID
		if req.VisibleRoles == "" {
			content.VisibleRoles = category.RoleScope
		}
	}
	if req.FilePath != "" {
		content.FilePath = req.FilePath
	}
	if req.CoverURL != "" {
		content.CoverURL = req.CoverURL
	}
	if req.Summary != "" {
		content.Summary = req.Summary
	}
	if req.VisibleRoles != "" {
		content.VisibleRoles = req.VisibleRoles
	}
	if req.DurationSeconds > 0 {
		content.DurationSeconds = req.DurationSeconds
	}
	if req.Status != "" {
		content.Status = req.Status
		if req.Status == "published" && content.PublishAt == nil {
			now := time.Now()
			content.PublishAt = &now
		}
	}

	if err := s.contents.Update(content); err != nil {
		return nil, err
	}
	return content, nil
}

// AdminListContents lists contents for admin.
func (s *ContentService) AdminListContents(adminID uint, filter dto.AdminListContentRequest) ([]model.Content, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}
	return s.contents.ListAdmin(filter.CategoryID, filter.Type, filter.Status)
}

// ListPublished returns contents for user.
func (s *ContentService) ListPublished(userID uint, categoryID uint, contentType string) ([]model.Content, error) {
	roleFilter, _, err := s.resolveUserRole(userID)
	if err != nil {
		return nil, err
	}
	return s.contents.ListPublishedByRole(roleFilter, categoryID, contentType)
}

// GetPublishedDetail returns published content if visible to user.
func (s *ContentService) GetPublishedDetail(userID, id uint) (*model.Content, error) {
	roleFilter, user, err := s.resolveUserRole(userID)
	if err != nil {
		return nil, err
	}

	content, err := s.contents.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user.Role != model.RoleAdmin {
		if content.Status != "published" {
			return nil, errors.New("内容未发布")
		}
		if roleFilter != "" && content.VisibleRoles != "both" && content.VisibleRoles != roleFilter {
			return nil, errors.New("无权查看该内容")
		}
	}
	return content, nil
}

func (s *ContentService) ensureAdmin(userID uint) error {
	user, err := s.users.FindByID(userID)
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

func (s *ContentService) resolveUserRole(userID uint) (string, *model.User, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return "", nil, err
	}
	if user.Role == model.RoleAdmin {
		return "", user, nil
	}
	return user.Role, user, nil
}
