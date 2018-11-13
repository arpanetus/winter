package ws

import "github.com/steplems/winter/core"

var (
	SimpleWebSocket = core.NewWinterSocket(func(socket *core.Socket) {
		socket.On("message", func(data interface{}) {
			core.WebSocketLogger.Info(data)
			socket.Emit("mess", "Socket response")
		})
	})
)
