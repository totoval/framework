package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/request/websocket"

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

	Handle(httpMethod, relativePath string, handlers ...request.HandlerFunc) routeEnder
	Any(relativePath string, handlers ...request.HandlerFunc) routeEnder
	GET(relativePath string, handlers ...request.HandlerFunc) routeEnder
	POST(relativePath string, handlers ...request.HandlerFunc) routeEnder
	DELETE(relativePath string, handlers ...request.HandlerFunc) routeEnder
	PATCH(relativePath string, handlers ...request.HandlerFunc) routeEnder
	PUT(relativePath string, handlers ...request.HandlerFunc) routeEnder
	OPTIONS(relativePath string, handlers ...request.HandlerFunc) routeEnder
	HEAD(relativePath string, handlers ...request.HandlerFunc) routeEnder

	Websocket(relativePath string, wsHandler websocket.Handler, handlers ...request.HandlerFunc) routeEnder

	StaticFile(relativePath, filepath string) gin.IRoutes
	Static(relativePath, root string) gin.IRoutes
	StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes
}

type group struct {
	engineHash request.EngineHash

	*gin.RouterGroup
}

type routeEnder interface {
	policy.RoutePolicier

	Name(routeName string) policy.RoutePolicier
}

func (g *group) Handle(httpMethod, relativePath string, handlers ...request.HandlerFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute(httpMethod, g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.Handle(httpMethod, relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) Any(relativePath string, handlers ...request.HandlerFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("Any", g, g.clearPath(relativePath), func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.Any(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) GET(relativePath string, handlers ...request.HandlerFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("GET", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.GET(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) POST(relativePath string, handlers ...request.HandlerFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("POST", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.POST(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) DELETE(relativePath string, handlers ...request.HandlerFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("DELETE", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.DELETE(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) PATCH(relativePath string, handlers ...request.HandlerFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("PATCH", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.PATCH(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) PUT(relativePath string, handlers ...request.HandlerFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("PUT", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.PUT(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) OPTIONS(relativePath string, handlers ...request.HandlerFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("OPTIONS", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.OPTIONS(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) HEAD(relativePath string, handlers ...request.HandlerFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newRoute("HEAD", g, relativePath, func(innerHandlers ...request.HandlerFunc) { g.RouterGroup.HEAD(relativePath, request.ConvertHandlers(innerHandlers)...) }, handlers...)
}

func (g *group) StaticFile(relativePath, filepath string) gin.IRoutes {
	relativePath = g.clearPath(relativePath)
	return g.RouterGroup.StaticFile(relativePath, filepath)
}

func (g *group) Static(relativePath, root string) gin.IRoutes {
	relativePath = g.clearPath(relativePath)
	return g.RouterGroup.Static(relativePath, root)
}

func (g *group) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	relativePath = g.clearPath(relativePath)
	return g.RouterGroup.StaticFS(relativePath, fs)
}

const httpMethodWebsocket = "WS"

func (g *group) Websocket(relativePath string, wsHandler websocket.Handler, handlers ...request.HandlerFunc) routeEnder {
	relativePath = g.clearPath(relativePath)
	return newWsRoute(httpMethodWebsocket, g, relativePath, func(wsHandler websocket.Handler, innerHandlers ...request.HandlerFunc) {
		innerGinHandlers := append(request.ConvertHandlers(innerHandlers), websocket.ConvertHandler(wsHandler))
		g.RouterGroup.GET(relativePath, innerGinHandlers...)
	}, wsHandler, handlers...)
}

func (g *group) AddGroup(relativePath string, routeGrouper RouteGrouper, handlers ...request.HandlerFunc) {
	ginGroup := g.RouterGroup.Group(relativePath, request.ConvertHandlers(handlers)...)
	routeGrouper.Group(&group{engineHash: g.engineHash, RouterGroup: ginGroup})
}

func (g *group) clearPath(relativePath string) string {
	if relativePath == "" {
		return relativePath
	}

	basePath := g.RouterGroup.BasePath()
	if basePath[len(basePath)-1:] == "/" && relativePath[:1] == "/" {
		return relativePath[1:]
	}
	if basePath[len(basePath)-1:] != "/" && relativePath[:1] != "/" {
		return "/" + relativePath
	}
	return relativePath
}
