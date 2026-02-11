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

var _ UsersHandler = (*usersHandler)(nil)

// UsersHandler defining the handler interface
type UsersHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type usersHandler struct {
	iDao dao.UsersDao
}

// NewUsersHandler creating the handler interface
func NewUsersHandler() UsersHandler {
	return &usersHandler{
		iDao: dao.NewUsersDao(
			database.GetDB(), // db driver is mysql
			cache.NewUsersCache(database.GetCacheType()),
		),
	}
}

// Create a new users
// @Summary Create a new users
// @Description Creates a new users entity using the provided data in the request body.
// @Tags users
// @Accept json
// @Produce json
// @Param data body types.CreateUsersRequest true "users information"
// @Success 200 {object} types.CreateUsersReply{}
// @Router /api/v1/users [post]
// @Security BearerAuth
func (h *usersHandler) Create(c *gin.Context) {
	form := &types.CreateUsersRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	users := &model.Users{}
	err = copier.Copy(users, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, users)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": users.ID})
}

// DeleteByID delete a users by id
// @Summary Delete a users by id
// @Description Deletes a existing users identified by the given id in the path.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteUsersByIDReply{}
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth
func (h *usersHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getUsersIDFromPath(c)
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

// UpdateByID update a users by id
// @Summary Update a users by id
// @Description Updates the specified users by given id in the path, support partial update.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateUsersByIDRequest true "users information"
// @Success 200 {object} types.UpdateUsersByIDReply{}
// @Router /api/v1/users/{id} [put]
// @Security BearerAuth
func (h *usersHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getUsersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateUsersByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	users := &model.Users{}
	err = copier.Copy(users, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, users)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a users by id
// @Summary Get a users by id
// @Description Gets detailed information of a users specified by the given id in the path.
// @Tags users
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUsersByIDReply{}
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth
func (h *usersHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getUsersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	users, err := h.iDao.GetByID(ctx, id)
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

	data := &types.UsersObjDetail{}
	err = copier.Copy(data, users)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"users": data})
}

// List get a paginated list of userss by custom conditions
// @Summary Get a paginated list of userss by custom conditions
// @Description Returns a paginated list of users based on query filters, including page number and size.
// @Tags users
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListUserssReply{}
// @Router /api/v1/users/list [post]
// @Security BearerAuth
func (h *usersHandler) List(c *gin.Context) {
	form := &types.ListUserssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUserss(userss)
	if err != nil {
		response.Error(c, ecode.ErrListUsers)
		return
	}

	response.Success(c, gin.H{
		"userss": data,
		"total":  total,
	})
}

func getUsersIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertUsers(users *model.Users) (*types.UsersObjDetail, error) {
	data := &types.UsersObjDetail{}
	err := copier.Copy(data, users)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertUserss(fromValues []*model.Users) ([]*types.UsersObjDetail, error) {
	toValues := []*types.UsersObjDetail{}
	for _, v := range fromValues {
		data, err := convertUsers(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
