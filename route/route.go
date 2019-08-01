package route

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/totoval/framework/app"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/http/middleware"
	"github.com/totoval/framework/policy"
	"github.com/totoval/framework/request"
	"github.com/totoval/framework/request/websocket"
)

type route struct {
	bindFunc          func(handlers ...request.HandlerFunc)
	handlers          []request.HandlerFunc
	relativePath      string
	httpMethod        string
	basicPath         string
	prefixHandlersNum int
	wsHandler         websocket.Handler
}

const maxRouteMapLength = 1000

type engineRoute struct {
	Lock sync.RWMutex
	data map[request.EngineHash]chan *route
}

func newEngineRoute() *engineRoute {
	return &engineRoute{
		data: make(map[request.EngineHash]chan *route),
	}
}
func (erm *engineRoute) init(eh request.EngineHash) {
	if erm.data[eh] == nil {
		erm.data[eh] = make(chan *route, maxRouteMapLength)
	}
}
func (erm *engineRoute) Get(eh request.EngineHash) chan *route {
	erm.Lock.RLock()
	defer erm.Lock.RUnlock()
	return erm.data[eh]
}
func (erm *engineRoute) Set(eh request.EngineHash, r *route) {
	erm.Lock.Lock()
	defer erm.Lock.Unlock()
	erm.init(eh)
	erm.data[eh] <- r
}
func (erm *engineRoute) Close(eh request.EngineHash) {
	close(erm.data[eh])
}

var engineRouteMap *engineRoute

func init() {
	engineRouteMap = newEngineRoute()
}

func newRoute(httpMethod string, g *group, relativePath string, bindFunc func(ginHandlers ...request.HandlerFunc), handlers ...request.HandlerFunc) *route {
	r := route{httpMethod: httpMethod, prefixHandlersNum: len(g.RouterGroup.Handlers), basicPath: g.RouterGroup.BasePath(), relativePath: relativePath, bindFunc: bindFunc, handlers: handlers}

	engineRouteMap.Set(g.engineHash, &r)

	return &r
}
func newWsRoute(httpMethod string, g *group, relativePath string, bindWsFunc func(wsHandler websocket.Handler, ginHandlers ...request.HandlerFunc), wsHandler websocket.Handler, handlers ...request.HandlerFunc) *route {
	r := route{wsHandler: wsHandler, httpMethod: httpMethod, prefixHandlersNum: len(g.RouterGroup.Handlers), basicPath: g.RouterGroup.BasePath(), relativePath: relativePath, bindFunc: func(handlers ...request.HandlerFunc) {
		// normalHandler -> wshandler
		httpHs := append(handlers, func(context request.Context) {
			//@todo websocket handler placeholder
		})
		bindWsFunc(wsHandler, httpHs...)
	}, handlers: handlers}

	engineRouteMap.Set(g.engineHash, &r)

	return &r
}

func (r *route) Can(policies policy.Policier, action policy.Action) {
	r.handlers = append([]request.HandlerFunc{middleware.Policy(policies, action)}, r.handlers...)
}
func (r *route) handlerNum() int {
	return r.prefixHandlersNum + len(r.handlers)
}
func (r *route) absolutePath() string {
	return fmt.Sprintf("%s%s", r.basicPath, r.relativePath)
}
func (r *route) lastHandlerName() string {
	if r.httpMethod == httpMethodWebsocket {
		h := reflect.ValueOf(r.wsHandler).Elem().Type()
		return h.PkgPath() + "." + h.Name()
	}
	if len(r.handlers) > 0 {
		return runtime.FuncForPC(reflect.ValueOf(r.handlers[len(r.handlers)-1]).Pointer()).Name()
	}
	return "nil"
}

func Bind(engine *request.Engine) {
	hash := engine.Hash()
	defer engineRouteMap.Close(hash)

	for r := range engineRouteMap.Get(hash) {
		r.bindFunc(r.handlers...)

		if app.GetMode() != app.ModeProduction {
			log.Info(fmt.Sprintf("%-6s %-30s --> %s (%d handlers)\n", r.httpMethod, r.absolutePath(), r.lastHandlerName(), r.handlerNum()))
		}

		if len(engineRouteMap.Get(hash)) <= 0 {
			break
		}
	}
}
