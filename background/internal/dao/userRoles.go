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

var _ UserRolesDao = (*userRolesDao)(nil)

// UserRolesDao defining the dao interface
type UserRolesDao interface {
	Create(ctx context.Context, table *model.UserRoles) error
	DeleteByUserID(ctx context.Context, userID uint64) error
	UpdateByUserID(ctx context.Context, table *model.UserRoles) error
	GetByUserID(ctx context.Context, userID uint64) (*model.UserRoles, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.UserRoles, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.UserRoles) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, userID uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.UserRoles) error
}

type userRolesDao struct {
	db    *gorm.DB
	cache cache.UserRolesCache // if nil, the cache is not used.
	sfg   *singleflight.Group    // if cache is nil, the sfg is not used.
}

// NewUserRolesDao creating the dao interface
func NewUserRolesDao(db *gorm.DB, xCache cache.UserRolesCache) UserRolesDao {
	if xCache == nil {
		return &userRolesDao{db: db}
	}
	return &userRolesDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *userRolesDao) deleteCache(ctx context.Context, userID uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, userID)
	}
	return nil
}

// Create a new userRoles, insert the record and the userID value is written back to the table
func (d *userRolesDao) Create(ctx context.Context, table *model.UserRoles) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByUserID delete a userRoles by userID
func (d *userRolesDao) DeleteByUserID(ctx context.Context, userID uint64) error {
	err := d.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.UserRoles{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, userID)

	return nil
}

// UpdateByUserID update a userRoles by userID
func (d *userRolesDao) UpdateByUserID(ctx context.Context, table *model.UserRoles) error {
	err := d.updateDataByUserID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.UserID)

	return err
}

func (d *userRolesDao) updateDataByUserID(ctx context.Context, db *gorm.DB, table *model.UserRoles) error {
		if table.UserID < 1 {
		return errors.New("userID cannot be 0")
	}


	update := map[string]interface{}{}
	
	if table.UserID != 0 {
		update["user_id"] = table.UserID
	}
	if table.RoleID != 0 {
		update["role_id"] = table.RoleID
	}
	

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByUserID get a userRoles by userID
func (d *userRolesDao) GetByUserID(ctx context.Context, userID uint64) (*model.UserRoles, error) {
	// no cache
	if d.cache == nil {
		record := &model.UserRoles{}
		err := d.db.WithContext(ctx).Where("user_id = ?", userID).First(record).Error
		return record, err
	}

	// get from cache
	record, err := d.cache.Get(ctx, userID)
	if err == nil {
		return record, nil
	}

	// get from database
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same userID, prevent high concurrent simultaneous access to database
				val, err, _ := d.sfg.Do(utils.Uint64ToStr(userID), func() (interface{}, error) {

			table := &model.UserRoles{}
			err = d.db.WithContext(ctx).Where("user_id = ?", userID).First(table).Error
			if err != nil {
				// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, database.ErrRecordNotFound) {
					if err = d.cache.SetPlaceholder(ctx, userID); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("userID", userID))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, userID, table, cache.UserRolesExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("userID", userID))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.UserRoles)
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

// GetByColumns get a paginated list of userRoles by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *userRolesDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.UserRoles, int64, error) {
	if params.Sort == "" {
		params.Sort = "-user_id"
	}
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.UserRolesColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.UserRoles{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.UserRoles{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *userRolesDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.UserRoles) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.UserID, err
}

// DeleteByTx delete a record by userID in the database using the provided transaction
func (d *userRolesDao) DeleteByTx(ctx context.Context, tx *gorm.DB, userID uint64) error {
	err := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.UserRoles{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, userID)

	return nil
}

// UpdateByTx update a record by userID in the database using the provided transaction
func (d *userRolesDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.UserRoles) error {
	err := d.updateDataByUserID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.UserID)

	return err
}
