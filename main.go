package main

import (
	"log"

	"github.com/d0kur0/webm-api/cmd"
)

var version = "unknown"

func main() {
	cmd.Version = version
	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
