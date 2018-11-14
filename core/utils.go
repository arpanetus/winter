package core

import (
	"encoding/json"
	"net/http"
	"time"
)

func GetTrackedTime(cb func(t time.Duration)) func() {
	start := time.Now()
	return func() {
		cb(time.Since(start))
	}
}

func TrackTime() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		return time.Since(start)
	}
}

func Sender(json interface{}) Resolver {
	return func(ctx *Context) Response {
		ctx.JSON(json)
		return NullResponse()
	}
}

func SendResponse(response Response) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(response.Status)
		json.NewEncoder(res).Encode(response)
	}
}
