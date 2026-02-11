package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateUsersRequest request params
type CreateUsersRequest struct {
	UserName   string `json:"userName" binding:""`
	Password   string `json:"password" binding:""`
	UserGender string `json:"userGender" binding:""`
	NickName   string `json:"nickName" binding:""`
	UserPhone  string `json:"userPhone" binding:""`
	UserEmail  string `json:"userEmail" binding:""`
	Status     string `json:"status" binding:""`
}

// UpdateUsersByIDRequest request params
type UpdateUsersByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	UserName   string `json:"userName" binding:""`
	Password   string `json:"password" binding:""`
	UserGender string `json:"userGender" binding:""`
	NickName   string `json:"nickName" binding:""`
	UserPhone  string `json:"userPhone" binding:""`
	UserEmail  string `json:"userEmail" binding:""`
	Status     string `json:"status" binding:""`
}

// UsersObjDetail detail
type UsersObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	UserName   string     `json:"userName"`
	Password   string     `json:"password"`
	UserGender string     `json:"userGender"`
	NickName   string     `json:"nickName"`
	UserPhone  string     `json:"userPhone"`
	UserEmail  string     `json:"userEmail"`
	Status     string     `json:"status"`
}

// CreateUsersReply only for api docs
type CreateUsersReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteUsersByIDReply only for api docs
type DeleteUsersByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateUsersByIDReply only for api docs
type UpdateUsersByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetUsersByIDReply only for api docs
type GetUsersByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users UsersObjDetail `json:"users"`
	} `json:"data"` // return data
}

// ListUserssRequest request params
type ListUserssRequest struct {
	query.Params
}

// ListUserssReply only for api docs
type ListUserssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Userss []UsersObjDetail `json:"userss"`
	} `json:"data"` // return data
}
