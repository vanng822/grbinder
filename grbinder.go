package grbinder

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vanng822/gorlock/v2"
)

// InitSupported binds GET: /path_new for the web, init the form
type InitSupported interface {
	InitHandler(*gin.Context)
}

// CreateSupported binds POST: /path create item
type CreateSupported interface {
	CreateHandler(*gin.Context)
}

// ListSupported binds GET: /path list item
type ListSupported interface {
	ListHandler(*gin.Context)
}

// TakeSupported binds GET: /path/:id get item
type TakeSupported interface {
	TakeHandler(*gin.Context)
}

// UpdateSupported binds PUT:  /path/:id change item
type UpdateSupported interface {
	UpdateHandler(*gin.Context)
}

// DeleteSupported binds DELETE: /path/:id delete item
type DeleteSupported interface {
	DeleteHandler(*gin.Context)
}

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
}

func defaultEntityLockOptions() *entityLockOptions {
	return &entityLockOptions{
		EnableLock:     false,
		Name:           "",
		LockTakeAction: false,
		Locker:         entityLock,
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
		options.Locker = locker
	}
}

func WithEntityLockTakeAction(LockTakeAction bool) Option {
	return func(options *entityLockOptions) {
		options.LockTakeAction = LockTakeAction
	}
}

func lockEntityAndHandle(ctx *gin.Context, options *entityLockOptions, handler func(*gin.Context)) {
	// Lock the entity
	var id = ctx.Param("id")
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

func crud(group *gin.RouterGroup, handler any, options *entityLockOptions) {
	if handler, ok := handler.(CreateSupported); ok {
		group.POST("", handler.CreateHandler)
	}
	if handler, ok := handler.(ListSupported); ok {
		group.GET("", handler.ListHandler)
	}
	if handler, ok := handler.(TakeSupported); ok {
		if options.EnableLock && options.LockTakeAction {
			group.GET("/:id", func(ctx *gin.Context) {
				lockEntityAndHandle(ctx, options, handler.TakeHandler)
			})
		} else {
			group.GET("/:id", handler.TakeHandler)
		}
	}
	if handler, ok := handler.(UpdateSupported); ok {
		if options.EnableLock {
			group.PUT("/:id", func(ctx *gin.Context) {
				lockEntityAndHandle(ctx, options, handler.UpdateHandler)
			})
		} else {
			group.PUT("/:id", handler.UpdateHandler)
		}
	}

	if handler, ok := handler.(DeleteSupported); ok {
		if options.EnableLock {
			group.DELETE("/:id", func(ctx *gin.Context) {
				lockEntityAndHandle(ctx, options, handler.DeleteHandler)
			})
		} else {
			group.DELETE("/:id", handler.DeleteHandler)
		}
	}
}

// CRUD set up 5 handlers for this group
// beside CRUD it includes list
func CRUD(group *gin.RouterGroup, handler any, options ...Option) {
	var opts = defaultEntityLockOptions()
	for _, option := range options {
		option(opts)
	}
	crud(group, handler, opts)
}

// CRUDI set up 6 handlers for this group
// beside CRUD it includes list and init form at path {path}/:id/new
// you may call GET path/0/new
func CRUDI(group *gin.RouterGroup, handler any, options ...Option) {

	if handler, ok := handler.(InitSupported); ok {
		group.GET("/:id/new", handler.InitHandler)
	}

	CRUD(group, handler, options...)
}

// POSTSupported binds POST: /path
type POSTSupported interface {
	POST(*gin.Context)
}

// GETSupported binds GET: /path
type GETSupported interface {
	GET(*gin.Context)
}

// PUTSupported binds PUT:  /path
type PUTSupported interface {
	PUT(*gin.Context)
}

// DELETESupported binds DELETE: /path
type DELETESupported interface {
	DELETE(*gin.Context)
}

// PATCHSupported binds PATCH: /path
type PATCHSupported interface {
	PATCH(*gin.Context)
}

// HEADSupported binds HEAD: /path
type HEADSupported interface {
	HEAD(*gin.Context)
}

// OPTIONSSupported binds OPTIONS: /path
type OPTIONSSupported interface {
	OPTIONS(*gin.Context)
}

// BindVerb bind funcs with http verbs to a route group
// be aware that all binds to the same route
func BindVerb(group *gin.RouterGroup, handler any) {
	if handler, ok := handler.(GETSupported); ok {
		group.GET("", handler.GET)
	}
	if handler, ok := handler.(POSTSupported); ok {
		group.POST("", handler.POST)
	}

	if handler, ok := handler.(PUTSupported); ok {
		group.PUT("", handler.PUT)
	}

	if handler, ok := handler.(DELETESupported); ok {
		group.DELETE("", handler.DELETE)
	}
	if handler, ok := handler.(PATCHSupported); ok {
		group.PATCH("", handler.PATCH)
	}
	if handler, ok := handler.(HEADSupported); ok {
		group.HEAD("", handler.HEAD)
	}
	if handler, ok := handler.(OPTIONSSupported); ok {
		group.OPTIONS("", handler.OPTIONS)
	}
}

// Mix runs BindVerb and CRUB. One has to make sure no collision!
func Mix(group *gin.RouterGroup, handler any) {
	BindVerb(group, handler)
	CRUD(group, handler)
}
