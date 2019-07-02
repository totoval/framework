package route

import (
	"github.com/gin-gonic/gin"
)

type RouteGrouper interface {
	Group(grp Grouper)
}

type Grouper interface {
	Can() *group
	AddGroup(relativePath string, routeGrouper RouteGrouper, handlers ...gin.HandlerFunc)
	gin.IRoutes
}
type group struct {
	*gin.RouterGroup
}

func (g *group) Can() *group {
	//@todo for policy module
	return g
}

func (g *group) AddGroup(relativePath string, routeGrouper RouteGrouper, handlers ...gin.HandlerFunc) {
	ginGroup := g.Group(relativePath, handlers...)
	routeGrouper.Group(&group{ginGroup})
}
