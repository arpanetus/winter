package main

import (
	"github.com/steplems/winter/core"
	"github.com/steplems/winter/examples/multiproto/ws"
	"time"
)

var (
	addr = "127.0.0.1:5548"
	log = core.NewLogger("multiproto")
)

func socketClient() {
	socket := core.NewWinterSocketClient("ws://" + addr + "/ws", nil)
	socket.Open(func(addr string) {
		log.Info("Connected to the server on", addr)

		socket.Emit("message", "Data from client")

		socket.On("message", func(data interface{}) {
			log.Info("Data from server:", data)
		})
	})
}

func main() {
	server := core.NewServer(addr)

	server.GracefulShutdown = true
	server.Debug = true
	server.HandleWebSocket("/ws", ws.SimpleWebSocket)

	go func() {
		time.Sleep(time.Second * 3)
		socketClient()
	}()

	server.Start()
}
