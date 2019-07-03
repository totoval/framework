package route

import (
	"github.com/gin-gonic/gin"

	policy_package "github.com/totoval/framework/policy"
)

type route struct {
	bindFunc func(handlers ...gin.HandlerFunc)
	handlers []gin.HandlerFunc
}

func newRoute(bindFunc func(handlers ...gin.HandlerFunc), handlers ...gin.HandlerFunc) *route {
	r := route{bindFunc: bindFunc, handlers: handlers}
	theList = append(theList, &r)
	return &r
}

func (rp *route) Can(policy policy_package.Policier, action policy_package.Action) {
	rp.handlers = append([]gin.HandlerFunc{policy_package.Middleware(policy, action)}, rp.handlers...)
}

var theList []*route

func Bind() {
	for _, f := range theList {
		f.bindFunc(f.handlers...)
	}
}
