package grbinder

import (
	"github.com/gin-gonic/gin"
)

// Mix runs BindVerb and CRUD. One has to make sure no collision!
// when using entity lock with own entity id lookup func
// one has to make sure the entity id lookup func works
// if it returns empty string, no lock will be applied
func Mix(group *gin.RouterGroup, handler any, options ...Option) {
	BindVerb(group, handler, options...)
	CRUD(group, handler, options...)
}
