package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	
	"github.com/paganjoshua/control.jynx.dev/webhooks/pkg/ansible"
)

func main() {
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/devl", devlHandler)
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var ansi ansible.Runner
	params := ansible.Params{
		User:	"jynx",
		Log:	"goansi.log",
		Cmd:	ansible.Cmd{
			AnsibleCommand: "playbook",
			Args: []string{	"-u", "jynx", "--tags", "deploy", "webserver.yml" },
		},
	}

	ansi = ansible.New(params)
	ansi.Run()
	w.WriteHeader(http.StatusOK)
}

func devlHandler (w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var payload struct {
		Ref string `json:"ref"`
	}
	err := dec.Decode(&payload)
	if err != nil {
		http.Error(w, "Error: Could Not Read JSON Body", http.StatusBadRequest)
		log.Println(err)
		return
	}
	re := regexp.MustCompile("devl")
	devl := re.MatchString(payload.Ref)
	if devl {
		updateImages()
	}
	w.WriteHeader(http.StatusOK)
}

func updateImages() {
	var ansi ansible.Runner
	params := ansible.Params{
		User:	"jynx",
		Log:	"goansi.log",
		Cmd:	ansible.Cmd{
			AnsibleCommand: "playbook",
			Args: []string{	"-u", "jynx", "registry.yml" },
		},
	}

	ansi = ansible.New(params)
	ansi.Run()
}
