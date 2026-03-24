/**
 * Report API agent — TypeScript example.
 *
 * Generates financial reports and resumes with result.
 *
 * Usage:
 *   export AXME_API_KEY="<agent-key>"
 *   npx tsx agent.ts
 */

import { AxmeClient } from "@axme/axme";

const AGENT_ADDRESS = "report-api-demo";

async function handleIntent(client: AxmeClient, intentId: string) {
  const intentData = await client.getIntent(intentId);
  const intent = intentData.intent ?? intentData;
  let payload = intent.payload ?? {};
  if (payload.parent_payload) {
    payload = payload.parent_payload;
  }

  const reportType = payload.report_type ?? "unknown";
  const fmt = payload.format ?? "pdf";
  const year = payload.fiscal_year ?? "2025";

  console.log(`  Generating ${reportType} report (${fmt}) for FY${year}...`);
  await new Promise((r) => setTimeout(r, 2000));

  const result = {
    action: "complete",
    report_url: `https://reports.example.com/FY${year}-${reportType}.${fmt}`,
    pages: 128,
    generated_at: new Date().toISOString(),
  };

  await client.resumeIntent(intentId, result, { ownerAgent: "report-api-demo" });
  console.log(`  Report ready: ${result.report_url}`);
}

async function main() {
  const apiKey = process.env.AXME_API_KEY;
  if (!apiKey) {
    console.error("Error: AXME_API_KEY not set.");
    process.exit(1);
  }

  const client = new AxmeClient({ apiKey });

  console.log(`Agent listening on ${AGENT_ADDRESS}...`);
  console.log("Waiting for intents (Ctrl+C to stop)\n");

  for await (const delivery of client.listen(AGENT_ADDRESS)) {
    const intentId = delivery.intent_id;
    const status = delivery.status;
    if (intentId && ["DELIVERED", "CREATED", "IN_PROGRESS"].includes(status)) {
      console.log(`[${status}] Intent received: ${intentId}`);
      try {
        await handleIntent(client, intentId);
      } catch (e) {
        console.error(`  Error: ${e}`);
      }
    }
  }
}

main().catch(console.error);
