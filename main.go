package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nantapop-kj/go-bulkify/config"
	"github.com/nantapop-kj/go-bulkify/internal/runner"
	"github.com/nantapop-kj/go-bulkify/payload"
)

func main() {
	configFile := flag.String("config", "config.json", "path to JSON config file")
	flag.Parse()

	cfg, err := config.LoadFromFile(*configFile)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cfg.BuildPayload = payload.BuildPayload

	fmt.Printf("🚀 Starting: url=%s  method=%s  workers=%d  total=%d\n", cfg.URL, cfg.Method, cfg.WorkerCount, cfg.TotalRecords)

	successCount, failCount := runner.Run(cfg)

	fmt.Printf("\n🏁 Done! Success: %d | Failed: %d | Total: %d\n", successCount, failCount, cfg.TotalRecords)
}
