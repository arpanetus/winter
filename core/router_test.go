package core

import (
	"bytes"
	"encoding/json"
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

func TestRouter_Guard(t *testing.T) {
	addr := "localhost:5060"
	reqAddr := "http://" + addr
	server := startServer(addr)

	//server.EnableDoc("/@doc")

	type Contacts struct {
		List []string `json:"list" winter:"max: 25"`
	}

	type User struct {
		Email string `json:"email" winter:"max: 40, min: 6, contains: '@'~'.'"`
		Password string `json:"password" winter:"extends: Email, contains: ''"`
		Username string `json:"username" winter:"max: 50"`
		Contacts Contacts `json:"contacts" winter:"omitempty"`
	}

	server.Set("/api", NewRouter(func(r *Router) {
		r.Post("/new", func(ctx *Context) Response {
			body, err := ctx.GetGuardBody()
			if err != nil {
				t.Log(err.Error())
				return NewErrorResponse(NewError(http.StatusBadRequest, err.Error()))
			}

			t.Log(body)

			return NewSuccessResponse(body)
		}).Guard(User{}, true).Doc("Creates new user")
	}))

	time.Sleep(time.Second)

	b, _ := json.Marshal(User{
		"shit@gmail.com",
		"55114411aAZ",
		"username",
		Contacts{
			[]string{"Godddmaaaamn"},
		},
	})

	res, err := http.Post(reqAddr + "/api/new", "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Error("Error trying to post user to guarded route:", err)
	}

	t.Log(res)

	//	/@doc/api/
	//	{
	//		"http://localhost:5060/api/users": "GET - Returns registered users"
	//		"http://localhost:5060/api/users/new/{id}": {
	//			"explanation": "Creates new user",x
	//			"params": {
	//				"id": "string - user id"
	//			},
	//			"body": {
	//				"name": "string - username"
	//			}
	//		}
	//	}

	server.NativeServer.Shutdown(nil)
}
