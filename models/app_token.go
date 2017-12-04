package models

import (
	"time"
)

// AppToken api访问key
type AppToken struct {
	BaseModel
	AppName      string
	AppKey       string
	Token        string
	RefreshToken string
	ExpiresIn    int
	RevokedAt    time.Time
}

// TableName 返回表名
func (AppToken) TableName() string {
	return "app_tokens"
}
