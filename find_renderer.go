package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/labstack/echo"
	"github.com/wellington/go-libsass"
	"net/http"
	"path"
	"log"
	"os"
)

func make_error(err error) (func (echo.Context) error) {
	return func(c echo.Context) error {
		return c.String(http.StatusBadRequest, fmt.Sprintf("%v", err))
	}
}

type TemplateRenderer struct {
	Template string
}

type config struct {
	Render string
	Template TemplateRenderer
}

func find_renderer(dir string, fs billy.Filesystem) (func(echo.Context) error) {
	var conf config
	fname := path.Join(dir, ".template")
	ops, err := toml.DecodeFile(fname, &conf)
	if err != nil {
		fmt.Printf("%v", err)
		return make_error(err);
	}
	if !ops.IsDefined("render") {
		return make_error(fmt.Errorf("Missing render field in %s", fname))
	}
	ret := make_error(fmt.Errorf("Renderer not found: %v", conf.Render))
	switch(conf.Render) {
	case "sass":
		ret = func(c echo.Context) error {
			f, err := os.Open(path.Join(dir, "index.scss"))
			if err != nil {
				return err
			}
			buf := bufio.NewReader(f)
			var b bytes.Buffer
			out_buf := bufio.NewWriter(&b)
			comp, err := libsass.New(out_buf, buf)
			if err != nil {
				return err
			}
			if err := comp.Run(); err != nil {
				log.Fatal(err)
			}
			out_buf.Flush()
			return c.String(http.StatusOK, b.String())
		}
	case "Template":
		ret = func(c echo.Context) error {
			template := "home.html"
			if ops.IsDefined("Template", "Template") {
				template = conf.Template.Template
			}
			return c.Render(http.StatusOK, template, "")
		}
	}
	return ret
}
