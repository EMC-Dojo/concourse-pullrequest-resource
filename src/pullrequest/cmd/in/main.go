package main

import (
	"encoding/json"
	"log"
	"os"

	r "pullrequest/resource"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("usage: %s <source directory>\n", os.Args[0])
		os.Exit(1)
	}

	req := r.NewInRequest()
	inputRequest(&req)

	destDir := os.Args[1]

	github, err := r.NewGithubClient(req.Source)
	if err != nil {
		log.Fatalf("constructing github client: %+v", err)
	}

	command := r.NewInCommand(github, os.Stderr)
	resp, err := command.Run(destDir, req)
	if err != nil {
		log.Fatalf("running command: %+v", err)
	}

	outputResponse(resp)
}

func inputRequest(req *r.InRequest) {
	err := json.NewDecoder(os.Stdin).Decode(req)
	if err != nil {
		log.Fatalf("reading request from stdin: %+v", err)
	}
}

func outputResponse(resp r.InResponse) {
	err := json.NewEncoder(os.Stdout).Encode(resp)
	if err != nil {
		log.Fatalf("writing response to stdout: %+v", err)
	}
}
