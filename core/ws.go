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
		winterChanSelect(socket)
	}
}

func WinterSocketClient(conn *Connection) *Socket {
	socket := newSocket(conn)
	go winterChanSelect(socket)
	return socket
}

func newSocket(conn *Connection) *Socket {
	return &Socket{
		conn: conn,
		events: map[string]*SocketResolver{},
		onOpen: func(addr string) {},
		onClose: func(err error) {},
	}
}

func winterChanSelect(socket *Socket) {
	socket.onOpen(<-socket.conn.Open)

	callEvent(socket, <-socket.conn.Message)

	select {
	case err := <-socket.conn.Close:
		socket.onClose(err)
	case err := <-socket.conn.CloseError:
		socket.onClose(err)
	case err := <-socket.conn.UnexpectedCloseError:
		socket.onClose(err)
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

func (s *Socket) On(event string, resolver SocketResolver) {
	s.events[event] = &resolver
}

func (s *Socket) Emit(event string, data ...interface{}) {
	var payload interface{}
	if len(data) > 0 {
		payload = data[0]
	} else {
		payload = ""
	}

	s.conn.JSON(EventMessage{
		Event: event,
		Payload: payload,
	})
}

func (s *Socket) Open(onOpen func(adrr string)) {
	s.onOpen = onOpen
}

func (s *Socket) Close(onClose func(err error)) {
	s.onClose = onClose
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
	messageChan := make(chan *Message, 1)
	closeChan := make(chan error, 1)
	openChan := make(chan string, 1)
	closeErrorChan := make(chan error, 1)
	unexpectedErrorChan := make(chan error, 1)

	connection := &Connection{
		Conn: conn,
		Message: messageChan,
		Open: openChan,
		Close: closeChan,
		CloseError: closeErrorChan,
		UnexpectedCloseError: unexpectedErrorChan,
	}

	go w.Resolver(connection)

	go func() {
		defer conn.Close()

		connection.Open <- conn.RemoteAddr().String()

		for {
			mt, message, err := conn.ReadMessage()

			if err != nil {
				if websocket.IsCloseError(err, connection.UnexpectedCloseErrorCodes...) {
					connection.CloseError <- err
				} else if websocket.IsUnexpectedCloseError(err, connection.CloseErrorCodes...) {
					connection.UnexpectedCloseError <- err
				} else {
					connection.Close <- err
				}

				break
			}

			connection.Message <- NewMessage(mt, message)
		}
		connection.Close <- nil
	}()

	return connection
}
