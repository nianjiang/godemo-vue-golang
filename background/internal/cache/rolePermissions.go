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
	rolePermissionsCachePrefixKey = "rolePermissions:"
	// RolePermissionsExpireTime expire time
	RolePermissionsExpireTime = 5 * time.Minute
)

var _ RolePermissionsCache = (*rolePermissionsCache)(nil)

// RolePermissionsCache cache interface
type RolePermissionsCache interface {
	Set(ctx context.Context, roleID uint64, data *model.RolePermissions, duration time.Duration) error
	Get(ctx context.Context, roleID uint64) (*model.RolePermissions, error)
	MultiGet(ctx context.Context, roleIDs []uint64) (map[uint64]*model.RolePermissions, error)
	MultiSet(ctx context.Context, data []*model.RolePermissions, duration time.Duration) error
	Del(ctx context.Context, roleID uint64) error
	SetPlaceholder(ctx context.Context, roleID uint64) error
	IsPlaceholderErr(err error) bool
}

// rolePermissionsCache define a cache struct
type rolePermissionsCache struct {
	cache cache.Cache
}

// NewRolePermissionsCache new a cache
func NewRolePermissionsCache(cacheType *database.CacheType) RolePermissionsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.RolePermissions{}
		})
		return &rolePermissionsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.RolePermissions{}
		})
		return &rolePermissionsCache{cache: c}
	}

	return nil // no cache
}

// GetRolePermissionsCacheKey cache key
func (c *rolePermissionsCache) GetRolePermissionsCacheKey(roleID uint64) string {
	return rolePermissionsCachePrefixKey + utils.Uint64ToStr(roleID)
}

// Set write to cache
func (c *rolePermissionsCache) Set(ctx context.Context, roleID uint64, data *model.RolePermissions, duration time.Duration) error {
	if data == nil {
		return nil
	}
	cacheKey := c.GetRolePermissionsCacheKey(roleID)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *rolePermissionsCache) Get(ctx context.Context, roleID uint64) (*model.RolePermissions, error) {
	var data *model.RolePermissions
	cacheKey := c.GetRolePermissionsCacheKey(roleID)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *rolePermissionsCache) MultiSet(ctx context.Context, data []*model.RolePermissions, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetRolePermissionsCacheKey(v.RoleID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is roleID value
func (c *rolePermissionsCache) MultiGet(ctx context.Context, roleIDs []uint64) (map[uint64]*model.RolePermissions, error) {
	var keys []string
	for _, v := range roleIDs {
		cacheKey := c.GetRolePermissionsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.RolePermissions)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.RolePermissions)
	for _, roleID := range roleIDs {
		val, ok := itemMap[c.GetRolePermissionsCacheKey(roleID)]
		if ok {
			retMap[roleID] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *rolePermissionsCache) Del(ctx context.Context, roleID uint64) error {
	cacheKey := c.GetRolePermissionsCacheKey(roleID)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *rolePermissionsCache) SetPlaceholder(ctx context.Context, roleID uint64) error {
	cacheKey := c.GetRolePermissionsCacheKey(roleID)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *rolePermissionsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
