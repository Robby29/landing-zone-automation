package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var counter int
var mutex = &sync.Mutex{}

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, strconv.Itoa(counter))
	mutex.Unlock()
}

func main() {
	DownloadLZTemplate()
	ConnectGithub()
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	// http.HandleFunc("/increment", incrementCounter)

	// http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hi")
	// })

	log.Fatal(http.ListenAndServe(":8081", nil))

}

func DownloadLZTemplate() {
	cmd := exec.Command("git", "clone", "https://github.com/Robby29/google-cft-template.git")
	cmd.Dir = "/Users/robinsingh/Documents/Development/copied-github-repos-golang"
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func ConnectGithub() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "<enter access token here>"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// create a new private repository named "foo"
	repo := &github.Repository{
		Name:    github.String("test-repo-golang"),
		Private: github.Bool(true),
	}
	client.Repositories.Create(ctx, "", repo)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Number of repositories ", len(repos))

	//push code to gitrepo
	CopyCodeToGithub("test-repo-golang", "template-1")
}

func CopyCodeToGithub(reponame string, templatename string) {
	fmt.Println("Running copy command")
	cmdCopy := exec.Command("cp", "-pr", templatename, "/Users/robinsingh/Documents/Development/test-repo-golang")
	cmdCopy.Dir = "/Users/robinsingh/Documents/Development/copied-github-repos-golang/google-cft-template"
	if err := cmdCopy.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running git init command")
	cmdGitInit := exec.Command("git", "init", "-b", "main")
	cmdGitInit.Dir = "/Users/robinsingh/Documents/Development/test-repo-golang"
	if err := cmdGitInit.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running git add command")
	cmdGitAdd := exec.Command("git", "add", ".")
	cmdGitAdd.Dir = "/Users/robinsingh/Documents/Development/test-repo-golang"
	if err := cmdGitAdd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running git commit command")
	cmdGitCommit := exec.Command("git", "commit", "-m", `"Initial TF Setup using GCP LZ Automation"`)
	cmdGitCommit.Dir = "/Users/robinsingh/Documents/Development/test-repo-golang"
	if err := cmdGitCommit.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running git remote add command")
	cmdGitRemoteAdd := exec.Command("git", "remote", "add", "origin", "https://github.com/Robby29/test-repo-golang.git")
	cmdGitRemoteAdd.Dir = "/Users/robinsingh/Documents/Development/test-repo-golang"
	if err := cmdGitRemoteAdd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running git push command")
	cmdGitPush := exec.Command("git", "push", "-u", "https://<token>@github.com/Robby29/test-repo-golang.git", "main")
	cmdGitPush.Dir = "/Users/robinsingh/Documents/Development/test-repo-golang"
	if err := cmdGitPush.Run(); err != nil {
		log.Fatal(err)
	}
}
