package core

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	ansi_prefix = "\x1b["
	ansi_suffix = "m"
	ansi_clear = ansi_prefix + "0" + ansi_suffix

	tag_info = "INFO"
	tag_warn = "WARN"
	tag_error = "ERR!"

	cors = "Access-Control-Allow-"
	cors_origin = cors + "Origin"
	cors_credentials = cors + "Credentials"
	cors_methods = cors + "Methods"
	cors_headers = cors + "Headers"

	router_init_func_name = "Init"

	winter_logo = " __     __     __     __   __     ______   ______     ______   \n" +
		"/\\ \\  _ \\ \\   /\\ \\   /\\ \"-.\\ \\   /\\__  _\\ /\\  ___\\   /\\  == \\  \n" +
		"\\ \\ \\/ \".\\ \\  \\ \\ \\  \\ \\ \\-.  \\  \\/_/\\ \\/ \\ \\  __\\   \\ \\  __<  \n" +
		" \\ \\__/\".~\\_\\  \\ \\_\\  \\ \\_\\\\\"\\_\\    \\ \\_\\  \\ \\_____\\  \\ \\_\\ \\_\\\n" +
		"  \\/_/   \\/_/   \\/_/   \\/_/ \\/_/     \\/_/   \\/_____/   \\/_/ /_/ \n"
)

var (
	MainLogger = NewLogger("main")
	requestLogger = NewLogger("request")
	routerLogger = NewLogger("router")
)

// logger.go
type (
	ILogger interface {
		Log(mess ...interface{})
		Err(mess ...interface{})
		Info(mess ...interface{})
		Warn(mess ...interface{})
	}
	Logger struct {
		Name string
	}
)

// server.go
type (
	IServer interface {
		Start()
		StartTLS()
	}

	Server struct {
		*Router
		Addr string
		Headers ServerHeaders
		CORS ServerCORSHeaders
	}

	ServerConfig struct {
	}

	ServerHeaders struct {
		headerMap map[string]string
	}

	ServerCORSHeaders struct {
		*ServerHeaders
	}
)

// router.go
type (
	IRouter interface {
		GetHandler() *mux.Router

		Set(path string, router interface{})
		SetWebSocket(path string, wsRouter interface{})
		SetHandler(path string, handler http.Handler)

		All(path string, resolver Resolver)
		Get(path string, resolver Resolver)
		Put(path string, resolver Resolver)
		Post(path string, resolver Resolver)
		Delete(path string, resolver Resolver)
		Handle(path string, resolver Resolver, methods ...string)

		Use(resolver MiddlewareResolver)
	}
	Router struct {
		mux *mux.Router
	}
)

// ws.go
type (
	IWebSocketRouter interface {
	}
	WebSocketRouter struct {
	}
)

// context.go
type (
	IContext interface {
		Send(msg []byte)
		JSON(msg interface{})
		Status(code int) *Context
	}
	Context struct {
		Response http.ResponseWriter
		Request *http.Request
		TrackTime func() time.Duration
	}

	IMiddlewareContext interface {
		IContext
		Next()
	}
	MiddlewareContext struct {
		*Context
		Handler http.Handler
	}

	Resolver func(ctx *Context)
	MiddlewareResolver func(ctx *MiddlewareContext)
)
type (
	IWebSocketContext interface {
	}
	WebSocketContext struct {
	}
)

// error.go
type (
	IError interface {
		Send(ctx *Context)
		SetMessage(mess string)
		SetStatus(status int)
	}
	Error struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	BindError struct {
		Code int
		*Error
	}

	IErrorMap interface {
		Get(code int) *Error
		Set(code int, err Error)
	}
	ErrorMap map[int]*Error
)
