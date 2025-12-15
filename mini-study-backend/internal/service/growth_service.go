package service

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
)

const (
	GrowthStatusPending  = "pending"
	GrowthStatusApproved = "approved"
	GrowthStatusRejected = "rejected"
)

// GrowthService 处理成长圈业务逻辑。
type GrowthService struct {
	posts *repository.GrowthPostRepository
	users *repository.UserRepository
	audit *AuditService
}

// NewGrowthService 创建成长圈服务。
func NewGrowthService(posts *repository.GrowthPostRepository, users *repository.UserRepository, audit *AuditService) *GrowthService {
	return &GrowthService{posts: posts, users: users, audit: audit}
}

func (s *GrowthService) ensureAdmin(userID uint) error {
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

func (s *GrowthService) ensureManager(userID uint) (*model.User, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user.Role != model.RoleManager {
		return nil, errors.New("仅店长可发布成长圈")
	}
	return user, nil
}

// CreatePost 店长创建成长圈动态。
func (s *GrowthService) CreatePost(creatorID uint, req dto.CreateGrowthPostRequest) (*dto.GrowthPostResponse, error) {
	user, err := s.ensureManager(creatorID)
	if err != nil {
		return nil, err
	}

	if len(req.ImagePaths) > 9 {
		return nil, errors.New("每条动态最多上传9张图片")
	}

	imgJSON, err := json.Marshal(req.ImagePaths)
	if err != nil {
		return nil, err
	}

	post := &model.GrowthPost{
		CreatorID:  user.ID,
		Content:    req.Content,
		ImagePaths: string(imgJSON),
		Status:     GrowthStatusPending,
	}

	if err := s.posts.Create(post); err != nil {
		return nil, err
	}
	if s.audit != nil {
		_ = s.audit.Record(creatorID, "create_growth_post", "growth_posts", post.Content, "success")
	}
	return s.toResponse(post), nil
}

// ListPublic 返回所有已审核通过的成长圈动态。
func (s *GrowthService) ListPublic(keyword string) ([]dto.GrowthPostResponse, error) {
	posts, err := s.posts.ListPublic(keyword)
	if err != nil {
		return nil, err
	}
	return s.toResponses(posts), nil
}

// ListMine 返回当前用户自己的成长圈动态列表。
func (s *GrowthService) ListMine(userID uint, query dto.GrowthMyListQuery) ([]dto.GrowthPostResponse, error) {
	posts, err := s.posts.ListByCreator(userID, query.Keyword, query.Status)
	if err != nil {
		return nil, err
	}
	return s.toResponses(posts), nil
}

// AdminList 返回管理员视角的成长圈动态列表。
func (s *GrowthService) AdminList(adminID uint, query dto.AdminGrowthListQuery) ([]dto.GrowthPostResponse, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}
	posts, err := s.posts.AdminList(query.Keyword, query.Status)
	if err != nil {
		return nil, err
	}
	return s.toResponses(posts), nil
}

// Approve 审核通过某条成长圈动态。
func (s *GrowthService) Approve(adminID, postID uint) (*dto.GrowthPostResponse, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}
	post, err := s.posts.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post.Status != GrowthStatusApproved {
		post.Status = GrowthStatusApproved
		now := time.Now()
		post.ApprovedAt = &now
		if err := s.posts.Update(post); err != nil {
			return nil, err
		}
		if s.audit != nil {
			_ = s.audit.Record(adminID, "approve_growth_post", "growth_posts", post.Content, "success")
		}
	}
	return s.toResponse(post), nil
}

// Reject 将某条成长圈动态标记为拒绝。
func (s *GrowthService) Reject(adminID, postID uint) (*dto.GrowthPostResponse, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}
	post, err := s.posts.FindByID(postID)
	if err != nil {
		return nil, err
	}
	if post.Status != GrowthStatusRejected {
		post.Status = GrowthStatusRejected
		post.ApprovedAt = nil
		if err := s.posts.Update(post); err != nil {
			return nil, err
		}
		if s.audit != nil {
			_ = s.audit.Record(adminID, "reject_growth_post", "growth_posts", post.Content, "success")
		}
	}
	return s.toResponse(post), nil
}

// Delete 删除成长圈动态（店长删除未通过的自己的，管理员可删除所有）。
func (s *GrowthService) Delete(userID, postID uint) error {
	user, err := s.users.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}
	post, err := s.posts.FindByID(postID)
	if err != nil {
		return err
	}

	// 管理员可以删除任意动态
	if user.Role == model.RoleAdmin {
		if err := s.posts.Delete(post); err != nil {
			return err
		}
		if s.audit != nil {
			_ = s.audit.Record(userID, "delete_growth_post", "growth_posts", post.Content, "success")
		}
		return nil
	}

	// 店长只能删除自己发布且未通过审核的动态
	if user.Role == model.RoleManager && post.CreatorID == user.ID {
		if post.Status == GrowthStatusApproved {
			return errors.New("已通过审核的动态仅管理员可删除")
		}
		if err := s.posts.Delete(post); err != nil {
			return err
		}
		if s.audit != nil {
			_ = s.audit.Record(userID, "delete_own_growth_post", "growth_posts", post.Content, "success")
		}
		return nil
	}

	return errors.New("无权删除该动态")
}

func (s *GrowthService) toResponses(posts []model.GrowthPost) []dto.GrowthPostResponse {
	resp := make([]dto.GrowthPostResponse, 0, len(posts))
	for i := range posts {
		resp = append(resp, *s.toResponse(&posts[i]))
	}
	return resp
}

func (s *GrowthService) toResponse(post *model.GrowthPost) *dto.GrowthPostResponse {
	var images []string
	if post.ImagePaths != "" {
		_ = json.Unmarshal([]byte(post.ImagePaths), &images)
	}

	publisherName := ""
	publisherRole := ""
	if post.Creator.ID != 0 {
		publisherName = post.Creator.Name
		publisherRole = string(post.Creator.Role)
	}

	return &dto.GrowthPostResponse{
		ID:            post.ID,
		Content:       post.Content,
		ImagePaths:    images,
		Status:        post.Status,
		PublisherID:   post.CreatorID,
		PublisherName: publisherName,
		PublisherRole: publisherRole,
		CreatedAt:     post.CreatedAt,
		ApprovedAt:    post.ApprovedAt,
	}
}
