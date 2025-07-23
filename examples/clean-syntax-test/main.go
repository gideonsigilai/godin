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

	log.Println("Starting clean syntax test on :8082")
	log.Println("Visit http://localhost:8082 to see your app")
	if err := app.Serve(":8082"); err != nil {
		log.Fatal(err)
	}
}

// This is exactly the syntax you wanted - no widgets. prefix needed!
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

func IncrementHandler(ctx *Context) Widget {
	counter++
	return Text{
		Data: fmt.Sprintf("%d", counter),
	}
}
