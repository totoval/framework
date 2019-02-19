package route

import "github.com/gin-gonic/gin"

type RouteGrouper interface {
	Register(group *gin.RouterGroup)
}
type RouteVersioner interface {
	Register(router *gin.Engine)
	noAuth(group *gin.RouterGroup)
	auth(group *gin.RouterGroup)
}

func RegisterRouteGroup(g RouteGrouper, group *gin.RouterGroup) {
	g.Register(group)
}
