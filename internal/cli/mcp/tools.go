package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetActionParametersArgs defines the parameters for the get_action_parameters tool.
type GetActionParametersArgs struct {
	ActionRef string `json:"actionRef" jsonschema:"GitHub Action reference (e.g., 'actions/checkout@v5')"`
}

// handleGetActionParameters handles the get_action_parameters tool call.
func (m *MCPServer) handleGetActionParameters(ctx context.Context, req *mcp.CallToolRequest, args GetActionParametersArgs) (*mcp.CallToolResult, any, error) {
	// Validate input
	if args.ActionRef == "" {
		return nil, nil, fmt.Errorf("actionRef is required")
	}

	// Fetch and parse action parameters
	params, err := m.actionsService.GetActionParameters(args.ActionRef)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get action parameters: %w", err)
	}

	// Format as text output
	textOutput := fmt.Sprintf("Action: %s\n\n", args.ActionRef)

	// Add name and description if available
	if name, ok := params["name"].(string); ok {
		textOutput += fmt.Sprintf("Name: %s\n", name)
	}
	if desc, ok := params["description"].(string); ok {
		textOutput += fmt.Sprintf("Description: %s\n", desc)
	}

	// Add inputs summary if available
	if inputs, ok := params["inputs"].(map[string]interface{}); ok {
		textOutput += fmt.Sprintf("\nInputs: %d defined\n", len(inputs))
	}

	// Add outputs summary if available
	if outputs, ok := params["outputs"].(map[string]interface{}); ok {
		textOutput += fmt.Sprintf("Outputs: %d defined\n", len(outputs))
	}

	textOutput += "\nFull action.yml structure returned in structured data."

	// Return response with both text and structured data
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: textOutput,
			},
		},
	}, params, nil
}
