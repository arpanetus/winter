package main

import (
	"github.com/steplems/winter/core"
)

type ApiRouter struct {
	*core.Router
}

func (a *ApiRouter) Init() {
	a.All("/", func(ctx *core.Context) {
		ctx.JSON("Api router /api/")
	})
}

func main() {
	server := core.NewServer(":5539")

	server.Get("/", func(ctx *core.Context) {
		core.MainLogger.Info("Wow, what a fancy logger")
		ctx.JSON("I love cookies")
	})

	server.Set("/api", &ApiRouter{})

	server.Start()
}
