package model

type RolePermissions struct {
	RoleID       uint64 `gorm:"column:role_id;type:bigint(20) unsigned;primary_key" json:"roleID"`
	PermissionID uint64 `gorm:"column:permission_id;type:bigint(20) unsigned;not null" json:"permissionID"`
}

// RolePermissionsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var RolePermissionsColumnNames = map[string]bool{
	"role_id":       true,
	"permission_id": true,
}


