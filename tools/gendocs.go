package main

import (
	"log"
	
	"github.com/johanhellman/alpaca-broker-cli/cmd/broker"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(broker.RootCmd(), "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
