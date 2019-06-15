package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/google/go-github/v26/github"
)

func handleWebhook(c echo.Context) error{
	r := c.Request()
	//// TODO: Add secret key validation!
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nill {
		return err
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return err
	}
	// check the type of action
	switch event.(type) {
	case *github.PushEvent:
		fmt.Printf("Got the push request!\n")
	default:
		err = fmt.Errorf("unknown event type %s\n", github.WebHookType(r))
	}
	return err
}
