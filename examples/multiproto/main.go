package main

import "github.com/steplems/winter/core"

func main() {
	server := core.NewServer(":5549")
	server.Start()
}
