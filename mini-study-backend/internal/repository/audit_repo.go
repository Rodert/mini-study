package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// AuditRepository 审计日志仓储层。
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository 创建审计仓库实例。
func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// Create 写入一条审计日志。
func (r *AuditRepository) Create(entry *model.AuditLog) error {
	if err := r.db.Create(entry).Error; err != nil {
		return errors.Wrap(err, "create audit log")
	}
	return nil
}
