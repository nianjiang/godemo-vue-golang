package model

import (
	"time"
)

type Permissions struct {
	ID          uint64     `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deletedAt"`
	Name        string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Code        string     `gorm:"column:code;type:varchar(255);not null" json:"code"`
	Description string     `gorm:"column:description;type:text" json:"description"`
}

// PermissionsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var PermissionsColumnNames = map[string]bool{
	"id":          true,
	"created_at":  true,
	"updated_at":  true,
	"deleted_at":  true,
	"name":        true,
	"code":        true,
	"description": true,
}


