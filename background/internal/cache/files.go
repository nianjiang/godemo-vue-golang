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
	filesCachePrefixKey = "files:"
	// FilesExpireTime expire time
	FilesExpireTime = 5 * time.Minute
)

var _ FilesCache = (*filesCache)(nil)

// FilesCache cache interface
type FilesCache interface {
	Set(ctx context.Context, id uint64, data *model.Files, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Files, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Files, error)
	MultiSet(ctx context.Context, data []*model.Files, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// filesCache define a cache struct
type filesCache struct {
	cache cache.Cache
}

// NewFilesCache new a cache
func NewFilesCache(cacheType *database.CacheType) FilesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Files{}
		})
		return &filesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Files{}
		})
		return &filesCache{cache: c}
	}

	return nil // no cache
}

// GetFilesCacheKey cache key
func (c *filesCache) GetFilesCacheKey(id uint64) string {
	return filesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *filesCache) Set(ctx context.Context, id uint64, data *model.Files, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetFilesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *filesCache) Get(ctx context.Context, id uint64) (*model.Files, error) {
	var data *model.Files
	cacheKey := c.GetFilesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *filesCache) MultiSet(ctx context.Context, data []*model.Files, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetFilesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *filesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Files, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetFilesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Files)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Files)
	for _, id := range ids {
		val, ok := itemMap[c.GetFilesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *filesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetFilesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *filesCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetFilesCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *filesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
