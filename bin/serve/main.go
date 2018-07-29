package main

import (
	"log"

	"github.com/halverneus/static-file-server/cli"
)

func main() {
	if err := cli.Execute(); nil != err {
		log.Fatalf("Error: %v\n", err)
	}
}
