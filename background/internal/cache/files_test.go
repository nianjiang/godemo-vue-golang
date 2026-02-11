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

func newFilesCache() *gotest.Cache {
	record1 := &model.Files{}
	record1.ID = 1
	record2 := &model.Files{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewFilesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_filesCache_Set(t *testing.T) {
	c := newFilesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Files)
	err := c.ICache.(FilesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(FilesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_filesCache_Get(t *testing.T) {
	c := newFilesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Files)
	err := c.ICache.(FilesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(FilesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(FilesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_filesCache_MultiGet(t *testing.T) {
	c := newFilesCache()
	defer c.Close()

	var testData []*model.Files
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Files))
	}

	err := c.ICache.(FilesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(FilesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Files))
	}
}

func Test_filesCache_MultiSet(t *testing.T) {
	c := newFilesCache()
	defer c.Close()

	var testData []*model.Files
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Files))
	}

	err := c.ICache.(FilesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_filesCache_Del(t *testing.T) {
	c := newFilesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Files)
	err := c.ICache.(FilesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_filesCache_SetCacheWithNotFound(t *testing.T) {
	c := newFilesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Files)
	err := c.ICache.(FilesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(FilesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewFilesCache(t *testing.T) {
	c := NewFilesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewFilesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewFilesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
