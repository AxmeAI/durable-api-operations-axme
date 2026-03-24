"""Report API agent - generates reports and resumes with result."""

import os, sys, time
sys.stdout.reconfigure(line_buffering=True)
from axme import AxmeClient, AxmeClientConfig

AGENT_ADDRESS = "report-api-demo"

def handle_intent(client, intent_id):
    intent_data = client.get_intent(intent_id)
    intent = intent_data.get("intent", intent_data)
    payload = intent.get("payload", {})
    if "parent_payload" in payload:
        payload = payload["parent_payload"]

    report_type = payload.get("report_type", "unknown")
    fmt = payload.get("format", "pdf")
    year = payload.get("fiscal_year", "2025")

    print(f"  Generating {report_type} report ({fmt}) for FY{year}...")
    time.sleep(2)

    result = {
        "action": "complete",
        "report_url": f"https://reports.example.com/FY{year}-{report_type}.{fmt}",
        "pages": 128,
        "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
    }
    client.resume_intent(intent_id, result)
    print(f"  Report ready: {result['report_url']}")

def main():
    api_key = os.environ.get("AXME_API_KEY", "")
    if not api_key:
        print("Error: AXME_API_KEY not set."); sys.exit(1)
    client = AxmeClient(AxmeClientConfig(api_key=api_key))
    print(f"Agent listening on {AGENT_ADDRESS}...")
    print("Waiting for intents (Ctrl+C to stop)\n")
    for delivery in client.listen(AGENT_ADDRESS):
        intent_id = delivery.get("intent_id", "")
        status = delivery.get("status", "")
        if intent_id and status in ("DELIVERED", "CREATED", "IN_PROGRESS"):
            print(f"[{status}] Intent received: {intent_id}")
            try:
                handle_intent(client, intent_id)
            except Exception as e:
                print(f"  Error: {e}")

if __name__ == "__main__":
    main()
