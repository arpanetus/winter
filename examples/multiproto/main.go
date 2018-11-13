package main

import (
	"github.com/steplems/winter/core"
	"github.com/steplems/winter/examples/multiproto/ws"
)

func main() {
	server := core.NewServer(":5548")
	server.GracefulShutdown = true

	server.HandleWebSocket("/ws", ws.SimpleWebSocket)

	server.Start()
}
