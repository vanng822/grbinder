package grbinder

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
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

func TestLockEntityConcurrentSuccess(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("PUT", testParamPath, nil)
	params := make([]gin.Param, 0)
	params = append(params, gin.Param{Key: "id", Value: "123"})
	c.Params = params

	var opts = defaultEntityLockOptions()
	opts.Name = "testing"
	opts.EnableLock = true

	var counter atomic.Int32
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		lockEntityAndHandle(c, opts, func(c *gin.Context) {
			val, err := client.Get(context.Background(), "grbinder.entity_lock:testing.123").Result()
			assert.NoError(t, err)
			assert.NotNil(t, val)
			time.Sleep(1 * time.Second)
			counter.Add(1)
		})
	}()

	go func() {
		defer wg.Done()

		lockEntityAndHandle(c, opts, func(c *gin.Context) {
			val, err := client.Get(context.Background(), "grbinder.entity_lock:testing.123").Result()
			assert.NoError(t, err)
			assert.NotNil(t, val)
			time.Sleep(1 * time.Second)
			counter.Add(1)
		})
	}()

	wg.Wait()

	assert.Equal(t, int32(2), counter.Load())
}

func TestLockEntityConcurrentRejected(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("PUT", testParamPath, nil)
	params := make([]gin.Param, 0)
	params = append(params, gin.Param{Key: "id", Value: "123"})
	c.Params = params

	var opts = defaultEntityLockOptions()
	opts.Name = "testing"
	opts.EnableLock = true

	var counter atomic.Int32
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		lockEntityAndHandle(c, opts, func(c *gin.Context) {
			val, err := client.Get(context.Background(), "grbinder.entity_lock:testing.123").Result()
			assert.NoError(t, err)
			assert.NotNil(t, val)
			time.Sleep(2200 * time.Millisecond)
			counter.Add(1)
		})
	}()

	go func() {
		defer wg.Done()

		lockEntityAndHandle(c, opts, func(c *gin.Context) {
			val, err := client.Get(context.Background(), "grbinder.entity_lock:testing.123").Result()
			assert.NoError(t, err)
			assert.NotNil(t, val)
			time.Sleep(2200 * time.Millisecond)
			counter.Add(1)
		})
	}()

	wg.Wait()

	assert.Equal(t, int32(1), counter.Load())
}

func TestLockEntityIdLockup(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("PUT", testParamPath, nil)
	params := make([]gin.Param, 0)
	params = append(params, gin.Param{Key: "id", Value: "123"}, gin.Param{Key: "custom", Value: "456"})
	c.Params = params

	var opts = defaultEntityLockOptions()
	opts.Name = "testing"
	opts.EnableLock = true
	opts.EntityIdLookup = func(ctx *gin.Context) string {
		return ctx.Param("custom")
	}

	lockEntityAndHandle(c, opts, func(c *gin.Context) {
		val, err := client.Get(context.Background(), "grbinder.entity_lock:testing.456").Result()
		assert.NoError(t, err)
		assert.NotNil(t, val)
	})
}
