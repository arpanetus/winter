package core

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"time"
)

func NewRouter(init func(r *Router)) interface{} {
	return &struct {
		*Router
		Init func(r *Router)
	}{Init: init}
}

func NewCoreRouter() *Router {
	return &Router{
		mux: mux.NewRouter(),
		Errors: NewErrorMap(),
	}
}

func (r *Router) GetHandler() *mux.Router {
	return r.mux
}

func (r *Router) Get(path string, resolver Resolver) {
	r.Handle(path, resolver, http.MethodGet)
}

func (r *Router) Put(path string, resolver Resolver) {
	r.Handle(path, resolver, http.MethodPut)
}

func (r *Router) Post(path string, resolver Resolver) {
	r.Handle(path, resolver, http.MethodPost)
}

func (r *Router) Delete(path string, resolver Resolver) {
	r.Handle(path, resolver, http.MethodDelete)
}

func (r *Router) All(path string, resolver Resolver) {
	r.Handle(path, resolver)
}

func (r *Router) Handle(path string, resolver Resolver, methods ...string) {
	handlerFunc := r.mux.HandleFunc(path, r.resolver(resolver))
	if len(methods) > 0 {
		handlerFunc.Methods(methods...)
	}
}

func (r *Router) Use(middlewareResolver MiddlewareResolver) {
	r.mux.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			profiler := TrackTime()
			middlewareResolver(r.getMiddlewareContext(res, req, handler, profiler))
		})
	})
}

func (r *Router) Set(path string, router interface{}) {
	routerPrefix := r.mux.PathPrefix(path).Subrouter()
	newPrefixedRouter := &Router{routerPrefix, NewErrorMap()}

	routerValue := reflect.ValueOf(router).Elem()
	field := routerValue.FieldByName("Router")

	if field.IsValid() {
		field.Set(reflect.ValueOf(newPrefixedRouter))
	} else {
		return
	}

	routerType := reflect.TypeOf(router)
	method, ok := routerType.MethodByName(router_init_func_name)
	if !ok {
		simpleMethod := routerValue.FieldByName(router_init_func_name)

		if simpleMethod != reflect.Zero(reflect.TypeOf(simpleMethod)).Interface() {
			simpleMethod.Call([]reflect.Value{reflect.ValueOf(newPrefixedRouter)})
		} else {
			routerLogger.Warn("Missing Init method in router!")
		}
		return
	}

	method.Func.Call([]reflect.Value{reflect.ValueOf(router)})
}

func (r *Router) SetHandler(path string, handler http.Handler) {
	r.mux.Handle(path, handler)
}

func (r *Router) resolver(resolver Resolver) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		profiler := TrackTime()

		response := resolver(r.getContext(res, req, profiler))

		if response != (Response{}) {
			res.WriteHeader(response.Status)
			json.NewEncoder(res).Encode(response)
		}
	}
}

func (r *Router) getContext(res http.ResponseWriter, req *http.Request, executionTracker func() time.Duration) *Context {
	return &Context{
		Request: req,
		Response: res,
		TrackTime: executionTracker,
	}
}

func (r *Router) getMiddlewareContext(res http.ResponseWriter, req *http.Request, handler http.Handler, executionTracker func() time.Duration) *MiddlewareContext {
	return &MiddlewareContext{
		Context: r.getContext(res, req, executionTracker),
		Handler: handler,
	}
}
