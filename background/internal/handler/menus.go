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

var _ MenusHandler = (*menusHandler)(nil)

// MenusHandler defining the handler interface
type MenusHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type menusHandler struct {
	iDao dao.MenusDao
}

// NewMenusHandler creating the handler interface
func NewMenusHandler() MenusHandler {
	return &menusHandler{
		iDao: dao.NewMenusDao(
			database.GetDB(), // db driver is mysql
			cache.NewMenusCache(database.GetCacheType()),
		),
	}
}

// Create a new menus
// @Summary Create a new menus
// @Description Creates a new menus entity using the provided data in the request body.
// @Tags menus
// @Accept json
// @Produce json
// @Param data body types.CreateMenusRequest true "menus information"
// @Success 200 {object} types.CreateMenusReply{}
// @Router /api/v1/menus [post]
// @Security BearerAuth
func (h *menusHandler) Create(c *gin.Context) {
	form := &types.CreateMenusRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	menus := &model.Menus{}
	err = copier.Copy(menus, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateMenus)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, menus)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": menus.ID})
}

// DeleteByID delete a menus by id
// @Summary Delete a menus by id
// @Description Deletes a existing menus identified by the given id in the path.
// @Tags menus
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteMenusByIDReply{}
// @Router /api/v1/menus/{id} [delete]
// @Security BearerAuth
func (h *menusHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getMenusIDFromPath(c)
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

// UpdateByID update a menus by id
// @Summary Update a menus by id
// @Description Updates the specified menus by given id in the path, support partial update.
// @Tags menus
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateMenusByIDRequest true "menus information"
// @Success 200 {object} types.UpdateMenusByIDReply{}
// @Router /api/v1/menus/{id} [put]
// @Security BearerAuth
func (h *menusHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getMenusIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateMenusByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	menus := &model.Menus{}
	err = copier.Copy(menus, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDMenus)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, menus)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a menus by id
// @Summary Get a menus by id
// @Description Gets detailed information of a menus specified by the given id in the path.
// @Tags menus
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetMenusByIDReply{}
// @Router /api/v1/menus/{id} [get]
// @Security BearerAuth
func (h *menusHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getMenusIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	menus, err := h.iDao.GetByID(ctx, id)
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

	data := &types.MenusObjDetail{}
	err = copier.Copy(data, menus)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDMenus)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"menus": data})
}

// List get a paginated list of menuss by custom conditions
// @Summary Get a paginated list of menuss by custom conditions
// @Description Returns a paginated list of menus based on query filters, including page number and size.
// @Tags menus
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListMenussReply{}
// @Router /api/v1/menus/list [post]
// @Security BearerAuth
func (h *menusHandler) List(c *gin.Context) {
	form := &types.ListMenussRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	menuss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertMenuss(menuss)
	if err != nil {
		response.Error(c, ecode.ErrListMenus)
		return
	}

	response.Success(c, gin.H{
		"menuss": data,
		"total":        total,
	})
}

func getMenusIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertMenus(menus *model.Menus) (*types.MenusObjDetail, error) {
	data := &types.MenusObjDetail{}
	err := copier.Copy(data, menus)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertMenuss(fromValues []*model.Menus) ([]*types.MenusObjDetail, error) {
	toValues := []*types.MenusObjDetail{}
	for _, v := range fromValues {
		data, err := convertMenus(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
