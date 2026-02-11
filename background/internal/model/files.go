package model

import (
	"time"
)

type Files struct {
	ID        uint64     `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deletedAt"`
	Filename  string     `gorm:"column:filename;type:varchar(255);not null" json:"filename"`
	URL       string     `gorm:"column:url;type:varchar(255);not null" json:"url"`
	Size      int64      `gorm:"column:size;type:bigint(20)" json:"size"`
	MimeType  string     `gorm:"column:mime_type;type:varchar(100)" json:"mimeType"`
	UserID    uint64     `gorm:"column:user_id;type:bigint(20) unsigned" json:"userID"`
}

// FilesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var FilesColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"filename":   true,
	"url":        true,
	"size":       true,
	"mime_type":  true,
	"user_id":    true,
}
