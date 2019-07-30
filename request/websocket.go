package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/toto"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsHandlerFunc func(WsContext)

func ConvertWsHandlers(handlers []HandlerFunc) (ginHandlers []gin.HandlerFunc) {
	for i, h := range handlers {
		handler := h // must new a variable for `range's val`, or the `val` in anonymous funcs will be overwrited in every loop

		if i == len(handlers)-1 { // websocket handler only could be bind at the last handler
			ginHandlers = append(ginHandlers, func(c *gin.Context) {
				ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
				if err != nil {
					_ = log.Error(err, toto.V{"msg": "Failed to set websocket upgrade"})
					c.JSON(http.StatusUnprocessableEntity, toto.V{"error": err})
					return
				}
				totovalContext := websocketContext{Context: c, CommonContext: &CommonContext{c}, ws: ws}
				handler(&totovalContext)
			})
			return
		}

		ginHandlers = append(ginHandlers, func(c *gin.Context) {
			totovalContext := httpContext{Context: c, CommonContext: &CommonContext{c}}
			handler(&totovalContext)
		})
	}
	return
}
