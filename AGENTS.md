# AGENTS.md

## Project Overview

GitHub Actions Utils CLI is an MCP (Model Context Protocol) server that provides tools for working with GitHub Actions. It exposes a single tool `get_action_parameters` that fetches and parses action.yml files from GitHub Actions repositories.

**Key characteristics:**

- Simple Go CLI application
- Single MCP tool for fetching GitHub Actions metadata
- No database, no API server, no complex dependencies
- Designed for use with AI coding agents (Claude, Cursor, Cline, etc.)

## Setup Commands

```bash
# Initialize project (installs dependencies)
make init

# Install dependencies only
make install

# Build the CLI binary
make build
# Output: dist/github-actions-utils-cli

# Run the MCP server
make mcp
# Or directly: ./dist/github-actions-utils-cli mcp
```

## Development Commands

```bash
# Build (faster for local development)
make build

# Run tests
make test

# Format code (runs go fmt and dprint)
make format

# Run static analysis (vet, staticcheck, govulncheck)
make analyze

# Update dependencies
make upgrade-deps
```

## Project Structure

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
├── go.mod                # Go dependencies
├── Makefile             # Build commands
└── README.md            # User documentation
```

## Code Style Guidelines

### Go Conventions

- **Package naming**: Short, lowercase, single word (e.g., `github`, `mcp`, `cmd`)
- **File naming**: Use snake_case (e.g., `actions.go`, `server.go`)
- **Import order**: Standard library, external deps, internal packages (separated by blank lines)
- **Error handling**: Always check and wrap errors with context using `fmt.Errorf("context: %w", err)`

### Go Import Organization

```go
import (
    // 1. Standard library
    "context"
    "fmt"
    "log/slog"

    // 2. External dependencies
    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"

    // 3. Internal packages
    "github.com/techprimate/github-actions-utils-cli/internal/github"
)
```

### MCP Server Patterns

- **Silent logging**: MCP uses stdin/stdout for JSON-RPC, so always use `slog.New(slog.NewTextHandler(io.Discard, nil))` for MCP server logger
- **Tool handlers**: Implement in `internal/cli/mcp/tools.go` with clear argument structs
- **Error handling**: Return descriptive errors from tool handlers, they're sent to the client

## Testing Instructions

```bash
# Run all tests
make test

# Run with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/github/

# Run with race detection
go test -race ./...

# Test MCP server manually
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./dist/github-actions-utils-cli mcp
```

### Manual Testing GitHub Actions Fetcher

```bash
# Build first
make build

# Test fetching an action
# Example: create a simple test file
cat > test_fetch.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "github.com/techprimate/github-actions-utils-cli/internal/github"
)

func main() {
    svc := github.NewActionsService()
    params, err := svc.GetActionParameters("actions/checkout@v5")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Fetched %v inputs\n", len(params["inputs"].(map[string]interface{})))
}
EOF

go run test_fetch.go
rm test_fetch.go
```

## Commit Message Format

This project uses [Conventional Commits 1.0.0](https://www.conventionalcommits.org/).

**Format:**

```
<type>([optional scope]): <description>

[optional body]

[optional footer(s)]
```

**Types:**

- `feat:` - New feature (MINOR version bump)
- `fix:` - Bug fix (PATCH version bump)
- `build:` - Build system or dependencies
- `chore:` - Routine tasks, maintenance
- `ci:` - CI configuration changes
- `docs:` - Documentation only
- `style:` - Code style (formatting, etc.)
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `test:` - Adding or updating tests

**Breaking changes:**

- Add `!` after type: `feat!:` or `feat(mcp)!:`
- Or use footer: `BREAKING CHANGE: description`

**Examples:**

```
feat(mcp): add get_action_parameters tool
fix: handle 404 errors when fetching action.yml
docs: update README with installation instructions
refactor(github): simplify YAML parsing logic
```

**Important:** Never mention AI assistants or Claude in commit messages.

## Pull Request Guidelines

### Branch Naming

- `feature/` - New features
- `fix/` - Bug fixes
- `refactor/` - Code refactoring
- `docs/` - Documentation only
- `chore/` - Maintenance tasks

### Creating PRs

```bash
# Create feature branch
git checkout -b feature/add-new-tool

# Make changes and commit
git add -A
git commit -m "feat: add new MCP tool"

# Push branch
git push -u origin feature/add-new-tool

# Create PR (simplest - uses commit messages)
gh pr create --fill

# Or open browser to add detailed description
gh pr create --web
```

### PR Checklist

Before submitting:

- [ ] Code builds successfully (`make build`)
- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`make format`)
- [ ] Static analysis passes (`make analyze`)
- [ ] Commit messages follow Conventional Commits
- [ ] PR title follows Conventional Commits format

### PR Size Guidelines

- **Small**: 1-100 lines (ideal)
- **Medium**: 100-500 lines (good)
- **Large**: 500-1000 lines (needs extra context)
- **Too Large**: 1000+ lines (consider splitting)

## CI/CD Workflows

All workflows use GitHub Actions:

- **analyze.yml** - Runs `make analyze` on every PR/push to main
- **build.yml** - Runs `make build` to verify compilation
- **test.yml** - Runs `make test` for all tests
- **format.yml** - Checks code formatting with `make format`
- **release.yml** - Creates GitHub releases with signed binaries (triggered by tags or main branch)
- **build-binaries.yml** - Reusable workflow that builds for all platforms and signs macOS binaries
- **build-cli-docker.yml** - Builds and publishes Docker images to GitHub Container Registry

**Platforms built:**

- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon) - **code signed and notarized**
- Windows: amd64

## Docker Deployment

The project provides Docker images for easy deployment and containerized usage.

**Docker commands:**

```bash
# Build Linux binaries for Docker
make build-linux

# Build Docker image locally
make docker-build

# Test Docker image
make docker-test

# Run Docker container
make docker-run
```

**Image details:**

- **Registries**:
  - `ghcr.io/techprimate/github-actions-utils-cli` (GitHub Container Registry)
  - `docker.io/techprimate/github-actions-utils-cli` (Docker Hub)
- **Base image**: `buildpack-deps:bookworm`
- **Platforms**: linux/amd64, linux/arm64

**Using the Docker image:**

```bash
# Pull latest image (from GitHub Container Registry)
docker pull ghcr.io/techprimate/github-actions-utils-cli:latest

# Or from Docker Hub
docker pull docker.io/techprimate/github-actions-utils-cli:latest

# Run MCP server
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | \
  docker run -i --rm ghcr.io/techprimate/github-actions-utils-cli:latest mcp

# Use with MCP clients (Claude, Cursor, etc.)
# Configure in MCP client settings:
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

**Publishing images:**

Docker images are automatically built and published by the `build-cli-docker.yml` workflow:

- On **main branch push**: Tagged as `latest`
- On **version tags** (v1.0.0): Tagged as `1.0.0`, `1.0`, and `1`
- On **pull requests**: Built but not pushed (validation only)

**Manual Docker build:**

```bash
# Build for specific platform
docker build --platform linux/amd64 -t github-actions-utils-cli:latest .

# Build multi-platform (requires buildx)
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t github-actions-utils-cli:latest \
  --load \
  .
```

## Sentry Integration

This project uses Sentry for error tracking and monitoring.

**Disable telemetry:**

```bash
export TELEMETRY_ENABLED=false
```

**Sentry DSN:** `https://445c4c2185068fa980b83ddbe4bf1fd7@o188824.ingest.us.sentry.io/4510306572828672`

### MCP Tracing

The project includes automatic tracing for MCP tool calls following [OpenTelemetry MCP Semantic Conventions](https://github.com/open-telemetry/semantic-conventions/pull/2083).

**Key Features:**

- Automatic span creation for tool calls
- Detailed attributes (method name, tool name, arguments, results)
- Error capture and correlation
- Compatible with Sentry performance monitoring

**Implementation:**

All MCP tools are automatically wrapped with Sentry tracing using the `WithSentryTracing` wrapper:

```go
mcp.AddTool(server, &mcp.Tool{
    Name:        "my_tool",
    Description: "Tool description",
}, WithSentryTracing("my_tool", m.handleMyTool))
```

**Span Attributes:**

Each tool call creates a span with:

- Operation: `mcp.server`
- Name: `tools/call {tool_name}`
- Attributes: `mcp.method.name`, `mcp.tool.name`, `mcp.request.argument.*`, `mcp.tool.result.*`

See `docs/MCP_TRACING.md` for complete documentation on span conventions, attributes, and examples.

## Security Considerations

- **No credentials in code**: Never commit API keys or certificates
- **GitHub secrets**: Required secrets for releases:
  - `DEVELOPER_ID_P12_BASE64` - Apple Developer certificate
  - `DEVELOPER_ID_PASSWORD` - Certificate password
  - `APPLE_NOTARIZATION_APPLE_ID_PASSWORD` - App-specific password
- **Code signing**: macOS binaries are signed with Developer ID certificate
- **Notarization**: All macOS binaries are notarized by Apple

## Common Tasks

### Adding a New MCP Tool

1. Define arguments struct in `internal/cli/mcp/tools.go`:

```go
type NewToolArgs struct {
    Parameter string `json:"parameter" jsonschema:"Description"`
}
```

2. Implement handler in `internal/cli/mcp/tools.go`:

```go
func (m *MCPServer) handleNewTool(ctx context.Context, req *mcp.CallToolRequest, args NewToolArgs) (*mcp.CallToolResult, any, error) {
    // Implementation
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{Text: "Result"},
        },
    }, data, nil
}
```

3. Register tool in `internal/cli/mcp/server.go` with Sentry tracing:

```go
mcp.AddTool(server, &mcp.Tool{
    Name:        "new_tool",
    Description: "Description of what the tool does",
}, WithSentryTracing("new_tool", m.handleNewTool))
```

4. Test the tool:

```bash
make build
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./dist/github-actions-utils-cli mcp | jq
```

### Updating Dependencies

```bash
# Update all dependencies
make upgrade-deps

# Or manually
go get -u ./...
go mod tidy

# Verify everything still works
make test
make build
```

### Creating a Release

Releases are automated via GitHub Actions:

**Stable release:**

```bash
# Tag the commit
git tag v1.2.3
git push origin v1.2.3

# Workflow automatically:
# 1. Builds binaries for all platforms
# 2. Signs and notarizes macOS binaries
# 3. Creates GitHub release with binaries
```

**Pre-release (latest):**

```bash
# Push to main branch
git push origin main

# Workflow automatically:
# 1. Creates/updates "latest" tag
# 2. Builds and publishes binaries
# 3. Marks as pre-release
```

## Debugging Tips

### MCP Server Not Working

1. Check if server starts:

```bash
./dist/github-actions-utils-cli mcp
# Should not output anything (uses stdio)
```

2. Test with simple request:

```bash
printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}\n' | ./dist/github-actions-utils-cli mcp
```

3. Check MCP client configuration:

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

### GitHub Actions Fetching Fails

1. Check if URL is correct:

```bash
# Format: https://raw.githubusercontent.com/{owner}/{repo}/refs/tags/{version}/action.yml
curl -I https://raw.githubusercontent.com/actions/checkout/refs/tags/v5/action.yml
```

2. Verify action reference parsing:

```bash
# Should split into: owner=actions, repo=checkout, version=v5
echo "actions/checkout@v5"
```

## File Organization

### When to Create New Packages

✅ Create a new package when:

- Code has a distinct responsibility
- Code will be reused across features
- Package would have 5+ related files
- Clear interface boundary exists

❌ Don't create a package when:

- It would only have 1-2 files
- It's just for organization (use subdirectories)
- No clear interface boundary

### Test File Placement

Place test files next to the code they test:

```
internal/github/
├── actions.go
└── actions_test.go
```

## Large Dataset Handling

This project fetches action.yml files from GitHub. These are typically small (< 10KB), but if you need to handle larger responses:

1. Stream the response instead of loading into memory
2. Add timeout handling for slow responses
3. Consider caching frequently accessed actions (future enhancement)

## Future Enhancements

Potential features to consider:

- Cache for frequently accessed actions
- Support for fetching from branches (not just tags)
- Tool to list all versions of an action
- Tool to search for actions by name/description
- Batch fetching of multiple actions
- Local cache management commands
