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
		Resolver: resolver,
		Upgrader: defaultUpgrader,
	}
}

func NewMessage(messageType int, message []byte) *Message {
	return &Message{
		Type: messageType,
		Value: message,
	}
}

func WinterSocket(resolver WinterSocketResolver) WebSocketResolver {
	return func(conn *Connection) {
		socket := &Socket{
			conn: conn,
			events: map[string]*SocketResolver{},
			OnCloseError: func(err error) {},
			OnUnexpectedCloseError: func(err error) {},
		}

		resolver(socket)

		select {
		case message := <-conn.Message:
			callEvent(socket, message)
		case err := <-conn.CloseError:
			socket.OnCloseError(err)
		case err := <-conn.UnexpectedCloseError:
			socket.OnUnexpectedCloseError(err)
		}
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

func (w *Socket) Send(messageType int, message []byte) {
	w.conn.Conn.WriteMessage(messageType, message)
}

func (w *Socket) JSON(json interface{}) {
	w.conn.Conn.WriteJSON(json)
}

func (w *WebSocket) resolver(res http.ResponseWriter, req *http.Request) {
	conn, err := w.Upgrader.Upgrade(res, req, nil)
	if err != nil {
		WebSocketLogger.Err("Error while trying to get new Connection:", err)
		return
	}

	w.dialer(conn)
}

func (w *WebSocket) dialer(conn *websocket.Conn) {
	defer conn.Close()

	messageChan := make(chan *Message)
	closeErrorChan := make(chan error)
	unexpectedErrorChan := make(chan error)

	connection := &Connection{
		Conn: conn,
		Message: messageChan,
		CloseError: closeErrorChan,
		UnexpectedCloseError: unexpectedErrorChan,
	}

	go w.Resolver(connection)

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
}
