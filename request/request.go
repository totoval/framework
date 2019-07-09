package request

import (
	"github.com/gin-gonic/gin"
)

type Engine struct {
	*gin.Engine
}

func New() *Engine {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		//log.Info(fmt.Sprintf("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers))
	}
	return &Engine{Engine: gin.New()}

}

func (e *Engine) Use(handlers ...HandlerFunc) {
	e.Engine.Use(ConvertHandlers(handlers)...)
}
func (e *Engine) UseGin(handlerFunc ...gin.HandlerFunc) {
	e.Engine.Use(handlerFunc...)
}

func (e *Engine) GinEngine() *gin.Engine {
	return e.Engine
}
