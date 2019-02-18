package route

import "github.com/gin-gonic/gin"

type RouteGrouper interface {
	Register(group *gin.RouterGroup)
}
type RouteVersioner interface {
	Register(router *gin.Engine)
}

func RegisterRouteGroup(g RouteGrouper, group *gin.RouterGroup) {
	g.Register(group)
}
