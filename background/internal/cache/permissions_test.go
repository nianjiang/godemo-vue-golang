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

func newPermissionsCache() *gotest.Cache {
	record1 := &model.Permissions{}
	record1.ID = 1
	record2 := &model.Permissions{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewPermissionsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_permissionsCache_Set(t *testing.T) {
	c := newPermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Permissions)
	err := c.ICache.(PermissionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(PermissionsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_permissionsCache_Get(t *testing.T) {
	c := newPermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Permissions)
	err := c.ICache.(PermissionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(PermissionsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(PermissionsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_permissionsCache_MultiGet(t *testing.T) {
	c := newPermissionsCache()
	defer c.Close()

	var testData []*model.Permissions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Permissions))
	}

	err := c.ICache.(PermissionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(PermissionsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Permissions))
	}
}

func Test_permissionsCache_MultiSet(t *testing.T) {
	c := newPermissionsCache()
	defer c.Close()

	var testData []*model.Permissions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Permissions))
	}

	err := c.ICache.(PermissionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_permissionsCache_Del(t *testing.T) {
	c := newPermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Permissions)
	err := c.ICache.(PermissionsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_permissionsCache_SetCacheWithNotFound(t *testing.T) {
	c := newPermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Permissions)
	err := c.ICache.(PermissionsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(PermissionsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewPermissionsCache(t *testing.T) {
	c := NewPermissionsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewPermissionsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewPermissionsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
