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
	client := axme.NewClient(axme.Config{
		APIKey: os.Getenv("AXME_API_KEY"),
	})

	ctx := context.Background()

	// Submit a durable operation — replaces POST /reports/generate → 202 + job_id
	intentID, err := client.SendIntent(ctx, axme.SendIntentRequest{
		IntentType: "report.generate.v1",
		ToAgent:    "agent://myorg/production/report-service",
		Payload: map[string]interface{}{
			"report_type": "quarterly",
			"format":      "pdf",
			"quarter":     "Q1-2026",
		},
	})
	if err != nil {
		log.Fatalf("send intent: %v", err)
	}
	fmt.Printf("Intent submitted: %s\n", intentID)

	// Wait for completion — no polling, no webhooks
	result, err := client.WaitFor(ctx, intentID)
	if err != nil {
		log.Fatalf("wait: %v", err)
	}
	fmt.Printf("Final status: %s\n", result.Status)
}
