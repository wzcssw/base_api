package mw

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

// TemplateRenderer 模板渲染器
type TemplateRenderer struct {
	Templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	// fmt.Printf("\nI am in Render \n")
	// fmt.Printf("\ndata is : %v\n", data)
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
		viewContext["csrfKey"] = c.Get("csrf")
		viewContext["csrfToken"] = echo.HeaderXCSRFToken
	}
	// fmt.Printf("\nafter data Reverse,name is %s\n", name)
	result := t.Templates.ExecuteTemplate(w, name, data)
	// fmt.Printf("\nfinish render!!!!!!!!!,result is %s\n", result)
	return result
}
