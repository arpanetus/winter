package main

import (
	"github.com/steplems/winter/core"
)

var (
	cookieLogger = core.NewLogger("cookie")

	apiRouter = core.NewRouter(func(r *core.Router) {
		cookieLogger.Info("Executed when used by some other router")

		r.Get("/", core.Sender("Some shit"))

		r.Get("/simple", func(ctx *core.Context) {
		})
	})
)

func main() {
	server := core.NewServer(":5539")

	server.Set("/api", apiRouter)

	server.Get("/", func(ctx *core.Context) {
		cookieLogger.Info("Wow, what a fancy logger")
		ctx.JSON("I love cookies")
	})

	// Now checkout http://localhost:5539/ and /api/
	server.Start()
}
