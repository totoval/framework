package route

import (
	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/http/middleware"
)

type version struct {
	engine *gin.Engine
	group  *gin.RouterGroup
	prefix string
}

func NewVersion(engine *gin.Engine, prefix string) *version {
	ver := &version{engine: engine, prefix: prefix}
	ver.group = ver.engine.Group(prefix)
	return ver
}

func (v *version) Auth(relativePath string, groupFunc func(grp Grouper), handlers ...gin.HandlerFunc) {
	ginGroup := v.group.Group(relativePath, append([]gin.HandlerFunc{middleware.AuthRequired(),}, handlers...)...)
	groupFunc(&group{ginGroup})
}

func (v *version) NoAuth(relativePath string, groupFunc func(grp Grouper), handlers ...gin.HandlerFunc) {
	ginGroup := v.group.Group(relativePath, handlers...)
	groupFunc(&group{ginGroup})
}
