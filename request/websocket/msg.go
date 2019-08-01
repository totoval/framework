package websocket

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"

	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/zone"
)

type Msg struct {
	msgType int
	data    *[]byte
	err     error
}

func (wm *Msg) SetJSON(data interface{}) {
	wm.msgType = websocket.TextMessage
	dataBytes, err := json.Marshal(data)
	wm.data = &dataBytes
	wm.err = err
}
func (wm *Msg) SetProtoBuf(data proto.Message) {
	wm.msgType = websocket.BinaryMessage
	dataBytes, err := proto.Marshal(data)
	wm.data = &dataBytes
	wm.err = err
}
func (wm *Msg) SetString(data string) {
	dataBytes := []byte(data)
	wm.msgType = websocket.TextMessage
	wm.data = &dataBytes
	wm.err = nil
}
func (wm *Msg) SendDone() {
	wm.msgType = -1
	wm.data = nil
	wm.err = nil
}
func (wm *Msg) isDone() bool {
	if wm.msgType == -1 {
		return true
	}
	return false
}

func (wm *Msg) SetByte(msgType int, msg *[]byte) {
	wm.msgType = msgType
	wm.data = msg
	wm.err = nil
}
func (wm *Msg) Type() int {
	return wm.msgType
}
func (wm *Msg) Error() error {
	return wm.err
}

func (wm *Msg) ProtoBuf(dataPtr proto.Message) error {
	return proto.Unmarshal(*wm.data, dataPtr)
}
func (wm *Msg) JSON(dataPtr interface{}) error {
	return json.Unmarshal(*wm.data, dataPtr)
}
func (wm *Msg) String() string {
	return string(*wm.data)
}
func (wm *Msg) Byte() *[]byte {
	return wm.data
}
func (wm *Msg) send(ws *websocket.Conn, wsHandler Handler) error {
	if err := ws.SetWriteDeadline(zone.Now().Add(wsHandler.WriteTimeout())); err != nil {
		return err
	}

	return ws.WriteMessage(wm.msgType, *wm.data)
}
func (wm *Msg) scan(ws *websocket.Conn, wsHandler Handler) error {
	if err := ws.SetReadDeadline(zone.Now().Add(wsHandler.ReadTimeout())); err != nil {
		return err
	}

	msgType, msg, err := ws.ReadMessage()
	if err != nil {
		return log.Error(err)
	}
	wm.msgType = msgType
	wm.data = &msg
	wm.err = err
	return nil
}
