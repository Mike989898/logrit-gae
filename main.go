package main

import (
	"html/template"
	"io"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"path/filepath"
	"os"
	"gopkg.in/src-d/go-git.v4"
	"strings"
)

const repo_to_watch string = "https://github.com/chathaway-codes/logrit-gae-sample.git"

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	_, err := git.PlainClone("sample/", false, &git.CloneOptions{
		URL: repo_to_watch,
		Progress: os.Stdout,
	})
	if err != nil {
		panic(err)
	}

	root := "gs://logrit-gae-test/sample/"
	// Setup templates
	t := &Template{
		templates: template.Must(template.ParseGlob(root + "_templates/*.html")),
	}
	e.Renderer = t

	// Load folders, except for /_templates and /.git
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() && !strings.HasPrefix(path, root + ".git") &&
			!strings.HasPrefix(path, root + "_templates") {
			e.GET(path[len(root):] + "/", find_renderer(path))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))
}
