package model

import (
	"time"
)

type Users struct {
	ID         uint64     `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt  *time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt  *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deletedAt"`
	UserName   string     `gorm:"column:user_name;type:varchar(255);not null" json:"userName"`
	Password   string     `gorm:"column:password;type:varchar(255);not null" json:"password"`
	UserGender string     `gorm:"column:user_gender;type:varchar(10)" json:"userGender"`
	NickName   string     `gorm:"column:nick_name;type:varchar(255)" json:"nickName"`
	UserPhone  string     `gorm:"column:user_phone;type:varchar(20)" json:"userPhone"`
	UserEmail  string     `gorm:"column:user_email;type:varchar(255)" json:"userEmail"`
	Status     string     `gorm:"column:status;type:varchar(10)" json:"status"`
}

// UsersColumnNames Whitelist for custom query fields to prevent sql injection attacks
var UsersColumnNames = map[string]bool{
	"id":          true,
	"created_at":  true,
	"updated_at":  true,
	"deleted_at":  true,
	"user_name":   true,
	"password":    true,
	"user_gender": true,
	"nick_name":   true,
	"user_phone":  true,
	"user_email":  true,
	"status":      true,
}


