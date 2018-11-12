package core

import (
	"bufio"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

const (
	ansi_prefix = "\x1b["
	ansi_suffix = "m"
	ansi_clear = ansi_prefix + "0" + ansi_suffix

	tag_info = "INFO"
	tag_warn = "WARN"
	tag_error = "ERR!"
	tag_note = "NOTE"

	cors = "Access-Control-Allow-"
	cors_origin = cors + "Origin"
	cors_credentials = cors + "Credentials"
	cors_methods = cors + "Methods"
	cors_headers = cors + "Headers"

	router_init_func_name = "Init"

	bad_os = "windows"

	winter_logo = " __     __     __     __   __     ______   ______     ______   \n" +
		"/\\ \\  _ \\ \\   /\\ \\   /\\ \"-.\\ \\   /\\__  _\\ /\\  ___\\   /\\  == \\  \n" +
		"\\ \\ \\/ \".\\ \\  \\ \\ \\  \\ \\ \\-.  \\  \\/_/\\ \\/ \\ \\  __\\   \\ \\  __<  \n" +
		" \\ \\__/\".~\\_\\  \\ \\_\\  \\ \\_\\\\\"\\_\\    \\ \\_\\  \\ \\_____\\  \\ \\_\\ \\_\\\n" +
		"  \\/_/   \\/_/   \\/_/   \\/_/ \\/_/     \\/_/   \\/_____/   \\/_/ /_/ \n"
)

var (
	MainLogger = NewLogger("main")
	RequestLogger = NewLogger("request")
	RouterLogger = NewLogger("router")
)

// server.go
type (
	IServer interface {
		Start()
		StartTLS(certPath, keyPath string)
		SetRootRouter(router *Router)

		OnStart(onStart func(addr string))
		OnError(onErr func(err error))
		OnShutdown(onShutdown func(err error))
	}
	Server struct {
		*Router

		Addr string

		Debug bool
		GracefulShutdown bool

		Headers ServerHeaders
		CORS ServerCORSHeaders

		NativeServer *http.Server

		onStart func(addr string)
		onError func(err error)
		onShutdown func(err error)
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
		Errors *ErrorMap
	}
)

// context.go
type (
	IContext interface {
		Send(msg []byte)
		JSON(msg interface{})
		Status(code int) *Context

		GetParams() map[string]string
		GetParam(key string) (string, bool)
		GetBody(body interface{}) error

		SendError(err *Error)
		SendSuccess(message interface{})
		SendResponse(status int, message interface{})
	}
	Context struct {
		Response http.ResponseWriter
		Request *http.Request
		TrackTime func() time.Duration
	}

	IMiddlewareContext interface {
		IContext
		Next()
		NewNext() Response
	}
	MiddlewareContext struct {
		*Context
		handler http.Handler
	}

	Response struct {
		Status 	int 		`json:"status,omitempty"`
		Message interface{} `json:"message,omitempty"`
	}

	Resolver func(ctx *Context) Response
	MiddlewareResolver func(ctx *MiddlewareContext) Response
)

// error.go
type (
	IError interface {
		Send(ctx *Context)
		SetMessage(mess interface{})
		SetStatus(status int)
	}
	Error struct {
		*Response
	}
	BindError struct {
		Code int
		*Error
	}

	IErrorMap interface {
		Get(code int) *Error
		Set(code int, err *Error)
	}
	ErrorMap map[int]*Error
)

// logger.go
type (
	ILogger interface {
		Log(mess ...interface{})
		Logf(format string, mess ...interface{})
		Err(mess ...interface{})
		Errf(format string, mess ...interface{})
		Info(mess ...interface{})
		Infof(format string, mess ...interface{})
		Warn(mess ...interface{})
		Warnf(format string, mess ...interface{})
		Note(mess ...interface{})
		Notef(format string, mess ...interface{})
	}
	Logger struct {
		Name string
		logIntoFile bool
		logFile *os.File

		fileWriter *bufio.Writer
		writer func(a ...interface{}) (n int, err error)
		writerf func(format string, a ...interface{}) (n int, err error)
		writerln func(a ...interface{}) (n int, err error)
	}
)
