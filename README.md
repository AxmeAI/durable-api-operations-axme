# Durable API Operations

Your API returns `202 Accepted` and a `job_id`. Your client polls every 2 seconds. Redis tracks state. A cleanup cron catches orphaned jobs. You build this same stack every time you need an async endpoint.

**There is a better way.** Submit an intent, subscribe to lifecycle events via SSE, and let the platform handle delivery, retries, and timeouts.

> **Alpha** · Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
> [cloud.axme.ai](https://cloud.axme.ai) · [hello@axme.ai](mailto:hello@axme.ai)

---

## The Problem

Every async API reinvents the same infrastructure:

```
Client → POST /reports/generate → 202 Accepted + job_id
Client → GET /reports/status/abc  (every 2s... 50 requests)
Client → GET /reports/status/abc  → 200 { "status": "done", "url": "..." }
```

What you end up building:
- **Polling endpoint** — clients hammer your server for status updates
- **Redis job table** — tracks state, needs TTLs, needs cleanup
- **Webhook callbacks** — signature verification, retry logic, dead letter queues
- **Orphan detection** — cron job to catch jobs that never completed
- **Timeout handling** — client gives up, server keeps working, no coordination

---

## The Solution: Intent Lifecycle

```
Client → send_intent("generate report") → intent_id
Client → observe(intent_id)  ← real-time SSE stream

Server: CREATED → SUBMITTED → DELIVERED → IN_PROGRESS → COMPLETED
```

One call to submit. One stream to watch. The platform handles retries, timeouts, and delivery guarantees.

---

## Quick Start

### Python

```bash
pip install axme
export AXME_API_KEY="your-key"   # Get one: axme login
```

```python
from axme import AxmeClient, AxmeClientConfig
import os

client = AxmeClient(AxmeClientConfig(api_key=os.environ["AXME_API_KEY"]))

# FastAPI endpoint returns intent_id instead of job_id
intent_id = client.send_intent({
    "intent_type": "report.generate.v1",
    "to_agent": "agent://myorg/production/report-service",
    "payload": {"report_type": "quarterly", "format": "pdf"},
})

print(f"Submitted: {intent_id}")

# Subscribe to lifecycle events via SSE — no polling
for event in client.observe(intent_id):
    print(f"  [{event['status']}] {event['event_type']}")
    if event["status"] in ("COMPLETED", "FAILED", "TIMED_OUT"):
        break
```

### TypeScript

```bash
npm install @axme/axme
```

```typescript
import { AxmeClient } from "@axme/axme";

const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

const intentId = await client.sendIntent({
  intentType: "report.generate.v1",
  toAgent: "agent://myorg/production/report-service",
  payload: { reportType: "quarterly", format: "pdf" },
});

console.log(`Submitted: ${intentId}`);

const result = await client.waitFor(intentId);
console.log(`Done: ${result.status}`);
```

---

## More Languages

Full implementations in all 5 languages:

| Language | Directory | Install |
|----------|-----------|---------|
| [Python](python/) | `python/` | `pip install axme` |
| [TypeScript](typescript/) | `typescript/` | `npm install @axme/axme` |
| [Go](go/) | `go/` | `go get github.com/AxmeAI/axme-sdk-go` |
| [Java](java/) | `java/` | Maven Central: `ai.axme:axme-sdk` |
| [.NET](dotnet/) | `dotnet/` | `dotnet add package Axme.Sdk` |

---

## Before / After

### Before: Polling API + Redis + Cleanup (200+ lines)

```python
@app.post("/reports/generate")
async def generate(req):
    job_id = str(uuid4())
    redis.set(f"job:{job_id}", json.dumps({"status": "pending", "created": time.time()}))
    queue.enqueue(generate_report, job_id, req.params)
    return JSONResponse({"job_id": job_id, "status": "pending"}, status_code=202)

@app.get("/reports/status/{job_id}")
async def status(job_id):
    state = redis.get(f"job:{job_id}")
    if not state:
        raise HTTPException(404, "Job not found or expired")
    return json.loads(state)

@app.post("/webhooks/report-done")
async def webhook_callback(req):
    if not verify_signature(req.headers["x-signature"], req.body):
        raise HTTPException(401)
    # Process... retry on failure... dead letter if exhausted...
    pass

# Plus: Redis TTL cleanup, orphan cron, timeout sweeper, webhook retry queue...
```

### After: AXME Durable Operations (15 lines)

```python
from axme import AxmeClient, AxmeClientConfig

client = AxmeClient(AxmeClientConfig(api_key=os.environ["AXME_API_KEY"]))

intent_id = client.send_intent({
    "intent_type": "report.generate.v1",
    "to_agent": "agent://myorg/production/report-service",
    "payload": {"report_type": "quarterly", "format": "pdf"},
})

for event in client.observe(intent_id):
    print(f"[{event['status']}] {event['event_type']}")
    if event["status"] in ("COMPLETED", "FAILED", "TIMED_OUT"):
        break
```

No Redis. No polling endpoint. No webhook handler. No cleanup cron. No orphan detector.

---

## How It Works

```
┌────────────┐  POST /generate   ┌────────────────┐   deliver    ┌──────────────┐
│            │ returns intent_id │                │ ──────────>  │              │
│   API      │ ────────────────> │   AXME Cloud   │              │   Report     │
│   Client   │                   │   (platform)   │ <─ resume()  │   Service    │
│            │ <─ observe(SSE) ─ │                │  with result │   (agent)    │
│            │  lifecycle events │   retries,     │              │              │
└────────────┘                   │   timeouts     │              │  processes   │
                                 │                │              │  the report  │
                                 └────────────────┘              └──────────────┘
```

1. FastAPI endpoint receives request, submits an **intent** via AXME SDK
2. Returns `intent_id` to the client immediately (no 202 + job_id)
3. Platform **delivers** the intent to the report service agent
4. Report service processes and **resumes** with result
5. Client **observes** lifecycle events via SSE — no polling needed

---

## Run the Full Example

### Prerequisites

```bash
# Install CLI (one-time)
curl -fsSL https://raw.githubusercontent.com/AxmeAI/axme-cli/main/install.sh | sh
# Open a new terminal, or run the "source" command shown by the installer

# Log in
axme login

# Install Python SDK
pip install axme
```

### Terminal 1 - submit the intent

```bash
axme scenarios apply scenario.json
# Note the intent_id in the output
```

### Terminal 2 - start the agent

Get the agent key after scenario apply:

```bash
# macOS
cat ~/Library/Application\ Support/axme/scenario-agents.json | grep -A2 report-api-demo

# Linux
cat ~/.config/axme/scenario-agents.json | grep -A2 report-api-demo
```

Then run the agent in your language of choice:

```bash
# Python (SSE stream listener)
AXME_API_KEY=<agent-key> python agent.py

# TypeScript (SSE stream listener, requires Node 20+)
cd typescript && npm install
AXME_API_KEY=<agent-key> npx tsx agent.ts

# Go (SSE stream listener)
cd go && go run ./cmd/agent/

# Java (processes a single intent by ID)
cd java/agent && mvn compile
AXME_API_KEY=<agent-key> mvn -q exec:java -Dexec.mainClass="Agent" -Dexec.args="<step-intent-id>"

# .NET (processes a single intent by ID)
cd dotnet/agent && dotnet run -- <step-intent-id>
```

### Verify

```bash
axme intents get <intent_id>
# lifecycle_status: COMPLETED
```

---

## Related

- [AXME](https://github.com/AxmeAI/axme) — project overview
- [AXP Spec](https://github.com/AxmeAI/axme-spec) — open Intent Protocol specification
- [AXME Examples](https://github.com/AxmeAI/axme-examples) — 20+ runnable examples across 5 languages
- [AXME CLI](https://github.com/AxmeAI/axme-cli) — manage intents, agents, scenarios from the terminal

---

Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
