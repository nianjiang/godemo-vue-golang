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
	permissionsCachePrefixKey = "permissions:"
	// PermissionsExpireTime expire time
	PermissionsExpireTime = 5 * time.Minute
)

var _ PermissionsCache = (*permissionsCache)(nil)

// PermissionsCache cache interface
type PermissionsCache interface {
	Set(ctx context.Context, id uint64, data *model.Permissions, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Permissions, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Permissions, error)
	MultiSet(ctx context.Context, data []*model.Permissions, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// permissionsCache define a cache struct
type permissionsCache struct {
	cache cache.Cache
}

// NewPermissionsCache new a cache
func NewPermissionsCache(cacheType *database.CacheType) PermissionsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Permissions{}
		})
		return &permissionsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Permissions{}
		})
		return &permissionsCache{cache: c}
	}

	return nil // no cache
}

// GetPermissionsCacheKey cache key
func (c *permissionsCache) GetPermissionsCacheKey(id uint64) string {
	return permissionsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *permissionsCache) Set(ctx context.Context, id uint64, data *model.Permissions, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetPermissionsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *permissionsCache) Get(ctx context.Context, id uint64) (*model.Permissions, error) {
	var data *model.Permissions
	cacheKey := c.GetPermissionsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *permissionsCache) MultiSet(ctx context.Context, data []*model.Permissions, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetPermissionsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *permissionsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Permissions, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetPermissionsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Permissions)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Permissions)
	for _, id := range ids {
		val, ok := itemMap[c.GetPermissionsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *permissionsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetPermissionsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *permissionsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetPermissionsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *permissionsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
