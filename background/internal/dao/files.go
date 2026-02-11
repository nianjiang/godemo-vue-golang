package dao

import (
	"context"
	"errors"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"godemo/internal/cache"
	"godemo/internal/database"
	"godemo/internal/model"
)

var _ FilesDao = (*filesDao)(nil)

// FilesDao defining the dao interface
type FilesDao interface {
	Create(ctx context.Context, table *model.Files) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Files) error
	GetByID(ctx context.Context, id uint64) (*model.Files, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Files, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Files) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Files) error
}

type filesDao struct {
	db    *gorm.DB
	cache cache.FilesCache    // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewFilesDao creating the dao interface
func NewFilesDao(db *gorm.DB, xCache cache.FilesCache) FilesDao {
	if xCache == nil {
		return &filesDao{db: db}
	}
	return &filesDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *filesDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new files, insert the record and the id value is written back to the table
func (d *filesDao) Create(ctx context.Context, table *model.Files) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a files by id
func (d *filesDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Files{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a files by id, support partial update
func (d *filesDao) UpdateByID(ctx context.Context, table *model.Files) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *filesDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Files) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Filename != "" {
		update["filename"] = table.Filename
	}
	if table.URL != "" {
		update["url"] = table.URL
	}
	if table.Size != 0 {
		update["size"] = table.Size
	}
	if table.MimeType != "" {
		update["mime_type"] = table.MimeType
	}
	if table.UserID != 0 {
		update["user_id"] = table.UserID
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a files by id
func (d *filesDao) GetByID(ctx context.Context, id uint64) (*model.Files, error) {
	// no cache
	if d.cache == nil {
		record := &model.Files{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		return record, err
	}

	// get from cache
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	// get from database
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { //nolint
			table := &model.Files{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				if errors.Is(err, database.ErrRecordNotFound) {
					// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, id, table, cache.FilesExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Files)
		if !ok {
			return nil, database.ErrRecordNotFound
		}
		return table, nil
	}

	if d.cache.IsPlaceholderErr(err) {
		return nil, database.ErrRecordNotFound
	}

	return nil, err
}

// GetByColumns get a paginated list of filess by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *filesDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Files, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.FilesColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Files{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Files{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *filesDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Files) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *filesDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Files{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *filesDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Files) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
