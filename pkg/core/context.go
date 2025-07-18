package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

// Context provides request context and utilities for handlers
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	App      *App
	vars     map[string]string
	params   map[string]interface{}
	handlers map[string]Handler
}

// NewContext creates a new request context
func NewContext(w http.ResponseWriter, r *http.Request, app *App) *Context {
	return &Context{
		Request:  r,
		Response: w,
		App:      app,
		vars:     mux.Vars(r),
		params:   make(map[string]interface{}),
		handlers: make(map[string]Handler),
	}
}

// Param gets a URL parameter by name
func (c *Context) Param(name string) string {
	return c.vars[name]
}

// ParamInt gets a URL parameter as integer
func (c *Context) ParamInt(name string) (int, error) {
	value := c.vars[name]
	return strconv.Atoi(value)
}

// Query gets a query parameter by name
func (c *Context) Query(name string) string {
	return c.Request.URL.Query().Get(name)
}

// QueryInt gets a query parameter as integer
func (c *Context) QueryInt(name string) (int, error) {
	value := c.Query(name)
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

// FormValue gets a form value by name
func (c *Context) FormValue(name string) string {
	return c.Request.FormValue(name)
}

// JSON binds request body to a struct
func (c *Context) JSON(v interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// Set stores a value in the context
func (c *Context) Set(key string, value interface{}) {
	c.params[key] = value
}

// Get retrieves a value from the context
func (c *Context) Get(key string) interface{} {
	return c.params[key]
}

// GetString retrieves a string value from the context
func (c *Context) GetString(key string) string {
	if value, ok := c.params[key].(string); ok {
		return value
	}
	return ""
}

// GetInt retrieves an integer value from the context
func (c *Context) GetInt(key string) int {
	if value, ok := c.params[key].(int); ok {
		return value
	}
	return 0
}

// GetBool retrieves a boolean value from the context
func (c *Context) GetBool(key string) bool {
	if value, ok := c.params[key].(bool); ok {
		return value
	}
	return false
}

// Header gets a request header value
func (c *Context) Header(name string) string {
	return c.Request.Header.Get(name)
}

// SetHeader sets a response header
func (c *Context) SetHeader(name, value string) {
	c.Response.Header().Set(name, value)
}

// IsHTMX returns true if the request is from HTMX
func (c *Context) IsHTMX() bool {
	return c.Header("HX-Request") == "true"
}

// HTMXTarget returns the HTMX target element ID
func (c *Context) HTMXTarget() string {
	return c.Header("HX-Target")
}

// HTMXTrigger returns the HTMX trigger element ID
func (c *Context) HTMXTrigger() string {
	return c.Header("HX-Trigger")
}

// HTMXCurrentURL returns the current URL from HTMX
func (c *Context) HTMXCurrentURL() string {
	return c.Header("HX-Current-URL")
}

// Method returns the HTTP method
func (c *Context) Method() string {
	return c.Request.Method
}

// URL returns the request URL
func (c *Context) URL() string {
	return c.Request.URL.String()
}

// UserAgent returns the user agent string
func (c *Context) UserAgent() string {
	return c.Header("User-Agent")
}

// RemoteAddr returns the client's remote address
func (c *Context) RemoteAddr() string {
	return c.Request.RemoteAddr
}

// IsSecure returns true if the request is HTTPS
func (c *Context) IsSecure() bool {
	return c.Request.TLS != nil
}

// Redirect sends a redirect response
func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.Response, c.Request, url, code)
}

// Error sends an error response
func (c *Context) Error(message string, code int) {
	http.Error(c.Response, message, code)
}

// WriteJSON writes a JSON response
func (c *Context) WriteJSON(data interface{}) error {
	c.SetHeader("Content-Type", "application/json")
	return json.NewEncoder(c.Response).Encode(data)
}

// WriteHTML writes an HTML response
func (c *Context) WriteHTML(html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Response.Write([]byte(html))
}

// TemplateData represents data for template rendering
type TemplateData struct {
	Title   string
	Content template.HTML // Use template.HTML to prevent escaping
	CSS     template.CSS  // Use template.CSS for CSS content
	JS      template.JS   // Use template.JS for JavaScript content
}

// RenderTemplate renders a widget using the base HTML template
func (c *Context) RenderTemplate(widget Widget, title string) {
	// Render the widget content
	content := widget.Render(c)

	// Prepare template data
	data := TemplateData{
		Title:   title,
		Content: template.HTML(content),
	}

	// Load and parse the base template
	templatePath := filepath.Join("pkg", "templates", "base.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// Fallback to simple HTML if template fails
		c.WriteHTML(content)
		return
	}

	// Execute template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		// Fallback to simple HTML if template execution fails
		c.WriteHTML(content)
		return
	}

	// Write the complete HTML document
	c.WriteHTML(buf.String())
}

// WriteText writes a plain text response
func (c *Context) WriteText(text string) {
	c.SetHeader("Content-Type", "text/plain")
	c.Response.Write([]byte(text))
}

// RegisterHandler registers a handler function and returns a unique ID
func (c *Context) RegisterHandler(handler Handler) string {
	// Generate a unique ID for the handler
	handlerID := fmt.Sprintf("handler_%d", len(c.handlers))
	c.handlers[handlerID] = handler

	// Register the handler with the app's router
	c.App.router.HandleFunc("/handlers/"+handlerID, func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, c.App)
		widget := handler(ctx)
		if widget != nil {
			html := widget.Render(ctx)
			ctx.WriteHTML(html)
		}
	}).Methods("GET", "POST", "PUT", "DELETE")

	return handlerID
}
