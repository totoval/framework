package route

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/totoval/framework/app"
	"github.com/totoval/framework/helpers/log"
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

func newRoute(httpMethod string, prefixHandlersNum int, basicPath string, relativePath string, bindFunc func(ginHandlers ...request.HandlerFunc), handlers ...request.HandlerFunc) *route {
	r := route{httpMethod: httpMethod, prefixHandlersNum: prefixHandlersNum, basicPath: basicPath, relativePath: relativePath, bindFunc: bindFunc, handlers: handlers}
	theList = append(theList, &r)
	return &r
}

func (r *route) Can(policies policy.Policier, action policy.Action) {
	r.handlers = append([]request.HandlerFunc{policy.Middleware(policies, action)}, r.handlers...)
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

var theList []*route

func Bind() {
	for _, r := range theList {
		r.bindFunc(r.handlers...)

		if app.GetMode() != app.ModeProduction {
			log.Info(fmt.Sprintf("%-6s %-30s --> %s (%d handlers)\n", r.httpMethod, r.absolutePath(), r.lastHandlerName(), r.handlerNum()))
		}
	}
}
