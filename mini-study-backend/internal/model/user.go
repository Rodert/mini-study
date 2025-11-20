package model

// User represents an application user.
// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// User 用户表
type User struct {
	Base
	WorkNo       string `gorm:"size:50;uniqueIndex;not null;comment:工号" json:"work_no"`
	Phone        string `gorm:"size:20;comment:手机号" json:"phone"`
	PasswordHash string `gorm:"size:255;not null;comment:密码哈希" json:"-"`
	Role         Role   `gorm:"size:16;default:'employee';comment:角色(employee员工/manager店长/admin管理员)" json:"role"`
	Name         string `gorm:"size:100;comment:姓名" json:"name"`
	Status       bool   `gorm:"default:true;comment:状态(启用/禁用)" json:"status"`
}
