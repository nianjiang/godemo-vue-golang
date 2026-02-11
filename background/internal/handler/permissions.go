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

var _ PermissionsHandler = (*permissionsHandler)(nil)

// PermissionsHandler defining the handler interface
type PermissionsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type permissionsHandler struct {
	iDao dao.PermissionsDao
}

// NewPermissionsHandler creating the handler interface
func NewPermissionsHandler() PermissionsHandler {
	return &permissionsHandler{
		iDao: dao.NewPermissionsDao(
			database.GetDB(), // db driver is mysql
			cache.NewPermissionsCache(database.GetCacheType()),
		),
	}
}

// Create a new permissions
// @Summary Create a new permissions
// @Description Creates a new permissions entity using the provided data in the request body.
// @Tags permissions
// @Accept json
// @Produce json
// @Param data body types.CreatePermissionsRequest true "permissions information"
// @Success 200 {object} types.CreatePermissionsReply{}
// @Router /api/v1/permissions [post]
// @Security BearerAuth
func (h *permissionsHandler) Create(c *gin.Context) {
	form := &types.CreatePermissionsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	permissions := &model.Permissions{}
	err = copier.Copy(permissions, form)
	if err != nil {
		response.Error(c, ecode.ErrCreatePermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, permissions)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": permissions.ID})
}

// DeleteByID delete a permissions by id
// @Summary Delete a permissions by id
// @Description Deletes a existing permissions identified by the given id in the path.
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeletePermissionsByIDReply{}
// @Router /api/v1/permissions/{id} [delete]
// @Security BearerAuth
func (h *permissionsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getPermissionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update a permissions by id
// @Summary Update a permissions by id
// @Description Updates the specified permissions by given id in the path, support partial update.
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdatePermissionsByIDRequest true "permissions information"
// @Success 200 {object} types.UpdatePermissionsByIDReply{}
// @Router /api/v1/permissions/{id} [put]
// @Security BearerAuth
func (h *permissionsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getPermissionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdatePermissionsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	permissions := &model.Permissions{}
	err = copier.Copy(permissions, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDPermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, permissions)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a permissions by id
// @Summary Get a permissions by id
// @Description Gets detailed information of a permissions specified by the given id in the path.
// @Tags permissions
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetPermissionsByIDReply{}
// @Router /api/v1/permissions/{id} [get]
// @Security BearerAuth
func (h *permissionsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getPermissionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	permissions, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.PermissionsObjDetail{}
	err = copier.Copy(data, permissions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDPermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"permissions": data})
}

// List get a paginated list of permissionss by custom conditions
// @Summary Get a paginated list of permissionss by custom conditions
// @Description Returns a paginated list of permissions based on query filters, including page number and size.
// @Tags permissions
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListPermissionssReply{}
// @Router /api/v1/permissions/list [post]
// @Security BearerAuth
func (h *permissionsHandler) List(c *gin.Context) {
	form := &types.ListPermissionssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	permissionss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertPermissionss(permissionss)
	if err != nil {
		response.Error(c, ecode.ErrListPermissions)
		return
	}

	response.Success(c, gin.H{
		"permissionss": data,
		"total":        total,
	})
}

func getPermissionsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertPermissions(permissions *model.Permissions) (*types.PermissionsObjDetail, error) {
	data := &types.PermissionsObjDetail{}
	err := copier.Copy(data, permissions)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertPermissionss(fromValues []*model.Permissions) ([]*types.PermissionsObjDetail, error) {
	toValues := []*types.PermissionsObjDetail{}
	for _, v := range fromValues {
		data, err := convertPermissions(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
