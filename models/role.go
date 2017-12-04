package models

// Role 角色
type Role struct {
	BaseModel
	ID   int
	Name string
}

// TableName 返回表名
func (Role) TableName() string {
	return "roles"
}
