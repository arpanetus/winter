package main

import (
	"github.com/steplems/winter/core"
	"github.com/steplems/winter/examples/multiproto/ws"
)

var rootRouter = core.NewRouter(func(r *core.Router) {
	r.HandleWebSocket("/ws", ws.SimpleWebSocket)
})

func main() {
	server := core.NewServer(":5549")
	server.SetRootRouter(rootRouter)
	server.Start()
}
