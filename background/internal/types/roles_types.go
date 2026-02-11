package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.


// CreateRolesRequest request params
type CreateRolesRequest struct {
	RoleName  string `json:"roleName" binding:""`
	RoleCode  string `json:"roleCode" binding:""`
	RoleDesc  string `json:"roleDesc" binding:""`
	Status  string `json:"status" binding:""`
}

// UpdateRolesByIDRequest request params
type UpdateRolesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	RoleName  string `json:"roleName" binding:""`
	RoleCode  string `json:"roleCode" binding:""`
	RoleDesc  string `json:"roleDesc" binding:""`
	Status  string `json:"status" binding:""`
}

// RolesObjDetail detail
type RolesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	RoleName  string `json:"roleName"`
	RoleCode  string `json:"roleCode"`
	RoleDesc  string `json:"roleDesc"`
	Status  string `json:"status"`
}


// CreateRolesReply only for api docs
type CreateRolesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteRolesByIDReply only for api docs
type DeleteRolesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateRolesByIDReply only for api docs
type UpdateRolesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetRolesByIDReply only for api docs
type GetRolesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Roles RolesObjDetail `json:"roles"`
	} `json:"data"` // return data
}

// ListRolessRequest request params
type ListRolessRequest struct {
	query.Params
}

// ListRolessReply only for api docs
type ListRolessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Roless []RolesObjDetail `json:"roless"`
	} `json:"data"` // return data
}
