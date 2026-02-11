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

var _ RolePermissionsDao = (*rolePermissionsDao)(nil)

// RolePermissionsDao defining the dao interface
type RolePermissionsDao interface {
	Create(ctx context.Context, table *model.RolePermissions) error
	DeleteByRoleID(ctx context.Context, roleID uint64) error
	UpdateByRoleID(ctx context.Context, table *model.RolePermissions) error
	GetByRoleID(ctx context.Context, roleID uint64) (*model.RolePermissions, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.RolePermissions, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.RolePermissions) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, roleID uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.RolePermissions) error
}

type rolePermissionsDao struct {
	db    *gorm.DB
	cache cache.RolePermissionsCache // if nil, the cache is not used.
	sfg   *singleflight.Group        // if cache is nil, the sfg is not used.
}

// NewRolePermissionsDao creating the dao interface
func NewRolePermissionsDao(db *gorm.DB, xCache cache.RolePermissionsCache) RolePermissionsDao {
	if xCache == nil {
		return &rolePermissionsDao{db: db}
	}
	return &rolePermissionsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *rolePermissionsDao) deleteCache(ctx context.Context, roleID uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, roleID)
	}
	return nil
}

// Create a new rolePermissions, insert the record and the roleID value is written back to the table
func (d *rolePermissionsDao) Create(ctx context.Context, table *model.RolePermissions) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByRoleID delete a rolePermissions by roleID
func (d *rolePermissionsDao) DeleteByRoleID(ctx context.Context, roleID uint64) error {
	err := d.db.WithContext(ctx).Where("role_id = ?", roleID).Delete(&model.RolePermissions{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, roleID)

	return nil
}

// UpdateByRoleID update a rolePermissions by roleID
func (d *rolePermissionsDao) UpdateByRoleID(ctx context.Context, table *model.RolePermissions) error {
	err := d.updateDataByRoleID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.RoleID)

	return err
}

func (d *rolePermissionsDao) updateDataByRoleID(ctx context.Context, db *gorm.DB, table *model.RolePermissions) error {
	if table.RoleID < 1 {
		return errors.New("roleID cannot be 0")
	}

	update := map[string]interface{}{}

	if table.RoleID != 0 {
		update["role_id"] = table.RoleID
	}
	if table.PermissionID != 0 {
		update["permission_id"] = table.PermissionID
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByRoleID get a rolePermissions by roleID
func (d *rolePermissionsDao) GetByRoleID(ctx context.Context, roleID uint64) (*model.RolePermissions, error) {
	// no cache
	if d.cache == nil {
		record := &model.RolePermissions{}
		err := d.db.WithContext(ctx).Where("role_id = ?", roleID).First(record).Error
		return record, err
	}

	// get from cache
	record, err := d.cache.Get(ctx, roleID)
	if err == nil {
		return record, nil
	}

	// get from database
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same roleID, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(roleID), func() (interface{}, error) {

			table := &model.RolePermissions{}
			err = d.db.WithContext(ctx).Where("role_id = ?", roleID).First(table).Error
			if err != nil {
				// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, database.ErrRecordNotFound) {
					if err = d.cache.SetPlaceholder(ctx, roleID); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("roleID", roleID))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, roleID, table, cache.RolePermissionsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("roleID", roleID))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.RolePermissions)
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

// GetByColumns get a paginated list of rolePermissions by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *rolePermissionsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.RolePermissions, int64, error) {
	if params.Sort == "" {
		params.Sort = "-role_id"
	}
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.RolePermissionsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.RolePermissions{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.RolePermissions{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *rolePermissionsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.RolePermissions) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.RoleID, err
}

// DeleteByTx delete a record by roleID in the database using the provided transaction
func (d *rolePermissionsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, roleID uint64) error {
	err := tx.WithContext(ctx).Where("role_id = ?", roleID).Delete(&model.RolePermissions{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, roleID)

	return nil
}

// UpdateByTx update a record by roleID in the database using the provided transaction
func (d *rolePermissionsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.RolePermissions) error {
	err := d.updateDataByRoleID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.RoleID)

	return err
}
