## 🚧 Latest Development Build

**This is a pre-release build from the latest `main` branch.**

- **Commit**: {{COMMIT_SHA}}
- **Built**: {{BUILD_DATE}}
- **Version**: {{VERSION}}

⚠️ **Warning**: This is an unstable development build. For production use, download a stable release instead.

### Installation

#### macOS

```bash
# Apple Silicon (M1/M2/M3)
curl -L -o github-actions-utils-cli https://github.com/{{REPOSITORY}}/releases/download/latest/github-actions-utils-cli-darwin-arm64
chmod +x github-actions-utils-cli
sudo mv github-actions-utils-cli /usr/local/bin/

# Intel
curl -L -o github-actions-utils-cli https://github.com/{{REPOSITORY}}/releases/download/latest/github-actions-utils-cli-darwin-amd64
chmod +x github-actions-utils-cli
sudo mv github-actions-utils-cli /usr/local/bin/
```

#### Linux

```bash
# AMD64
curl -L -o github-actions-utils-cli https://github.com/{{REPOSITORY}}/releases/download/latest/github-actions-utils-cli-linux-amd64
chmod +x github-actions-utils-cli
sudo mv github-actions-utils-cli /usr/local/bin/

# ARM64
curl -L -o github-actions-utils-cli https://github.com/{{REPOSITORY}}/releases/download/latest/github-actions-utils-cli-linux-arm64
chmod +x github-actions-utils-cli
sudo mv github-actions-utils-cli /usr/local/bin/
```

#### Windows

Download `github-actions-utils-cli-windows-amd64.exe` and add it to your PATH.

### What's New?

See the [commit history](https://github.com/{{REPOSITORY}}/commits/main) for recent changes.

### Checksums

See `checksums.txt` for SHA256 checksums of all binaries.
