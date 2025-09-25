package grbinder

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vanng822/gorlock/v2"
)

var entityLock gorlock.Gorlock

func init() {
	entityLock = gorlock.NewDefault().WithSettings(&gorlock.Settings{
		KeyPrefix:     "grbinder.entity_lock",
		LockTimeout:   30 * time.Second,
		RetryTimeout:  2 * time.Second,
		RetryInterval: 15 * time.Millisecond,
		LockWaiting:   true,
	})
}

type Locker interface {
	Lock(key string) (bool, error)
	Unlock(key string) error
}

type entityLockOptions struct {
	EnableLock     bool
	Name           string
	LockTakeAction bool
	Locker         Locker
	EntityIdLookup EntityIdLookup
}

// EntityIdLookup defines a func to lookup entity id from context
// if it returns empty string, no lock will be applied
// if the context is aborted in the func, the handler will not be called
// if abort then a status code should be set in the context
// this is not middleware, so don't call .Next()
type EntityIdLookup func(*gin.Context) string

func defaultEntityLockOptions() *entityLockOptions {
	return &entityLockOptions{
		EnableLock:     false,
		Name:           "",
		LockTakeAction: false,
		Locker:         entityLock,
		EntityIdLookup: func(ctx *gin.Context) string {
			return ctx.Param("id")
		},
	}
}

type Option func(*entityLockOptions)

func WithEntityLockEnable(enable bool) Option {
	return func(options *entityLockOptions) {
		options.EnableLock = enable
	}
}

func WithEntityLockName(name string) Option {
	return func(options *entityLockOptions) {
		options.Name = name
	}
}

func WithEntityLockLocker(locker Locker) Option {
	return func(options *entityLockOptions) {
		if locker == nil {
			panic("locker cannot be nil")
		}
		options.Locker = locker
	}
}

func WithEntityLockTakeAction(LockTakeAction bool) Option {
	return func(options *entityLockOptions) {
		options.LockTakeAction = LockTakeAction
	}
}

func WithEntityLookup(entityIdLookup EntityIdLookup) Option {
	return func(options *entityLockOptions) {
		if entityIdLookup == nil {
			panic("entityIdLockup cannot be nil")
		}
		options.EntityIdLookup = entityIdLookup
	}
}

func lockEntityAndHandle(ctx *gin.Context, options *entityLockOptions, handler func(*gin.Context)) {
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
