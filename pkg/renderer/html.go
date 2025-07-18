package renderer

import (
	"fmt"
	"strings"
)

// HTMLRenderer handles HTML generation for widgets
type HTMLRenderer struct {
	indentLevel int
	buffer      strings.Builder
}

// NewHTMLRenderer creates a new HTML renderer
func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{
		indentLevel: 0,
	}
}

// RenderElement renders an HTML element with attributes and content
func (r *HTMLRenderer) RenderElement(tag string, attributes map[string]string, content string, selfClosing bool) string {
	r.buffer.Reset()
	
	// Opening tag
	r.writeIndent()
	r.buffer.WriteString("<")
	r.buffer.WriteString(tag)
	
	// Attributes
	for key, value := range attributes {
		if value != "" {
			r.buffer.WriteString(fmt.Sprintf(` %s="%s"`, key, escapeHTML(value)))
		}
	}
	
	if selfClosing {
		r.buffer.WriteString(" />")
		return r.buffer.String()
	}
	
	r.buffer.WriteString(">")
	
	// Content
	if content != "" {
		if strings.Contains(content, "\n") {
			r.buffer.WriteString("\n")
			r.indentLevel++
			r.writeIndentedContent(content)
			r.indentLevel--
			r.writeIndent()
		} else {
			r.buffer.WriteString(content)
		}
	}
	
	// Closing tag
	r.buffer.WriteString("</")
	r.buffer.WriteString(tag)
	r.buffer.WriteString(">")
	
	return r.buffer.String()
}

// RenderContainer renders a container element with children
func (r *HTMLRenderer) RenderContainer(tag string, attributes map[string]string, children []string) string {
	r.buffer.Reset()
	
	// Opening tag
	r.writeIndent()
	r.buffer.WriteString("<")
	r.buffer.WriteString(tag)
	
	// Attributes
	for key, value := range attributes {
		if value != "" {
			r.buffer.WriteString(fmt.Sprintf(` %s="%s"`, key, escapeHTML(value)))
		}
	}
	
	r.buffer.WriteString(">")
	
	// Children
	if len(children) > 0 {
		r.buffer.WriteString("\n")
		r.indentLevel++
		
		for _, child := range children {
			r.writeIndentedContent(child)
		}
		
		r.indentLevel--
		r.writeIndent()
	}
	
	// Closing tag
	r.buffer.WriteString("</")
	r.buffer.WriteString(tag)
	r.buffer.WriteString(">")
	
	return r.buffer.String()
}

// RenderText renders escaped text content
func (r *HTMLRenderer) RenderText(text string) string {
	return escapeHTML(text)
}

// RenderRawHTML renders raw HTML content (use with caution)
func (r *HTMLRenderer) RenderRawHTML(html string) string {
	return html
}

// BuildAttributes builds HTML attributes from a map
func (r *HTMLRenderer) BuildAttributes(attrs map[string]string) string {
	if len(attrs) == 0 {
		return ""
	}
	
	var parts []string
	for key, value := range attrs {
		if value != "" {
			parts = append(parts, fmt.Sprintf(`%s="%s"`, key, escapeHTML(value)))
		}
	}
	
	return " " + strings.Join(parts, " ")
}

// AddClass adds a CSS class to existing classes
func (r *HTMLRenderer) AddClass(existing, newClass string) string {
	if existing == "" {
		return newClass
	}
	if newClass == "" {
		return existing
	}
	return existing + " " + newClass
}

// MergeAttributes merges multiple attribute maps
func (r *HTMLRenderer) MergeAttributes(attrs ...map[string]string) map[string]string {
	result := make(map[string]string)
	
	for _, attrMap := range attrs {
		for key, value := range attrMap {
			if key == "class" {
				result[key] = r.AddClass(result[key], value)
			} else {
				result[key] = value
			}
		}
	}
	
	return result
}

// writeIndent writes the current indentation
func (r *HTMLRenderer) writeIndent() {
	for i := 0; i < r.indentLevel; i++ {
		r.buffer.WriteString("  ")
	}
}

// writeIndentedContent writes content with proper indentation
func (r *HTMLRenderer) writeIndentedContent(content string) {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i > 0 {
			r.buffer.WriteString("\n")
		}
		if strings.TrimSpace(line) != "" {
			r.writeIndent()
			r.buffer.WriteString(line)
		}
	}
	r.buffer.WriteString("\n")
}

// escapeHTML escapes HTML special characters
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// Common HTML generation helpers

// RenderDiv renders a div element
func (r *HTMLRenderer) RenderDiv(attributes map[string]string, content string) string {
	return r.RenderElement("div", attributes, content, false)
}

// RenderSpan renders a span element
func (r *HTMLRenderer) RenderSpan(attributes map[string]string, content string) string {
	return r.RenderElement("span", attributes, content, false)
}

// RenderButton renders a button element
func (r *HTMLRenderer) RenderButton(attributes map[string]string, content string) string {
	return r.RenderElement("button", attributes, content, false)
}

// RenderInput renders an input element
func (r *HTMLRenderer) RenderInput(attributes map[string]string) string {
	return r.RenderElement("input", attributes, "", true)
}

// RenderImg renders an img element
func (r *HTMLRenderer) RenderImg(attributes map[string]string) string {
	return r.RenderElement("img", attributes, "", true)
}

// RenderLink renders an anchor element
func (r *HTMLRenderer) RenderLink(attributes map[string]string, content string) string {
	return r.RenderElement("a", attributes, content, false)
}

// RenderForm renders a form element
func (r *HTMLRenderer) RenderForm(attributes map[string]string, children []string) string {
	return r.RenderContainer("form", attributes, children)
}

// RenderList renders a ul or ol element
func (r *HTMLRenderer) RenderList(listType string, attributes map[string]string, items []string) string {
	var children []string
	for _, item := range items {
		children = append(children, r.RenderElement("li", nil, item, false))
	}
	return r.RenderContainer(listType, attributes, children)
}

// RenderTable renders a table with headers and rows
func (r *HTMLRenderer) RenderTable(attributes map[string]string, headers []string, rows [][]string) string {
	var children []string
	
	// Header
	if len(headers) > 0 {
		var headerCells []string
		for _, header := range headers {
			headerCells = append(headerCells, r.RenderElement("th", nil, header, false))
		}
		thead := r.RenderContainer("thead", nil, []string{
			r.RenderContainer("tr", nil, headerCells),
		})
		children = append(children, thead)
	}
	
	// Body
	if len(rows) > 0 {
		var bodyRows []string
		for _, row := range rows {
			var cells []string
			for _, cell := range row {
				cells = append(cells, r.RenderElement("td", nil, cell, false))
			}
			bodyRows = append(bodyRows, r.RenderContainer("tr", nil, cells))
		}
		tbody := r.RenderContainer("tbody", nil, bodyRows)
		children = append(children, tbody)
	}
	
	return r.RenderContainer("table", attributes, children)
}
