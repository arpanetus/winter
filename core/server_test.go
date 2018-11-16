package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func startServer(addr string) *Server {
	server := NewServer(addr)
	server.OnStart(func(addr string) {})
	server.OnError(func(err error) {})
	server.OnShutdown(func(signal string) {})
	go server.Start()
	return server
}

func TestNewServer(t *testing.T) {
	addr := "localhost:5050"
	jsonValue := "Main /"
	resp := Response{}

	server := startServer(addr)

	server.Get("/", Sender(NewResponse(http.StatusOK, jsonValue)))

	time.Sleep(time.Second)
	res, err := http.Get("http://" + addr)
	if err != nil {
		t.Error("Error trying to request server main route", err)
		return
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("Error trying to read value from body", err)
		return
	}

	err = json.Unmarshal(b, &resp)

	if err != nil || resp.Message != jsonValue {
		t.Error("Error trying to Unmarshal json", err)
		return
	}

	server.NativeServer.Shutdown(nil)
}

func TestServerCORS(t *testing.T) {
	addr := "localhost:5051"
	jsonValue := "Main /"
	contentTypeHeader := "Content-Type"
	contentType := "application/json"

	server := startServer(addr)

	server.Headers.Add(contentTypeHeader, contentType)

	server.CORS.Origin("*")
	server.CORS.Methods([]string{"GET"})

	server.Get("/", Sender(NewResponse(http.StatusOK, jsonValue)))

	time.Sleep(time.Second)
	res, err := http.Get("http://" + addr)
	if err != nil {
		t.Error("Error trying to request server main route", err)
		return
	}

	val, ok := res.Header[cors_methods]
	if !ok || val[0] != server.CORS.Get(cors_methods) {
		t.Error("No cors header represented")
		return
	}

	val, ok = res.Header[contentTypeHeader]
	if !ok || val[0] != server.Headers.Get(contentTypeHeader) {
		t.Error("No Content-Type header represented")
		return
	}

	server.NativeServer.Shutdown(nil)
}
