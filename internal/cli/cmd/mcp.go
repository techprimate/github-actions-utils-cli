package cmd

import (
	"io"
	"log/slog"

	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/techprimate/github-actions-utils-cli/internal/cli/mcp"
	"github.com/techprimate/github-actions-utils-cli/internal/github"
)

var MCPCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run MCP server for agent integration",
	Long: `Runs an MCP (Model Context Protocol) server that exposes GitHub Actions utilities as tools.

The server communicates over stdin/stdout and provides this tool:
  - get_action_parameters: Fetch and parse GitHub Action action.yml files

This allows AI agents to programmatically retrieve information about GitHub Actions,
including their inputs, outputs, and configuration.

Example MCP client configuration:
{
  "mcpServers": {
    "github-actions-utils": {
      "command": "github-actions-utils-cli",
      "args": ["mcp"]
    }
  }
}`,
	RunE: runMCP,
}

func init() {
	rootCmd.AddCommand(MCPCmd)
}

func runMCP(cmd *cobra.Command, args []string) error {
	// MCP uses stdio for JSON-RPC, so we need to silence the logger
	// to avoid interfering with the protocol
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Create GitHub Actions service
	actionsService := github.NewActionsService()

	// Create MCP server wrapper with silent logger
	mcpSrv := mcp.NewMCPServer(actionsService, silentLogger)

	// Create go-sdk MCP server
	server := mcp_sdk.NewServer(&mcp_sdk.Implementation{
		Name:    "github-actions-utils",
		Version: version,
	}, nil)

	// Register all tools
	mcpSrv.RegisterTools(server)

	// Run server on stdio (logging disabled to keep stdio clean for JSON-RPC)
	return server.Run(cmd.Context(), &mcp_sdk.StdioTransport{})
}
