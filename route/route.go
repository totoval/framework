package route

import (
	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/policy"
)

type route struct {
	bindFunc     func(handlers ...gin.HandlerFunc)
	handlers     []gin.HandlerFunc
	relativePath string
}

func parseParamKey(relativePath string) {

}

func newRoute(relativePath string, bindFunc func(handlers ...gin.HandlerFunc), handlers ...gin.HandlerFunc) *route {
	r := route{relativePath: relativePath, bindFunc: bindFunc, handlers: handlers}
	theList = append(theList, &r)
	return &r
}

func (rp *route) Can(policies policy.Policier, action policy.Action) {
	rp.handlers = append([]gin.HandlerFunc{policy.Middleware(policies, action)}, rp.handlers...)
}

var theList []*route

func Bind() {
	for _, f := range theList {
		f.bindFunc(f.handlers...)
	}
}
