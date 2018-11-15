package main

import (
	"github.com/steplems/winter/core"
	"github.com/steplems/winter/examples/multiproto/ws"
)

var (
	log = core.NewLogger("example")
	addr = "localhost:5549"
	server = core.NewServer(addr)
)

func main() {
	server.GracefulShutdown = true
	server.HandleWebSocket("/ws", ws.SocketServer)

	server.OnStart(func(addr string) {
		core.NewWebSocketClient("ws://" + addr + "/ws", nil, ws.SocketClient)
	})

	server.Start()
}
