package core

import (
	"net/http"
	"testing"
	"time"
)

func TestRouter_Set(t *testing.T) {
	addr := "localhost:5052"
	reqAddr := "http://" + addr

	r1 := NewRouter(func(r *Router) {
		r.Get("/get", func(ctx *Context) Response {
			return NewSuccessResponse("/r1/get")
		})
	})
	r := NewRouter(func(r *Router) {
		r.Set("/r1", r1)
	})

	server := startServer(addr)
	defer server.NativeServer.Shutdown(nil)

	server.Set("/r", r)

	time.Sleep(time.Second)
	res, err := http.Get(reqAddr + "/r/r1/get")
	if err != nil {
		t.Error("Error trying to request /r/r1/get route")
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Error("Router didn't set correctly")
		return
	}

	server.NativeServer.Shutdown(nil)
}

func TestRouter_SetHandler(t *testing.T) {
	addr := "localhost:5053"
	reqAddr := "http://" + addr

	mux := http.NewServeMux()
	mux.HandleFunc("/", SendResponse(NewSuccessResponse("")))

	r := NewRouter(func(r *Router) {
		r.SetHandler("/mux", mux)
	})

	server := startServer(addr)
	server.SetRootRouter(r)

	time.Sleep(time.Second)

	res, err := http.Get(reqAddr + "/mux")
	if err != nil {
		t.Error("Error trying to request new handler route", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Error("Error trying to execute http mux router")
		return
	}

	server.NativeServer.Shutdown(nil)
}

func TestRouterDoc(t *testing.T)  {
	addr := "localhost:5060"
	reqAddr := "http://" + addr
	server := startServer(addr)

	server.EnableDoc("/@doc")

	server.Set("/api", NewRouter(func(r *Router) {
		r.Doc(NewDoc("Main API router"))

		r.Get("/users", Sender(NewSuccessResponse([]struct{
			name string
			email string
		}{
			{
				"alex",
				"alex@mail.com",
			},
		}))).Doc(NewDoc("Returns registered users"))
	}))

	time.Sleep(time.Second)

	server.NativeServer.Shutdown()
}
