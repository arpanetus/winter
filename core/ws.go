package core

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

func NewWebSocket(resolver WebSocketResolver, upgrader ...*websocket.Upgrader) *WebSocket {
	defaultUpgrader := &websocket.Upgrader{}
	if len(upgrader) > 0 {
		defaultUpgrader = upgrader[0]
	}

	return &WebSocket{
		upgrader: NewUpgrader(defaultUpgrader),
		resolver: resolver,
		connection: &Connection{
			RoomPath: main_room_path,
			events: map[string]map[string]*EventResolver{},
			onClose: func() {},
			onError: func(err error) {},
			onMessage: func(message Message) {},
		},
		Headers: nil,
	}
}

func NewWebSocketClient(url string, headers http.Header, onOpen func(conn *Connection), upgrader ...*websocket.Upgrader) {
	defaultUpgrader := &websocket.Upgrader{}
	if len(upgrader) > 0 {
		defaultUpgrader = upgrader[0]
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		WebSocketLogger.Err("Failed to connect:", err)
		return
	}

	ws := NewWebSocket(func(conn *Connection) {}, NewUpgrader(defaultUpgrader))
	ws.Headers = headers
	onOpen(ws.newDialerConnection(conn))
	ws.dial()
}

func NewUpgrader(upgrader *websocket.Upgrader) *websocket.Upgrader {
	upgrader.Error = func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		SendResponse(NewResponse(status, reason.Error()))(w, r)
	}
	return upgrader
}

func (w *WebSocket) GetResolver() Resolver {
	return func(ctx *Context) Response {
		w.handler(ctx.Response, ctx.Request)
		return NullResponse()
	}
}

func (w *WebSocket) GetHandlerFunc() http.HandlerFunc {
	return w.handler
}

func (w *WebSocket) handler(res http.ResponseWriter, req *http.Request) {
	conn, _ := w.upgrader.Upgrade(res, req, w.Headers)
	w.newDialerConnection(conn)
	w.resolver(w.connection)
	w.dial()
}

func (w *WebSocket) newDialerConnection(conn *websocket.Conn) *Connection {
	w.connection.ExtendedConnection = conn
	return w.connection
}

func (w *WebSocket) dial() {
	defer w.connection.ExtendedConnection.Close()

	for {
		mt, mess, err := w.connection.ExtendedConnection.ReadMessage()
		if err != nil {
			w.connection.onError(err)
			break
		}

		message := Message{mt, mess}

		w.connection.onMessage(message)
		w.connection.callEvent(message)
	}

	w.connection.onClose()
}

func (c *Connection) OnMessage(onMessage func(message Message)) {
	c.onMessage = onMessage
}

func (c *Connection) OnError(onError func(err error)) {
	c.onError = onError
}

func (c *Connection) OnClose(onClose func()) {
	c.onClose = onClose
}

func (c *Connection) Send(mt int, data []byte) {
	c.ExtendedConnection.WriteMessage(mt, data)
}

func (c *Connection) JSON(mess interface{}) {
	c.ExtendedConnection.WriteJSON(mess)
}

func (c *Connection) On(event string, resolver EventResolver) {
	roomMap := c.events[c.RoomPath]
	if len(roomMap) == 0 {
		c.events[c.RoomPath] = map[string]*EventResolver{
			event: &resolver,
		}
	}

	c.events[c.RoomPath][event] = &resolver
}

func (c *Connection) Emit(event string, data ...interface{}) {
	var payload interface{}
	if len(data) > 0 {
		payload = data[0]
	} else {
		payload = ""
	}

	c.JSON(c.newEventMessage(event, payload))
}

func (c Connection) Room(name string, resolver ...WebSocketResolver) *Connection {
	c.RoomPath = name
	if len(resolver) > 0 { resolver[0](&c) }
	return &c
}

func (c *Connection) newEventMessage(event string, payload interface{}) EventMessage {
	return EventMessage{
		Room: c.RoomPath,
		Event: event,
		Payload: payload,
	}
}

func (c *Connection) callEvent(message Message) {
	eventMessage := EventMessage{}
	err := json.Unmarshal(message.Data, &eventMessage)
	if err != nil {
		return
	}

	if len(c.events[eventMessage.Room]) > 0 {
		fnc, ok := c.events[eventMessage.Room][eventMessage.Event]
		if ok {
			(*fnc)(eventMessage.Payload)
		}
	}
}
