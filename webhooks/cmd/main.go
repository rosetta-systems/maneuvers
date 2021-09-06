package main

import (
	"log"
	"net/http"

	"github.com/paganjoshua/control.jynx.dev/webhooks/pkg/ansible"
)

func main() {
	http.HandleFunc("/", homeHandler)
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var ansi ansible.Runner
	params := ansible.Params{
		User:	"jynx",
		Log:	"goansi.log",
		Cmd:	ansible.Cmd{
			AnsibleCommand: "playbook",
			Args: []string{	"-u", "jynx", "--tags", "deploy", "--check", "workstation.yml" },
		},
	}

	ansi = ansible.New(params)
	ansi.Run()
	w.WriteHeader(http.StatusOK)
}
