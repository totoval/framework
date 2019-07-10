package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/policy"
	"github.com/totoval/framework/request"
)

type RouteGrouper interface {
	Group(grp Grouper)
}

type Grouper interface {
	AddGroup(relativePath string, routeGrouper RouteGrouper, handlers ...request.HandlerFunc)
	iRoutes
}
type iRoutes interface {
	//Use(...request.HandlerFunc) gin.IRoutes

	Handle(httpMethod, relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier
	Any(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier
	GET(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier
	POST(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier
	DELETE(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier
	PATCH(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier
	PUT(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier
	OPTIONS(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier
	HEAD(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier

	StaticFile(relativePath, filepath string) gin.IRoutes
	Static(relativePath, root string) gin.IRoutes
	StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes
}

type group struct {
	*gin.RouterGroup
}

func (g *group) Handle(httpMethod, relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier {
	return newRoute(httpMethod, g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.Handle(httpMethod, relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) Any(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier {
	return newRoute("Any", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.Any(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) GET(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier {
	return newRoute("GET", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.GET(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) POST(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier {
	return newRoute("POST", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.POST(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) DELETE(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier {
	return newRoute("DELETE", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.DELETE(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) PATCH(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier {
	return newRoute("PATCH", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.PATCH(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) PUT(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier {
	return newRoute("PUT", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.PUT(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) OPTIONS(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier {
	return newRoute("OPTIONS", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.OPTIONS(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) HEAD(relativePath string, handlers ...request.HandlerFunc) policy.RoutePolicier {
	return newRoute("HEAD", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.HEAD(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
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

func (g *group) AddGroup(relativePath string, routeGrouper RouteGrouper, handlers ...request.HandlerFunc) {
	ginGroup := g.RouterGroup.Group(relativePath, request.ConvertHandlers(handlers)...)
	routeGrouper.Group(&group{RouterGroup: ginGroup})
}
