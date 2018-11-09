package main

import (
	"github.com/steplems/winter/core"
	"net/http"
)

var (
	// Create a new logger with 'cookie' tag.
	// Now every log made with cookieLogger will be tagged like so 'cookie'.
	cookieLogger = core.NewLogger("cookie")

	// NewRouters creates router that will be initiated when its attached to the root router (Server).
	// You can create multiple router and attach them to each other
	apiRouter = core.NewRouter(func(r *core.Router) {
		cookieLogger.Info("Executed when used by some other router")

		// Create your own errors attached to a router or create new ErrorMap
		r.Errors.Set(1, core.NewError(http.StatusInternalServerError, "This route is not working right now"))

		// Creating new handler with util resolver factory Sender
		// Sender returns resolver with null response and sends given interface as json
		r.Get("/", core.Sender("Good util to send json if router needs to only echo"))

		// Resolvers are function that will be called when request hits exact method and path
		// In order to respond to the request you must return the Response or use ctx
		r.Get("/simple", func(ctx *core.Context) core.Response {
			// Returns Response struct with default 200 status and given interface message
			return core.NewSuccessResponse("Cool")
		})

		r.Get("/error", func(ctx *core.Context) core.Response {
			return core.NewErrorResponse(r.Errors.Get(1))
		})

		r.Get("/default_err", func(ctx *core.Context) core.Response {
			// You still have default http errors in every ErrorMap
			return core.NewErrorResponse(r.Errors.Get(http.StatusInternalServerError))
		})
	})
)

func main() {
	// Creating new server with given addr.
	server := core.NewServer(":5539")

	// Debug mode sets middleware to the main router and logs every request with its execution time.
	server.Debug = true

	// GracefulShutdown allows your server to shutdown properly.
	// Let’s imagine you have a http server with Winter which connected to a database,
	// and every time server gets called then it send request to database to get/set a data
	// which also will send to client by response.
	// Imagine you need to shutdown the server, the easiest way to do that is <Ctrl>+C
	// and server will be killed, but wait, what if your server did not finish all the requests,
	// what if some client connections is closed because server is killed
	// and can not handle the requests.
	//
	// By default graceful shutdown is false
	server.GracefulShutdown = true

	// Setting apiRouter to the root router at `/api` path
	server.Set("/api", apiRouter)

	server.Get("/", func(ctx *core.Context) core.Response {
		// Asynchronously sends a message in json formats
		ctx.JSON("I love cookies")

		// Returns Response with empty body. Empty Responses are ignored by handler
		return core.NullResponse()
	})

	// Now checkout http://localhost:5539/
	server.Start()
}
