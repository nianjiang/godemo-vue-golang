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

func newRolesCache() *gotest.Cache {
	record1 := &model.Roles{}
	record1.ID = 1
	record2 := &model.Roles{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewRolesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_rolesCache_Set(t *testing.T) {
	c := newRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Roles)
	err := c.ICache.(RolesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(RolesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_rolesCache_Get(t *testing.T) {
	c := newRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Roles)
	err := c.ICache.(RolesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RolesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(RolesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_rolesCache_MultiGet(t *testing.T) {
	c := newRolesCache()
	defer c.Close()

	var testData []*model.Roles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Roles))
	}

	err := c.ICache.(RolesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RolesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Roles))
	}
}

func Test_rolesCache_MultiSet(t *testing.T) {
	c := newRolesCache()
	defer c.Close()

	var testData []*model.Roles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Roles))
	}

	err := c.ICache.(RolesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_rolesCache_Del(t *testing.T) {
	c := newRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Roles)
	err := c.ICache.(RolesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_rolesCache_SetCacheWithNotFound(t *testing.T) {
	c := newRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Roles)
	err := c.ICache.(RolesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(RolesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewRolesCache(t *testing.T) {
	c := NewRolesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewRolesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewRolesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
