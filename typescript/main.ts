/**
 * Durable API operations — TypeScript example.
 *
 * Submit a report generation intent, subscribe to SSE lifecycle events,
 * receive completion without polling or webhooks.
 *
 * Usage:
 *   npm install @axme/axme
 *   export AXME_API_KEY="your-key"
 *   npx tsx main.ts
 */

import { AxmeClient } from "@axme/axme";

async function main() {
  const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

  // Submit a durable operation — replaces POST /reports/generate → 202 + job_id
  const intentId = await client.sendIntent({
    intentType: "report.generate.v1",
    toAgent: "agent://myorg/production/report-service",
    payload: {
      reportType: "quarterly",
      format: "pdf",
      quarter: "Q1-2026",
    },
  });
  console.log(`Intent submitted: ${intentId}`);

  // Wait for completion — no polling, no webhooks
  const result = await client.waitFor(intentId);
  console.log(`Final status: ${result.status}`);
}

main().catch(console.error);
