package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// menus business-level http error codes.
// the menusNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	menusNO       = 26
	menusName     = "menus"
	menusBaseCode = errcode.HCode(menusNO)

	ErrCreateMenus     = errcode.NewError(menusBaseCode+1, "failed to create "+menusName)
	ErrDeleteByIDMenus = errcode.NewError(menusBaseCode+2, "failed to delete "+menusName)
	ErrUpdateByIDMenus = errcode.NewError(menusBaseCode+3, "failed to update "+menusName)
	ErrGetByIDMenus    = errcode.NewError(menusBaseCode+4, "failed to get "+menusName+" details")
	ErrListMenus       = errcode.NewError(menusBaseCode+5, "failed to list of "+menusName)

	// error codes are globally unique, adding 1 to the previous error code
)
