package main

import (
	"flag"
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
	fmt.Println(scriptDir)
}

func main() {
	flagListen := flag.String("l", ":4567", "interface an port to listen on.")
	flag.Parse()
	gitHooks := webhooks.WebhookListener{}
	gitHooks.OnPush = githubHook
	http.HandleFunc("/gh", gitHooks.GetHttpListener())
	http.ListenAndServe(*flagListen, nil)
}

func githubHook(event *webhooks.PushEvent, _ *webhooks.WebhookContext) {
	repo := strings.Replace(event.Repository.FullName, "/", ".", -1)
	refPath := strings.Split(event.Ref, "/")
	ref := refPath[len(refPath)-1]
	runScriptIfExists(fmt.Sprintf("gh-%s", repo))
	runScriptIfExists(fmt.Sprintf("gh-%s~%s", repo, ref))
}

func runScriptIfExists(name string) {
	filename := filepath.Join(scriptDir, name+scriptExt)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("Script does not exist: %s. Skipping.\n", filename)
		return
	}
	cmd := exec.Command(filename)
	log.Printf("Executing %s...\n", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing %s: %s.", filename, err.Error())
	}
	log.Println(string(output))
}
