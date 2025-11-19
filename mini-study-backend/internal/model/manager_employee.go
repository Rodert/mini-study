package model

// ManagerEmployee represents mapping between managers and employees.
type ManagerEmployee struct {
	Base
	ManagerID  uint `gorm:"not null" json:"manager_id"`
	EmployeeID uint `gorm:"not null" json:"employee_id"`
}
