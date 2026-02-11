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
	rolesCachePrefixKey = "roles:"
	// RolesExpireTime expire time
	RolesExpireTime = 5 * time.Minute
)

var _ RolesCache = (*rolesCache)(nil)

// RolesCache cache interface
type RolesCache interface {
	Set(ctx context.Context, id uint64, data *model.Roles, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Roles, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Roles, error)
	MultiSet(ctx context.Context, data []*model.Roles, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// rolesCache define a cache struct
type rolesCache struct {
	cache cache.Cache
}

// NewRolesCache new a cache
func NewRolesCache(cacheType *database.CacheType) RolesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Roles{}
		})
		return &rolesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Roles{}
		})
		return &rolesCache{cache: c}
	}

	return nil // no cache
}

// GetRolesCacheKey cache key
func (c *rolesCache) GetRolesCacheKey(id uint64) string {
	return rolesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *rolesCache) Set(ctx context.Context, id uint64, data *model.Roles, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetRolesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *rolesCache) Get(ctx context.Context, id uint64) (*model.Roles, error) {
	var data *model.Roles
	cacheKey := c.GetRolesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *rolesCache) MultiSet(ctx context.Context, data []*model.Roles, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetRolesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *rolesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Roles, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetRolesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Roles)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Roles)
	for _, id := range ids {
		val, ok := itemMap[c.GetRolesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *rolesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetRolesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *rolesCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetRolesCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *rolesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
