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

// ListByEmployeeIDs 批量查询员工与店长的关联关系。
func (r *ManagerEmployeeRepository) ListByEmployeeIDs(employeeIDs []uint) ([]model.ManagerEmployee, error) {
	if len(employeeIDs) == 0 {
		return nil, nil
	}

	var relations []model.ManagerEmployee
	if err := r.db.Where("employee_id IN ?", employeeIDs).Find(&relations).Error; err != nil {
		return nil, errors.Wrap(err, "list manager employee relations")
	}
	return relations, nil
}

// ListEmployeeIDsByManager returns employee IDs managed by a manager.
func (r *ManagerEmployeeRepository) ListEmployeeIDsByManager(managerID uint) ([]uint, error) {
	var ids []uint
	if err := r.db.
		Model(&model.ManagerEmployee{}).
		Select("employee_id").
		Where("manager_id = ?", managerID).
		Pluck("employee_id", &ids).Error; err != nil {
		return nil, errors.Wrap(err, "list employees by manager")
	}
	return ids, nil
}
