package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

// ActionsService provides functionality to fetch and parse GitHub Actions.
type ActionsService struct {
	httpClient *http.Client
}

// NewActionsService creates a new ActionsService.
func NewActionsService() *ActionsService {
	return &ActionsService{
		httpClient: &http.Client{},
	}
}

// ActionRef represents a parsed GitHub Action reference.
type ActionRef struct {
	Owner   string
	Repo    string
	Version string
}

// ParseActionRef parses an action reference string like "owner/repo@version".
// Examples:
//   - "actions/checkout@v5" -> {Owner: "actions", Repo: "checkout", Version: "v5"}
//   - "actions/setup-node@v4" -> {Owner: "actions", Repo: "setup-node", Version: "v4"}
func ParseActionRef(ref string) (*ActionRef, error) {
	if ref == "" {
		return nil, fmt.Errorf("action reference cannot be empty")
	}

	// Split by @ to separate repo from version
	parts := strings.Split(ref, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid action reference format: expected 'owner/repo@version', got '%s'", ref)
	}

	repoPath := parts[0]
	version := parts[1]

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

	return &ActionRef{
		Owner:   owner,
		Repo:    repo,
		Version: version,
	}, nil
}

// FetchActionYAML fetches the action.yml file from GitHub's raw content CDN.
// It constructs the URL in the format:
// https://raw.githubusercontent.com/{owner}/{repo}/refs/tags/{version}/action.yml
func (s *ActionsService) FetchActionYAML(owner, repo, version string) ([]byte, error) {
	// Construct URL to raw action.yml on GitHub
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/refs/tags/%s/action.yml",
		owner, repo, version)

	// Make HTTP GET request
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch action.yml: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("action.yml not found at %s (status: 404) - verify the action reference and version", url)
		}
		return nil, fmt.Errorf("failed to fetch action.yml from %s (status: %d)", url, resp.StatusCode)
	}

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read action.yml response: %w", err)
	}

	return data, nil
}

// ParseActionYAML parses YAML data into a map that can be JSON-encoded.
// This converts the action.yml structure into a JSON-compatible format.
func ParseActionYAML(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}

	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return result, nil
}

// GetActionParameters fetches and parses a GitHub Action's action.yml file.
// It takes an action reference (e.g., "actions/checkout@v5") and returns
// the parsed action.yml content as a JSON-compatible map.
func (s *ActionsService) GetActionParameters(actionRef string) (map[string]interface{}, error) {
	// Parse the action reference
	ref, err := ParseActionRef(actionRef)
	if err != nil {
		return nil, fmt.Errorf("invalid action reference: %w", err)
	}

	// Fetch the action.yml file
	yamlData, err := s.FetchActionYAML(ref.Owner, ref.Repo, ref.Version)
	if err != nil {
		return nil, err
	}

	// Parse YAML to map
	parsed, err := ParseActionYAML(yamlData)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

// GetActionParametersJSON is a convenience method that returns the action
// parameters as a JSON string instead of a map.
func (s *ActionsService) GetActionParametersJSON(actionRef string) (string, error) {
	params, err := s.GetActionParameters(actionRef)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return string(jsonData), nil
}
