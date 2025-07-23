package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

func main() {
	app := core.New()

	// Debug endpoint to see raw HTML output
	app.GET("/", func(ctx *core.Context) core.Widget {
		return homePage()
	})

	// Raw HTML endpoint to see what's being generated
	app.GET("/debug-html", func(ctx *core.Context) core.Widget {
		button := widgets.Button{
			Text: "Debug Button",
			Type: "primary",
			OnPressed: func() {
				fmt.Println("DEBUG BUTTON CLICKED!")
			},
		}
		
		html := button.Render(ctx)
		fmt.Printf("Generated HTML: %s\n", html)
		
		return widgets.Container{
			Child: widgets.Text{
				Data: fmt.Sprintf("Generated HTML: %s", html),
			},
		}
	})

	// Simple HTTP handler to test if basic requests work
	http.HandleFunc("/simple-test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("SIMPLE HTTP HANDLER CALLED!")
		w.Write([]byte("Simple handler works!"))
	})

	log.Println("Starting debug server on :8091")
	log.Println("Visit http://localhost:8091/debug-html to see generated HTML")
	log.Println("Visit http://localhost:8091/simple-test to test basic HTTP")
	if err := app.Serve(":8091"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func homePage() core.Widget {
	return widgets.Container{
		Style: "padding: 20px;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				widgets.Text{
					Data: "Debug HTML Output",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{24}[0],
					},
				},
				
				widgets.Button{
					Text: "Test Button",
					Type: "primary",
					OnPressed: func() {
						fmt.Println("TEST BUTTON CLICKED!")
					},
				},
			},
		},
	}
}
