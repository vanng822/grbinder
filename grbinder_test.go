package grbinder

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var (
	testParamPath = "/123"
	client        = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // use default DB
	})
)

func TestLockEntity(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("PUT", testParamPath, nil)
	params := make([]gin.Param, 0)
	params = append(params, gin.Param{Key: "id", Value: "123"})
	c.Params = params

	var opts = defaultEntityLockOptions()
	opts.Name = "testing"
	opts.EnableLock = true

	lockEntityAndHandle(c, opts, func(c *gin.Context) {
		val, err := client.Get(context.Background(), "grbinder.entity_lock:testing.123").Result()
		assert.NoError(t, err)
		assert.NotNil(t, val)
	})
}

func TestLockEntityConcurrent(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("PUT", testParamPath, nil)
	params := make([]gin.Param, 0)
	params = append(params, gin.Param{Key: "id", Value: "123"})
	c.Params = params

	var opts = defaultEntityLockOptions()
	opts.Name = "testing"
	opts.EnableLock = true

	go lockEntityAndHandle(c, opts, func(c *gin.Context) {
		val, err := client.Get(context.Background(), "grbinder.entity_lock:testing.123").Result()
		assert.NoError(t, err)
		assert.NotNil(t, val)
		time.Sleep(1 * time.Second)
	})

	go lockEntityAndHandle(c, opts, func(c *gin.Context) {
		val, err := client.Get(context.Background(), "grbinder.entity_lock:testing.123").Result()
		assert.NoError(t, err)
		assert.NotNil(t, val)
		time.Sleep(1 * time.Second)
	})
}
