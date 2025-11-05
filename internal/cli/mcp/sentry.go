package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCP Attribute Constants
// Based on OpenTelemetry MCP Semantic Conventions
// See: https://github.com/open-telemetry/semantic-conventions/pull/2083

const (
	// Core MCP Attributes
	AttrMCPMethodName      = "mcp.method.name"
	AttrMCPRequestID       = "mcp.request.id"
	AttrMCPSessionID       = "mcp.session.id"
	AttrMCPTransport       = "mcp.transport"
	AttrNetworkTransport   = "network.transport"
	AttrNetworkProtocolVer = "network.protocol.version"

	// Tool-specific Attributes
	AttrMCPToolName               = "mcp.tool.name"
	AttrMCPToolResultIsError      = "mcp.tool.result.is_error"
	AttrMCPToolResultContentCount = "mcp.tool.result.content_count"
	AttrMCPToolResultContent      = "mcp.tool.result.content"

	// Request Arguments Prefix
	AttrMCPRequestArgumentPrefix = "mcp.request.argument"

	// Sentry-specific Values
	OpMCPServer          = "mcp.server"
	OriginMCPFunction    = "auto.function.mcp_server"
	SourceMCPRoute       = "route"
	TransportStdio       = "stdio"
	NetworkTransportPipe = "pipe"
	JSONRPCVersion       = "2.0"
)

// WithSentryTracing wraps an MCP tool handler with Sentry tracing.
// It creates spans following OpenTelemetry MCP semantic conventions and
// captures tool execution results and errors.
//
// Example usage:
//
//	mcp.AddTool(server, &mcp.Tool{
//	    Name:        "my_tool",
//	    Description: "Does something useful",
//	}, WithSentryTracing("my_tool", func(ctx context.Context, req *mcp.CallToolRequest, args MyToolArgs) (*mcp.CallToolResult, any, error) {
//	    return m.handleMyTool(ctx, req, args)
//	}))
func WithSentryTracing[In, Out any](toolName string, handler mcp.ToolHandlerFor[In, Out]) mcp.ToolHandlerFor[In, Out] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args In) (*mcp.CallToolResult, Out, error) {
		// Create span for tool execution
		span := sentry.StartSpan(ctx, OpMCPServer)
		defer span.Finish()

		// Set span name following MCP conventions: "tools/call {tool_name}"
		span.Description = fmt.Sprintf("tools/call %s", toolName)

		// Set common MCP attributes
		span.SetData(AttrMCPMethodName, "tools/call")
		span.SetData(AttrMCPToolName, toolName)
		span.SetData(AttrMCPTransport, TransportStdio)
		span.SetData(AttrNetworkTransport, NetworkTransportPipe)
		span.SetData(AttrNetworkProtocolVer, JSONRPCVersion)

		// Set Sentry-specific attributes
		span.SetData("sentry.origin", OriginMCPFunction)
		span.SetData("sentry.source", SourceMCPRoute)

		// Extract and set request ID if available
		if req != nil {
			// The CallToolRequest may have metadata we can extract
			// For now, we'll use reflection to check if there's an ID field
			setRequestMetadata(span, req)
		}

		// Extract and set tool arguments
		setToolArguments(span, args)

		// Execute the handler with the span's context
		ctx = span.Context()
		result, data, err := handler(ctx, req, args)

		// Capture error if present
		if err != nil {
			span.Status = sentry.SpanStatusInternalError
			span.SetData(AttrMCPToolResultIsError, true)

			// Capture the error to Sentry with context
			hub := sentry.GetHubFromContext(ctx)
			if hub == nil {
				hub = sentry.CurrentHub()
			}
			hub.CaptureException(err)
		} else {
			span.Status = sentry.SpanStatusOK
			span.SetData(AttrMCPToolResultIsError, false)

			// Extract result metadata
			if result != nil {
				setResultMetadata(span, result)
			}
		}

		return result, data, err
	}
}

// setRequestMetadata extracts metadata from the CallToolRequest
func setRequestMetadata(span *sentry.Span, req *mcp.CallToolRequest) {
	// Use reflection to safely check for an ID field
	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Try to find common ID/request ID fields
	if val.Kind() == reflect.Struct {
		// Check for ID field
		if idField := val.FieldByName("ID"); idField.IsValid() {
			switch idField.Kind() {
			case reflect.String:
				if id := idField.String(); id != "" {
					span.SetData(AttrMCPRequestID, id)
				}
			case reflect.Int, reflect.Int64:
				if id := idField.Int(); id != 0 {
					span.SetData(AttrMCPRequestID, fmt.Sprintf("%d", id))
				}
			}
		}

		// Check for SessionID field
		if sessionField := val.FieldByName("SessionID"); sessionField.IsValid() && sessionField.Kind() == reflect.String {
			if sessionID := sessionField.String(); sessionID != "" {
				span.SetData(AttrMCPSessionID, sessionID)
			}
		}
	}
}

// setToolArguments extracts tool arguments and sets them as span attributes
func setToolArguments(span *sentry.Span, args any) {
	if args == nil {
		return
	}

	val := reflect.ValueOf(args)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		// Get JSON tag name or use field name
		jsonTag := fieldType.Tag.Get("json")
		fieldName := fieldType.Name
		if jsonTag != "" {
			// Split on comma to handle tags like "json:field,omitempty"
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				fieldName = parts[0]
			}
		}

		// Convert field name to lowercase for attribute
		attrKey := fmt.Sprintf("%s.%s", AttrMCPRequestArgumentPrefix, strings.ToLower(fieldName))

		// Set the value based on type
		switch field.Kind() {
		case reflect.String:
			if value := field.String(); value != "" {
				span.SetData(attrKey, value)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			span.SetData(attrKey, field.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			span.SetData(attrKey, field.Uint())
		case reflect.Float32, reflect.Float64:
			span.SetData(attrKey, field.Float())
		case reflect.Bool:
			span.SetData(attrKey, field.Bool())
		default:
			// For complex types, serialize to JSON
			if field.CanInterface() {
				if jsonBytes, err := json.Marshal(field.Interface()); err == nil {
					span.SetData(attrKey, string(jsonBytes))
				}
			}
		}
	}
}

// setResultMetadata extracts result metadata and sets span attributes
func setResultMetadata(span *sentry.Span, result *mcp.CallToolResult) {
	if result == nil {
		return
	}

	// Count content items
	contentCount := len(result.Content)
	span.SetData(AttrMCPToolResultContentCount, contentCount)

	// If there's content, serialize it for the span
	// Note: We only capture metadata about the content, not the full content
	// to avoid potentially large payloads
	if contentCount > 0 {
		contentTypes := make([]string, 0, contentCount)
		for _, content := range result.Content {
			// Extract content type information
			if content != nil {
				contentTypes = append(contentTypes, getContentType(content))
			}
		}

		if len(contentTypes) > 0 {
			// Store content types as JSON array string
			if typesJSON, err := json.Marshal(contentTypes); err == nil {
				span.SetData(AttrMCPToolResultContent, string(typesJSON))
			}
		}
	}
}

// getContentType returns the type of content
func getContentType(content mcp.Content) string {
	switch c := content.(type) {
	case *mcp.TextContent:
		return "text"
	case *mcp.ImageContent:
		return "image"
	case *mcp.AudioContent:
		return "audio"
	case *mcp.ResourceLink:
		return "resource_link"
	case *mcp.EmbeddedResource:
		return "embedded_resource"
	default:
		return fmt.Sprintf("%T", c)
	}
}
