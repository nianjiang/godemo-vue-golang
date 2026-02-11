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

var _ RolesDao = (*rolesDao)(nil)

// RolesDao defining the dao interface
type RolesDao interface {
	Create(ctx context.Context, table *model.Roles) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Roles) error
	GetByID(ctx context.Context, id uint64) (*model.Roles, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Roles, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Roles) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Roles) error
}

type rolesDao struct {
	db    *gorm.DB
	cache cache.RolesCache    // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewRolesDao creating the dao interface
func NewRolesDao(db *gorm.DB, xCache cache.RolesCache) RolesDao {
	if xCache == nil {
		return &rolesDao{db: db}
	}
	return &rolesDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *rolesDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new roles, insert the record and the id value is written back to the table
func (d *rolesDao) Create(ctx context.Context, table *model.Roles) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a roles by id
func (d *rolesDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Roles{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a roles by id, support partial update
func (d *rolesDao) UpdateByID(ctx context.Context, table *model.Roles) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *rolesDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Roles) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.RoleName != "" {
		update["role_name"] = table.RoleName
	}
	if table.RoleCode != "" {
		update["role_code"] = table.RoleCode
	}
	if table.RoleDesc != "" {
		update["role_desc"] = table.RoleDesc
	}
	if table.Status != "" {
		update["status"] = table.Status
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a roles by id
func (d *rolesDao) GetByID(ctx context.Context, id uint64) (*model.Roles, error) {
	// no cache
	if d.cache == nil {
		record := &model.Roles{}
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
			table := &model.Roles{}
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
			if err = d.cache.Set(ctx, id, table, cache.RolesExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Roles)
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

// GetByColumns get a paginated list of roless by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *rolesDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Roles, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.RolesColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Roles{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Roles{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *rolesDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Roles) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *rolesDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Roles{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *rolesDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Roles) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
