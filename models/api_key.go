package models

import (
	"time"
	"tx_base_api/util"

	"github.com/jinzhu/gorm"
)

// APIKey api访问key
type APIKey struct {
	BaseModel
	AccessToken int
	RoleID      int
	LoginType   string
	LoginID     uint
	ExpiresAt   time.Time
	Active      bool
}

// TableName 返回表名
func (APIKey) TableName() string {
	return "api_keys"
}

// IsExpired 是否过期
func (apikey *APIKey) IsExpired() bool {
	return time.Now().After(apikey.ExpiresAt)
}

// AfterSave 保存以后
func (apikey *APIKey) AfterSave(scope *gorm.Scope) error {
	scope.SetColumn("AccessToken", util.RandomString(32))
	scope.SetColumn("ExpiresAt", time.Now().AddDate(0, 0, 30))
	return nil
}
