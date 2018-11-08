package core

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func NewServer(addr string) *Server {
	return &Server{
		Router: NewCoreRouter(),
		Addr: addr,
		Debug: false,
		Headers: ServerHeaders{
			map[string]string{},
		},
		CORS: ServerCORSHeaders{
			&ServerHeaders{
				map[string]string{},
			},
		},
	}
}

func (s *Server) Start() {
	s.serverInit()

	srv := &http.Server{
		Addr: s.Addr,
		Handler: s.processRouterByDefault(),
	}
	srv.ListenAndServe()
}

func (s *Server) StartTLS(certPath string, keyPath string) {
	s.serverInit()

	srv := &http.Server{
		Addr: s.Addr,
		Handler: s.processRouterByDefault(),
	}
	srv.ListenAndServeTLS(certPath, keyPath)
}

func (s *Server) serverInit() {
	fmt.Println(winter_logo)
	MainLogger.Info("Your server is running on " + s.Addr)
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
