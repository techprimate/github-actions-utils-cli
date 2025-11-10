# GitHub Actions Utils MCP Server

Connect AI tools directly to GitHub Actions. This MCP server gives AI agents, assistants, and chatbots the ability to fetch and analyze GitHub Action definitions, explore available inputs and outputs, and understand action configurations—all through natural language interactions.

### Use Cases

- **Workflow Development**: Quickly discover GitHub Action parameters while building CI/CD workflows
- **Action Discovery**: Explore available inputs, outputs, and configuration options for any public GitHub Action
- **Documentation**: Automatically generate documentation for actions used in your workflows
- **Migration & Updates**: Understand parameter changes when upgrading action versions
- **Validation**: Verify action configurations before deploying workflows

Built for developers who want to enhance their AI tools with GitHub Actions context, from simple action queries to complex workflow generation.

---

## Installation

### Prerequisites

1. A compatible MCP host application:
   - VS Code 1.101+ with GitHub Copilot
   - Claude Desktop
   - Cursor IDE
   - Windsurf IDE
   - Any MCP-compatible client

2. One of the following:
   - Pre-built binary (recommended)
   - Docker
   - Go 1.25+ (for building from source)

### Quick Install

#### Option 1: Pre-built Binary (Recommended)

**Download the latest release:**

Visit [GitHub Releases](https://github.com/techprimate/github-actions-utils-cli/releases/latest) and download the appropriate binary for your platform:

- **macOS**: `github-actions-utils-cli-darwin-amd64` (Intel) or `github-actions-utils-cli-darwin-arm64` (Apple Silicon)
- **Linux**: `github-actions-utils-cli-linux-amd64` or `github-actions-utils-cli-linux-arm64`
- **Windows**: `github-actions-utils-cli-windows-amd64.exe`

**Install:**

```bash
# macOS/Linux - move to a directory in your PATH
sudo mv github-actions-utils-cli-* /usr/local/bin/github-actions-utils-cli
sudo chmod +x /usr/local/bin/github-actions-utils-cli

# Verify installation
github-actions-utils-cli --version
```

#### Option 2: Docker

```bash
# Pull the image
docker pull ghcr.io/techprimate/github-actions-utils-cli:latest

# Test it works
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | docker run -i --rm ghcr.io/techprimate/github-actions-utils-cli:latest mcp
```

#### Option 3: Build from Source

```bash
# Clone and build
git clone https://github.com/techprimate/github-actions-utils-cli.git
cd github-actions-utils-cli
make build

# Binary will be at ./dist/github-actions-utils-cli
```

---

## Configuration

### VS Code with GitHub Copilot

**Prerequisites**: VS Code 1.101+ with GitHub Copilot installed

**Setup**:

1. Open VS Code settings (Cmd/Ctrl + ,)
2. Search for "MCP"
3. Click "Edit in settings.json"
4. Add the server configuration:

<table>
<tr><th>Using Binary</th><th>Using Docker</th></tr>
<tr valign=top>
<td>

```json
{
  "github.copilot.chat.mcp.servers": {
    "github-actions-utils": {
      "command": "/usr/local/bin/github-actions-utils-cli",
      "args": ["mcp"]
    }
  }
}
```

</td>
<td>

```json
{
  "github.copilot.chat.mcp.servers": {
    "github-actions-utils": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "ghcr.io/techprimate/github-actions-utils-cli:latest",
        "mcp"
      ]
    }
  }
}
```

</td>
</tr>
</table>

5. Restart VS Code
6. Open Copilot Chat and toggle "Agent mode" to activate MCP tools

**Troubleshooting**:

- Ensure VS Code is version 1.101 or later
- Verify GitHub Copilot extension is installed and activated
- Check that the binary path is correct: `which github-actions-utils-cli`

### Claude Desktop

**Prerequisites**: [Claude Desktop](https://claude.ai/download) installed

**Setup**:

1. Locate your Claude configuration file:
   - **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
   - **Linux**: `~/.config/Claude/claude_desktop_config.json`

2. Add the server configuration:

<table>
<tr><th>Using Binary</th><th>Using Docker</th></tr>
<tr valign=top>
<td>

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

</td>
<td>

```json
{
  "mcpServers": {
    "github-actions-utils": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "ghcr.io/techprimate/github-actions-utils-cli:latest",
        "mcp"
      ]
    }
  }
}
```

</td>
</tr>
</table>

3. Restart Claude Desktop
4. Look for the 🔌 icon in the bottom right to verify the server is connected

**Troubleshooting**:

- Verify the config file is valid JSON
- Check Claude Desktop logs for connection errors
- Ensure the binary path is absolute, not relative

### Cursor IDE

**Prerequisites**: [Cursor IDE](https://cursor.sh/) installed

**Setup**:

1. Open Cursor Settings (Cmd/Ctrl + ,)
2. Navigate to "Cursor Settings" → "MCP"
3. Click "Edit Config"
4. Add the server configuration:

<table>
<tr><th>Using Binary</th><th>Using Docker</th></tr>
<tr valign=top>
<td>

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

</td>
<td>

```json
{
  "mcpServers": {
    "github-actions-utils": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "ghcr.io/techprimate/github-actions-utils-cli:latest",
        "mcp"
      ]
    }
  }
}
```

</td>
</tr>
</table>

5. Restart Cursor
6. The MCP tools will be available in the AI chat

### Windsurf IDE

**Prerequisites**: [Windsurf IDE](https://codeium.com/windsurf) installed

**Setup**:

1. Open Windsurf Settings
2. Navigate to "MCP Servers"
3. Add new server configuration:

<table>
<tr><th>Using Binary</th><th>Using Docker</th></tr>
<tr valign=top>
<td>

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

</td>
<td>

```json
{
  "mcpServers": {
    "github-actions-utils": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "ghcr.io/techprimate/github-actions-utils-cli:latest",
        "mcp"
      ]
    }
  }
}
```

</td>
</tr>
</table>

4. Restart Windsurf

### Other MCP Clients

For other MCP-compatible clients, use the standard MCP configuration format:

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

Refer to your MCP client's documentation for the specific configuration file location.

---

## Available Tools

### get_action_parameters

Fetches and parses a GitHub Action's `action.yml` or `action.yaml` file, returning complete metadata about inputs, outputs, and configuration.

**Parameters:**

| Parameter   | Type   | Required | Description                                                |
| ----------- | ------ | -------- | ---------------------------------------------------------- |
| `actionRef` | string | Yes      | GitHub Action reference in the format `owner/repo@version` |

**Examples:**

```
Can you show me the parameters for actions/checkout@v5?

What inputs does actions/setup-node@v4 accept?

Explain the outputs of docker/build-push-action@v6
```

**Response Structure:**

```json
{
  "name": "Checkout",
  "description": "Checkout a Git repository at a particular version",
  "inputs": {
    "repository": {
      "description": "Repository name with owner. For example, actions/checkout",
      "required": false,
      "default": "${{ github.repository }}"
    },
    "ref": {
      "description": "The branch, tag or SHA to checkout",
      "required": false
    },
    "token": {
      "description": "Personal access token (PAT) used to fetch the repository",
      "required": false,
      "default": "${{ github.token }}"
    }
  },
  "outputs": {
    "ref": {
      "description": "The branch, tag or SHA that was checked out"
    }
  },
  "runs": {
    "using": "node24",
    "main": "dist/index.js"
  },
  "branding": {
    "icon": "download",
    "color": "blue"
  }
}
```

### get_readme

Fetches the README.md file from a GitHub repository, useful for understanding how to use actions or exploring their documentation.

**Parameters:**

| Parameter | Type   | Required | Description                                                                                             |
| --------- | ------ | -------- | ------------------------------------------------------------------------------------------------------- |
| `repoRef` | string | Yes      | GitHub repository reference in the format `owner/repo[@ref]`. If no ref is provided, defaults to `main` |

**Examples:**

```
Can you get the README for actions/checkout?

Show me the documentation for docker/build-push-action@v6

What does the README say about github/github-mcp-server?
```

**Response:**

Returns the full README content as markdown text.

---

## Example Workflows

### Discovering Action Parameters

**User**: "What are all the inputs for actions/setup-python@v5?"

**AI Response**: Uses `get_action_parameters` to fetch and explain all available inputs including `python-version`, `cache`, `architecture`, etc.

### Building a Workflow

**User**: "Help me create a workflow that checks out code, sets up Node.js 20, and runs tests"

**AI**: Uses `get_action_parameters` for `actions/checkout@v5` and `actions/setup-node@v4` to understand the correct parameters and generate a complete workflow file.

### Comparing Action Versions

**User**: "What changed between actions/upload-artifact@v3 and @v4?"

**AI**: Fetches both versions using `get_action_parameters` and highlights the differences in inputs, outputs, and behavior.

### Exploring New Actions

**User**: "Show me how to use aws-actions/configure-aws-credentials"

**AI**: Uses `get_readme` to fetch documentation and `get_action_parameters` to understand all configuration options.

---

## Telemetry

This project uses Sentry for error tracking and performance monitoring to help improve the tool.

**Disable telemetry:**

```bash
export TELEMETRY_ENABLED=false
```

Add this to your shell profile (`.bashrc`, `.zshrc`, etc.) to make it permanent.

---

## Development

### Prerequisites

- Go 1.25 or later
- Make
- [dprint](https://dprint.dev/) (for formatting)

### Setup

```bash
# Clone the repository
git clone https://github.com/techprimate/github-actions-utils-cli.git
cd github-actions-utils-cli

# Install dependencies
make init

# Build
make build

# Run tests
make test

# Format code
make format

# Run static analysis (vet, staticcheck, govulncheck)
make analyze
```

### Project Structure

```
.
├── cmd/cli/              # CLI entry point with main.go
├── internal/
│   ├── cli/
│   │   ├── cmd/          # Cobra commands (root, mcp)
│   │   └── mcp/          # MCP server and tool handlers
│   ├── github/           # GitHub Actions fetcher and parser
│   └── logging/          # Multi-handler for Sentry integration
├── .github/workflows/    # CI/CD pipelines
├── docs/                 # Documentation
├── Makefile             # Build commands
└── README.md            # This file
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Test specific package
go test ./internal/github/

# Run with race detection
go test -race ./...
```

### Manual MCP Testing

**Important**: MCP servers communicate via JSON-RPC over stdin/stdout and are designed to be used by MCP clients, not directly from the command line. The server is working correctly when no output appears - it's waiting for proper JSON-RPC formatted input from an MCP client.

To verify the server is working:

```bash
# Build the CLI
make build

# Test that the server starts (it should produce an error about invalid message format)
echo '{}' | ./dist/github-actions-utils-cli mcp
# Output: Error: invalid message version tag ""; expected "2.0"
# This confirms the server is reading stdin and validating JSON-RPC messages

# Verify binary is executable
./dist/github-actions-utils-cli --version
```

**To actually test the tools**, configure the server in an MCP client (VS Code, Claude Desktop, Cursor, etc.) using the configuration examples above, then use natural language queries:

- "Show me the parameters for actions/checkout@v5"
- "Get the README for github/github-mcp-server"

The MCP client will handle the JSON-RPC communication and present results in a user-friendly format.

### Adding New Tools

See [AGENTS.md](./AGENTS.md) for detailed instructions on adding new MCP tools.

Quick overview:

1. Define arguments struct in `internal/cli/mcp/tools.go`
2. Implement handler function
3. Register tool in `internal/cli/mcp/server.go`
4. Add tests
5. Update README

---

## Contributing

Contributions are welcome! Please see our contributing guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes using [Conventional Commits](https://www.conventionalcommits.org/)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add support for private actions
fix: handle action.yaml extension
docs: update installation instructions
```

---

## Security

### Reporting Security Issues

Please report security vulnerabilities to [security@techprimate.com](mailto:security@techprimate.com).

### Security Best Practices

- This tool only fetches publicly accessible GitHub Action definitions
- No authentication tokens are required or stored
- All network requests go through GitHub's public raw content CDN
- Sentry telemetry can be disabled via environment variable

---

## License

This project is licensed under the terms of the MIT License. See [LICENSE](./LICENSE) for details.

---

## Support

- 📖 [Documentation](./docs/)
- 💬 [GitHub Discussions](https://github.com/techprimate/github-actions-utils-cli/discussions)
- 🐛 [Issue Tracker](https://github.com/techprimate/github-actions-utils-cli/issues)
- 🌟 [Star on GitHub](https://github.com/techprimate/github-actions-utils-cli)

---

**Made with ❤️ by [techprimate](https://github.com/techprimate)**
