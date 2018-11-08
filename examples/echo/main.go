package main

import (
	"github.com/steplems/winter/core"
)

var (
	cookieLogger = core.NewLogger("cookie")

	apiRouter = core.NewRouter(func(r *core.Router) {
		cookieLogger.Info("Executed when used by some other router")

		r.Get("/", core.Sender("Some shit"))

		r.Get("/simple", func(ctx *core.Context) core.Response {
			return core.NewSuccessResponse("Cool")
		})
	})
)

func main() {
	server := core.NewServer(":5539")

	server.Set("/api", apiRouter)

	server.Get("/", func(ctx *core.Context) core.Response {
		cookieLogger.Info("Wow, what a fancy logger")

		ctx.JSON("I love cookies")

		return core.NullResponse()
	})

	// Now checkout http://localhost:5539/ and /api/
	server.Start()
}
