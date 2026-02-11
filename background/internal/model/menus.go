package model

import (
	"time"
)

type Menus struct {
	ID        uint64     `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deletedAt"`
	Name      string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Path      string     `gorm:"column:path;type:varchar(255);not null" json:"path"`
	Icon      string     `gorm:"column:icon;type:varchar(255)" json:"icon"`
	ParentID  uint64     `gorm:"column:parent_id;type:bigint(20) unsigned" json:"parentID"`
	Order     int        `gorm:"column:order;type:int(11)" json:"order"`
}

// MenusColumnNames Whitelist for custom query fields to prevent sql injection attacks
var MenusColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"name":       true,
	"path":       true,
	"icon":       true,
	"parent_id":  true,
	"order":      true,
}
