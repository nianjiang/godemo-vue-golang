package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// rolePermissions business-level http error codes.
// the rolePermissionsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	rolePermissionsNO       = 69
	rolePermissionsName     = "rolePermissions"
	rolePermissionsBaseCode = errcode.HCode(rolePermissionsNO)

	ErrCreateRolePermissions         = errcode.NewError(rolePermissionsBaseCode+1, "failed to create "+rolePermissionsName)
	ErrDeleteByRoleIDRolePermissions = errcode.NewError(rolePermissionsBaseCode+2, "failed to delete "+rolePermissionsName)
	ErrUpdateByRoleIDRolePermissions = errcode.NewError(rolePermissionsBaseCode+3, "failed to update "+rolePermissionsName)
	ErrGetByRoleIDRolePermissions    = errcode.NewError(rolePermissionsBaseCode+4, "failed to get "+rolePermissionsName+" details")
	ErrListRolePermissions           = errcode.NewError(rolePermissionsBaseCode+5, "failed to list of "+rolePermissionsName)

	// error codes are globally unique, adding 1 to the previous error code
)
