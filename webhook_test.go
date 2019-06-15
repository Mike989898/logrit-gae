package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
  "gopkg.in/src-d/go-billy.v4/memfs"
)

const repo_to_watch string = "https://github.com/Mike989898/logrit-gae-sample.git"

func set_up_http_server(ctx context.Context) {
  e := echo.New()
  //add the hook to service POST request from github when a push happens
  e.Add(http.MethodPost, "/webhook", handleWebhook)
  e.Use(middleware.Logger())
  e.Logger.Fatal(e.Start(":1234"))
}

func main() {
  // Create a new context
  ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)
  fs := memfs.New()

  go set_up_http_server(ctx)
  //set up local host port forwarding
  exec.CommandContext(ctx, "/bin/sh", "ssh -R logrit.serveo.net:80:localhost:1234 serveo.net")
  //make a small change to the render repo
  r, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
    URL: repo,
    Progress: os.Stdout,
  })
  if err != nill {
    panic(err)
  }
  _, err := fs.Create("testing.txt")
  if err != nil {
    return err
  }
  w, err := r.Worktree()
  w.Add("testing.txt")
  err = r.Push(&git.PushOptions{})
  if err != nill {
    panic(err)
  }
  //clean up our repo
  err = fs.Remove("testing.txt")
  if err != nil {
    return err
  }
  //cancel background procs
  cancel()
}
