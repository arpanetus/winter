package ws

import (
	"github.com/steplems/winter/core"
	"net/http"
)

var (
	logServer = core.NewLogger("socket server")
	logClient = core.NewLogger("socket client")

	SocketServer = core.NewWebSocket(func(conn *core.Connection) {
		conn.OnMessage(func(message core.Message) {
			logServer.Info("Every new message will execute this func")
		})

		conn.On("token", func(data interface{}) {
			logServer.Info("I guess this token is correct")
			conn.Emit("token", http.StatusOK)
		})

		conn.Room("chat", func(conn *core.Connection) {
			conn.On("message", func(data interface{}) {
				conn.Emit("message", data)
			})
		})

		conn.OnClose(func() {
			logServer.Warn("Connection with client closed")
		})
	})
)

func SocketClient(conn *core.Connection) {
	chatRoom := conn.Room("chat", func(conn *core.Connection) {
		conn.On("message", func(data interface{}) {
			logClient.Info("Wow, ive got new message:", data)
		})
	})

	conn.Emit("token", "iGuessItsToken")

	conn.On("token", func(data interface{}) {
		chatRoom.Emit("message", "sup")
	})
}
