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

var _ UserRolesHandler = (*userRolesHandler)(nil)

// UserRolesHandler defining the handler interface
type UserRolesHandler interface {
	Create(c *gin.Context)
	DeleteByUserID(c *gin.Context)
	UpdateByUserID(c *gin.Context)
	GetByUserID(c *gin.Context)
	List(c *gin.Context)
}

type userRolesHandler struct {
	iDao dao.UserRolesDao
}

// NewUserRolesHandler creating the handler interface
func NewUserRolesHandler() UserRolesHandler {
	return &userRolesHandler{
		iDao: dao.NewUserRolesDao(
			database.GetDB(), // db driver is mysql
			cache.NewUserRolesCache(database.GetCacheType()),
		),
	}
}

// Create a new userRoles
// @Summary Create a new userRoles
// @Description Creates a new userRoles entity using the provided data in the request body.
// @Tags userRoles
// @Accept json
// @Produce json
// @Param data body types.CreateUserRolesRequest true "userRoles information"
// @Success 200 {object} types.CreateUserRolesReply{}
// @Router /api/v1/userRoles [post]
// @Security BearerAuth
func (h *userRolesHandler) Create(c *gin.Context) {
	form := &types.CreateUserRolesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	userRoles := &model.UserRoles{}
	err = copier.Copy(userRoles, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUserRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, userRoles)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"userID": userRoles.UserID})
}

// DeleteByUserID delete a userRoles by userID
// @Summary Delete a userRoles by userID
// @Description Deletes a existing userRoles identified by the given userID in the path.
// @Tags userRoles
// @Accept json
// @Produce json
// @Param userID path string true "userID"
// @Success 200 {object} types.DeleteUserRolesByUserIDReply{}
// @Router /api/v1/userRoles/{userID} [delete]
// @Security BearerAuth
func (h *userRolesHandler) DeleteByUserID(c *gin.Context) {
	userID, isAbort := getUserRolesUserIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByUserID(ctx, userID)
	if err != nil {
		logger.Error("DeleteByUserID error", logger.Err(err), logger.Any("userID", userID), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByUserID update a userRoles by userID
// @Summary Update a userRoles by userID
// @Description Updates the specified userRoles by given userID in the path, support partial update.
// @Tags userRoles
// @Accept json
// @Produce json
// @Param userID path string true "userID"
// @Param data body types.UpdateUserRolesByUserIDRequest true "userRoles information"
// @Success 200 {object} types.UpdateUserRolesByUserIDReply{}
// @Router /api/v1/userRoles/{userID} [put]
// @Security BearerAuth
func (h *userRolesHandler) UpdateByUserID(c *gin.Context) {
	userID, isAbort := getUserRolesUserIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateUserRolesByUserIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.UserID = userID

	userRoles := &model.UserRoles{}
	err = copier.Copy(userRoles, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByUserIDUserRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByUserID(ctx, userRoles)
	if err != nil {
		logger.Error("UpdateByUserID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByUserID get a userRoles by userID
// @Summary Get a userRoles by userID
// @Description Gets detailed information of a userRoles specified by the given userID in the path.
// @Tags userRoles
// @Param userID path string true "userID"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserRolesByUserIDReply{}
// @Router /api/v1/userRoles/{userID} [get]
// @Security BearerAuth
func (h *userRolesHandler) GetByUserID(c *gin.Context) {
	userID, isAbort := getUserRolesUserIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userRoles, err := h.iDao.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByUserID not found", logger.Err(err), logger.Any("userID", userID), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByUserID error", logger.Err(err), logger.Any("userID", userID), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.UserRolesObjDetail{}
	err = copier.Copy(data, userRoles)
	if err != nil {
		response.Error(c, ecode.ErrGetByUserIDUserRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"userRoles": data})
}

// List get a paginated list of userRoles by custom conditions
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
// @Summary Get a paginated list of userRoles by custom conditions
// @Description Returns a paginated list of userRoles based on query filters, including page number and size.
// @Tags userRoles
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListUserRolesReply{}
// @Router /api/v1/userRoles/list [post]
// @Security BearerAuth
func (h *userRolesHandler) List(c *gin.Context) {
	form := &types.ListUserRolesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userRoles, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUserRoles(userRoles)
	if err != nil {
		response.Error(c, ecode.ErrListUserRoles)
		return
	}

	response.Success(c, gin.H{
		"userRoles": data,
		"total":     total,
	})
}

func getUserRolesUserIDFromPath(c *gin.Context) (uint64, bool) {
	userIDStr := c.Param("userID")

	userID, err := utils.StrToUint64E(userIDStr)
	if err != nil || userIDStr == "" {
		logger.Warn("StrToUint64E error: ", logger.String("userIDStr", userIDStr), middleware.GCtxRequestIDField(c))
		return 0, true
	}
	return userID, false

}

func convertUserRoles1(userRoles *model.UserRoles) (*types.UserRolesObjDetail, error) {
	data := &types.UserRolesObjDetail{}
	err := copier.Copy(data, userRoles)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertUserRoles(fromValues []*model.UserRoles) ([]*types.UserRolesObjDetail, error) {
	toValues := []*types.UserRolesObjDetail{}
	for _, v := range fromValues {
		data, err := convertUserRoles1(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
