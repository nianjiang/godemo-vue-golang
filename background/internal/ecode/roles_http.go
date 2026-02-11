package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// roles business-level http error codes.
// the rolesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	rolesNO       = 47
	rolesName     = "roles"
	rolesBaseCode = errcode.HCode(rolesNO)

	ErrCreateRoles     = errcode.NewError(rolesBaseCode+1, "failed to create "+rolesName)
	ErrDeleteByIDRoles = errcode.NewError(rolesBaseCode+2, "failed to delete "+rolesName)
	ErrUpdateByIDRoles = errcode.NewError(rolesBaseCode+3, "failed to update "+rolesName)
	ErrGetByIDRoles    = errcode.NewError(rolesBaseCode+4, "failed to get "+rolesName+" details")
	ErrListRoles       = errcode.NewError(rolesBaseCode+5, "failed to list of "+rolesName)

	// error codes are globally unique, adding 1 to the previous error code
)
