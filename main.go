package main

import (
	"html/template"
	"io"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"strings"
	"fmt"
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
	fs := memfs.New()

	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: repo_to_watch,
		Progress: os.Stdout,
	})
	if err != nil {
		panic(err)
	}

	// Setup templates
	tmpls, err := fs.ReadDir("/_templates/")
	templates := template.New("")
	if err == nil {
		var read_dir func(os.FileInfo, string) error
		read_dir = func(f os.FileInfo, parent string) error {
			new_path := parent + f.Name()
			fmt.Printf("Processing %v\n", new_path)
			if f.IsDir() {
				files, err := fs.ReadDir(new_path)
				if err != nil {
					return err
				}
				for _, file := range files {
					err = read_dir(file, new_path)
					if err != nil {
						return err
					}
				}
			} else {
				file, err := fs.Open("/_templates/" + new_path)
				if err != nil {
					return err
				}
				buff := make([]byte, f.Size())
				file.Read(buff)
				_, err = templates.New(new_path).Parse(string(buff))
				if err != nil {
					return err
				}
			}
			return nil
		}
		for _, file := range tmpls {
			err = read_dir(file, "")
		}
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf(templates.DefinedTemplates())
	t := &Template{
		templates: templates,
	}
	e.Renderer = t

	// Load folders, except for /_templates and /.git
	var walk_dirs func(path string) error
	walk_dirs = func(path string) error {
		info, err := fs.Stat(path)
		if err != nil {
			return err
		}
		if info.IsDir() && !strings.HasPrefix(path, ".git") &&
			!strings.HasPrefix(path, "_templates") {
			path = "/" + info.Name()
			fmt.Printf("\nScanning %v\n", path)
			e.GET(path + "/", find_renderer(path, fs))
			files, err := fs.ReadDir(path)
			if err != nil {
				return err
			}
			for _, file := range files {
				err := walk_dirs(path + file.Name())
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}
		return nil
	}
	err = walk_dirs("")
	if err != nil {
		panic(err)
	}

	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))
}
