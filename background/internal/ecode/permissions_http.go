package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// permissions business-level http error codes.
// the permissionsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	permissionsNO = 68
	permissionsName     = "permissions"
	permissionsBaseCode = errcode.HCode(permissionsNO)

	ErrCreatePermissions     = errcode.NewError(permissionsBaseCode+1, "failed to create "+permissionsName)
	ErrDeleteByIDPermissions = errcode.NewError(permissionsBaseCode+2, "failed to delete "+permissionsName)
	ErrUpdateByIDPermissions = errcode.NewError(permissionsBaseCode+3, "failed to update "+permissionsName)
	ErrGetByIDPermissions    = errcode.NewError(permissionsBaseCode+4, "failed to get "+permissionsName+" details")
	ErrListPermissions       = errcode.NewError(permissionsBaseCode+5, "failed to list of "+permissionsName)

	// error codes are globally unique, adding 1 to the previous error code
)
