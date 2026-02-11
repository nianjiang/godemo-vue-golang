package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"godemo/internal/database"
	"godemo/internal/model"
)

const (
	// cache prefix key, must end with a colon
	userRolesCachePrefixKey = "userRoles:"
	// UserRolesExpireTime expire time
	UserRolesExpireTime = 5 * time.Minute
)

var _ UserRolesCache = (*userRolesCache)(nil)

// UserRolesCache cache interface
type UserRolesCache interface {
	Set(ctx context.Context, userID uint64, data *model.UserRoles, duration time.Duration) error
	Get(ctx context.Context, userID uint64) (*model.UserRoles, error)
	MultiGet(ctx context.Context, userIDs []uint64) (map[uint64]*model.UserRoles, error)
	MultiSet(ctx context.Context, data []*model.UserRoles, duration time.Duration) error
	Del(ctx context.Context, userID uint64) error
	SetPlaceholder(ctx context.Context, userID uint64) error
	IsPlaceholderErr(err error) bool
}

// userRolesCache define a cache struct
type userRolesCache struct {
	cache cache.Cache
}

// NewUserRolesCache new a cache
func NewUserRolesCache(cacheType *database.CacheType) UserRolesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserRoles{}
		})
		return &userRolesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserRoles{}
		})
		return &userRolesCache{cache: c}
	}

	return nil // no cache
}

// GetUserRolesCacheKey cache key
func (c *userRolesCache) GetUserRolesCacheKey(userID uint64) string {
	return userRolesCachePrefixKey + utils.Uint64ToStr(userID)
}

// Set write to cache
func (c *userRolesCache) Set(ctx context.Context, userID uint64, data *model.UserRoles, duration time.Duration) error {
	if data == nil {
		return nil
	}
	cacheKey := c.GetUserRolesCacheKey(userID)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *userRolesCache) Get(ctx context.Context, userID uint64) (*model.UserRoles, error) {
	var data *model.UserRoles
	cacheKey := c.GetUserRolesCacheKey(userID)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *userRolesCache) MultiSet(ctx context.Context, data []*model.UserRoles, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserRolesCacheKey(v.UserID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is userID value
func (c *userRolesCache) MultiGet(ctx context.Context, userIDs []uint64) (map[uint64]*model.UserRoles, error) {
	var keys []string
	for _, v := range userIDs {
		cacheKey := c.GetUserRolesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.UserRoles)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.UserRoles)
	for _, userID := range userIDs {
		val, ok := itemMap[c.GetUserRolesCacheKey(userID)]
		if ok {
			retMap[userID] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *userRolesCache) Del(ctx context.Context, userID uint64) error {
	cacheKey := c.GetUserRolesCacheKey(userID)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *userRolesCache) SetPlaceholder(ctx context.Context, userID uint64) error {
	cacheKey := c.GetUserRolesCacheKey(userID)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *userRolesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
