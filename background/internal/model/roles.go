package model

import (
	"time"
)

type Roles struct {
	ID        uint64     `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deletedAt"`
	RoleName  string     `gorm:"column:role_name;type:varchar(255);not null" json:"roleName"`
	RoleCode  string     `gorm:"column:role_code;type:varchar(255);not null" json:"roleCode"`
	RoleDesc  string     `gorm:"column:role_desc;type:text" json:"roleDesc"`
	Status    string     `gorm:"column:status;type:varchar(10)" json:"status"`
}

// RolesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var RolesColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"role_name":  true,
	"role_code":  true,
	"role_desc":  true,
	"status":     true,
}


