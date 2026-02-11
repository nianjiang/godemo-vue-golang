package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"godemo/internal/cache"
	"godemo/internal/dao"
	"godemo/internal/database"
	"godemo/internal/ecode"
	"godemo/internal/model"
	"godemo/internal/types"
)

var _ RolePermissionsHandler = (*rolePermissionsHandler)(nil)

// RolePermissionsHandler defining the handler interface
type RolePermissionsHandler interface {
	Create(c *gin.Context)
	DeleteByRoleID(c *gin.Context)
	UpdateByRoleID(c *gin.Context)
	GetByRoleID(c *gin.Context)
	List(c *gin.Context)
}

type rolePermissionsHandler struct {
	iDao dao.RolePermissionsDao
}

// NewRolePermissionsHandler creating the handler interface
func NewRolePermissionsHandler() RolePermissionsHandler {
	return &rolePermissionsHandler{
		iDao: dao.NewRolePermissionsDao(
			database.GetDB(), // db driver is mysql
			cache.NewRolePermissionsCache(database.GetCacheType()),
		),
	}
}

// Create a new rolePermissions
// @Summary Create a new rolePermissions
// @Description Creates a new rolePermissions entity using the provided data in the request body.
// @Tags rolePermissions
// @Accept json
// @Produce json
// @Param data body types.CreateRolePermissionsRequest true "rolePermissions information"
// @Success 200 {object} types.CreateRolePermissionsReply{}
// @Router /api/v1/rolePermissions [post]
// @Security BearerAuth
func (h *rolePermissionsHandler) Create(c *gin.Context) {
	form := &types.CreateRolePermissionsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	rolePermissions := &model.RolePermissions{}
	err = copier.Copy(rolePermissions, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateRolePermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, rolePermissions)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"roleID": rolePermissions.RoleID})
}

// DeleteByRoleID delete a rolePermissions by roleID
// @Summary Delete a rolePermissions by roleID
// @Description Deletes a existing rolePermissions identified by the given roleID in the path.
// @Tags rolePermissions
// @Accept json
// @Produce json
// @Param roleID path string true "roleID"
// @Success 200 {object} types.DeleteRolePermissionsByRoleIDReply{}
// @Router /api/v1/rolePermissions/{roleID} [delete]
// @Security BearerAuth
func (h *rolePermissionsHandler) DeleteByRoleID(c *gin.Context) {
	roleID, isAbort := getRolePermissionsRoleIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByRoleID(ctx, roleID)
	if err != nil {
		logger.Error("DeleteByRoleID error", logger.Err(err), logger.Any("roleID", roleID), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByRoleID update a rolePermissions by roleID
// @Summary Update a rolePermissions by roleID
// @Description Updates the specified rolePermissions by given roleID in the path, support partial update.
// @Tags rolePermissions
// @Accept json
// @Produce json
// @Param roleID path string true "roleID"
// @Param data body types.UpdateRolePermissionsByRoleIDRequest true "rolePermissions information"
// @Success 200 {object} types.UpdateRolePermissionsByRoleIDReply{}
// @Router /api/v1/rolePermissions/{roleID} [put]
// @Security BearerAuth
func (h *rolePermissionsHandler) UpdateByRoleID(c *gin.Context) {
	roleID, isAbort := getRolePermissionsRoleIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRolePermissionsByRoleIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.RoleID = roleID

	rolePermissions := &model.RolePermissions{}
	err = copier.Copy(rolePermissions, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByRoleIDRolePermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByRoleID(ctx, rolePermissions)
	if err != nil {
		logger.Error("UpdateByRoleID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByRoleID get a rolePermissions by roleID
// @Summary Get a rolePermissions by roleID
// @Description Gets detailed information of a rolePermissions specified by the given roleID in the path.
// @Tags rolePermissions
// @Param roleID path string true "roleID"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRolePermissionsByRoleIDReply{}
// @Router /api/v1/rolePermissions/{roleID} [get]
// @Security BearerAuth
func (h *rolePermissionsHandler) GetByRoleID(c *gin.Context) {
	roleID, isAbort := getRolePermissionsRoleIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	rolePermissions, err := h.iDao.GetByRoleID(ctx, roleID)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByRoleID not found", logger.Err(err), logger.Any("roleID", roleID), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByRoleID error", logger.Err(err), logger.Any("roleID", roleID), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.RolePermissionsObjDetail{}
	err = copier.Copy(data, rolePermissions)
	if err != nil {
		response.Error(c, ecode.ErrGetByRoleIDRolePermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"rolePermissions": data})
}

// List get a paginated list of rolePermissions by custom conditions
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
// @Summary Get a paginated list of rolePermissions by custom conditions
// @Description Returns a paginated list of rolePermissions based on query filters, including page number and size.
// @Tags rolePermissions
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListRolePermissionsReply{}
// @Router /api/v1/rolePermissions/list [post]
// @Security BearerAuth
func (h *rolePermissionsHandler) List(c *gin.Context) {
	form := &types.ListRolePermissionsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	rolePermissions, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRolePermissions(rolePermissions)
	if err != nil {
		response.Error(c, ecode.ErrListRolePermissions)
		return
	}

	response.Success(c, gin.H{
		"rolePermissions": data,
		"total":           total,
	})
}

func getRolePermissionsRoleIDFromPath(c *gin.Context) (uint64, bool) {
	roleIDStr := c.Param("roleID")

	roleID, err := utils.StrToUint64E(roleIDStr)
	if err != nil || roleIDStr == "" {
		logger.Warn("StrToUint64E error: ", logger.String("roleIDStr", roleIDStr), middleware.GCtxRequestIDField(c))
		return 0, true
	}
	return roleID, false

}

func convertRolePermissions1(rolePermissions *model.RolePermissions) (*types.RolePermissionsObjDetail, error) {
	data := &types.RolePermissionsObjDetail{}
	err := copier.Copy(data, rolePermissions)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertRolePermissions(fromValues []*model.RolePermissions) ([]*types.RolePermissionsObjDetail, error) {
	toValues := []*types.RolePermissionsObjDetail{}
	for _, v := range fromValues {
		data, err := convertRolePermissions1(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
