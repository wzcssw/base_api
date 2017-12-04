package models

// UserRole 用户角色
type UserRole struct {
	BaseModel
	User   User `gorm:"ForeignKey:UserID"`
	UserID int
	Role   Role `gorm:"ForeignKey:RoleID"`
	RoleID int
}

// TableName 返回表名
func (UserRole) TableName() string {
	return "users_roles"
}
