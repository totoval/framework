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
)

type route struct {
	bindFunc          func(handlers ...request.HandlerFunc)
	handlers          []request.HandlerFunc
	relativePath      string
	httpMethod        string
	basicPath         string
	prefixHandlersNum int
}

const maxRouteMapLength = 1000

type engineRoute struct {
	Lock sync.RWMutex
	data map[versionHash]chan *route
}

func newEngineRoute() *engineRoute {
	return &engineRoute{
		data: make(map[versionHash]chan *route),
	}
}
func (erm *engineRoute) init(vh versionHash) {
	if erm.data[vh] == nil {
		erm.data[vh] = make(chan *route, maxRouteMapLength)
	}
}
func (erm *engineRoute) Get(vh versionHash) chan *route {
	erm.Lock.RLock()
	defer erm.Lock.RUnlock()
	return erm.data[vh]
}
func (erm *engineRoute) Set(vh versionHash, r *route) {
	erm.Lock.Lock()
	defer erm.Lock.Unlock()
	erm.init(vh)
	erm.data[vh] <- r
}
func (erm *engineRoute) Close(vh versionHash) {
	close(erm.data[vh])
}

var engineRouteMap *engineRoute

func init() {
	engineRouteMap = newEngineRoute()
}

func newRoute(httpMethod string, g *group, relativePath string, bindFunc func(ginHandlers ...request.HandlerFunc), handlers ...request.HandlerFunc) *route {
	r := route{httpMethod: httpMethod, prefixHandlersNum: len(g.RouterGroup.Handlers), basicPath: g.RouterGroup.BasePath(), relativePath: relativePath, bindFunc: bindFunc, handlers: handlers}

	engineRouteMap.Set(g.versionHash, &r)

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
	if len(r.handlers) > 0 {
		return runtime.FuncForPC(reflect.ValueOf(r.handlers[len(r.handlers)-1]).Pointer()).Name()
	}
	return "nil"
}

func Bind(engine *request.Engine) {
	hash := engineHash(engine)
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
