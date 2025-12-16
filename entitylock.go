package grbinder

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vanng822/gorlock/v2"
)

var entityLock atomic.Pointer[Locker]

func InitDefaultLocker() {
	gorlock.InitDefaultRedisClient()
	var locker Locker = gorlock.NewDefault().WithSettings(&gorlock.Settings{
		KeyPrefix:     "grbinder.entity_lock",
		LockTimeout:   30 * time.Second,
		RetryTimeout:  2 * time.Second,
		RetryInterval: 15 * time.Millisecond,
		LockWaiting:   true,
	})

	entityLock.Store(&locker)
}

func SetDefaultLocker(locker Locker) {
	entityLock.Store(&locker)
}

func hasDefaultEntityLocker() bool {
	return entityLock.Load() != nil
}

func getDefaultEntityLocker() Locker {
	if locker := entityLock.Load(); locker != nil {
		return *locker
	}
	panic("Must run InitDefaultLocker or Set a default Locker first")
}

type Locker interface {
	Lock(key string) (bool, error)
	Unlock(key string) error
}

type EntityLockOptions struct {
	EnableLock     bool
	Name           string
	LockTakeAction bool
	Locker         Locker
	EntityIdLookup EntityIdLookup
}

func (o *EntityLockOptions) With(options ...Option) *EntityLockOptions {
	clone := *o
	for _, option := range options {
		option(&clone)
	}
	return &clone
}

// EntityIdLookup defines a func to lookup entity id from context
// if it returns empty string, no lock will be applied
// if the context is aborted in the func, the handler will not be called
// if abort then a status code should be set in the context
// this is not middleware, so don't call .Next()
type EntityIdLookup func(*gin.Context) string

type EntityLockSupported interface {
	EntityLockOptions() *EntityLockOptions
}

func DefaultEntityLockOptions() *EntityLockOptions {
	return &EntityLockOptions{
		EnableLock:     false,
		Name:           "",
		LockTakeAction: false,
		Locker:         getDefaultEntityLocker(),
		EntityIdLookup: func(ctx *gin.Context) string {
			return ctx.Param("id")
		},
	}
}

func nilEntityLockOptions() *EntityLockOptions {
	return &EntityLockOptions{
		EnableLock:     false,
		Name:           "",
		LockTakeAction: false,
		Locker:         nil,
		EntityIdLookup: nil,
	}
}

type Option func(*EntityLockOptions)

func WithEntityLockEnable(enable bool) Option {
	return func(options *EntityLockOptions) {
		options.EnableLock = enable
	}
}

func WithEntityLockName(name string) Option {
	return func(options *EntityLockOptions) {
		options.Name = name
	}
}

func WithEntityLockLocker(locker Locker) Option {
	return func(options *EntityLockOptions) {
		if locker == nil {
			panic("locker cannot be nil")
		}
		options.Locker = locker
	}
}

func WithEntityLockTakeAction(LockTakeAction bool) Option {
	return func(options *EntityLockOptions) {
		options.LockTakeAction = LockTakeAction
	}
}

func WithEntityLookup(entityIdLookup EntityIdLookup) Option {
	return func(options *EntityLockOptions) {
		if entityIdLookup == nil {
			panic("entityIdLockup cannot be nil")
		}
		options.EntityIdLookup = entityIdLookup
	}
}

func lockEntityAndHandle(ctx *gin.Context, options *EntityLockOptions, handler func(*gin.Context)) {
	// Lock the entity
	var id = options.EntityIdLookup(ctx)
	// support abort in EntityIdLookup
	// for example, if the id is not valid or the current user has no access to the entity
	// the EntityIdLookup can abort the context
	// so we need to check if the context is aborted here
	// and return directly
	if ctx.IsAborted() {
		return
	}
	// don't have id, don't lock
	if id == "" {
		handler(ctx)
		return
	}
	var name = options.Name
	if name == "" {
		name = ctx.FullPath()
	}
	var key = fmt.Sprintf("%s.%s", name, id)
	var locked, err = options.Locker.Lock(key)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !locked {
		ctx.AbortWithStatus(http.StatusConflict)
		return
	}

	defer options.Locker.Unlock(key)
	handler(ctx)
}
