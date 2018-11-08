package main

import (
	"github.com/steplems/winter/core"
)

var mainRouter = core.NewRouter(func(r *core.Router) {
	r.Errors.Set(1, core.NewError(400, "Bad user object"))

	r.Get("/", func(ctx *core.Context) core.Response {
		return core.NullResponse()
	})
})

func main() {
	server := core.NewServer(":5000")
	server.Debug = true
	server.Set("", mainRouter)

	server.Start()
}
