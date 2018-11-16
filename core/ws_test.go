package core

import (
	"github.com/gorilla/websocket"
	"testing"
	"time"
)

func TestNewWebSocket(t *testing.T) {
	addr := "localhost:5150"
	wsAddr := "ws://localhost:5150/ws"
	tokenMessage := "test"

	ws := NewWebSocket(func(conn *Connection) {
		conn.OnMessage(func(message Message) {
			t.Log("Received message from client")

			if string(message.Data) != tokenMessage {
				t.Error("Incorrect token")
				return
			}

			t.Log("Token is correct, sending it back")
			conn.Send(websocket.TextMessage, []byte(tokenMessage))
		})

		conn.OnError(func(err error) {
			t.Error("Error trying to connect with client", err)
		})
	}, &websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,
	})

	server := startServer(addr)
	server.HandleWebSocket("/ws", ws)

	time.Sleep(time.Second)

	go NewWebSocketClient(wsAddr, nil, func(conn *Connection) {
		t.Log("Connected, sending token")

		conn.Send(websocket.TextMessage, []byte(tokenMessage))
		conn.OnMessage(func(message Message) {
			t.Log("New message", message)
			if string(message.Data) != tokenMessage {
				t.Error("Incorrect token")
			}
		})

		conn.OnError(func(err error) {
			t.Error("Error trying to connect with server", err)
		})
	})

	time.Sleep(time.Second * 3)

	server.NativeServer.Shutdown(nil)
}
