package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// ManagerEmployeeRepository handles manager-employee relations.
type ManagerEmployeeRepository struct {
	db *gorm.DB
}

// NewManagerEmployeeRepository constructs a ManagerEmployeeRepository.
func NewManagerEmployeeRepository(db *gorm.DB) *ManagerEmployeeRepository {
	return &ManagerEmployeeRepository{db: db}
}

// CreateRelations creates manager-employee relations in batch.
func (r *ManagerEmployeeRepository) CreateRelations(employeeID uint, managerIDs []uint) error {
	if len(managerIDs) == 0 {
		return nil
	}

	var rows []model.ManagerEmployee
	for _, mid := range managerIDs {
		rows = append(rows, model.ManagerEmployee{
			ManagerID:  mid,
			EmployeeID: employeeID,
		})
	}

	if err := r.db.Create(&rows).Error; err != nil {
		return errors.Wrap(err, "create manager_employee relations")
	}
	return nil
}

// ReplaceRelations replaces employee's manager bindings with provided manager IDs.
func (r *ManagerEmployeeRepository) ReplaceRelations(employeeID uint, managerIDs []uint) error {
	if err := r.db.Where("employee_id = ?", employeeID).Delete(&model.ManagerEmployee{}).Error; err != nil {
		return errors.Wrap(err, "clear existing relations")
	}
	return r.CreateRelations(employeeID, managerIDs)
}
