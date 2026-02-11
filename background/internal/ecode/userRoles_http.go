package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// userRoles business-level http error codes.
// the userRolesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	userRolesNO       = 94
	userRolesName     = "userRoles"
	userRolesBaseCode = errcode.HCode(userRolesNO)

	ErrCreateUserRoles         = errcode.NewError(userRolesBaseCode+1, "failed to create "+userRolesName)
	ErrDeleteByUserIDUserRoles = errcode.NewError(userRolesBaseCode+2, "failed to delete "+userRolesName)
	ErrUpdateByUserIDUserRoles = errcode.NewError(userRolesBaseCode+3, "failed to update "+userRolesName)
	ErrGetByUserIDUserRoles    = errcode.NewError(userRolesBaseCode+4, "failed to get "+userRolesName+" details")
	ErrListUserRoles           = errcode.NewError(userRolesBaseCode+5, "failed to list of "+userRolesName)

	// error codes are globally unique, adding 1 to the previous error code
)
