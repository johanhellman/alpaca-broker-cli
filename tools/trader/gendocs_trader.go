package main

import "log"
import trader "github.com/johanhellman/alpaca-broker-cli/cmd/trader"
import "github.com/spf13/cobra/doc"

func main() {
	err := doc.GenMarkdownTree(trader.RootCmd, "./docs/trader")
	if err != nil {
		log.Fatal(err)
	}
}
