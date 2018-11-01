package main

import "github.com/steplems/winter/core"

type WsRouter struct {
	*core.WebSocketRouter
}

func (w *WsRouter) Open(ctx *core.WebSocketContext) {
}

func main() {
	server := core.NewServer(":5000")
	server.SetWebSocket("/ws", &WsRouter{})
	server.Start()
}
