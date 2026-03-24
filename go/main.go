// Durable API operations — Go example.
//
// Submit a report generation intent, wait for completion
// without polling or webhooks.
//
// Usage:
//
//	export AXME_API_KEY="your-key"
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

func main() {
	client, err := axme.NewClient(axme.ClientConfig{
		APIKey: os.Getenv("AXME_API_KEY"),
	})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	// Submit a durable operation — replaces POST /reports/generate -> 202 + job_id
	intentID, err := client.SendIntent(ctx, map[string]any{
		"intent_type": "report.generate.v1",
		"to_agent":    "agent://myorg/production/report-service",
		"report_type": "quarterly",
		"format":      "pdf",
		"quarter":     "Q1-2026",
	}, axme.RequestOptions{})
	if err != nil {
		log.Fatalf("send intent: %v", err)
	}
	fmt.Printf("Intent submitted: %s\n", intentID)

	// Wait for completion — no polling, no webhooks
	result, err := client.WaitFor(ctx, intentID, axme.ObserveOptions{})
	if err != nil {
		log.Fatalf("wait: %v", err)
	}
	fmt.Printf("Final status: %v\n", result["status"])
}
