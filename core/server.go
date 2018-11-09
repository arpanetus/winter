package core

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
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
		onError: func(err error) {
			MainLogger.Err(err)
		},
		onStart: func(addr string) {
			fmt.Println(winter_logo)
			MainLogger.Info("Your server is running on " + addr)
		},
		onShutdown: func(err error) {
			if err != nil {
				MainLogger.Err("Server shutdown with error", err)
				return
			}
			MainLogger.Warn("Server shutdown")
		},
	}
}

func (s *Server) OnStart(onStart func(addr string)) {
	s.onStart = onStart
}

func (s *Server) OnError(onErr func(err error)) {
	s.onError = onErr
}

func (s *Server) OnShutdown(onShutdown func(err error)) {
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

func (s *Server) gracefulShutdown(useTLS bool, certPath, keyPath string) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go s.start(useTLS, certPath, keyPath)

	<-stop

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)

	s.onShutdown(s.NativeServer.Shutdown(ctx))
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

func (s *Server) processRouterByDefault() *mux.Router {
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

func (s *Server) loggingMiddleware(ctx *MiddlewareContext) {
	ctx.Next()
	requestLogger.Info(ctx.Request.Method, ctx.Request.RequestURI,
		"ms -", float32(ctx.TrackTime().Nanoseconds()) / float32(1000000))
}

func (s *Server) headerSetterMiddleware(ctx *MiddlewareContext) {
	for key, value := range s.Headers.GetMap() {
		ctx.Response.Header().Set(key, value)
	}
	ctx.Next()
}

func (s *Server) corsMiddleware(ctx *MiddlewareContext) {
	for key, value := range s.CORS.GetMap() {
		ctx.Response.Header().Set(key, value)
	}
	ctx.Next()
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
