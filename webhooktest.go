package main

import (
	"context"
	"os"
	"fmt"
	"net/http"
	"time"
	"os/exec"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	gghttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
  "gopkg.in/src-d/go-billy.v4/memfs"
)

const repo_fetch string = "https://github.com/Mike989898/logrit-gae-sample.git"

func set_up_http_server(ctx context.Context) {
  e := echo.New()
  //add the hook to service POST request from github when a push happens
  e.Add(http.MethodPost, "/webhook", handleWebhook)
	e.Use(middleware.Logger())
  e.Logger.Fatal(e.Start(":1234"))
}

func commit_and_push(cmsg string, w *git.Worktree, repo *git.Repository, uname string, pass string) {
	fmt.Println("Commiting and Pushing")
	_, err := w.Commit(cmsg, &git.CommitOptions {
		Author: &object.Signature{
			Name:  "Test Script",
			Email: "test@script.org",
			When:  time.Now(),
		},
	})
	if err != nil {
		panic(err)
	}
	err = repo.Push(&git.PushOptions{
			Auth: &gghttp.BasicAuth{
			Username: uname,
			Password: pass,
		},
	})
	if err != nil {
		panic(err)
	}
}
func main() {
  // Create a new context
	usrname := os.Getenv("GITUSR")
	if len(usrname) == 0 {
		fmt.Println("Please set the enviroment variable for your github username.")
		return
	}
	pass := os.Getenv("GITPASS")
	if len(pass) == 0 {
		fmt.Println("Please set the enviroment variable for your github password.")
		return
	}

  ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)
  fs := memfs.New()
	fmt.Println("Starting Server...")
  go set_up_http_server(ctx)
  defer cancel()
  //set up local host port forwarding
	fmt.Println("Setting up port fowarding")
	//// TODO: Probably make the server name configurable
  exec.CommandContext(ctx, "/bin/sh", "ssh -R logrit.serveo.net:80:localhost:1234 serveo.net")
  //make a small change to the render repo
  repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions {
    URL: repo_fetch,
    Progress: os.Stdout,
  })
  if err != nil {
    panic(err)
  }
	fmt.Println("Creating test file in repo.")
  _, err = fs.Create("testing.txt")
  if err != nil {
    panic(err)
  }
  w, err := repo.Worktree()
  w.Add("testing.txt")
	commit_and_push("Adding test file", w, repo, usrname, pass)
  //clean up our repo
	fmt.Println("Removing test file in repo.")
  err = fs.Remove("testing.txt")
	w.Add("testing.txt")
  if err != nil {
    panic(err)
  }
	commit_and_push("Removing test file", w, repo, usrname, pass)
}
