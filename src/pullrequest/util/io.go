package util

import (
	"encoding/json"
	"log"
	"os"
)

// InputRequest is
func InputRequest(request interface{}) {
	err := json.NewDecoder(os.Stdin).Decode(request)
	if err != nil {
		log.Fatalf("reading request from stdin: %+v", err)
	}
}

// OutputArrayResponse is
func OutputArrayResponse(response []interface{}) {
	err := json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		log.Fatalf("writing response to stdout: %+v", err)
	}
}

// OutputResponse is
func OutputResponse(resp interface{}) {
	err := json.NewEncoder(os.Stdout).Encode(resp)
	if err != nil {
		log.Fatalf("writing response to stdout: %+v", err)
	}
}
