package ws

import (
	"github.com/steplems/winter/core"
)

var (
	log = core.NewLogger("simple socket")
	SimpleWebSocket = core.NewWinterSocket(func(socket *core.Socket) {
		socket.Close(func(err error) {
			log.Err(err)
		})

		socket.On("message", func(data interface{}) {
			core.WebSocketLogger.Info(data)
			socket.Emit("message", "Socket response")
		})
	})
)
