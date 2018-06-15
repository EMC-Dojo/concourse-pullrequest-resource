package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("usage: %s <source directory>\n", os.Args[0])
		os.Exit(1)
	}

}
