package service

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
)

const (
	pointSourceContentCompletion = "content_completion"
	contentCompletionPoints      = 1
)

// PointService handles awarding and querying user points.
type PointService struct {
	repo  *repository.PointRepository
	users *repository.UserRepository
}

// NewPointService creates a PointService.
func NewPointService(pointRepo *repository.PointRepository, userRepo *repository.UserRepository) *PointService {
	return &PointService{
		repo:  pointRepo,
		users: userRepo,
	}
}

// AwardContentCompletion gives points when a user completes a content.
func (s *PointService) AwardContentCompletion(userID uint, content *model.Content) error {
	referenceID := fmt.Sprintf("content:%d", content.ID)
	exists, err := s.repo.ExistsByReference(userID, referenceID, pointSourceContentCompletion)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	txn := &model.PointTransaction{
		UserID:      userID,
		Change:      contentCompletionPoints,
		Source:      pointSourceContentCompletion,
		ReferenceID: referenceID,
		ContentID:   &content.ID,
		Description: fmt.Sprintf("完成学习内容《%s》", content.Title),
	}
	return s.repo.AddTransaction(txn)
}

// GetTotalsMap returns point totals for users.
func (s *PointService) GetTotalsMap(userIDs []uint) (map[uint]int64, error) {
	return s.repo.GetTotalsByUserIDs(userIDs)
}

// AdminUserPointDetails returns point summary and transactions for a user.
func (s *PointService) AdminUserPointDetails(adminID, targetUserID uint, query dto.PointTransactionListQuery) (*dto.UserPointDetailResponse, error) {
	if _, err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	target, err := s.users.FindByID(targetUserID)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.GetTotalByUserID(targetUserID)
	if err != nil {
		return nil, err
	}

	page := query.Page
	if page == 0 {
		page = 1
	}
	size := query.PageSize
	if size == 0 {
		size = 20
	}

	transactions, totalCount, err := s.repo.ListTransactionsByUser(targetUserID, page, size)
	if err != nil {
		return nil, err
	}

	resp := dto.UserPointDetailResponse{
		User: dto.UserResponse{
			ID:     target.ID,
			WorkNo: target.WorkNo,
			Phone:  target.Phone,
			Name:   target.Name,
			Role:   target.Role,
			Status: target.Status,
		},
		TotalPoints:  total,
		Transactions: make([]dto.PointTransactionResponse, 0, len(transactions)),
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: size,
			Total:    totalCount,
		},
	}

	for _, txn := range transactions {
		resp.Transactions = append(resp.Transactions, dto.PointTransactionResponse{
			ID:          txn.ID,
			UserID:      txn.UserID,
			Change:      txn.Change,
			Source:      txn.Source,
			ReferenceID: txn.ReferenceID,
			ContentID:   txn.ContentID,
			Description: txn.Description,
			Memo:        txn.Memo,
			CreatedAt:   txn.CreatedAt,
		})
	}

	return &resp, nil
}

// AdminListAllPoints returns paginated list of all users with their points.
func (s *PointService) AdminListAllPoints(adminID uint, query dto.AdminListPointsQuery) (*dto.AdminListPointsResponse, error) {
	if _, err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	page := query.Page
	if page == 0 {
		page = 1
	}
	size := query.PageSize
	if size == 0 {
		size = 20
	}

	items, total, err := s.repo.ListAllUserPointsWithUserInfo(query.Role, query.Keyword, page, size)
	if err != nil {
		return nil, err
	}

	result := &dto.AdminListPointsResponse{
		Items: make([]dto.UserPointListItem, 0, len(items)),
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: size,
			Total:    total,
		},
	}

	for _, item := range items {
		result.Items = append(result.Items, dto.UserPointListItem{
			UserResponse: dto.UserResponse{
				ID:     item.ID,
				WorkNo: item.WorkNo,
				Name:   item.Name,
				Phone:  item.Phone,
				Role:   item.Role,
				Status: item.Status,
			},
			Points: item.Points,
		})
	}

	return result, nil
}

func (s *PointService) ensureAdmin(userID uint) (*model.User, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user.Role != model.RoleAdmin {
		return nil, errors.New("无权限，仅管理员可操作")
	}
	return user, nil
}
