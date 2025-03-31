package main

import (
	"log"
	"os"

	"gogomate/internal/cli"
	"gogomate/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cli := cli.NewCLI(cfg)

	if err := cli.Run(os.Args); err != nil {
		log.Fatalf("failed to run cli: %v", err)
	}
}
