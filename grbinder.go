package grbinder

import (
	"github.com/gin-gonic/gin"
)

// Mix runs BindVerb and CRUB. One has to make sure no collision!
func Mix(group *gin.RouterGroup, handler any) {
	BindVerb(group, handler)
	CRUD(group, handler)
}
