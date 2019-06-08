package main

import (
	"html/template"
	"io"
	"net/http"
	"github.com/labstack/echo"
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Static("/static", "static")
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/hello", func(c echo.Context) error {
		return c.Render(http.StatusOK, "hello", "World")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
