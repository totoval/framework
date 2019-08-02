package websocket

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/helpers/zone"
	request_http "github.com/totoval/framework/request/http"
)

func ConvertHandler(wsHandler Handler) gin.HandlerFunc {
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var pingPeriod = (wsHandler.ReadTimeout() * 9) / 10
	return func(c *gin.Context) {
		totovalContext := request_http.ConvertContext(c)

		// create connectionHub
		hub := &connectionHub{} //@todo interface
		hub.init(totovalContext)
		//@todo add Hub to hub list, for broadcast

		////@todo every handler struct has share'd it's Context
		//typeof := reflect.TypeOf(wsHandler)
		//
		//ptr := reflect.New(typeof).Elem()
		//val := reflect.New(typeof.Elem())
		//ptr.Set(val)
		//newHandler := ptr.Interface().(Handler)
		//debug.Dump(newHandler)

		ws, err := wsUpgrader.Upgrade(totovalContext.Writer(), totovalContext.Request(), nil)
		if err != nil {
			_ = log.Error(err, toto.V{"msg": "Failed to set websocket upgrade"})
			totovalContext.JSON(http.StatusUnprocessableEntity, toto.V{"error": err})
			return
		}
		// close ws connection
		defer func() {
			if err := ws.Close(); err != nil {
				_ = log.Error(err)
			}
			hub.close()
		}()

		// handle panic
		defer func() {
			if _err := recover(); _err != nil {
				if __err, ok := _err.(error); ok {
					err = log.Error(__err)
					return
				}
				err = log.Error(errors.New(fmt.Sprint(_err)))
				return
			}
		}()

		// setMaxSize
		ws.SetReadLimit(wsHandler.MaxMessageSize())

		// set handlers
		setPingPongCloseHandlers(ws, wsHandler, hub)

		// --------------------------
		// ---- websocket process ----
		// --------------------------

		// send msg
		go func() {
			pingTicker := zone.NewTicker(pingPeriod)
			defer pingTicker.Stop()

			for {
				if err := wsHandler.Loop(hub); err != nil {
					return
				}

				select {
				// send message
				case msg, ok := <-hub.getChan():
					if !ok {
						return
					}
					if msg.isDone() {
						return
					}
					if err := msg.send(ws, wsHandler); err != nil {
						_ = log.Error(err)
						return
					}

				// send ping
				case <-pingTicker.C:
					ping := Msg{
						msgType: websocket.PingMessage,
						data:    &[]byte{},
						err:     nil,
					}
					if err := ping.send(ws, wsHandler); err != nil {
						_ = log.Error(err)
						return
					}

				default:
					continue
				}
			}
		}()

		// read msg
		for {
			msg := &Msg{}
			if msg.scan(ws, wsHandler) != nil {
				_ = log.Error(err)
				return
			}

			switch msg.msgType {
			// TextMessage denotes a text data message. The text message payload is
			// interpreted as UTF-8 encoded text data.
			case websocket.TextMessage:
				fallthrough
			// BinaryMessage denotes a binary data message.
			case websocket.BinaryMessage:
				wsHandler.OnMessage(hub, msg)
			// PingMessage denotes a ping control message. The optional message payload
			// is UTF-8 encoded text.
			case websocket.PingMessage: // defined at SetPingHandler
				// send pong
				pong := Msg{
					msgType: websocket.PongMessage,
					data:    &[]byte{},
					err:     nil,
				}
				if err := pong.send(ws, wsHandler); err != nil {
					_ = log.Error(err)
					return
				}
			// PongMessage denotes a pong control message. The optional message payload
			// is UTF-8 encoded text.
			case websocket.PongMessage: // defined at SetPongHandler
				// check deadline
				if msg.err != nil {
					_ = log.Error(err)
					return
				}
			case websocket.CloseMessage: // defined at SetCloseHandler
				closeMsg := Msg{
					msgType: websocket.CloseMessage,
					data:    &[]byte{},
					err:     nil,
				}
				if err := closeMsg.send(ws, wsHandler); err != nil {
					_ = log.Error(err)
					return
				}

				return
			default:
				log.Warn("No websocket handler on this msgType", toto.V{"msgType": msg.msgType})
			}

		}
	}
}

func setPingPongCloseHandlers(ws *websocket.Conn, wsHandler Handler, hub Hub) {
	ws.SetCloseHandler(func(code int, text string) error {
		wsHandler.OnClose(hub, code, text)
		return nil
	})
	ws.SetPingHandler(func(appData string) error {
		wsHandler.OnPing(hub, appData)
		return nil
	})
	ws.SetPongHandler(func(appData string) error {
		wsHandler.OnPing(hub, appData)
		return nil
	})
}
