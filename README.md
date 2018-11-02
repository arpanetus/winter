# Winter Framework
Winter is framework made for developing scalable web applications.

# Examples
You can find some examples in /winter/examples dir

# Usage
Simple usage:
```
server := core.NewServer(":<addr>")
server.Get("/", func(ctx *core.Context) {
    core.MainLogger.Info("Wow, what a fancy logger")
    ctx.JSON("I love cookies")
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

# TODO
* Documentation
* Winter CLI
* MultipleProtocol/RPC Usage (WS, Twirp, gRPC, HTTP/2)

# License
MIT
