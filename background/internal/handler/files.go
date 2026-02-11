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

var _ FilesHandler = (*filesHandler)(nil)

// FilesHandler defining the handler interface
type FilesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type filesHandler struct {
	iDao dao.FilesDao
}

// NewFilesHandler creating the handler interface
func NewFilesHandler() FilesHandler {
	return &filesHandler{
		iDao: dao.NewFilesDao(
			database.GetDB(), // db driver is mysql
			cache.NewFilesCache(database.GetCacheType()),
		),
	}
}

// Create a new files
// @Summary Create a new files
// @Description Creates a new files entity using the provided data in the request body.
// @Tags files
// @Accept json
// @Produce json
// @Param data body types.CreateFilesRequest true "files information"
// @Success 200 {object} types.CreateFilesReply{}
// @Router /api/v1/files [post]
// @Security BearerAuth
func (h *filesHandler) Create(c *gin.Context) {
	form := &types.CreateFilesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	files := &model.Files{}
	err = copier.Copy(files, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateFiles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, files)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": files.ID})
}

// DeleteByID delete a files by id
// @Summary Delete a files by id
// @Description Deletes a existing files identified by the given id in the path.
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteFilesByIDReply{}
// @Router /api/v1/files/{id} [delete]
// @Security BearerAuth
func (h *filesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getFilesIDFromPath(c)
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

// UpdateByID update a files by id
// @Summary Update a files by id
// @Description Updates the specified files by given id in the path, support partial update.
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateFilesByIDRequest true "files information"
// @Success 200 {object} types.UpdateFilesByIDReply{}
// @Router /api/v1/files/{id} [put]
// @Security BearerAuth
func (h *filesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getFilesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateFilesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	files := &model.Files{}
	err = copier.Copy(files, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDFiles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, files)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a files by id
// @Summary Get a files by id
// @Description Gets detailed information of a files specified by the given id in the path.
// @Tags files
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetFilesByIDReply{}
// @Router /api/v1/files/{id} [get]
// @Security BearerAuth
func (h *filesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getFilesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	files, err := h.iDao.GetByID(ctx, id)
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

	data := &types.FilesObjDetail{}
	err = copier.Copy(data, files)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDFiles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"files": data})
}

// List get a paginated list of filess by custom conditions
// @Summary Get a paginated list of filess by custom conditions
// @Description Returns a paginated list of files based on query filters, including page number and size.
// @Tags files
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListFilessReply{}
// @Router /api/v1/files/list [post]
// @Security BearerAuth
func (h *filesHandler) List(c *gin.Context) {
	form := &types.ListFilessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	filess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertFiless(filess)
	if err != nil {
		response.Error(c, ecode.ErrListFiles)
		return
	}

	response.Success(c, gin.H{
		"filess": data,
		"total":  total,
	})
}

func getFilesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertFiles(files *model.Files) (*types.FilesObjDetail, error) {
	data := &types.FilesObjDetail{}
	err := copier.Copy(data, files)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertFiless(fromValues []*model.Files) ([]*types.FilesObjDetail, error) {
	toValues := []*types.FilesObjDetail{}
	for _, v := range fromValues {
		data, err := convertFiles(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
