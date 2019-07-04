package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/policy"
)

type RouteGrouper interface {
	Group(grp Grouper)
}

type Grouper interface {
	AddGroup(relativePath string, routeGrouper RouteGrouper, handlers ...gin.HandlerFunc)
	iRoutes
}
type iRoutes interface {
	//Use(...gin.HandlerFunc) gin.IRoutes

	Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier
	Any(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier
	GET(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier
	POST(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier
	DELETE(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier
	PATCH(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier
	PUT(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier
	OPTIONS(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier
	HEAD(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier

	StaticFile(relativePath, filepath string) gin.IRoutes
	Static(relativePath, root string) gin.IRoutes
	StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes
}

type group struct {
	*gin.RouterGroup
}

func (g *group) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier {
	return newRoute(relativePath, func(handlers ...gin.HandlerFunc) { g.RouterGroup.Handle(httpMethod, relativePath, handlers...) }, handlers...)
}

func (g *group) Any(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier {
	return newRoute(relativePath, func(handlers ...gin.HandlerFunc) { g.RouterGroup.Any(relativePath, handlers...) }, handlers...)
}

func (g *group) GET(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier {
	return newRoute(relativePath, func(handlers ...gin.HandlerFunc) { g.RouterGroup.GET(relativePath, handlers...) }, handlers...)
}

func (g *group) POST(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier {
	return newRoute(relativePath, func(handlers ...gin.HandlerFunc) { g.RouterGroup.POST(relativePath, handlers...) }, handlers...)
}

func (g *group) DELETE(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier {
	return newRoute(relativePath, func(handlers ...gin.HandlerFunc) { g.RouterGroup.DELETE(relativePath, handlers...) }, handlers...)
}

func (g *group) PATCH(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier {
	return newRoute(relativePath, func(handlers ...gin.HandlerFunc) { g.RouterGroup.PATCH(relativePath, handlers...) }, handlers...)
}

func (g *group) PUT(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier {
	return newRoute(relativePath, func(handlers ...gin.HandlerFunc) { g.RouterGroup.PUT(relativePath, handlers...) }, handlers...)
}

func (g *group) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier {
	return newRoute(relativePath, func(handlers ...gin.HandlerFunc) { g.RouterGroup.OPTIONS(relativePath, handlers...) }, handlers...)
}

func (g *group) HEAD(relativePath string, handlers ...gin.HandlerFunc) policy.RoutePolicier {
	return newRoute(relativePath, func(handlers ...gin.HandlerFunc) { g.RouterGroup.HEAD(relativePath, handlers...) }, handlers...)
}

//func (g *group) StaticFile(relativePath, filepath string) gin.IRoutes {
//	return g.RouterGroup.StaticFile(relativePath, filepath)
//}
//
//func (g *group) Static(relativePath, root string) gin.IRoutes {
//	return g.RouterGroup.Static(relativePath, root)
//}
//
//func (g *group) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
//	return g.RouterGroup.StaticFS(relativePath, fs)
//}

func (g *group) AddGroup(relativePath string, routeGrouper RouteGrouper, handlers ...gin.HandlerFunc) {
	ginGroup := g.Group(relativePath, handlers...)
	routeGrouper.Group(&group{ginGroup})
}
