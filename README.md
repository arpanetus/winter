# Winter Framework
Winter is a framework made in order to rapidly develope scalable web applications.

# Third-party libraries
This framework uses other libraries such as:
* [gorilla/mux](https://github.com/gorilla/mux)
* [gorilla/websocket](https://github.com/gorilla/websocket)

***Many thanks to the creators of these libs!***

# Examples
You can find some examples in `/winter/examples` dir

# Usage
Simple usage:
```
server := core.NewServer(":<addr>")
server.Get("/", func(ctx *core.Context) core.Response {
    core.MainLogger.Info("Wow, what a fancy logger")
    return core.NewSuccessResponse("I love cookies")
})
server.Start()
```

Creating new routers:
```
// You need to create a new struct that extends Router
type SuperRouter struct {
    *core.Router
}

// and implement the Init method. It'll execute when you "Set" router
func (s *SuperRouter) Init() {
    // Sender(json interface{}) is just an utils function to fast response with json
    s.All("/", core.Sender("Works gud"))
}

func main() {
    server := core.NewServer(":<addr>")
    server.Set("/api", &SuperRouter{})
    server.Start()
}
```

Creating router with more simple way:
```
var superRouter = core.NewRouter(func(r *core.Router) {
    r.All("/", core.Sender("Simple Router")))
})

// ...

// Set it as always
server.Set("/api", superRouter)
```

# TODO
* Documentation
* Winter CLI
* MultipleProtocol/RPC Usage (WS, Twirp, gRPC, HTTP/2)

# License
MIT
