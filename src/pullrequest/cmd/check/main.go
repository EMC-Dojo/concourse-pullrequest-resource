package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	r "pullrequest/resource"
	"pullrequest/util"
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	req := r.NewCheckRequest()
	util.InputRequest(&req)

	github, err := r.NewGithubClient(req.Source)
	if err != nil {
		log.Fatalf("contstructing github client: %+v", err)
	}

	command := r.NewCheckCommand(github)
	resp, err := command.Run(req)
	if err != nil {
		log.Fatalf("running command: %+v", err)
	}

	respInterface := make([]interface{}, len(resp))
	for i, version := range resp {
		respInterface[i] = version
	}

	util.OutputArrayResponse(respInterface)
}
