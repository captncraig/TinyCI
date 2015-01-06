package main

import (
	"fmt"
	"github.com/captncraig/github-webhooks"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var scriptDir string
var scriptExt = ".sh"

func init() {
	scriptDir = os.Getenv("TINYCI-SCRIPT-DIR")

	if scriptDir == "" {
		var err error
		scriptDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		scriptDir = filepath.Join(scriptDir, "scripts")
	}
	if runtime.GOOS == "windows" {
		scriptExt = ".bat"
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
	runScriptIfExists(fmt.Sprintf("%s", repo))
	runScriptIfExists(fmt.Sprintf("%s:%s", repo, ref))
}

func runScriptIfExists(name string) {
	filename := filepath.Join(scriptDir, name+scriptExt)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("Script does not exist: %s. Skipping\n", filename)
		return
	}
	cmd := exec.Command(filename)
	output, err := cmd.CombinedOutput()
	fmt.Println(output, err)
}
