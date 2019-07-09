package request

import (
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}
type HandlerFunc func(*Context)

func ConvertHandlers(handlers []HandlerFunc) (ginHandlers []gin.HandlerFunc) {
	for _, h := range handlers {
		handler := h // must new a variable for `range's val`, or the `val` in anonymous funcs will be overwrite every loop

		//debug.Dump(runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name())

		ginHandlers = append(ginHandlers, func(c *gin.Context) {
			totovalContext := Context{Context: c}
			handler(&totovalContext)
		})
	}
	return
}
