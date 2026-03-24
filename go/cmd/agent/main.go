// Report API agent — Go example.
//
// Generates financial reports and resumes with result.
//
// Usage:
//
//	export AXME_API_KEY="<agent-key>"
//	go run agent.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

const agentAddress = "report-api-demo"

func handleIntent(ctx context.Context, client *axme.Client, intentID string) error {
	intentData, err := client.GetIntent(ctx, intentID, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("get intent: %w", err)
	}

	intent, _ := intentData["intent"].(map[string]any)
	if intent == nil {
		intent = intentData
	}
	payload, _ := intent["payload"].(map[string]any)
	if payload == nil {
		payload = map[string]any{}
	}
	if pp, ok := payload["parent_payload"].(map[string]any); ok {
		payload = pp
	}

	reportType, _ := payload["report_type"].(string)
	if reportType == "" {
		reportType = "unknown"
	}
	format, _ := payload["format"].(string)
	if format == "" {
		format = "pdf"
	}
	year, _ := payload["fiscal_year"].(string)
	if year == "" {
		year = "2025"
	}

	fmt.Printf("  Generating %s report (%s) for FY%s...\n", reportType, format, year)
	time.Sleep(2 * time.Second)

	result := map[string]any{
		"action":       "complete",
		"report_url":   fmt.Sprintf("https://reports.example.com/FY%s-%s.%s", year, reportType, format),
		"pages":        128,
		"generated_at": time.Now().UTC().Format(time.RFC3339),
	}

	_, err = client.ResumeIntent(ctx, intentID, result, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("resume intent: %w", err)
	}
	fmt.Printf("  Report ready: %s\n", result["report_url"])
	return nil
}

func main() {
	apiKey := os.Getenv("AXME_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: AXME_API_KEY not set.")
	}

	client, err := axme.NewClient(axme.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	fmt.Printf("Agent listening on %s...\n", agentAddress)
	fmt.Println("Waiting for intents (Ctrl+C to stop)")

	intents, errCh := client.Listen(ctx, agentAddress, axme.ListenOptions{})

	go func() {
		for err := range errCh {
			log.Printf("Listen error: %v", err)
		}
	}()

	for delivery := range intents {
		intentID, _ := delivery["intent_id"].(string)
		status, _ := delivery["status"].(string)
		if intentID == "" {
			continue
		}
		if status == "DELIVERED" || status == "CREATED" || status == "IN_PROGRESS" {
			fmt.Printf("[%s] Intent received: %s\n", status, intentID)
			if err := handleIntent(ctx, client, intentID); err != nil {
				fmt.Printf("  Error: %v\n", err)
			}
		}
	}
}
