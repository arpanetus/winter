package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	addr := "localhost:5055"
	resAddr := "http://" + addr
	contentType := "application/json"

	server := startServer(addr)

	server.Post("/{id}", func(ctx *Context) Response {
		body := map[string]interface{}{}
		err := ctx.GetBody(&body)
		if err != nil {
			return NewErrorResponse(HTTPErrors.Get(http.StatusBadRequest))
		}

		if body["message"] != "test" {
			return NewErrorResponse(HTTPErrors.Get(http.StatusBadRequest))
		}

		id, ok := ctx.GetParam("id")
		if !ok {
			return NewErrorResponse(HTTPErrors.Get(http.StatusBadRequest))
		}

		if id != "123" {
			return NewErrorResponse(HTTPErrors.Get(http.StatusBadRequest))
		}

		return NewSuccessResponse("")
	})

	time.Sleep(time.Second)

	b, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		"test",
	})
	if err != nil {
		t.Error("Error trying to Marshal json struct to bytes")
		return
	}

	res, err := http.Post(resAddr + "/123", contentType, bytes.NewBuffer(b))
	if err != nil {
		t.Error("Error trying to POST server main route")
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Error("Bad request")
		return
	}

	server.NativeServer.Shutdown(nil)
}

func TestMiddlewareContext(t *testing.T) {
	addr := "localhost:5056"
	resAddr := "http://" + addr
	xHeader := "X-Best-Framework"
	xHeaderValue := "winter"

	server := startServer(addr)

	server.Get("/", Sender(NewSuccessResponse("")))
	server.Use(func(ctx *MiddlewareContext) Response {
		if val, ok := ctx.Request.Header[xHeader]; !ok || val[0] != xHeaderValue {
			return NewErrorResponse(HTTPErrors.Get(http.StatusBadRequest))
		}
		return ctx.NewNext()
	})

	time.Sleep(time.Second)

	req, _ := http.NewRequest(http.MethodGet, resAddr, nil)
	req.Header.Add(xHeader, xHeaderValue)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Error("Error trying to get server main route")
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Error("Bad request")
		return
	}

	server.NativeServer.Shutdown(nil)
}
