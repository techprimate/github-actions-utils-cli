package github

import (
	"encoding/json"
	"fmt"
	"net/http"

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

// ParseActionRef parses an action reference string like "owner/repo@version".
// The version part is required for actions.
// Examples:
//   - "actions/checkout@v5" -> {Owner: "actions", Repo: "checkout", Version: "v5"}
//   - "actions/setup-node@v4" -> {Owner: "actions", Repo: "setup-node", Version: "v4"}
func ParseActionRef(ref string) (*Ref, error) {
	return ParseRef(ref, true, "")
}

// FetchActionYAML fetches the action.yml or action.yaml file from GitHub's raw content CDN.
// It tries both common action file names in order of preference.
// It constructs the URL using tags format:
// https://raw.githubusercontent.com/{owner}/{repo}/refs/tags/{version}/action.yml
func (s *ActionsService) FetchActionYAML(owner, repo, version string) ([]byte, error) {
	// Try common action filenames in order of preference
	actionFilenames := []string{"action.yml", "action.yaml"}
	urlPath := fmt.Sprintf("refs/tags/%s", version)

	var lastErr error
	for _, filename := range actionFilenames {
		data, err := s.FetchRawFile(owner, repo, urlPath, filename)
		if err == nil {
			return data, nil
		}
		lastErr = err
	}

	// If we get here, none of the action files were found
	if lastErr != nil {
		return nil, fmt.Errorf("action.yml or action.yaml not found for %s/%s@%s: %w", owner, repo, version, lastErr)
	}
	return nil, fmt.Errorf("action.yml or action.yaml not found for %s/%s@%s", owner, repo, version)
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

// ParseRepoRef parses a repository reference string like "owner/repo@ref".
// The ref can be a tag, branch name, or commit SHA.
// If no ref is provided (e.g., "owner/repo"), it defaults to "main".
// Examples:
//   - "actions/checkout@v5" -> {Owner: "actions", Repo: "checkout", Version: "v5"}
//   - "owner/repo@main" -> {Owner: "owner", Repo: "repo", Version: "main"}
//   - "owner/repo" -> {Owner: "owner", Repo: "repo", Version: "main"}
func ParseRepoRef(ref string) (*Ref, error) {
	return ParseRef(ref, false, "main")
}

// FetchReadme fetches the README.md file from GitHub's raw content CDN.
// It tries multiple common README filenames in order of preference.
// The ref can be a branch name, tag, or commit SHA.
func (s *ActionsService) FetchReadme(owner, repo, ref string) (string, error) {
	// Try common README filenames in order of preference
	readmeNames := []string{"README.md", "readme.md", "Readme.md", "README", "readme"}
	urlPath := fmt.Sprintf("refs/heads/%s", ref)

	var lastErr error
	for _, filename := range readmeNames {
		data, err := s.FetchRawFile(owner, repo, urlPath, filename)
		if err == nil {
			return string(data), nil
		}
		lastErr = err
	}

	// If we get here, none of the README files were found
	if lastErr != nil {
		return "", fmt.Errorf("README not found in repository %s/%s@%s: %w", owner, repo, ref, lastErr)
	}
	return "", fmt.Errorf("README not found in repository %s/%s@%s", owner, repo, ref)
}

// GetReadme fetches a README.md file from a GitHub repository.
// It takes a repository reference (e.g., "owner/repo@main" or "owner/repo") and returns
// the README content as a string. If no ref is provided, it defaults to "main".
func (s *ActionsService) GetReadme(repoRef string) (string, error) {
	// Parse the repository reference
	ref, err := ParseRepoRef(repoRef)
	if err != nil {
		return "", fmt.Errorf("invalid repository reference: %w", err)
	}

	// Fetch the README file
	content, err := s.FetchReadme(ref.Owner, ref.Repo, ref.Version)
	if err != nil {
		return "", err
	}

	return content, nil
}
