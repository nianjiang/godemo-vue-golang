package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/gotest"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"godemo/internal/database"
	"godemo/internal/model"
)

func newMenusCache() *gotest.Cache {
	record1 := &model.Menus{}
	record1.ID = 1
	record2 := &model.Menus{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewMenusCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_menusCache_Set(t *testing.T) {
	c := newMenusCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Menus)
	err := c.ICache.(MenusCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(MenusCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_menusCache_Get(t *testing.T) {
	c := newMenusCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Menus)
	err := c.ICache.(MenusCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(MenusCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(MenusCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_menusCache_MultiGet(t *testing.T) {
	c := newMenusCache()
	defer c.Close()

	var testData []*model.Menus
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Menus))
	}

	err := c.ICache.(MenusCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(MenusCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Menus))
	}
}

func Test_menusCache_MultiSet(t *testing.T) {
	c := newMenusCache()
	defer c.Close()

	var testData []*model.Menus
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Menus))
	}

	err := c.ICache.(MenusCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_menusCache_Del(t *testing.T) {
	c := newMenusCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Menus)
	err := c.ICache.(MenusCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_menusCache_SetCacheWithNotFound(t *testing.T) {
	c := newMenusCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Menus)
	err := c.ICache.(MenusCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(MenusCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewMenusCache(t *testing.T) {
	c := NewMenusCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewMenusCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewMenusCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
