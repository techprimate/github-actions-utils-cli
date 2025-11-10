package github

import (
	"testing"
)

func TestParseRef(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		requireVersion bool
		defaultVersion string
		wantOwner      string
		wantRepo       string
		wantVersion    string
		wantErr        bool
	}{
		{
			name:           "valid reference with version",
			input:          "actions/checkout@v5",
			requireVersion: false,
			defaultVersion: "main",
			wantOwner:      "actions",
			wantRepo:       "checkout",
			wantVersion:    "v5",
			wantErr:        false,
		},
		{
			name:           "valid reference without version, defaults to provided",
			input:          "owner/repo",
			requireVersion: false,
			defaultVersion: "main",
			wantOwner:      "owner",
			wantRepo:       "repo",
			wantVersion:    "main",
			wantErr:        false,
		},
		{
			name:           "valid reference without version, defaults to develop",
			input:          "owner/repo",
			requireVersion: false,
			defaultVersion: "develop",
			wantOwner:      "owner",
			wantRepo:       "repo",
			wantVersion:    "develop",
			wantErr:        false,
		},
		{
			name:           "valid reference with branch name",
			input:          "owner/repo@main",
			requireVersion: false,
			defaultVersion: "main",
			wantOwner:      "owner",
			wantRepo:       "repo",
			wantVersion:    "main",
			wantErr:        false,
		},
		{
			name:           "valid reference with commit SHA",
			input:          "owner/repo@abc123def456",
			requireVersion: false,
			defaultVersion: "main",
			wantOwner:      "owner",
			wantRepo:       "repo",
			wantVersion:    "abc123def456",
			wantErr:        false,
		},
		{
			name:           "reference with trailing newline",
			input:          "owner/repo@v1\n",
			requireVersion: false,
			defaultVersion: "main",
			wantOwner:      "owner",
			wantRepo:       "repo",
			wantVersion:    "v1",
			wantErr:        false,
		},
		{
			name:           "reference with whitespace",
			input:          "  owner/repo@v2  ",
			requireVersion: false,
			defaultVersion: "main",
			wantOwner:      "owner",
			wantRepo:       "repo",
			wantVersion:    "v2",
			wantErr:        false,
		},
		{
			name:           "complex repo name with hyphens",
			input:          "techprimate/github-actions-utils-cli@v1.0.0",
			requireVersion: false,
			defaultVersion: "main",
			wantOwner:      "techprimate",
			wantRepo:       "github-actions-utils-cli",
			wantVersion:    "v1.0.0",
			wantErr:        false,
		},
		{
			name:           "empty string",
			input:          "",
			requireVersion: false,
			defaultVersion: "main",
			wantErr:        true,
		},
		{
			name:           "only whitespace",
			input:          "   \n\t  ",
			requireVersion: false,
			defaultVersion: "main",
			wantErr:        true,
		},
		{
			name:           "missing version when required",
			input:          "owner/repo",
			requireVersion: true,
			defaultVersion: "",
			wantErr:        true,
		},
		{
			name:           "missing repo",
			input:          "owner@v1",
			requireVersion: false,
			defaultVersion: "main",
			wantErr:        true,
		},
		{
			name:           "invalid format - too many slashes",
			input:          "owner/group/repo@v1",
			requireVersion: false,
			defaultVersion: "main",
			wantErr:        true,
		},
		{
			name:           "invalid format - multiple @ symbols",
			input:          "owner/repo@v1@extra",
			requireVersion: false,
			defaultVersion: "main",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRef(tt.input, tt.requireVersion, tt.defaultVersion)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseRef() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseRef() unexpected error: %v", err)
				return
			}

			if got.Owner != tt.wantOwner {
				t.Errorf("ParseRef() Owner = %v, want %v", got.Owner, tt.wantOwner)
			}
			if got.Repo != tt.wantRepo {
				t.Errorf("ParseRef() Repo = %v, want %v", got.Repo, tt.wantRepo)
			}
			if got.Version != tt.wantVersion {
				t.Errorf("ParseRef() Version = %v, want %v", got.Version, tt.wantVersion)
			}
		})
	}
}
