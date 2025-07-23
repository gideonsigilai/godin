package main

import (
	. "github.com/gideonsigilai/godin/pkg/godin"
)

func main() {
	app := New()

	// Enable WebSocket for real-time state updates
	app.WebSocket().Enable("/ws")

	// Routes
	app.GET("/", HomeHandler)

}

func HomeHandler(ctx *Context) Widget {
	return Container{
		Child: Text{
			Data: "Hello World Init ",
		},
	}
}
