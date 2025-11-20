package model

// TableName 指定表名
func (ManagerEmployee) TableName() string {
	return "manager_employees"
}

// ManagerEmployee represents mapping between managers and employees.
type ManagerEmployee struct {
	Base
	ManagerID  uint `gorm:"not null;comment:店长ID" json:"manager_id"`
	EmployeeID uint `gorm:"not null;comment:员工ID" json:"employee_id"`
}
