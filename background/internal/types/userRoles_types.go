package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateUserRolesRequest request params
type CreateUserRolesRequest struct {
	UserID uint64 `json:"userID" binding:""`
	RoleID uint64 `json:"roleID" binding:""`
}

// UpdateUserRolesByUserIDRequest request params
type UpdateUserRolesByUserIDRequest struct {
	UserID uint64 `json:"userID" binding:""`
	RoleID uint64 `json:"roleID" binding:""`
}

// UserRolesObjDetail detail
type UserRolesObjDetail struct {
	UserID uint64 `json:"userID"`
	RoleID uint64 `json:"roleID"`
}

// CreateUserRolesReply only for api docs
type CreateUserRolesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserID uint64 `json:"userID"`
	} `json:"data"` // return data
}

// DeleteUserRolesByUserIDReply only for api docs
type DeleteUserRolesByUserIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateUserRolesByUserIDReply only for api docs
type UpdateUserRolesByUserIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetUserRolesByUserIDReply only for api docs
type GetUserRolesByUserIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserRoles UserRolesObjDetail `json:"userRoles"`
	} `json:"data"` // return data
}

// ListUserRolesRequest request params
type ListUserRolesRequest struct {
	query.Params
}

// ListUserRolesReply only for api docs
type ListUserRolesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserRoles []UserRolesObjDetail `json:"userRoles"`
	} `json:"data"` // return data
}
