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

// CRUD set up 5 handlers for this group
// beside CRUD it includes list
func CRUD(group *gin.RouterGroup, handler interface{}) {
	if handler, ok := handler.(CreateSupported); ok {
		group.POST("", handler.CreateHandler)
	}
	if handler, ok := handler.(ListSupported); ok {
		group.GET("", handler.ListHandler)
	}
	if handler, ok := handler.(TakeSupported); ok {
		group.GET("/:id", handler.TakeHandler)
	}
	if handler, ok := handler.(UpdateSupported); ok {
		group.PUT("/:id", handler.UpdateHandler)
	}
	if handler, ok := handler.(DeleteSupported); ok {
		group.DELETE("/:id", handler.DeleteHandler)
	}
}

// CRUDI set up 6 handlers for this group
// beside CRUD it includes list and init form at path {path}_new
func CRUDI(path string, handler interface{}, router *gin.Engine, groupHandlers ...gin.HandlerFunc) {
	group := router.Group(path, groupHandlers...)

	if handler, ok := handler.(InitSupported); ok {
		router.GET(path+"_new", handler.InitHandler)
	}
	if handler, ok := handler.(CreateSupported); ok {
		group.POST("", handler.CreateHandler)
	}
	if handler, ok := handler.(ListSupported); ok {
		group.GET("", handler.ListHandler)
	}
	if handler, ok := handler.(TakeSupported); ok {
		group.GET("/:id", handler.TakeHandler)
	}
	if handler, ok := handler.(UpdateSupported); ok {
		group.PUT("/:id", handler.UpdateHandler)
	}
	if handler, ok := handler.(DeleteSupported); ok {
		group.DELETE("/:id", handler.DeleteHandler)
	}
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
func BindVerb(group *gin.RouterGroup, handler interface{}) {
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
