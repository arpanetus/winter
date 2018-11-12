package ws

import (
	"github.com/steplems/winter/core"
)

var SimpleWebSocket = core.NewWebSocket(core.WinterSocket(func(socket *core.Socket) {
	socket.On("message", func(data interface{}) {
		core.WebSocketLogger.Info(data)
		socket.JSON("Socket response")
	})
}))
