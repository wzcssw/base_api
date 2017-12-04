package models

import (
	"github.com/labstack/echo"
)

// User 用户
type User struct {
	BaseModel
	Name           string `gorm:"size:255"`
	Phone          string `gorm:"size:255"`
	Email          string `gorm:"size:255"`
	PasswordDigest string `gorm:"size:128"`
	IsDelete       bool
	Realname       string `gorm:"size:128"`
	Role           Role   `gorm:"ForeignKey:RolesID"`
	HospitalID     int
	RolesID        int
	UserRoles      []UserRole `gorm:"ForeignKey:UserID"`
}

// TableName 返回表名
func (User) TableName() string {
	return "users"
}

// Logon 验证登录
func (User) Logon(username, password string, c echo.Context) (bool, error) {
	if username == "zhaoyy" && password == "123" {
		return true, nil
	}
	return false, nil
}
