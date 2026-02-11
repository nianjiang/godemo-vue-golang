package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateMenusRequest request params
type CreateMenusRequest struct {
	Name     string `json:"name" binding:""`
	Path     string `json:"path" binding:""`
	Icon     string `json:"icon" binding:""`
	ParentID uint64 `json:"parentID" binding:""`
	Order    int    `json:"order" binding:""`
}

// UpdateMenusByIDRequest request params
type UpdateMenusByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Name     string `json:"name" binding:""`
	Path     string `json:"path" binding:""`
	Icon     string `json:"icon" binding:""`
	ParentID uint64 `json:"parentID" binding:""`
	Order    int    `json:"order" binding:""`
}

// MenusObjDetail detail
type MenusObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Icon      string     `json:"icon"`
	ParentID  uint64     `json:"parentID"`
	Order     int        `json:"order"`
}

// CreateMenusReply only for api docs
type CreateMenusReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteMenusByIDReply only for api docs
type DeleteMenusByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateMenusByIDReply only for api docs
type UpdateMenusByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetMenusByIDReply only for api docs
type GetMenusByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Menus MenusObjDetail `json:"menus"`
	} `json:"data"` // return data
}

// ListMenussRequest request params
type ListMenussRequest struct {
	query.Params
}

// ListMenussReply only for api docs
type ListMenussReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Menuss []MenusObjDetail `json:"menuss"`
	} `json:"data"` // return data
}
