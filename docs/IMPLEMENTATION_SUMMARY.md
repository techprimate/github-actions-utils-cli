# MCP Tracing Implementation Summary

This document summarizes the implementation of MCP (Model Context Protocol) tracing with Sentry for the GitHub Actions Utils CLI.

## Overview

Implemented automatic tracing for MCP tool calls following the OpenTelemetry MCP Semantic Conventions, based on the Sentry JavaScript SDK's MCP integration.

## What Was Implemented

### 1. Core Tracing Wrapper (`internal/cli/mcp/sentry.go`)

**File**: `internal/cli/mcp/sentry.go` (261 lines)

A comprehensive wrapper function that:

- Creates Sentry spans for MCP tool executions
- Extracts and sets MCP-specific attributes
- Captures tool arguments automatically
- Tracks tool results and errors
- Follows OpenTelemetry MCP semantic conventions

**Key Function**:

```go
func WithSentryTracing[In, Out any](toolName string, handler mcp.ToolHandlerFor[In, Out]) mcp.ToolHandlerFor[In, Out]
```

**Features**:

- Type-safe using Go generics
- Automatic argument extraction via reflection
- Result metadata capture
- Error capture and correlation
- Proper span status handling

### 2. Integration (`internal/cli/mcp/server.go`)

Updated tool registration to use the Sentry wrapper:

```go
mcp.AddTool(server, &mcp.Tool{
    Name:        "get_action_parameters",
    Description: "Fetch and parse a GitHub Action's action.yml file...",
}, WithSentryTracing("get_action_parameters", m.handleGetActionParameters))
```

### 3. Tests (`internal/cli/mcp/sentry_test.go`)

Comprehensive test suite covering:

- Successful tool executions
- Error handling and propagation
- Argument extraction
- Content type detection

All tests pass ✅

### 4. Documentation

Created three documentation files:

1. **`docs/MCP_TRACING.md`** (208 lines)
   - User-facing documentation
   - Usage examples
   - Span conventions
   - Attribute reference
   - Example span data

2. **`docs/ANALYSIS_SENTRY_MCP_INTEGRATION.md`** (259 lines)
   - Deep analysis of JavaScript implementation
   - Architecture breakdown
   - Component analysis
   - Comparison with Go implementation
   - Key learnings

3. **Updated `AGENTS.md`**
   - Added MCP Tracing section to Sentry Integration
   - Updated "Adding a New MCP Tool" instructions
   - References to documentation

## Attributes Captured

### Common Attributes (All Spans)

| Attribute                  | Example Value                |
| -------------------------- | ---------------------------- |
| `mcp.method.name`          | `"tools/call"`               |
| `mcp.tool.name`            | `"get_action_parameters"`    |
| `mcp.transport`            | `"stdio"`                    |
| `network.transport`        | `"pipe"`                     |
| `network.protocol.version` | `"2.0"`                      |
| `sentry.origin`            | `"auto.function.mcp_server"` |
| `sentry.source`            | `"route"`                    |

### Tool-Specific Attributes

- **Arguments**: `mcp.request.argument.*` (e.g., `mcp.request.argument.actionref`)
- **Results**:
  - `mcp.tool.result.is_error` (boolean)
  - `mcp.tool.result.content_count` (int)
  - `mcp.tool.result.content` (JSON array of content types)

### Optional Attributes

- `mcp.request.id` - Request identifier (if available)
- `mcp.session.id` - Session identifier (if available)

## Span Structure

**Operation**: `mcp.server`

**Name Pattern**: `tools/call {tool_name}`

**Examples**:

- `tools/call get_action_parameters`
- `tools/call my_custom_tool`

**Status**:

- `ok` - Successful execution
- `internal_error` - Tool returned an error

## Implementation Approach

### Go vs JavaScript Differences

| Aspect                 | JavaScript SDK                     | Go Implementation            |
| ---------------------- | ---------------------------------- | ---------------------------- |
| **Integration Point**  | Transport layer wrapping           | Handler-level wrapping       |
| **Complexity**         | Multi-layer (transport + handlers) | Single layer (handlers only) |
| **Type Safety**        | TypeScript interfaces              | Go generics                  |
| **Session Management** | Full session tracking              | Not implemented (stateless)  |
| **Coverage**           | All MCP messages                   | Tool calls only              |

### Why Handler-Level Wrapping?

The Go implementation uses a simpler approach:

1. **SDK Architecture**: The Go MCP SDK has strong type safety built-in
2. **Use Case**: CLI tool with simple stdio transport
3. **Stateless**: No need for session management
4. **Clean API**: `WithSentryTracing()` wrapper is intuitive
5. **Good Enough**: Captures essential observability data

### Could We Do Transport Wrapping?

Yes, but it would require:

- Wrapping the `mcp_sdk.Transport` interface
- Implementing custom transport type
- Managing request-response correlation
- More complexity without significant benefit for this use case

## Testing

### Unit Tests

All tests pass:

```bash
$ make test
ok  github.com/techprimate/github-actions-utils-cli/internal/cli/mcp  0.456s
```

**Test Coverage**:

- ✅ Successful tool execution with tracing
- ✅ Error handling and propagation
- ✅ Argument extraction from structs
- ✅ Content type detection

### Integration Test

Created `test_mcp_invocation.sh` to test end-to-end:

```bash
$ ./test_mcp_invocation.sh
✅ MCP server test completed
```

### Build & Analyze

All quality checks pass:

```bash
$ make build   # ✅ Builds successfully
$ make format  # ✅ Code formatted
$ make analyze # ✅ No issues found
$ make test    # ✅ All tests pass
```

## How to Use

### For New Tools

When adding a new tool, wrap the handler with `WithSentryTracing`:

```go
mcp.AddTool(server, &mcp.Tool{
    Name:        "my_new_tool",
    Description: "Does something useful",
}, WithSentryTracing("my_new_tool", m.handleMyNewTool))
```

That's it! The wrapper automatically:

- Creates spans
- Extracts arguments
- Captures results
- Handles errors

### Viewing in Sentry

Spans appear in Sentry with:

- **Performance** → **Traces**
- Operation: `mcp.server`
- Description: `tools/call {tool_name}`

Filter by:

- `mcp.tool.name` to see specific tools
- `mcp.tool.result.is_error:true` to find errors

## Architecture Decisions

### 1. Why Reflection for Arguments?

**Pros**:

- Works with any tool argument struct
- No boilerplate code needed
- Type-safe at compile time
- JSON tags automatically used

**Cons**:

- Slight runtime overhead (negligible)
- Cannot extract unexported fields (acceptable)

**Decision**: Benefits outweigh costs for observability.

### 2. Why Not PII Filtering?

**Reasoning**:

- This is a CLI tool, not a library
- Users control the environment
- Sentry DSN is configurable
- Can be added later if needed

**Mitigation**: Document that sensitive data may be captured.

### 3. Why Not Session Management?

**Reasoning**:

- Stdio transport is stateless
- Each invocation is independent
- Session tracking adds complexity
- Not needed for current use case

**Future**: Could add for HTTP/SSE transports.

## Comparison with Sentry JavaScript SDK

### Similarities ✅

- ✅ Follows same OpenTelemetry conventions
- ✅ Uses identical attribute names
- ✅ Same span operation and naming
- ✅ Captures arguments and results
- ✅ Error handling and correlation

### Differences 🔄

- 🔄 Simpler architecture (handler vs transport wrapping)
- 🔄 Go generics instead of TypeScript types
- 🔄 No session management (stateless)
- 🔄 No PII filtering (yet)
- 🔄 Tool calls only (no resources/prompts yet)

### JavaScript Features Not Implemented

1. **Transport Layer Wrapping**: Not needed for stdio
2. **Session Management**: Stateless design
3. **Resource/Prompt Spans**: Only have tool calls
4. **Notification Tracking**: Not applicable
5. **PII Filtering**: Can add if needed
6. **Result Content Capture**: Only capture metadata

## Future Enhancements

### High Priority

1. **PII Filtering**: Add `sendDefaultPii` option
2. **Resource Spans**: If we add resource handlers
3. **Prompt Spans**: If we add prompt handlers

### Medium Priority

4. **Session Tracking**: For HTTP/SSE transports
5. **Transport Wrapping**: For complete coverage
6. **Notification Spans**: For bidirectional communication

### Low Priority

7. **Full Result Capture**: With PII filtering
8. **Custom Attributes**: User-defined span attributes
9. **Sampling**: Control span sampling rate

## Files Changed/Added

### Added Files

- `internal/cli/mcp/sentry.go` (261 lines)
- `internal/cli/mcp/sentry_test.go` (200 lines)
- `docs/MCP_TRACING.md` (208 lines)
- `docs/ANALYSIS_SENTRY_MCP_INTEGRATION.md` (259 lines)
- `docs/IMPLEMENTATION_SUMMARY.md` (this file)
- `test_mcp_invocation.sh` (test script)

### Modified Files

- `internal/cli/mcp/server.go` (updated tool registration)
- `AGENTS.md` (added MCP Tracing section)

## Success Criteria

✅ **All criteria met**:

1. ✅ Follows OpenTelemetry MCP conventions
2. ✅ Creates spans with correct attributes
3. ✅ Captures tool arguments
4. ✅ Tracks results and errors
5. ✅ Type-safe implementation
6. ✅ Comprehensive tests
7. ✅ Complete documentation
8. ✅ All quality checks pass
9. ✅ Easy to use for new tools
10. ✅ Consistent with Sentry JS SDK

## Conclusion

The implementation successfully brings MCP tracing to the Go CLI, following the same conventions as the Sentry JavaScript SDK while adapting to Go's idioms and the project's simpler architecture.

**Key Achievements**:

- Clean, type-safe API
- Minimal boilerplate
- Comprehensive observability
- Well-documented
- Production-ready

**Next Steps**:

- Monitor spans in production
- Gather feedback
- Consider additional enhancements
- Keep aligned with OpenTelemetry conventions
