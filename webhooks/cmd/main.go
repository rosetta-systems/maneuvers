package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	
	"github.com/paganjoshua/control.jynx.dev/webhooks/pkg/ansible"
)

func main() {
	http.HandleFunc("/hybrid/site", siteHandler)
	http.HandleFunc("/hybrid/ctl", ctlHandler)
	http.HandleFunc("/ping", pingHandler)
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func siteHandler(w http.ResponseWriter, r *http.Request) {
	var ansi ansible.Runner
	params := ansible.Params{
		User:	"jynx",
		Log:	"goansi.log",
		Cmd:	ansible.Cmd{
			AnsibleCommand: "playbook",
			Args: []string{	"-u", "jynx", "-i", "production", "--tags", "deploy", "hybrid.yml" },
		},
	}

	ansi = ansible.New(params)
	ansi.Run()
	w.WriteHeader(http.StatusOK)
}

func ctlHandler(w http.ResponseWriter, r *http.Request) {
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
	
	w.WriteHeader(http.StatusOK)

	if strings.Contains(payload.Ref, "main") {
		var ansi ansible.Runner
		params := ansible.Params{
			User:	"jynx",
			Log:	"goansi.log",
			Cmd:	ansible.Cmd{
				AnsibleCommand: "playbook",
				Args: []string{	"-u", "jynx", "-i", "production", "--tags", "update-git", "hybrid.yml" },
			},
		}

		ansi = ansible.New(params)
		ansi.Run()

		params = ansible.Params{
			User:	"jynx",
			Log:	"goansi.log",
			Cmd:	ansible.Cmd{
				AnsibleCommand: "playbook",
				Args: []string{	"-u", "jynx", "-i", "production", "--skip-tags", "deploy,git", "hybrid.yml" },
			},
		}

		ansi = ansible.New(params)
		ansi.Run()
	}
	
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var ansi ansible.Runner
	params := ansible.Params{
		User:	"jynx",
		Log:	"goansi.log",
		Cmd:	ansible.Cmd{
			AnsibleCommand: "playbook",
			Args: []string{	"-u", "jynx", "-i", "production", "--tags", "deploy", "webserver.yml" },
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

	devl := strings.Contains(payload.Ref, "devl")
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

func pingHandler(w http.ResponseWriter, r *http.Request) {
//	var ansi ansible.Runner
//	params := ansible.Params{
//		User:	"jynx",
//		Log:	"goansi.log",
//		Cmd:	ansible.Cmd{
//			AnsibleCommand: "playbook",
//			Args: []string{	"-u", "jynx", "-v", "ping.yml" },
//		},
//	}
//
//	ansi = ansible.New(params)
//	ansi.Run()
	w.Write([]byte("\npong\n\n"))
}
