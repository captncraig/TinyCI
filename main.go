package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/captncraig/github-webhooks"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
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
	http.HandleFunc("/dh", dockerHubHook)
	go gitPoll()
	http.ListenAndServe(":4567", nil)
}

func githubHook(event *webhooks.PushEvent, _ *webhooks.WebhookContext) {
	repo := strings.Replace(event.Repository.FullName, "/", ".", -1)
	refPath := strings.Split(event.Ref, "/")
	ref := refPath[len(refPath)-1]
	runScriptIfExists(fmt.Sprintf("gh-%s", repo))
	runScriptIfExists(fmt.Sprintf("gh-%s~%s", repo, ref))
}

type DockerHubData struct {
	CallbackUrl string `json:"callback_url"`
	Repository  struct {
		Name string `json:"repo_name"`
	} `json:"repository"`
}

func dockerHubHook(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	data := DockerHubData{}
	json.Unmarshal(body, &data)
	runScriptIfExists(fmt.Sprintf("dh-%s", strings.Replace(data.Repository.Name, "/", ".", -1)))

	go func() {
		//wait for incoming request to finish before calling callback. The test on dockerhub is more consistent this way.
		time.Sleep(15 * time.Millisecond)
		resp, err := http.Post(data.CallbackUrl, "application/json", bytes.NewBuffer([]byte(`{"state": "success"}`)))
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(resp.StatusCode)
	}()
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

func gitPoll() {

}
