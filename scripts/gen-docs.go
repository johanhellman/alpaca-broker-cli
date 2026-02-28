package main

import (
	"log"
	"os"

	broker "github.com/johanhellman/alpaca-broker-cli/cmd/broker"
	trader "github.com/johanhellman/alpaca-broker-cli/cmd/trader"
	"github.com/spf13/cobra/doc"
)

func main() {
	// Generate Broker API Markdown Paths
	brokerDir := "./docs/broker"
	if err := os.MkdirAll(brokerDir, 0755); err != nil {
		log.Fatalf("failed to create broker docs directory: %v", err)
	}

	brokerRoot := broker.RootCmd()
	brokerRoot.DisableAutoGenTag = true
	err := doc.GenMarkdownTree(brokerRoot, brokerDir)
	if err != nil {
		log.Fatalf("failed to generate broker docs: %v", err)
	}
	log.Printf("Successfully generated Broker docs in %s", brokerDir)

	// Generate Trader API Markdown Paths
	traderDir := "./docs/trader"
	if err := os.MkdirAll(traderDir, 0755); err != nil {
		log.Fatalf("failed to create trader docs directory: %v", err)
	}

	traderRoot := trader.RootCmd
	traderRoot.DisableAutoGenTag = true
	err = doc.GenMarkdownTree(traderRoot, traderDir)
	if err != nil {
		log.Fatalf("failed to generate trader docs: %v", err)
	}
	log.Printf("Successfully generated Trader docs in %s", traderDir)
}
