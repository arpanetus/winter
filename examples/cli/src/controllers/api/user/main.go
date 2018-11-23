package main

import "github.com/steplems/winter/core"

var (
	Controller = core.NewController()
	_ = Controller.Post("/creat", func(ctx *core.Context) core.Response {
		return core.NewSuccessResponse("Hui")
	})
	_ = Controller.Get("/", func(ctx *core.Context) core.Response {
		return core.NewSuccessResponse("Hui")
	})
)
