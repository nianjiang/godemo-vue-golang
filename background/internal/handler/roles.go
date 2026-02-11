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

var _ RolesHandler = (*rolesHandler)(nil)

// RolesHandler defining the handler interface
type RolesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type rolesHandler struct {
	iDao dao.RolesDao
}

// NewRolesHandler creating the handler interface
func NewRolesHandler() RolesHandler {
	return &rolesHandler{
		iDao: dao.NewRolesDao(
			database.GetDB(), // db driver is mysql
			cache.NewRolesCache(database.GetCacheType()),
		),
	}
}

// Create a new roles
// @Summary Create a new roles
// @Description Creates a new roles entity using the provided data in the request body.
// @Tags roles
// @Accept json
// @Produce json
// @Param data body types.CreateRolesRequest true "roles information"
// @Success 200 {object} types.CreateRolesReply{}
// @Router /api/v1/roles [post]
// @Security BearerAuth
func (h *rolesHandler) Create(c *gin.Context) {
	form := &types.CreateRolesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	roles := &model.Roles{}
	err = copier.Copy(roles, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, roles)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": roles.ID})
}

// DeleteByID delete a roles by id
// @Summary Delete a roles by id
// @Description Deletes a existing roles identified by the given id in the path.
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteRolesByIDReply{}
// @Router /api/v1/roles/{id} [delete]
// @Security BearerAuth
func (h *rolesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getRolesIDFromPath(c)
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

// UpdateByID update a roles by id
// @Summary Update a roles by id
// @Description Updates the specified roles by given id in the path, support partial update.
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateRolesByIDRequest true "roles information"
// @Success 200 {object} types.UpdateRolesByIDReply{}
// @Router /api/v1/roles/{id} [put]
// @Security BearerAuth
func (h *rolesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getRolesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRolesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	roles := &model.Roles{}
	err = copier.Copy(roles, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, roles)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a roles by id
// @Summary Get a roles by id
// @Description Gets detailed information of a roles specified by the given id in the path.
// @Tags roles
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRolesByIDReply{}
// @Router /api/v1/roles/{id} [get]
// @Security BearerAuth
func (h *rolesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getRolesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roles, err := h.iDao.GetByID(ctx, id)
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

	data := &types.RolesObjDetail{}
	err = copier.Copy(data, roles)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"roles": data})
}

// List get a paginated list of roless by custom conditions
// @Summary Get a paginated list of roless by custom conditions
// @Description Returns a paginated list of roles based on query filters, including page number and size.
// @Tags roles
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListRolessReply{}
// @Router /api/v1/roles/list [post]
// @Security BearerAuth
func (h *rolesHandler) List(c *gin.Context) {
	form := &types.ListRolessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roless, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRoless(roless)
	if err != nil {
		response.Error(c, ecode.ErrListRoles)
		return
	}

	response.Success(c, gin.H{
		"roless": data,
		"total":  total,
	})
}

func getRolesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertRoles(roles *model.Roles) (*types.RolesObjDetail, error) {
	data := &types.RolesObjDetail{}
	err := copier.Copy(data, roles)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertRoless(fromValues []*model.Roles) ([]*types.RolesObjDetail, error) {
	toValues := []*types.RolesObjDetail{}
	for _, v := range fromValues {
		data, err := convertRoles(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
