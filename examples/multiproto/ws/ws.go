package ws

import (
	"github.com/steplems/winter/core"
)

var (
	log = core.NewLogger("simple socket")
	SimpleWebSocket = core.NewWinterSocket(func(socket *core.Socket) {
		socket.Open(func(addr string) {
			log.Info("Connection opened with", addr)
		})

		socket.Close(func(err error) {
			log.Warn("Connection closed")
		})

		socket.On("message", func(data interface{}) {
			log.Info(data)
			socket.Emit("message", "Socket response")
		})
	})
)
