package main

import (
	"fmt"
	"log"

	. "github.com/gideonsigilai/godin/pkg/godin"
)

var counter = 0

func main() {
	app := New()

	app.GET("/", HomeHandler)
	app.POST("/increment", IncrementHandler)
	app.POST("/decrement", DecrementHandler)

	log.Println("Starting direct import test on :8081")
	log.Println("Visit http://localhost:8081 to see your app")
	if err := app.Serve(":8081"); err != nil {
		log.Fatal(err)
	}
}

// HomeHandler demonstrates the new direct import style
func HomeHandler(ctx *Context) Widget {
	return Container{
		Child: Column{
			Children: []Widget{
				Text{
					Data: fmt.Sprintf("Count: %d", counter),
					TextStyle: &TextStyle{
						FontSize: &[]float64{24}[0],
					},
				},
				ElevatedButton{
					Child: Text{Data: "+"},
					OnPressed: func() {
						counter++
					},
				},
			},
		},
	}
}

// IncrementHandler increments the counter
func IncrementHandler(ctx *Context) Widget {
	counter++
	return Text{
		Data: fmt.Sprintf("%d", counter),
	}
}

// DecrementHandler decrements the counter
func DecrementHandler(ctx *Context) Widget {
	counter--
	return Text{
		Data: fmt.Sprintf("%d", counter),
	}
}
