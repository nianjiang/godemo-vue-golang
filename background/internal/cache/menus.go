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
	menusCachePrefixKey = "menus:"
	// MenusExpireTime expire time
	MenusExpireTime = 5 * time.Minute
)

var _ MenusCache = (*menusCache)(nil)

// MenusCache cache interface
type MenusCache interface {
	Set(ctx context.Context, id uint64, data *model.Menus, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Menus, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Menus, error)
	MultiSet(ctx context.Context, data []*model.Menus, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// menusCache define a cache struct
type menusCache struct {
	cache cache.Cache
}

// NewMenusCache new a cache
func NewMenusCache(cacheType *database.CacheType) MenusCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Menus{}
		})
		return &menusCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Menus{}
		})
		return &menusCache{cache: c}
	}

	return nil // no cache
}

// GetMenusCacheKey cache key
func (c *menusCache) GetMenusCacheKey(id uint64) string {
	return menusCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *menusCache) Set(ctx context.Context, id uint64, data *model.Menus, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetMenusCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *menusCache) Get(ctx context.Context, id uint64) (*model.Menus, error) {
	var data *model.Menus
	cacheKey := c.GetMenusCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *menusCache) MultiSet(ctx context.Context, data []*model.Menus, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetMenusCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *menusCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Menus, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetMenusCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Menus)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Menus)
	for _, id := range ids {
		val, ok := itemMap[c.GetMenusCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *menusCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetMenusCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *menusCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetMenusCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *menusCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
