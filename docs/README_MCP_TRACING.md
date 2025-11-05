# MCP Tracing - Quick Start

This is a quick reference for MCP tracing in the GitHub Actions Utils CLI.

## What is MCP Tracing?

Automatic instrumentation for MCP (Model Context Protocol) tool calls that creates Sentry spans following OpenTelemetry conventions.

## Quick Example

```go
// Register a tool with Sentry tracing
mcp.AddTool(server, &mcp.Tool{
    Name:        "my_tool",
    Description: "My awesome tool",
}, WithSentryTracing("my_tool", m.handleMyTool))
```

That's it! The tool is now automatically traced.

## What Gets Captured?

Every tool call creates a span with:

- **Operation**: `mcp.server`
- **Name**: `tools/call my_tool`
- **Attributes**:
  - Method name (`mcp.method.name`)
  - Tool name (`mcp.tool.name`)
  - All arguments (`mcp.request.argument.*`)
  - Result metadata (`mcp.tool.result.*`)
  - Transport info (`mcp.transport`, `network.transport`)
  - Error status (`mcp.tool.result.is_error`)

## Example Span in Sentry

```json
{
  "op": "mcp.server",
  "description": "tools/call get_action_parameters",
  "status": "ok",
  "data": {
    "mcp.method.name": "tools/call",
    "mcp.tool.name": "get_action_parameters",
    "mcp.transport": "stdio",
    "network.transport": "pipe",
    "mcp.request.argument.actionref": "actions/checkout@v4",
    "mcp.tool.result.is_error": false,
    "mcp.tool.result.content_count": 1
  }
}
```

## Benefits

✅ **Zero boilerplate**: One wrapper function, that's it\
✅ **Type-safe**: Uses Go generics\
✅ **Automatic**: Arguments and results captured automatically\
✅ **Standard**: Follows OpenTelemetry MCP conventions\
✅ **Production-ready**: Error capture, proper span lifecycle

## Documentation

- **User Guide**: See `docs/MCP_TRACING.md`
- **Analysis**: See `docs/ANALYSIS_SENTRY_MCP_INTEGRATION.md`
- **Implementation**: See `docs/IMPLEMENTATION_SUMMARY.md`

## Viewing Traces

In Sentry:

1. Go to **Performance** → **Traces**
2. Filter by operation: `mcp.server`
3. See tool calls with full context

## Disable Telemetry

```bash
export TELEMETRY_ENABLED=false
```

## Questions?

See the full documentation in `docs/MCP_TRACING.md` or check the implementation in `internal/cli/mcp/sentry.go`.
