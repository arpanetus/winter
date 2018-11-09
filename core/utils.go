package core

import (
	"runtime"
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

func Trace() (string, int, string) {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	return file, line, f.Name()
}

func Sender(json interface{}) Resolver {
	return func(ctx *Context) Response {
		ctx.JSON(json)
		return NullResponse()
	}
}
