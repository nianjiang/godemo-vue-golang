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

var _ MenusDao = (*menusDao)(nil)

// MenusDao defining the dao interface
type MenusDao interface {
	Create(ctx context.Context, table *model.Menus) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Menus) error
	GetByID(ctx context.Context, id uint64) (*model.Menus, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Menus, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Menus) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Menus) error
}

type menusDao struct {
	db    *gorm.DB
	cache cache.MenusCache    // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewMenusDao creating the dao interface
func NewMenusDao(db *gorm.DB, xCache cache.MenusCache) MenusDao {
	if xCache == nil {
		return &menusDao{db: db}
	}
	return &menusDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *menusDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new menus, insert the record and the id value is written back to the table
func (d *menusDao) Create(ctx context.Context, table *model.Menus) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a menus by id
func (d *menusDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Menus{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a menus by id, support partial update
func (d *menusDao) UpdateByID(ctx context.Context, table *model.Menus) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *menusDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Menus) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.Path != "" {
		update["path"] = table.Path
	}
	if table.Icon != "" {
		update["icon"] = table.Icon
	}
	if table.ParentID != 0 {
		update["parent_id"] = table.ParentID
	}
	if table.Order != 0 {
		update["order"] = table.Order
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a menus by id
func (d *menusDao) GetByID(ctx context.Context, id uint64) (*model.Menus, error) {
	// no cache
	if d.cache == nil {
		record := &model.Menus{}
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
			table := &model.Menus{}
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
			if err = d.cache.Set(ctx, id, table, cache.MenusExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Menus)
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

// GetByColumns get a paginated list of menuss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *menusDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Menus, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.MenusColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Menus{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Menus{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *menusDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Menus) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *menusDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Menus{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *menusDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Menus) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
