// Package cmd contains all CLI commands and their implementation.
// It uses the Cobra library for command-line interface construction.
//
// Command Structure:
//   - root: Base command with global flags
//   - mcp: Run MCP server for agent integration
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version information - set at build time via ldflags
	// Example: go build -ldflags "-X github.com/techprimate/github-actions-utils-cli/internal/cli/cmd.version=1.0.0"
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github-actions-utils-cli",
	Short: "MCP CLI for GitHub Actions utilities",
	Long: `GitHub Actions Utils CLI provides an MCP (Model Context Protocol) server
for interacting with GitHub Actions.

The primary use case is running as an MCP server that AI agents can use to
fetch and parse GitHub Action definitions (action.yml files).`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
