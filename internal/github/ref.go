package github

import (
	"fmt"
	"io"
	"strings"
)

// Ref represents a parsed GitHub reference (repository or action).
type Ref struct {
	Owner   string
	Repo    string
	Version string // Can be a tag, branch, commit SHA, or version
}

// ParseRef parses a GitHub reference string like "owner/repo@version".
// If requireVersion is true, the @version part is mandatory.
// If requireVersion is false and no @version is provided, defaultVersion is used.
//
// Examples:
//   - "actions/checkout@v5" -> {Owner: "actions", Repo: "checkout", Version: "v5"}
//   - "owner/repo@main" -> {Owner: "owner", Repo: "repo", Version: "main"}
//   - "owner/repo" with defaultVersion="main" -> {Owner: "owner", Repo: "repo", Version: "main"}
func ParseRef(ref string, requireVersion bool, defaultVersion string) (*Ref, error) {
	// Trim whitespace (including newlines, spaces, tabs)
	ref = strings.TrimSpace(ref)

	if ref == "" {
		return nil, fmt.Errorf("reference cannot be empty")
	}

	var repoPath, version string

	// Split by @ to separate repo from version
	if strings.Contains(ref, "@") {
		parts := strings.Split(ref, "@")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid reference format: expected 'owner/repo@version' or 'owner/repo', got '%s'", ref)
		}
		repoPath = parts[0]
		version = parts[1]
	} else {
		// No @ found
		if requireVersion {
			return nil, fmt.Errorf("invalid reference format: expected 'owner/repo@version', got '%s'", ref)
		}
		repoPath = ref
		version = defaultVersion
	}

	// Split repo path by / to get owner and repo
	repoParts := strings.Split(repoPath, "/")
	if len(repoParts) != 2 {
		return nil, fmt.Errorf("invalid repository path: expected 'owner/repo', got '%s'", repoPath)
	}

	owner := repoParts[0]
	repo := repoParts[1]

	if owner == "" || repo == "" || version == "" {
		return nil, fmt.Errorf("owner, repo, and version must all be non-empty")
	}

	return &Ref{
		Owner:   owner,
		Repo:    repo,
		Version: version,
	}, nil
}

// FetchRawFile fetches a file from GitHub's raw content CDN.
// The urlPath should specify the path type and version:
//   - For tags: "refs/tags/{version}"
//   - For branches: "refs/heads/{branch}"
//   - For commits: "{sha}"
func (s *ActionsService) FetchRawFile(owner, repo, urlPath, filename string) ([]byte, error) {
	// Construct URL to raw file on GitHub
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s",
		owner, repo, urlPath, filename)

	// Make HTTP GET request
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", filename, err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return nil, fmt.Errorf("%s not found at %s (status: 404)", filename, url)
		}
		return nil, fmt.Errorf("failed to fetch %s from %s (status: %d)", filename, url, resp.StatusCode)
	}

	// Read response body
	data, err := readAllBody(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s response: %w", filename, err)
	}

	return data, nil
}

// readAllBody is a helper to read all data from an io.Reader.
func readAllBody(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
