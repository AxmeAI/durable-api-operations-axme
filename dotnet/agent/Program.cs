// Report API agent — .NET example.
//
// Fetches an intent by ID, generates a financial report, and resumes with result.
//
// Usage:
//   export AXME_API_KEY="<agent-key>"
//   dotnet run -- <intent_id>

using Axme.Sdk;
using System.Text.Json.Nodes;

if (args.Length < 1)
{
    Console.Error.WriteLine("Usage: dotnet run -- <intent_id>");
    return 1;
}

var apiKey = Environment.GetEnvironmentVariable("AXME_API_KEY");
if (string.IsNullOrEmpty(apiKey))
{
    Console.Error.WriteLine("Error: AXME_API_KEY not set.");
    return 1;
}

var intentId = args[0];
var client = new AxmeClient(new AxmeClientConfig { ApiKey = apiKey });

Console.WriteLine($"Processing intent: {intentId}");

var intentData = await client.GetIntentAsync(intentId);
var intent = intentData["intent"]?.AsObject() ?? intentData;
var payload = intent["payload"]?.AsObject() ?? new JsonObject();
if (payload["parent_payload"] is JsonObject parentPayload)
{
    payload = parentPayload;
}

var reportType = payload["report_type"]?.ToString() ?? "unknown";
var fmt = payload["format"]?.ToString() ?? "pdf";
var year = payload["fiscal_year"]?.ToString() ?? "2025";

Console.WriteLine($"  Generating {reportType} report ({fmt}) for FY{year}...");
await Task.Delay(2000);

var reportUrl = $"https://reports.example.com/FY{year}-{reportType}.{fmt}";
var result = new JsonObject
{
    ["action"] = "complete",
    ["report_url"] = reportUrl,
    ["pages"] = 128,
    ["generated_at"] = DateTime.UtcNow.ToString("o")
};

await client.ResumeIntentAsync(intentId, result);
Console.WriteLine($"  Report ready: {reportUrl}");
return 0;
