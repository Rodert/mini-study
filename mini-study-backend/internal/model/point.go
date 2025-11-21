package model

// TableName specifies the user points table.
func (UserPoint) TableName() string {
	return "user_points"
}

// UserPoint keeps the latest point balance for a user.
type UserPoint struct {
	Base
	UserID uint  `gorm:"uniqueIndex;not null;comment:用户ID" json:"user_id"`
	Total  int64 `gorm:"not null;default:0;comment:积分总数" json:"total"`
}

// TableName specifies the point transactions table.
func (PointTransaction) TableName() string {
	return "point_transactions"
}

// PointTransaction records each point change for auditing.
type PointTransaction struct {
	Base
	UserID      uint   `gorm:"index;not null;comment:用户ID;uniqueIndex:idx_user_ref_source,priority:1" json:"user_id"`
	Change      int64  `gorm:"not null;comment:积分变动，正数为增加" json:"change"`
	Source      string `gorm:"size:32;not null;comment:积分来源;uniqueIndex:idx_user_ref_source,priority:3" json:"source"`
	ReferenceID string `gorm:"size:100;not null;default:'';uniqueIndex:idx_user_ref_source,priority:2;comment:业务关联ID" json:"reference_id"`
	ContentID   *uint  `gorm:"comment:关联内容ID" json:"content_id,omitempty"`
	Description string `gorm:"size:255;comment:描述信息" json:"description"`
	Memo        string `gorm:"size:255;comment:备注" json:"memo"`
}
