package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Mdb mysql连接
var Mdb *gorm.DB

// BaseModel 基础Model
type BaseModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
