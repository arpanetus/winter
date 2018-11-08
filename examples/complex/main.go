package main

import "github.com/steplems/winter/core"

var mainRouter = core.NewRouter(func(r *core.Router) {
	r.Get("/", func(ctx *core.Context) core.ResolverResponse {
	})
})

func main() {
	server := core.NewServer(":5000")
	server.Debug = true
	server.Set("", )

	server.Start()
}
