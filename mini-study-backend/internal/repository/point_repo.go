package repository

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// PointRepository handles user point persistence.
type PointRepository struct {
	db *gorm.DB
}

// NewPointRepository builds a PointRepository.
func NewPointRepository(db *gorm.DB) *PointRepository {
	return &PointRepository{db: db}
}

// UserPointWithUser joins user basic info with their point total.
type UserPointWithUser struct {
	ID        uint      `gorm:"column:id"`
	WorkNo    string    `gorm:"column:work_no"`
	Name      string    `gorm:"column:name"`
	Phone     string    `gorm:"column:phone"`
	Role      string    `gorm:"column:role"`
	Status    bool      `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Points    int64     `gorm:"column:points"`
}

// AddTransaction increments user points and records a transaction atomically.
func (r *PointRepository) AddTransaction(txn *model.PointTransaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var balance model.UserPoint
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", txn.UserID).
			First(&balance).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				balance = model.UserPoint{UserID: txn.UserID, Total: 0}
				if err := tx.Create(&balance).Error; err != nil {
					return errors.Wrap(err, "create user point balance")
				}
			} else {
				return errors.Wrap(err, "query user point balance")
			}
		}

		balance.Total += txn.Change
		if err := tx.Model(&balance).Update("total", balance.Total).Error; err != nil {
			return errors.Wrap(err, "update user point balance")
		}

		if err := tx.Create(txn).Error; err != nil {
			return errors.Wrap(err, "create point transaction")
		}
		return nil
	})
}

// ExistsByReference checks whether a transaction with the same reference already exists.
func (r *PointRepository) ExistsByReference(userID uint, referenceID, source string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.PointTransaction{}).
		Where("user_id = ? AND reference_id = ? AND source = ?", userID, referenceID, source).
		Count(&count).Error; err != nil {
		return false, errors.Wrap(err, "count point transactions by reference")
	}
	return count > 0, nil
}

// GetTotalsByUserIDs returns point totals for the given users.
func (r *PointRepository) GetTotalsByUserIDs(userIDs []uint) (map[uint]int64, error) {
	if len(userIDs) == 0 {
		return map[uint]int64{}, nil
	}

	var balances []model.UserPoint
	if err := r.db.Where("user_id IN ?", userIDs).Find(&balances).Error; err != nil {
		return nil, errors.Wrap(err, "list user point balances")
	}

	result := make(map[uint]int64, len(userIDs))
	for _, balance := range balances {
		result[balance.UserID] = balance.Total
	}
	return result, nil
}

// GetTotalByUserID returns single user's total points.
func (r *PointRepository) GetTotalByUserID(userID uint) (int64, error) {
	var balance model.UserPoint
	if err := r.db.Where("user_id = ?", userID).First(&balance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, errors.Wrap(err, "get user point balance")
	}
	return balance.Total, nil
}

// ListTransactionsByUser returns paginated transactions for the user.
func (r *PointRepository) ListTransactionsByUser(userID uint, page, pageSize int) ([]model.PointTransaction, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var total int64
	if err := r.db.Model(&model.PointTransaction{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "count point transactions")
	}

	var items []model.PointTransaction
	if total == 0 {
		return []model.PointTransaction{}, 0, nil
	}

	if err := r.db.Where("user_id = ?", userID).
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, errors.Wrap(err, "list point transactions")
	}

	return items, total, nil
}

// ListAllUserPointsWithUserInfo returns paginated list of users with their point totals, sorted by points descending.
func (r *PointRepository) ListAllUserPointsWithUserInfo(role, keyword string, page, pageSize int) ([]UserPointWithUser, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// Build base query for counting
	countQuery := r.db.Table("users")
	if role != "" {
		countQuery = countQuery.Where("role = ?", role)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		countQuery = countQuery.Where("work_no LIKE ? OR name LIKE ? OR phone LIKE ?", like, like, like)
	}

	// Count total users
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "count users with points")
	}

	// Get paginated results
	var items []UserPointWithUser
	if total == 0 {
		return []UserPointWithUser{}, 0, nil
	}

	// Build query with joins for listing
	query := r.db.Table("users u").
		Select("u.id, u.work_no, u.name, u.phone, u.role, u.status, u.created_at, COALESCE(up.total, 0) as points").
		Joins("LEFT JOIN user_points up ON u.id = up.user_id")

	// Apply filters
	if role != "" {
		query = query.Where("u.role = ?", role)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("u.work_no LIKE ? OR u.name LIKE ? OR u.phone LIKE ?", like, like, like)
	}

	// Execute query with pagination
	if err := query.
		Order("COALESCE(up.total, 0) DESC, u.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&items).Error; err != nil {
		return nil, 0, errors.Wrap(err, "list users with points")
	}

	return items, total, nil
}
