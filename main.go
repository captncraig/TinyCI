package main

import (
	"fmt"
	"github.com/captncraig/github-webhooks"
	"net/http"
	"strings"
)

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
	fmt.Printf("%s.sh\n", repo)
	fmt.Printf("%s:%s.sh\n", repo, ref)

}
