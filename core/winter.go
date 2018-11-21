package core

import (
	"bufio"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"reflect"
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

	main_room_path	= "main"

	shutdown_timeout = time.Second * 2

	winter_guard_tag 	= "winter"
	option_unrequired	= "omitempty"
	option_extends		= "extends"
	option_max			= "max"
	option_min			= "min"
	option_contains		= "contains"
	option_deprecate	= "deprecate"

	winter_logo = " __     __     __     __   __     ______   ______     ______   \n" +
		"/\\ \\  _ \\ \\   /\\ \\   /\\ \"-.\\ \\   /\\__  _\\ /\\  ___\\   /\\  == \\  \n" +
		"\\ \\ \\/ \".\\ \\  \\ \\ \\  \\ \\ \\-.  \\  \\/_/\\ \\/ \\ \\  __\\   \\ \\  __<  \n" +
		" \\ \\__/\".~\\_\\  \\ \\_\\  \\ \\_\\\\\"\\_\\    \\ \\_\\  \\ \\_____\\  \\ \\_\\ \\_\\\n" +
		"  \\/_/   \\/_/   \\/_/   \\/_/ \\/_/     \\/_/   \\/_____/   \\/_/ /_/ \n"
)

var (
	MainLogger 			= NewLogger("main")
	RequestLogger 		= NewLogger("request")
	RouterLogger 		= NewLogger("router")
	WebSocketLogger 	= NewLogger("ws")

	docEnabled			= false
	docPath				= ""
)

// server.go
type (
	IServer interface {
		Start()
		StartTLS(certPath, keyPath string)
		SetRootRouter(router interface{})

		OnStart(onStart func(addr string))
		OnError(onErr func(err error))
		OnShutdown(onShutdown func(signal string))

		EnableDoc(docPath string)
	}
	Server struct {
		*Router

		Addr string

		Debug bool
		GracefulShutdown bool

		ShutdownTimeout time.Duration

		Headers ServerHeaders
		CORS ServerCORSHeaders

		NativeServer *http.Server

		onStart func(addr string)
		onError func(err error)
		onShutdown func(signal string)
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

		All(path string, resolver Resolver) IRouterConfig
		Get(path string, resolver Resolver) IRouterConfig
		Put(path string, resolver Resolver) IRouterConfig
		Post(path string, resolver Resolver) IRouterConfig
		Delete(path string, resolver Resolver) IRouterConfig
		Handle(path string, resolver Resolver, methods ...string) IRouterConfig
		HandleWebSocket(path string, ws *WebSocket)

		Use(resolver MiddlewareResolver)
	}
	Router struct {
		*RouterConfig
		Errors 		*ErrorMap
		mux 		*mux.Router
		routerTree 	map[string]interface{}
	}

	SimpleRouter struct {
		*Router
		Init func(r *Router)
	}

	IRouterConfig interface {
		Doc(explanation string) IRouterConfig
		DocPath(path string) IRouterConfig
		Guard(config interface{}, passIfError ...bool) IRouterConfig
	}
	RouterConfig struct {
		guard 			interface{}
		guardConfig		*GuardConfigMap
		passIfError 	bool
	}

	GuardConfigMap map[string]*GuardConfig

	GuardConfig struct {
		FieldName 		string
		Type			reflect.Kind
		Omitempty		bool
		Options			map[string]interface{}
		Node 			*GuardConfigMap
	}
)

type (
	Doc struct {
		Protocol 	string					`json:"protocol,omitempty"`
		Method 		string 					`json:"method,omitempty"`
		Link 		string 					`json:"omitempty"`
		Params 		map[string]string 		`json:"params,omitempty"`
		Body		map[string]interface{}	`json:"body,omitempty"`
		handlerDoc 	bool
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
		GetGuardBody() (interface{}, error)

		SendError(err *Error)
		SendSuccess(message interface{})
		SendResponse(status int, message interface{})
	}
	Context struct {
		guard		interface{}
		Response 	http.ResponseWriter
		Request 	*http.Request
		TrackTime 	func() time.Duration
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


// ws.go
type (
	IWebSocket interface {
		GetResolver() Resolver
		GetHandlerFunc() http.HandlerFunc
	}
	WebSocket struct {
		Headers	 	http.Header
		upgrader 	*websocket.Upgrader
		resolver	WebSocketResolver
		connection  *Connection
	}

	WebSocketResolver func(conn *Connection)

	IConnection interface {
		OnMessage(onMessage func(message Message))
		OnError(onError func(err error))
		OnClose(onClose func())
		Send(mt int, data []byte)
		JSON(mess interface{})
		On(event string, resolver EventResolver)
		Emit(event string, data ...interface{})
		Room(name string, resolver ...WebSocketResolver) *Connection
	}
	Connection struct {
		ExtendedConnection 	*websocket.Conn
		RoomPath 			string
		events				EventHub
		onMessage 			func(message Message)
		onError				func(err error)
		onClose 			func()
	}

	Message struct {
		Type int
		Data []byte
	}

	EventMessage struct {
		Room 	string		`json:"room"`
		Event 	string		`json:"event"`
		Payload	interface{}	`json:"payload"`
	}

	EventResolver 	func(data interface{})
	EventHub 		map[string]map[string]*EventResolver
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
