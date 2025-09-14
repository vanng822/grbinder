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
