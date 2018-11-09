package core

import (
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
