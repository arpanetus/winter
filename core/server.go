package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func NewServer(addr string) *Server {
	return &Server{
		Router: NewCoreRouter(),
		Addr: addr,
		Debug: false,
		GracefulShutdown: false,
		ShutdownTimeout: shutdown_timeout,
		NativeServer: &http.Server{
			Addr: addr,
		},
		Headers: ServerHeaders{
			map[string]string{},
		},
		CORS: ServerCORSHeaders{
			&ServerHeaders{
				map[string]string{},
			},
		},
		onError: defaultOnError,
		onStart: defaultOnStart,
		onShutdown: defaultOnShutdown,
	}
}

func defaultOnError(err error)  {
	MainLogger.Err("Server closed with error:", err)
}

func defaultOnStart(addr string)  {
	MainLogger.Info("Your server is running on " + addr)
	fmt.Println(winter_logo)
}

func defaultOnShutdown(signal string)  {
	MainLogger.Warn("Server is shutting down with signal:", signal)
}

func (s *Server) OnStart(onStart func(addr string)) {
	s.onStart = onStart
}

func (s *Server) OnError(onErr func(err error)) {
	s.onError = onErr
}

func (s *Server) OnShutdown(onShutdown func(signal string)) {
	s.onShutdown = onShutdown
}

func (s *Server) Start() {
	s.NativeServer.Handler = s.processRouterByDefault()

	if s.GracefulShutdown {
		s.gracefulShutdown(false, "", "")
	} else {
		s.start(false, "", "")
	}
}

func (s *Server) StartTLS(certPath string, keyPath string) {
	s.NativeServer.Handler = s.processRouterByDefault()

	if s.GracefulShutdown {
		s.gracefulShutdown(true, certPath, keyPath)
	} else {
		s.start(true, certPath, keyPath)
	}
}

func (s *Server) SetRootRouter(router interface{}) {
	s.Set("", router)
}

func (s *Server) gracefulShutdown(useTLS bool, certPath, keyPath string) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go s.start(useTLS, certPath, keyPath)

	s.onShutdown((<-stop).String())

	ctx, _ := context.WithTimeout(context.Background(), shutdown_timeout)
	s.NativeServer.Shutdown(ctx)
	time.Sleep(shutdown_timeout)
	os.Exit(0)
}

func (s *Server) start(useTLS bool, certPath, keyPath string) {
	s.onStart(s.Addr)
	if useTLS {
		if err := s.NativeServer.ListenAndServeTLS(certPath, keyPath); err != nil {
			s.onError(err)
		}
	} else {
		if err := s.NativeServer.ListenAndServe(); err != nil {
			s.onError(err)
		}
	}
}

func (s *Server) processRouterByDefault() http.Handler {
	if len(s.CORS.headerMap) > 0 {
		s.Use(s.corsMiddleware)
	}
	if len(s.Headers.headerMap) > 0 {
		s.Use(s.headerSetterMiddleware)
	}
	if s.Debug {
		s.Use(s.loggingMiddleware)
	}

	return s.GetHandler()
}

func (s *Server) loggingMiddleware(ctx *MiddlewareContext) Response {
	ctx.Next()
	RequestLogger.Info(ctx.Request.Method, ctx.Request.RequestURI,
		"ms -", float32(ctx.TrackTime().Nanoseconds()) / float32(1000000))
	return NullResponse()
}

func (s *Server) headerSetterMiddleware(ctx *MiddlewareContext) Response {
	for key, value := range s.Headers.GetMap() {
		ctx.Response.Header().Set(key, value)
	}
	return ctx.NewNext()
}

func (s *Server) corsMiddleware(ctx *MiddlewareContext) Response {
	for key, value := range s.CORS.GetMap() {
		ctx.Response.Header().Set(key, value)
	}
	return ctx.NewNext()
}

func (s *ServerHeaders) Add(key, value string) {
	s.headerMap[key] = value
}

func (s *ServerHeaders) Get(key string) string {
	return s.headerMap[key]
}

func (s *ServerHeaders) GetMap() map[string]string {
	return s.headerMap
}

func (s *ServerCORSHeaders) Origin(value string) {
	s.Add(cors_origin, value)
}

func (s *ServerCORSHeaders) Credentials(value bool) {
	s.Add(cors_credentials, strconv.FormatBool(value))
}

func (s *ServerCORSHeaders) Methods(value []string) {
	headerValue := s.formatArrayHeader(value)
	s.Add(cors_methods, headerValue)
}

func (s *ServerCORSHeaders) Headers(value []string) {
	headerValue := s.formatArrayHeader(value)
	s.Add(cors_headers, headerValue)
}

func (s *ServerCORSHeaders) formatArrayHeader(value []string) string {
	headerValue := ""
	for _, n := range value {
		if headerValue == "" {
			headerValue = n
			continue
		}

		headerValue = headerValue + ", " + n
	}

	return headerValue
}
