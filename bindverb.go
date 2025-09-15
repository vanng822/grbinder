package grbinder

import (
	"github.com/gin-gonic/gin"
)

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
// When using entity lock with own entity id lookup func
// one has to make sure the entity id lookup func works
// if it returns empty string, no lock will be applied
func BindVerb(group *gin.RouterGroup, handler any, options ...Option) {
	var opts = defaultEntityLockOptions()
	for _, option := range options {
		option(opts)
	}
	if handler, ok := handler.(GETSupported); ok {
		if opts.EnableLock && opts.LockTakeAction {
			group.GET("", func(ctx *gin.Context) {
				lockEntityAndHandle(ctx, opts, handler.GET)
			})
		} else {
			group.GET("", handler.GET)
		}
	}
	if handler, ok := handler.(POSTSupported); ok {
		if opts.EnableLock {
			group.POST("", func(ctx *gin.Context) {
				lockEntityAndHandle(ctx, opts, handler.POST)
			})
		} else {
			group.POST("", handler.POST)
		}
	}

	if handler, ok := handler.(PUTSupported); ok {
		if opts.EnableLock {
			group.PUT("", func(ctx *gin.Context) {
				lockEntityAndHandle(ctx, opts, handler.PUT)
			})
		} else {
			group.PUT("", handler.PUT)
		}
	}

	if handler, ok := handler.(DELETESupported); ok {
		if opts.EnableLock {
			group.DELETE("", func(ctx *gin.Context) {
				lockEntityAndHandle(ctx, opts, handler.DELETE)
			})
		} else {
			group.DELETE("", handler.DELETE)
		}
	}
	if handler, ok := handler.(PATCHSupported); ok {
		if opts.EnableLock {
			group.PATCH("", func(ctx *gin.Context) {
				lockEntityAndHandle(ctx, opts, handler.PATCH)
			})
		} else {
			group.PATCH("", handler.PATCH)
		}
	}
	if handler, ok := handler.(HEADSupported); ok {
		group.HEAD("", handler.HEAD)
	}
	if handler, ok := handler.(OPTIONSSupported); ok {
		group.OPTIONS("", handler.OPTIONS)
	}
}
