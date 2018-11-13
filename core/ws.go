package core

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

func NewWinterSocket(resolver WinterSocketResolver, upgrader ...*websocket.Upgrader) *WebSocket {
	return NewWebSocket(WinterSocket(resolver), upgrader...)
}

func NewWinterSocketClient(url string, requestHeader http.Header) *Socket {
	return WinterSocketClient(NewWebSocketClient(url, requestHeader))
}

func NewWebSocket(resolver WebSocketResolver, upgrader ...*websocket.Upgrader) *WebSocket {
	defaultUpgrader := &websocket.Upgrader{}
	if len(upgrader) > 0 {
		defaultUpgrader = upgrader[0]
	}
	return &WebSocket{
		Resolver: resolver,
		Upgrader: defaultUpgrader,
	}
}

func NewWebSocketClient(url string, requestHeader http.Header) *Connection {
	conn, _, err := websocket.DefaultDialer.Dial(url, requestHeader)
	if err != nil {
		WebSocketLogger.Err("Couldn't connect to a websocket:", err)
		return &Connection{}
	}

	ws := WebSocket{}
	ws.Resolver = func(conn *Connection) {}

	return ws.dialer(conn)
}

func NewMessage(messageType int, message []byte) *Message {
	return &Message{
		Type: messageType,
		Value: message,
	}
}

func WinterSocket(resolver WinterSocketResolver) WebSocketResolver {
	return func(conn *Connection) {
		socket := newSocket(conn)
		resolver(socket)
		winterChanSelect(conn, socket)
	}
}

func WinterSocketClient(conn *Connection) *Socket {
	socket := newSocket(conn)
	go winterChanSelect(conn, socket)
	return socket
}

func newSocket(conn *Connection) *Socket {
	return &Socket{
		conn: conn,
		events: map[string]*SocketResolver{},
		OnCloseError: func(err error) {},
		OnUnexpectedCloseError: func(err error) {},
	}
}

func winterChanSelect(conn *Connection, socket *Socket) {
	select {
	case message := <-conn.Message:
		callEvent(socket, message)
	case err := <-conn.CloseError:
		socket.OnCloseError(err)
	case err := <-conn.UnexpectedCloseError:
		socket.OnUnexpectedCloseError(err)
	}
}

func callEvent(socket *Socket, message *Message) {
	eventMessage := EventMessage{}
	if err := json.Unmarshal(message.Value, &eventMessage); err != nil {
		return
	}

	for key, fn := range socket.events {
		if key == eventMessage.Event {
			(*fn)(eventMessage.Payload)
			break
		}
	}
}

func (w *Socket) On(event string, resolver SocketResolver) {
	w.events[event] = &resolver
}

func (w *Socket) Emit(event string, data ...interface{}) {
	w.conn.Conn.WriteJSON(EventMessage{
		Event: event,
		Payload: data[0],
	})
}

func (c *Connection) Send(messageType int, message []byte) {
	c.Conn.WriteMessage(messageType, message)
}

func (c *Connection) JSON(json interface{}) {
	c.Conn.WriteJSON(json)
}

func (w *WebSocket) GetHandlerFunc() http.HandlerFunc {
	return w.resolver
}

func (w *WebSocket) resolver(res http.ResponseWriter, req *http.Request) {
	conn, err := w.Upgrader.Upgrade(res, req, nil)
	if err != nil {
		WebSocketLogger.Err("Error while trying to get new Connection:", err)
		return
	}

	w.dialer(conn)
}

func (w *WebSocket) dialer(conn *websocket.Conn) *Connection {
	defer conn.Close()

	open := make(chan *websocket.Conn)
	messageChan := make(chan *Message)
	closeErrorChan := make(chan error)
	unexpectedErrorChan := make(chan error)

	connection := &Connection{
		Conn: conn,
		Open: open,
		Message: messageChan,
		CloseError: closeErrorChan,
		UnexpectedCloseError: unexpectedErrorChan,
	}

	go w.Resolver(connection)
	go func() {
		defer conn.Close()
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, connection.UnexpectedCloseErrorCodes...) {
					closeErrorChan <- err
				} else if websocket.IsUnexpectedCloseError(err, connection.CloseErrorCodes...) {
					unexpectedErrorChan <- err
				}
				break
			}

			messageChan <- NewMessage(mt, message)
		}
	}()

	return connection
}
