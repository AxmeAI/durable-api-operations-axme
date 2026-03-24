/*
 * Durable API operations — Java example.
 *
 * Submit a report generation intent, wait for completion
 * without polling or webhooks.
 *
 * Usage:
 *   export AXME_API_KEY="your-key"
 *   mvn compile exec:java -Dexec.mainClass="DurableApiOperations"
 */

import ai.axme.sdk.AxmeClient;
import ai.axme.sdk.AxmeClientConfig;
import java.util.Map;

public class DurableApiOperations {
    public static void main(String[] args) throws Exception {
        var client = new AxmeClient(
            AxmeClientConfig.builder()
                .apiKey(System.getenv("AXME_API_KEY"))
                .build()
        );

        // Submit a durable operation — replaces POST /reports/generate → 202 + job_id
        String intentId = client.sendIntent(Map.of(
            "intent_type", "report.generate.v1",
            "to_agent", "agent://myorg/production/report-service",
            "payload", Map.of(
                "report_type", "quarterly",
                "format", "pdf",
                "quarter", "Q1-2026"
            )
        ));
        System.out.println("Intent submitted: " + intentId);

        // Wait for completion — no polling, no webhooks
        var result = client.waitFor(intentId);
        System.out.println("Final status: " + result.getStatus());
    }
}
