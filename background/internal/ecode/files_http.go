package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// files business-level http error codes.
// the filesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	filesNO = 16
	filesName     = "files"
	filesBaseCode = errcode.HCode(filesNO)

	ErrCreateFiles     = errcode.NewError(filesBaseCode+1, "failed to create "+filesName)
	ErrDeleteByIDFiles = errcode.NewError(filesBaseCode+2, "failed to delete "+filesName)
	ErrUpdateByIDFiles = errcode.NewError(filesBaseCode+3, "failed to update "+filesName)
	ErrGetByIDFiles    = errcode.NewError(filesBaseCode+4, "failed to get "+filesName+" details")
	ErrListFiles       = errcode.NewError(filesBaseCode+5, "failed to list of "+filesName)

	// error codes are globally unique, adding 1 to the previous error code
)
