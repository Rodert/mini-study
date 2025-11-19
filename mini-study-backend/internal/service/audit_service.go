package service

import (
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
)

// AuditService encapsulates audit logging behavior.
type AuditService struct {
	repo *repository.AuditRepository
}

// NewAuditService builds an audit service.
func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

// Record stores an audit entry, ignoring persistence errors for flow safety.
func (s *AuditService) Record(actorID uint, action, target, payload, result string) error {
	entry := &model.AuditLog{
		ActorID: actorID,
		Action:  action,
		Target:  target,
		Payload: payload,
		Result:  result,
	}
	return s.repo.Create(entry)
}
