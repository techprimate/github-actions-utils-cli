package mcp

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MockArgs represents test arguments for a tool
type MockArgs struct {
	Name  string `json:"name" jsonschema:"The name parameter"`
	Count int    `json:"count" jsonschema:"The count parameter"`
}

func TestWithSentryTracing_Success(t *testing.T) {
	// Initialize Sentry with a test transport
	transport := &testTransport{}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:       "https://test@test.ingest.sentry.io/123456",
		Transport: transport,
	})
	if err != nil {
		t.Fatalf("Failed to initialize Sentry: %v", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Create a mock handler that succeeds
	mockHandler := func(ctx context.Context, req *mcp.CallToolRequest, args MockArgs) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Success"},
			},
		}, map[string]string{"status": "ok"}, nil
	}

	// Wrap with Sentry tracing
	wrappedHandler := WithSentryTracing("test_tool", mockHandler)

	// Execute the handler
	ctx := context.Background()
	args := MockArgs{Name: "test", Count: 42}
	result, data, err := wrappedHandler(ctx, &mcp.CallToolRequest{}, args)

	// Verify execution
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Error("Expected result, got nil")
	}
	if data == nil {
		t.Error("Expected data, got nil")
	}

	// Flush to ensure span is sent
	sentry.Flush(2 * time.Second)

	// Note: In a real test, you would verify the span was created with correct attributes
	// This requires either mocking the Sentry transport or using the test transport
}

func TestWithSentryTracing_Error(t *testing.T) {
	// Initialize Sentry with a test transport
	transport := &testTransport{}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:       "https://test@test.ingest.sentry.io/123456",
		Transport: transport,
	})
	if err != nil {
		t.Fatalf("Failed to initialize Sentry: %v", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Create a mock handler that fails
	expectedErr := errors.New("tool execution failed")
	mockHandler := func(ctx context.Context, req *mcp.CallToolRequest, args MockArgs) (*mcp.CallToolResult, any, error) {
		return nil, nil, expectedErr
	}

	// Wrap with Sentry tracing
	wrappedHandler := WithSentryTracing("test_tool_error", mockHandler)

	// Execute the handler
	ctx := context.Background()
	args := MockArgs{Name: "test", Count: 42}
	result, data, err := wrappedHandler(ctx, &mcp.CallToolRequest{}, args)

	// Verify error is propagated
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Errorf("Expected nil result on error, got: %v", result)
	}
	if data != nil {
		t.Errorf("Expected nil data on error, got: %v", data)
	}

	// Flush to ensure error is sent
	sentry.Flush(2 * time.Second)
}

func TestWithSentryTracing_ArgumentExtraction(t *testing.T) {
	// Initialize Sentry
	transport := &testTransport{}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:       "https://test@test.ingest.sentry.io/123456",
		Transport: transport,
	})
	if err != nil {
		t.Fatalf("Failed to initialize Sentry: %v", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Create a handler that just returns success
	mockHandler := func(ctx context.Context, req *mcp.CallToolRequest, args MockArgs) (*mcp.CallToolResult, any, error) {
		// Verify arguments were passed correctly
		if args.Name != "test_arg" {
			return nil, nil, errors.New("wrong name argument")
		}
		if args.Count != 123 {
			return nil, nil, errors.New("wrong count argument")
		}
		return &mcp.CallToolResult{}, nil, nil
	}

	// Wrap with Sentry tracing
	wrappedHandler := WithSentryTracing("test_args", mockHandler)

	// Execute with specific arguments
	ctx := context.Background()
	args := MockArgs{Name: "test_arg", Count: 123}
	_, _, err = wrappedHandler(ctx, &mcp.CallToolRequest{}, args)

	if err != nil {
		t.Errorf("Handler failed: %v", err)
	}

	// The span should have attributes:
	// - mcp.request.argument.name = "test_arg"
	// - mcp.request.argument.count = 123
	sentry.Flush(2 * time.Second)
}

func TestGetContentType(t *testing.T) {
	tests := []struct {
		name     string
		content  mcp.Content
		expected string
	}{
		{
			name:     "TextContent",
			content:  &mcp.TextContent{Text: "test"},
			expected: "text",
		},
		{
			name:     "ImageContent",
			content:  &mcp.ImageContent{Data: []byte("base64data"), MIMEType: "image/png"},
			expected: "image",
		},
		{
			name:     "AudioContent",
			content:  &mcp.AudioContent{Data: []byte("base64data"), MIMEType: "audio/mp3"},
			expected: "audio",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getContentType(tt.content)
			if result != tt.expected {
				t.Errorf("Expected content type %q, got %q", tt.expected, result)
			}
		})
	}
}

// testTransport is a no-op transport for testing
type testTransport struct{}

func (t *testTransport) Configure(options sentry.ClientOptions) {}

func (t *testTransport) SendEvent(event *sentry.Event) {}

func (t *testTransport) Flush(timeout time.Duration) bool {
	return true
}

func (t *testTransport) FlushWithContext(ctx context.Context) bool {
	return true
}

func (t *testTransport) Close() {}
