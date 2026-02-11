package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.


// CreateFilesRequest request params
type CreateFilesRequest struct {
	Filename  string `json:"filename" binding:""`
	URL  string `json:"url" binding:""`
	Size  int64 `json:"size" binding:""`
	MimeType  string `json:"mimeType" binding:""`
	UserID  uint64 `json:"userID" binding:""`
}

// UpdateFilesByIDRequest request params
type UpdateFilesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Filename  string `json:"filename" binding:""`
	URL  string `json:"url" binding:""`
	Size  int64 `json:"size" binding:""`
	MimeType  string `json:"mimeType" binding:""`
	UserID  uint64 `json:"userID" binding:""`
}

// FilesObjDetail detail
type FilesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	Filename  string `json:"filename"`
	URL  string `json:"url"`
	Size  int64 `json:"size"`
	MimeType  string `json:"mimeType"`
	UserID  uint64 `json:"userID"`
}


// CreateFilesReply only for api docs
type CreateFilesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteFilesByIDReply only for api docs
type DeleteFilesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateFilesByIDReply only for api docs
type UpdateFilesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetFilesByIDReply only for api docs
type GetFilesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Files FilesObjDetail `json:"files"`
	} `json:"data"` // return data
}

// ListFilessRequest request params
type ListFilessRequest struct {
	query.Params
}

// ListFilessReply only for api docs
type ListFilessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Filess []FilesObjDetail `json:"filess"`
	} `json:"data"` // return data
}
