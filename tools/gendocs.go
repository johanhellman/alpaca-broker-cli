package main

import (
	"log"
	
	"github.com/johanhellman/alpaca-broker-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd(), "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
