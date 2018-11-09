# Winter Framework
Winter is framework made for developing scalable web applications.

# Examples
You can find some examples in /winter/examples dir

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
// You need to create new struct that extends Router
type SuperRouter struct {
    *core.Router
}

// and implement Init method. It'll execute when you "Set" router
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

Creating router more simple way:
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
