package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

// Todo represents a todo item
type Todo struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoApp represents the todo application state
type TodoApp struct {
	todos  []Todo
	nextID int
	filter string // "all", "active", "completed"
}

var app *TodoApp

func main() {
	// Initialize todo app state
	app = &TodoApp{
		todos:  []Todo{},
		nextID: 1,
		filter: "all",
	}

	// Create Godin app
	godinApp := core.New()

	// Routes
	godinApp.GET("/", HomeHandler)
	godinApp.GET("/api/todos", TodosHandler)
	godinApp.POST("/api/todos", CreateTodoHandler)
	godinApp.PUT("/api/todos/{id}", UpdateTodoHandler)
	godinApp.DELETE("/api/todos/{id}", DeleteTodoHandler)
	godinApp.GET("/api/todos/filter/{filter}", FilterTodosHandler)

	log.Println("Starting Todo App on :8080")
	log.Fatal(godinApp.Serve(":8080"))
}

// HomeHandler renders the main todo application
func HomeHandler(ctx *core.Context) widgets.Widget {
	return widgets.Container{
		Style: "max-width: 800px; margin: 0 auto; padding: 20px; font-family: Arial, sans-serif;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// Header
				widgets.AppBar{
					Title: widgets.Text{
						Data: "Todo App",
						TextStyle: &widgets.TextStyle{
							Color: widgets.ColorWhite,
						},
					},
					Style: "background: #2196F3; color: white; padding: 16px; border-radius: 8px;",
				},

				// Spacer
				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Add todo form
				TodoForm(),

				// Spacer
				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Filter buttons
				FilterButtons(),

				// Spacer
				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Todo list
				TodoList(),

				// Spacer
				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Footer stats
				TodoStats(),
			},
		},
	}
}

// TodoForm renders the form to add new todos
func TodoForm() widgets.Widget {
	return widgets.Card{
		Style: "padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);",
		Child: widgets.Row{
			Children: []widgets.Widget{
				widgets.Expanded{
					Child: widgets.TextField{
						ID: "todo-input",
						Decoration: &widgets.InputDecoration{
							HintText: "What needs to be done?",
						},
						Style: "border: 1px solid #ddd; padding: 12px; border-radius: 4px; font-size: 16px;",
					},
				},
				// Spacer
				widgets.SizedBox{Width: &[]float64{10}[0]},
				widgets.Button{
					Text: "Add Todo",
					Type: "primary",
					OnClick: func() {
						// Handle add todo
					},
				},
			},
		},
	}
}

// FilterButtons renders the filter buttons
func FilterButtons() widgets.Widget {
	return widgets.Card{
		Style: "padding: 16px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);",
		Child: widgets.Row{
			MainAxisAlignment: widgets.MainAxisAlignmentCenter,
			Children: []widgets.Widget{
				FilterButton("all", "All"),
				// Spacer
				widgets.SizedBox{Width: &[]float64{10}[0]},
				FilterButton("active", "Active"),
				// Spacer
				widgets.SizedBox{Width: &[]float64{10}[0]},
				FilterButton("completed", "Completed"),
			},
		},
	}
}

// FilterButton renders a single filter button
func FilterButton(filter, label string) widgets.Widget {
	style := "padding: 8px 16px; border: 1px solid #ddd; border-radius: 4px; cursor: pointer; background: white;"
	if app.filter == filter {
		style = "padding: 8px 16px; border: 1px solid #2196F3; border-radius: 4px; cursor: pointer; background: #2196F3; color: white;"
	}

	return widgets.Button{
		Text:  label,
		Style: style,
		OnClick: func() {
			// Handle filter
		},
	}
}

// TodoList renders the list of todos
func TodoList() widgets.Widget {
	filteredTodos := getFilteredTodos()

	var todoItems []widgets.Widget
	for _, todo := range filteredTodos {
		todoItems = append(todoItems, TodoItem(todo))
	}

	if len(todoItems) == 0 {
		todoItems = append(todoItems, widgets.Text{
			Data:  "No todos found",
			Style: "text-align: center; color: #666; padding: 40px;",
		})
	}

	return widgets.Card{
		Style: "border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); overflow: hidden;",
		Child: widgets.Container{
			ID: "todo-list",
			Child: widgets.ListView{
				Children: todoItems,
				Style:    "min-height: 200px;",
			},
		},
	}
}

// TodoItem renders a single todo item
func TodoItem(todo Todo) widgets.Widget {
	textStyle := &widgets.TextStyle{}
	if todo.Completed {
		textStyle.Decoration = widgets.TextDecorationLineThrough
		textStyle.Color = widgets.Color("#666")
	}

	return widgets.ListTile{
		Style: "border-bottom: 1px solid #eee; padding: 12px 16px;",
		Leading: widgets.Checkbox{
			Value: &todo.Completed,
			OnChanged: func(checked bool) {
				// Handle todo completion toggle
			},
		},
		Title: widgets.Text{
			Data:      todo.Text,
			TextStyle: textStyle,
		},
		Trailing: widgets.Button{
			Text: "Delete",
			Type: "danger",
			OnClick: func() {
				// Handle todo deletion
			},
		},
	}
}

// TodoStats renders todo statistics
func TodoStats() widgets.Widget {
	total := len(app.todos)
	completed := 0
	for _, todo := range app.todos {
		if todo.Completed {
			completed++
		}
	}
	active := total - completed

	return widgets.Card{
		Style: "padding: 16px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);",
		Child: widgets.Row{
			MainAxisAlignment: widgets.MainAxisAlignmentSpaceBetween,
			Children: []widgets.Widget{
				widgets.Text{
					Data: fmt.Sprintf("Total: %d", total),
					TextStyle: &widgets.TextStyle{
						FontWeight: widgets.FontWeightBold,
					},
				},
				widgets.Text{
					Data: fmt.Sprintf("Active: %d", active),
					TextStyle: &widgets.TextStyle{
						Color:      widgets.Color("#2196F3"),
						FontWeight: widgets.FontWeightBold,
					},
				},
				widgets.Text{
					Data: fmt.Sprintf("Completed: %d", completed),
					TextStyle: &widgets.TextStyle{
						Color:      widgets.Color("#4CAF50"),
						FontWeight: widgets.FontWeightBold,
					},
				},
			},
		},
	}
}

// API Handlers

// TodosHandler returns the todo list
func TodosHandler(ctx *core.Context) widgets.Widget {
	return TodoList()
}

// CreateTodoHandler creates a new todo
func CreateTodoHandler(ctx *core.Context) widgets.Widget {
	text := ctx.FormValue("text")
	if text == "" {
		return TodoList()
	}

	todo := Todo{
		ID:        app.nextID,
		Text:      text,
		Completed: false,
		CreatedAt: time.Now(),
	}

	app.todos = append(app.todos, todo)
	app.nextID++

	// Broadcast update via WebSocket
	ctx.App.WebSocket().Broadcast("todos", app.todos)

	return TodoList()
}

// UpdateTodoHandler updates a todo
func UpdateTodoHandler(ctx *core.Context) widgets.Widget {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return TodoList()
	}

	completed := ctx.FormValue("completed") == "true"

	for i, todo := range app.todos {
		if todo.ID == id {
			app.todos[i].Completed = completed
			break
		}
	}

	// Broadcast update via WebSocket
	ctx.App.WebSocket().Broadcast("todos", app.todos)

	return TodoList()
}

// DeleteTodoHandler deletes a todo
func DeleteTodoHandler(ctx *core.Context) widgets.Widget {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return TodoList()
	}

	for i, todo := range app.todos {
		if todo.ID == id {
			app.todos = append(app.todos[:i], app.todos[i+1:]...)
			break
		}
	}

	// Broadcast update via WebSocket
	ctx.App.WebSocket().Broadcast("todos", app.todos)

	return TodoList()
}

// FilterTodosHandler filters todos
func FilterTodosHandler(ctx *core.Context) widgets.Widget {
	filter := ctx.Param("filter")
	app.filter = filter
	return TodoList()
}

// Helper functions

func getFilteredTodos() []Todo {
	switch app.filter {
	case "active":
		var active []Todo
		for _, todo := range app.todos {
			if !todo.Completed {
				active = append(active, todo)
			}
		}
		return active
	case "completed":
		var completed []Todo
		for _, todo := range app.todos {
			if todo.Completed {
				completed = append(completed, todo)
			}
		}
		return completed
	default:
		return app.todos
	}
}
