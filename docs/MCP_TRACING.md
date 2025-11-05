# MCP Tracing with Sentry

This document describes the MCP (Model Context Protocol) tracing integration with Sentry for the GitHub Actions Utils CLI.

## Overview

The MCP Server integration automatically instruments tool calls with Sentry spans, following the [OpenTelemetry MCP Semantic Conventions](https://github.com/open-telemetry/semantic-conventions/pull/2083). This provides comprehensive observability for MCP tool execution, including:

- Automatic span creation for tool calls
- Detailed attributes following MCP semantic conventions
- Error capture and correlation
- Tool result tracking

## Implementation

The implementation is based on the Sentry JavaScript SDK's MCP integration, adapted for Go. Key files:

- `internal/cli/mcp/sentry.go` - Tracing wrapper and attribute extraction
- `internal/cli/mcp/server.go` - Tool registration with tracing

## Usage

### Wrapping a Tool Handler

Use the `WithSentryTracing` wrapper when registering tools:

```go
mcp.AddTool(server, &mcp.Tool{
    Name:        "my_tool",
    Description: "Does something useful",
}, WithSentryTracing("my_tool", m.handleMyTool))
```

The wrapper:

1. Creates a span for the tool execution
2. Sets MCP-specific attributes
3. Captures tool arguments
4. Tracks results and errors
5. Reports to Sentry

### Example

See `internal/cli/mcp/server.go` for a complete example:

```go
func (m *MCPServer) RegisterTools(server *mcp.Server) {
    // Register get_action_parameters tool with Sentry tracing
    mcp.AddTool(server, &mcp.Tool{
        Name:        "get_action_parameters",
        Description: "Fetch and parse a GitHub Action's action.yml file...",
    }, WithSentryTracing("get_action_parameters", m.handleGetActionParameters))
}
```

## Span Conventions

All spans follow the OpenTelemetry MCP semantic conventions:

### Span Name

Tool call spans use the format: `tools/call {tool_name}`

Examples:

- `tools/call get_action_parameters`
- `tools/call my_custom_tool`

### Span Operation

All MCP tool spans use the operation: `mcp.server`

### Common Attributes

All spans include these attributes:

| Attribute                  | Type   | Description                  | Example                      |
| -------------------------- | ------ | ---------------------------- | ---------------------------- |
| `mcp.method.name`          | string | The MCP method name          | `"tools/call"`               |
| `mcp.tool.name`            | string | The tool being called        | `"get_action_parameters"`    |
| `mcp.transport`            | string | Transport method used        | `"stdio"`                    |
| `network.transport`        | string | OSI transport layer protocol | `"pipe"`                     |
| `network.protocol.version` | string | JSON-RPC version             | `"2.0"`                      |
| `sentry.origin`            | string | Sentry origin identifier     | `"auto.function.mcp_server"` |
| `sentry.source`            | string | Sentry source type           | `"route"`                    |

### Tool-Specific Attributes

#### Tool Arguments

Tool arguments are automatically extracted and set with the prefix `mcp.request.argument`:

```
mcp.request.argument.actionref = "actions/checkout@v5"
```

The argument names are:

- Extracted from JSON struct tags
- Converted to lowercase
- Prefixed with `mcp.request.argument.`

#### Tool Results

Result metadata is captured:

| Attribute                       | Type    | Description                        | Example    |
| ------------------------------- | ------- | ---------------------------------- | ---------- |
| `mcp.tool.result.is_error`      | boolean | Whether the tool returned an error | `false`    |
| `mcp.tool.result.content_count` | int     | Number of content items returned   | `1`        |
| `mcp.tool.result.content`       | string  | JSON array of content types        | `["text"]` |

### Request Metadata

If available, the following are extracted from the request:

| Attribute        | Type   | Description               |
| ---------------- | ------ | ------------------------- |
| `mcp.request.id` | string | Unique request identifier |
| `mcp.session.id` | string | MCP session identifier    |

## Span Status

Spans are marked with appropriate status:

- `ok` - Tool executed successfully
- `internal_error` - Tool returned an error

## Error Capture

When a tool handler returns an error:

1. The span status is set to `internal_error`
2. `mcp.tool.result.is_error` is set to `true`
3. The error is captured to Sentry with full context
4. The error is propagated to the MCP client

## Example Span Data

Here's an example of what a tool call span looks like in Sentry:

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
    "network.protocol.version": "2.0",
    "mcp.request.argument.actionref": "actions/checkout@v5",
    "mcp.tool.result.is_error": false,
    "mcp.tool.result.content_count": 1,
    "mcp.tool.result.content": "[\"text\"]",
    "sentry.origin": "auto.function.mcp_server",
    "sentry.source": "route"
  }
}
```

## Comparison with JavaScript SDK

This implementation closely follows the Sentry JavaScript SDK's MCP integration:

### Similarities

- Follows same OpenTelemetry MCP conventions
- Uses identical attribute names and values
- Implements same span creation patterns
- Captures results and errors similarly

### Differences

- **Language**: Go vs TypeScript
- **SDK Integration**: Direct wrapper vs transport interception
  - JS: Wraps transport layer to intercept all messages
  - Go: Wraps individual tool handlers (simpler, more idiomatic)
- **Type Safety**: Go uses generics for type-safe wrappers
- **Session Management**: Not yet implemented (stateless server)

### Why the Difference?

The Go MCP SDK has a different architecture:

- Tool handlers are registered directly with type safety
- No need to wrap transport layer for basic tool tracing
- Simpler approach that achieves the same observability goals

## Future Enhancements

Potential improvements to consider:

1. **Session Management**: Track client/server info across requests
2. **Transport Wrapping**: Intercept all MCP messages (not just tool calls)
3. **Resource Tracing**: Add spans for resource access
4. **Prompt Tracing**: Add spans for prompt requests
5. **Notification Tracing**: Track MCP notifications
6. **Result Content**: Optionally capture full result payloads (with PII filtering)

## References

- [OpenTelemetry MCP Semantic Conventions](https://github.com/open-telemetry/semantic-conventions/pull/2083)
- [MCP Specification](https://modelcontextprotocol.io/)
- [Sentry Go SDK](https://docs.sentry.io/platforms/go/)
- [Sentry JavaScript MCP Integration](https://github.com/getsentry/sentry-javascript/tree/develop/packages/core/src/integrations/mcp-server)
