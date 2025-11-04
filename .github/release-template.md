## GitHub Actions Utils CLI v{{VERSION}}

### Installation

#### macOS

```bash
# Apple Silicon (M1/M2/M3)
curl -L -o github-actions-utils-cli https://github.com/{{REPOSITORY}}/releases/download/v{{VERSION}}/github-actions-utils-cli-darwin-arm64
chmod +x github-actions-utils-cli
sudo mv github-actions-utils-cli /usr/local/bin/

# Intel
curl -L -o github-actions-utils-cli https://github.com/{{REPOSITORY}}/releases/download/v{{VERSION}}/github-actions-utils-cli-darwin-amd64
chmod +x github-actions-utils-cli
sudo mv github-actions-utils-cli /usr/local/bin/
```

#### Linux

```bash
# AMD64
curl -L -o github-actions-utils-cli https://github.com/{{REPOSITORY}}/releases/download/v{{VERSION}}/github-actions-utils-cli-linux-amd64
chmod +x github-actions-utils-cli
sudo mv github-actions-utils-cli /usr/local/bin/

# ARM64
curl -L -o github-actions-utils-cli https://github.com/{{REPOSITORY}}/releases/download/v{{VERSION}}/github-actions-utils-cli-linux-arm64
chmod +x github-actions-utils-cli
sudo mv github-actions-utils-cli /usr/local/bin/
```

#### Windows

Download `github-actions-utils-cli-windows-amd64.exe` and add it to your PATH.

### Usage

```bash
# Run as MCP server
github-actions-utils-cli mcp
```

See the [README](https://github.com/{{REPOSITORY}}/blob/main/README.md) for more details on configuring the MCP server with Claude Desktop, Cursor, or other MCP clients.

### Checksums

See `checksums.txt` for SHA256 checksums of all binaries.
