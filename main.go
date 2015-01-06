package main

import (
	"fmt"
	"github.com/captncraig/github-webhooks"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var scriptDir string

func init() {
	scriptDir = os.Getenv("TINYCI-SCRIPT-DIR")

	if scriptDir == "" {
		var err error
		scriptDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		scriptDir = path.Join(scriptDir, "scripts")
	}
}

func main() {

	gitHooks := webhooks.WebhookListener{}
	gitHooks.OnPush = githubHook
	http.HandleFunc("/gh", gitHooks.GetHttpListener())
	http.ListenAndServe(":4567", nil)
}

func githubHook(event *webhooks.PushEvent, _ *webhooks.WebhookContext) {
	repo := strings.Replace(event.Repository.FullName, "/", ".", -1)
	refPath := strings.Split(event.Ref, "/")
	ref := refPath[len(refPath)-1]
	anyBranchScript := path.Join(scriptDir, fmt.Sprintf("%s.sh", repo))
	fmt.Println(anyBranchScript)
	singleBranchScript := path.Join(scriptDir, fmt.Sprintf("%s:%s.sh", repo, ref))
	fmt.Println(singleBranchScript)
}
