package grbinder

import (
	"github.com/gin-gonic/gin"
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
// When using entity lock with own entity id lookup func
// one has to make sure the entity id lookup func works
// if it returns empty string, no lock will be applied
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
// When using entity lock with own entity id lookup func
// one has to make sure the entity id lookup func works
// if it returns empty string, no lock will be applied
func CRUDI(group *gin.RouterGroup, handler any, options ...Option) {

	if handler, ok := handler.(InitSupported); ok {
		group.GET("/:id/new", handler.InitHandler)
	}

	CRUD(group, handler, options...)
}
