// Durable API operations — .NET example.
//
// Submit a report generation intent, wait for completion
// without polling or webhooks.
//
// Usage:
//   export AXME_API_KEY="your-key"
//   dotnet run

using Axme.Sdk;

var client = new AxmeClient(new AxmeClientConfig
{
    ApiKey = Environment.GetEnvironmentVariable("AXME_API_KEY")!
});

// Submit a durable operation — replaces POST /reports/generate → 202 + job_id
var intentId = await client.SendIntentAsync(new
{
    intent_type = "report.generate.v1",
    to_agent = "agent://myorg/production/report-service",
    payload = new
    {
        report_type = "quarterly",
        format = "pdf",
        quarter = "Q1-2026"
    }
});
Console.WriteLine($"Intent submitted: {intentId}");

// Wait for completion — no polling, no webhooks
var result = await client.WaitForAsync(intentId);
Console.WriteLine($"Final status: {result.Status}");
