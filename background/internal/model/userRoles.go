package model

type UserRoles struct {
	UserID uint64 `gorm:"column:user_id;type:bigint(20) unsigned;primary_key" json:"userID"`
	RoleID uint64 `gorm:"column:role_id;type:bigint(20) unsigned;not null" json:"roleID"`
}

// UserRolesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var UserRolesColumnNames = map[string]bool{
	"user_id": true,
	"role_id": true,
}
