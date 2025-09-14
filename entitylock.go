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
	Acquire(key string) (bool, error)
	Unlock(key string) error
}

type entityLockOptions struct {
	EnableLock     bool
	Name           string
	LockTakeAction bool
	Locker         Locker
	EntityIdLookup EntityIdLookup
}

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
	var name = options.Name
	if name == "" {
		name = ctx.FullPath()
	}
	var key = fmt.Sprintf("%s.%s", name, id)
	var locked, err = options.Locker.Acquire(key)
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
