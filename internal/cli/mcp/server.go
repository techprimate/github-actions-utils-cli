package mcp

import (
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/techprimate/github-actions-utils-cli/internal/github"
)

// MCPServer wraps the ActionsService and provides MCP tool handlers.
// It uses dependency injection to receive its dependencies.
type MCPServer struct {
	actionsService *github.ActionsService
	logger         *slog.Logger
}

// NewMCPServer creates a new MCP server with the given dependencies.
func NewMCPServer(actionsService *github.ActionsService, logger *slog.Logger) *MCPServer {
	if logger == nil {
		logger = slog.Default()
	}
	return &MCPServer{
		actionsService: actionsService,
		logger:         logger,
	}
}

// RegisterTools registers all available tools with the MCP server.
func (m *MCPServer) RegisterTools(server *mcp.Server) {
	// Register get_action_parameters tool with Sentry tracing
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_action_parameters",
		Description: "Fetch and parse a GitHub Action's action.yml file. Returns the complete action.yml structure including inputs, outputs, runs configuration, and metadata.",
	}, WithSentryTracing("get_action_parameters", m.handleGetActionParameters))
}
