package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateRolePermissionsRequest request params
type CreateRolePermissionsRequest struct {
	RoleID       uint64 `json:"roleID" binding:""`
	PermissionID uint64 `json:"permissionID" binding:""`
}

// UpdateRolePermissionsByRoleIDRequest request params
type UpdateRolePermissionsByRoleIDRequest struct {
	RoleID       uint64 `json:"roleID" binding:""`
	PermissionID uint64 `json:"permissionID" binding:""`
}

// RolePermissionsObjDetail detail
type RolePermissionsObjDetail struct {
	RoleID       uint64 `json:"roleID"`
	PermissionID uint64 `json:"permissionID"`
}

// CreateRolePermissionsReply only for api docs
type CreateRolePermissionsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		RoleID uint64 `json:"roleID"`
	} `json:"data"` // return data
}

// DeleteRolePermissionsByRoleIDReply only for api docs
type DeleteRolePermissionsByRoleIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateRolePermissionsByRoleIDReply only for api docs
type UpdateRolePermissionsByRoleIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetRolePermissionsByRoleIDReply only for api docs
type GetRolePermissionsByRoleIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		RolePermissions RolePermissionsObjDetail `json:"rolePermissions"`
	} `json:"data"` // return data
}

// ListRolePermissionsRequest request params
type ListRolePermissionsRequest struct {
	query.Params
}

// ListRolePermissionsReply only for api docs
type ListRolePermissionsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		RolePermissions []RolePermissionsObjDetail `json:"rolePermissions"`
	} `json:"data"` // return data
}
