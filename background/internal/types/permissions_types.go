package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.


// CreatePermissionsRequest request params
type CreatePermissionsRequest struct {
	Name  string `json:"name" binding:""`
	Code  string `json:"code" binding:""`
	Description  string `json:"description" binding:""`
}

// UpdatePermissionsByIDRequest request params
type UpdatePermissionsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Name  string `json:"name" binding:""`
	Code  string `json:"code" binding:""`
	Description  string `json:"description" binding:""`
}

// PermissionsObjDetail detail
type PermissionsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	Name  string `json:"name"`
	Code  string `json:"code"`
	Description  string `json:"description"`
}


// CreatePermissionsReply only for api docs
type CreatePermissionsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeletePermissionsByIDReply only for api docs
type DeletePermissionsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdatePermissionsByIDReply only for api docs
type UpdatePermissionsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetPermissionsByIDReply only for api docs
type GetPermissionsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Permissions PermissionsObjDetail `json:"permissions"`
	} `json:"data"` // return data
}

// ListPermissionssRequest request params
type ListPermissionssRequest struct {
	query.Params
}

// ListPermissionssReply only for api docs
type ListPermissionssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Permissionss []PermissionsObjDetail `json:"permissionss"`
	} `json:"data"` // return data
}
