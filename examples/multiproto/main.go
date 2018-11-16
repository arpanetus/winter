package main

import (
	"github.com/steplems/winter/core"
	"github.com/steplems/winter/examples/multiproto/ws"
	"time"
)

var (
	log = core.NewLogger("example")
	addr = "localhost:5549"
	server = core.NewServer(addr)
)

func main() {
	server.GracefulShutdown = true
	server.HandleWebSocket("/ws", ws.SocketServer)

	go func() {
		time.Sleep(time.Second * 3)
		core.NewWebSocketClient("ws://" + addr + "/ws", nil, ws.SocketClient)
	}()

	server.Start()
}
