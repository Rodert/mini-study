package model

// Role represents user role in system.
// 使用类型别名，底层就是 string，既能复用常量，又不会和 string 不兼容。
type Role = string

const (
	RoleEmployee Role = "employee"
	RoleManager  Role = "manager"
	RoleAdmin    Role = "admin"
)
