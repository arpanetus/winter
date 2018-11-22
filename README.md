# Winter Framework
Winter is a framework made for developing scalable web applications.

# Other libraries
This framework uses other libraries such as:
* [gorilla/mux](https://github.com/gorilla/mux)
* [gorilla/websocket](https://github.com/gorilla/websocket)

***Many thanks to the creator of these libs!***

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
    - Init command - Init creates new project with given template name. If the default template is not found, then it finds it in the github repository
    - Run command - Run compiles plugins at the working dir (os.Getwd()) into tmp dir (.winter(?)) and runs app with given config in required winter.go file (watch core.App)
    - Build command - Build compiles app with given config in required winter.go file at the working dir plugins into one executable
    - Create Template Repositories
        - MVC Monolit
        - RPC Service
        - Micro Services
        - FastHTTP with WebSocket
* Server
    - FastHTTP server
    - C/C++ kore/proxygen servers (?)
* MultipleProtocol/RPC Usage (WS, Twirp, gRPC, HTTP/2(?))

# License
MIT
