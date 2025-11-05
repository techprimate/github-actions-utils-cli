# Analysis: Sentry JavaScript MCP Integration

This document provides a detailed analysis of how the Sentry JavaScript SDK implements MCP (Model Context Protocol) tracing, which was used as the basis for implementing the Go version.

## Overview

The Sentry JavaScript MCP integration (`@sentry/core/integrations/mcp-server`) provides comprehensive instrumentation for MCP servers following the [OpenTelemetry MCP Semantic Conventions](https://github.com/open-telemetry/semantic-conventions/pull/2083).

## Architecture Analysis

### High-Level Design

The JavaScript implementation uses a **wrapping pattern** at multiple layers:

```
┌─────────────────────────────────────────────┐
│  wrapMcpServerWithSentry()                  │
│  (Main entry point)                         │
└────────────────┬────────────────────────────┘
                 │
         ┌───────┴────────┐
         │                │
    ┌────▼─────┐    ┌────▼──────────┐
    │ Transport │    │ Handlers      │
    │ Wrapping  │    │ Wrapping      │
    └────┬─────┘    └────┬──────────┘
         │                │
    ┌────▼──────────┐     ├─── tool()
    │ onmessage     │     ├─── resource()
    │ send          │     └─── prompt()
    │ onclose       │
    │ onerror       │
    └───────────────┘
```

### Key Components

#### 1. **Transport Layer Wrapping** (`transport.ts`)

Intercepts all MCP messages at the transport level:

- **`wrapTransportOnMessage()`**: Creates spans for incoming requests
- **`wrapTransportSend()`**: Correlates responses with requests
- **`wrapTransportOnClose()`**: Cleans up pending spans
- **`wrapTransportError()`**: Captures transport errors

**Key Insight**: By wrapping the transport layer, the JS SDK can:

- Automatically detect request/notification types
- Create spans before handlers execute
- Correlate responses with their originating requests
- Track session lifecycle

#### 2. **Handler Wrapping** (`handlers.ts`)

Wraps individual MCP handler registration methods:

- `wrapToolHandlers()` - Wraps `server.tool()`
- `wrapResourceHandlers()` - Wraps `server.resource()`
- `wrapPromptHandlers()` - Wraps `server.prompt()`

**Purpose**: Provides error capture specific to each handler type.

#### 3. **Span Creation** (`spans.ts`)

Centralized span creation following conventions:

```typescript
function createMcpSpan(config: McpSpanConfig): unknown {
  const { type, message, transport, extra, callback } = config;
  
  // Build span name (e.g., "tools/call get_action_parameters")
  const spanName = createSpanName(method, target);
  
  // Collect attributes
  const attributes = {
    ...buildTransportAttributes(transport, extra),
    [MCP_METHOD_NAME_ATTRIBUTE]: method,
    ...buildTypeSpecificAttributes(type, message, params),
    ...buildSentryAttributes(type),
  };
  
  // Create and execute span
  return startSpan({ name: spanName, forceTransaction: true, attributes }, callback);
}
```

#### 4. **Attribute Extraction** (`attributeExtraction.ts`, `methodConfig.ts`)

Method-specific attribute extraction:

```typescript
const METHOD_CONFIGS: Record<string, MethodConfig> = {
  'tools/call': {
    targetField: 'name',
    targetAttribute: MCP_TOOL_NAME_ATTRIBUTE,
    captureArguments: true,
    argumentsField: 'arguments',
  },
  'resources/read': {
    targetField: 'uri',
    targetAttribute: MCP_RESOURCE_URI_ATTRIBUTE,
    captureUri: true,
  },
  // ...
};
```

#### 5. **Session Management** (`sessionManagement.ts`, `sessionExtraction.ts`)

Tracks session-level data:

- Extracts client/server info from `initialize` requests
- Stores per-transport session data
- Propagates session attributes to all spans in that session

#### 6. **Correlation** (`correlation.ts`)

Maps request IDs to their spans for result correlation:

```typescript
// Store span when request arrives
storeSpanForRequest(transport, requestId, span, method);

// Complete span when response is sent
completeSpanWithResults(transport, requestId, result);
```

## Span Lifecycle

### For a Tool Call

```
1. Client sends: {"jsonrpc":"2.0","id":1,"method":"tools/call","params":{...}}
   ↓
2. wrapTransportOnMessage() intercepts
   ↓
3. buildMcpServerSpanConfig() creates span config
   ↓
4. startInactiveSpan() creates span (not yet active)
   ↓
5. storeSpanForRequest() stores span for correlation
   ↓
6. withActiveSpan() executes handler within span
   ↓
7. Handler executes (potentially wrapped for error capture)
   ↓
8. wrapTransportSend() intercepts response
   ↓
9. extractToolResultAttributes() extracts result metadata
   ↓
10. completeSpanWithResults() finishes span
```

## Attribute Conventions

### Naming Pattern

All MCP attributes follow a consistent pattern:

```
mcp.{category}.{attribute}
```

Examples:

- `mcp.method.name` - Core protocol attribute
- `mcp.tool.name` - Tool-specific attribute
- `mcp.request.argument.{arg_name}` - Request arguments
- `mcp.tool.result.is_error` - Result metadata

### Attribute Categories

1. **Core Protocol**: `mcp.method.name`, `mcp.request.id`, `mcp.session.id`
2. **Transport**: `mcp.transport`, `network.transport`, `network.protocol.version`
3. **Client/Server**: `mcp.client.name`, `mcp.server.version`, etc.
4. **Method-Specific**: `mcp.tool.name`, `mcp.resource.uri`, `mcp.prompt.name`
5. **Arguments**: `mcp.request.argument.*`
6. **Results**: `mcp.tool.result.*`, `mcp.prompt.result.*`

## PII Filtering

The JavaScript SDK includes PII filtering (`piiFiltering.ts`):

```typescript
function filterMcpPiiFromSpanData(
  data: Record<string, unknown>,
  sendDefaultPii: boolean
): Record<string, unknown> {
  // Filter based on sendDefaultPii option
  // Removes: arguments, result content, client address, etc.
}
```

**When `sendDefaultPii` is false**, removes:

- Tool arguments (`mcp.request.argument.*`)
- Tool result content (`mcp.tool.result.content`, `mcp.tool.result.*`)
- Prompt arguments and results
- Client address/port
- Logging messages

## Error Handling

### Types of Errors Captured

1. **Tool Execution Errors**: Handler throws or returns error
2. **Protocol Errors**: JSON-RPC error responses (code -32603)
3. **Transport Errors**: Connection failures
4. **Validation Errors**: Invalid parameters or protocol violations
5. **Timeout Errors**: Long-running operations

Each error type is tagged appropriately for filtering in Sentry.

## Comparison: JavaScript vs Go Implementation

| Aspect                | JavaScript Implementation   | Go Implementation             |
| --------------------- | --------------------------- | ----------------------------- |
| **Approach**          | Wraps transport layer       | Wraps tool handlers directly  |
| **Complexity**        | High (multi-layer wrapping) | Low (single-layer wrapping)   |
| **Coverage**          | All MCP messages            | Tool calls only (currently)   |
| **Session Tracking**  | Full session management     | Not implemented (stateless)   |
| **Type Safety**       | TypeScript interfaces       | Go generics                   |
| **Integration Point** | `wrapMcpServerWithSentry()` | `WithSentryTracing()` wrapper |
| **Dependencies**      | Many internal modules       | Single sentry.go file         |

### Why the Go Implementation is Simpler

1. **SDK Architecture**: Go MCP SDK has different design
2. **Type System**: Go generics enable cleaner handler wrapping
3. **Use Case**: Simpler CLI tool vs full-featured SDK
4. **Stateless**: No session management needed for stdio transport

## Key Learnings

### What Works Well in JavaScript

1. **Transport wrapping** provides automatic instrumentation
2. **Span correlation** ensures proper request-response tracking
3. **Session management** enables rich contextual data
4. **PII filtering** protects sensitive information
5. **Comprehensive error capture** catches all failure modes

### What We Adapted for Go

1. **Simplified to handler-level wrapping** (good enough for tool calls)
2. **Used Go generics** for type-safe wrappers
3. **Focused on essential attributes** (no session management yet)
4. **Maintained naming conventions** for consistency
5. **Kept error capture** for production debugging

### Potential Future Enhancements

1. **Transport-level wrapping** to capture all messages
2. **Session tracking** for multi-request correlation
3. **PII filtering** with configuration options
4. **Resource/prompt spans** for complete coverage
5. **Notification tracking** for bidirectional communication

## Code Quality Observations

The JavaScript implementation demonstrates:

- ✅ **Excellent separation of concerns** (one file per responsibility)
- ✅ **Comprehensive documentation** (TSDoc comments)
- ✅ **Type safety** throughout
- ✅ **Extensive validation** for robustness
- ✅ **Defensive programming** (try-catch everywhere)
- ✅ **Consistent naming** following conventions
- ✅ **Testable design** (dependency injection)

## References

- [Sentry JS MCP Integration](https://github.com/getsentry/sentry-javascript/tree/develop/packages/core/src/integrations/mcp-server)
- [OpenTelemetry MCP Conventions](https://github.com/open-telemetry/semantic-conventions/pull/2083)
- [MCP Specification](https://modelcontextprotocol.io/)
