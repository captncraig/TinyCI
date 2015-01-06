package main

import (
	"fmt"
	"github.com/captncraig/github-webhooks"
	"net/http"
)

func main() {
	gitHooks := webhooks.WebhookListener{}
	gitHooks.OnPush = githubHook
	http.HandleFunc("/gh", gitHooks.GetHttpListener())
	http.ListenAndServe(":4567", nil)
}

func githubHook(event *webhooks.PushEvent, _ *webhooks.WebhookContext) {
	fmt.Println(event.Repository.Name)
}
