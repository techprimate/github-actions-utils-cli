# GitHub Actions Utils CLI

> MCP server for GitHub Actions utilities

A Model Context Protocol (MCP) server that provides tools for working with GitHub Actions. This allows AI agents and MCP clients to programmatically fetch and parse GitHub Action definitions.

## Features

- **get_action_parameters**: Fetch and parse any GitHub Action's `action.yml` file
- Returns complete action metadata including inputs, outputs, runs configuration, and description
- Works with any public GitHub Action repository
- Compatible with all MCP clients (Claude Desktop, Cline, etc.)

## Installation

### Download Pre-built Binary (Recommended)

Download the latest release for your platform as described in the [latest release notes](https://github.com/techprimate/github-actions-utils-cli/releases/tag/latest).

### Build from Source

```bash
# Clone the repository
git clone https://github.com/techprimate/github-actions-utils-cli.git
cd github-actions-utils-cli

# Build
make build

# Binary will be at ./dist/github-actions-utils-cli
```

## Usage

### MCP Server

The primary use case is running as an MCP server:

```bash
github-actions-utils-cli mcp
```

### MCP Client Configuration

#### Claude CLI

Add to your Claude CLI configuration using the `claude mcp` command:

```bash
claude mcp add --transport stdio github-actions-utils-cli github-actions-utils-cli mcp
```

#### Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "github-actions-utils": {
      "command": "/usr/local/bin/github-actions-utils-cli",
      "args": ["mcp"]
    }
  }
}
```

#### Cline (VS Code Extension)

Add to your Cline MCP settings:

```json
{
  "mcpServers": {
    "github-actions-utils": {
      "command": "/usr/local/bin/github-actions-utils-cli",
      "args": ["mcp"]
    }
  }
}
```

## Available Tools

### get_action_parameters

Fetches and parses a GitHub Action's `action.yml` file.

**Parameters:**

- `actionRef` (string, required): GitHub Action reference in the format `owner/repo@version`
  - Example: `actions/checkout@v5`
  - Example: `actions/setup-node@v4`

**Returns:**
Complete action.yml structure as JSON, including:

- `name`: Action name
- `description`: Action description
- `inputs`: All input parameters with descriptions, defaults, and whether they're required
- `outputs`: All output parameters with descriptions
- `runs`: Runtime configuration (node version, main entry point, etc.)
- `branding`: Icon and color information

**Example Usage in Claude:**

```
Can you show me the parameters for the actions/checkout@v5 action?
```

Claude will use the `get_action_parameters` tool and return structured information about all inputs and outputs.

**Example Response:**

```json
{
  "name": "Checkout",
  "description": "Checkout a Git repository at a particular version",
  "inputs": {
    "repository": {
      "description": "Repository name with owner. For example, actions/checkout",
      "default": "${{ github.repository }}"
    },
    "ref": {
      "description": "The branch, tag or SHA to checkout...",
      "default": ""
    },
    "token": {
      "description": "Personal access token (PAT) used to fetch the repository...",
      "default": "${{ github.token }}"
    }
    // ... more inputs
  },
  "runs": {
    "using": "node24",
    "main": "dist/index.js"
  }
}
```

## Use Cases

- **Building CI/CD Tools**: Get accurate information about available GitHub Actions
- **Documentation**: Automatically document which actions your workflows use
- **Validation**: Verify action parameters before using them in workflows
- **Discovery**: Explore available inputs and outputs for actions
- **Migration**: Understand action changes when updating versions

## Development

### Prerequisites

- Go 1.25 or later
- Make

### Setup

```bash
# Initialize project (installs dependencies)
make init

# Build
make build

# Run tests
make test

# Format code
make format

# Run static analysis
make analyze
```

### Project Structure

```
.
├── cmd/cli/              # CLI entry point
├── internal/
│   ├── cli/
│   │   ├── cmd/          # Cobra commands
│   │   └── mcp/          # MCP server and tools
│   ├── github/           # GitHub Actions service
│   └── logging/          # Logging utilities
├── .github/workflows/    # CI/CD pipelines
├── Makefile             # Development commands
└── README.md            # This file
```

### Testing the MCP Server

```bash
# Build the CLI
make build

# Test that MCP server responds
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./dist/github-actions-utils-cli mcp

# Test get_action_parameters tool
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_action_parameters","arguments":{"actionRef":"actions/checkout@v5"}}}' | ./dist/github-actions-utils-cli mcp | jq
```

## Telemetry

This project uses Sentry for error tracking and monitoring. You can disable telemetry by setting the environment variable:

```bash
export TELEMETRY_ENABLED=false
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.

## Related Projects

- [Model Context Protocol](https://github.com/modelcontextprotocol) - The MCP specification
- [GitHub Actions](https://github.com/features/actions) - GitHub's CI/CD platform
