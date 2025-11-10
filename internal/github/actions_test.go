package github

import (
	"testing"
)

func TestParseActionRef(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantOwner   string
		wantRepo    string
		wantVersion string
		wantErr     bool
	}{
		{
			name:        "valid action reference",
			input:       "actions/checkout@v5",
			wantOwner:   "actions",
			wantRepo:    "checkout",
			wantVersion: "v5",
			wantErr:     false,
		},
		{
			name:        "valid action with trailing newline",
			input:       "actions/checkout@v5\n",
			wantOwner:   "actions",
			wantRepo:    "checkout",
			wantVersion: "v5",
			wantErr:     false,
		},
		{
			name:        "valid action with leading and trailing whitespace",
			input:       "  actions/checkout@v5  ",
			wantOwner:   "actions",
			wantRepo:    "checkout",
			wantVersion: "v5",
			wantErr:     false,
		},
		{
			name:        "valid action with tabs and newlines",
			input:       "\t\nactions/checkout@v5\n\t",
			wantOwner:   "actions",
			wantRepo:    "checkout",
			wantVersion: "v5",
			wantErr:     false,
		},
		{
			name:        "complex action reference with whitespace",
			input:       "kula-app/wait-for-services-action@v1\n",
			wantOwner:   "kula-app",
			wantRepo:    "wait-for-services-action",
			wantVersion: "v1",
			wantErr:     false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only whitespace",
			input:   "   \n\t  ",
			wantErr: true,
		},
		{
			name:    "missing version",
			input:   "actions/checkout",
			wantErr: true,
		},
		{
			name:    "missing repo",
			input:   "actions@v5",
			wantErr: true,
		},
		{
			name:    "invalid format - too many slashes",
			input:   "actions/github/checkout@v5",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseActionRef(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseActionRef() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseActionRef() unexpected error: %v", err)
				return
			}

			if got.Owner != tt.wantOwner {
				t.Errorf("ParseActionRef() Owner = %v, want %v", got.Owner, tt.wantOwner)
			}
			if got.Repo != tt.wantRepo {
				t.Errorf("ParseActionRef() Repo = %v, want %v", got.Repo, tt.wantRepo)
			}
			if got.Version != tt.wantVersion {
				t.Errorf("ParseActionRef() Version = %v, want %v", got.Version, tt.wantVersion)
			}
		})
	}
}

func TestParseRepoRef(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantOwner   string
		wantRepo    string
		wantVersion string
		wantErr     bool
	}{
		{
			name:        "valid repo reference with tag",
			input:       "actions/checkout@v5",
			wantOwner:   "actions",
			wantRepo:    "checkout",
			wantVersion: "v5",
			wantErr:     false,
		},
		{
			name:        "valid repo reference with branch",
			input:       "owner/repo@main",
			wantOwner:   "owner",
			wantRepo:    "repo",
			wantVersion: "main",
			wantErr:     false,
		},
		{
			name:        "valid repo reference with commit SHA",
			input:       "owner/repo@abc123def456",
			wantOwner:   "owner",
			wantRepo:    "repo",
			wantVersion: "abc123def456",
			wantErr:     false,
		},
		{
			name:        "repo without ref defaults to main",
			input:       "owner/repo",
			wantOwner:   "owner",
			wantRepo:    "repo",
			wantVersion: "main",
			wantErr:     false,
		},
		{
			name:        "valid repo with trailing newline",
			input:       "owner/repo@main\n",
			wantOwner:   "owner",
			wantRepo:    "repo",
			wantVersion: "main",
			wantErr:     false,
		},
		{
			name:        "valid repo with whitespace",
			input:       "  owner/repo@develop  ",
			wantOwner:   "owner",
			wantRepo:    "repo",
			wantVersion: "develop",
			wantErr:     false,
		},
		{
			name:        "repo without ref and whitespace",
			input:       "  owner/repo\n",
			wantOwner:   "owner",
			wantRepo:    "repo",
			wantVersion: "main",
			wantErr:     false,
		},
		{
			name:        "complex repo name with hyphen",
			input:       "techprimate/github-actions-utils-cli@main",
			wantOwner:   "techprimate",
			wantRepo:    "github-actions-utils-cli",
			wantVersion: "main",
			wantErr:     false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only whitespace",
			input:   "   \n\t  ",
			wantErr: true,
		},
		{
			name:    "missing repo",
			input:   "owner@main",
			wantErr: true,
		},
		{
			name:    "invalid format - too many slashes",
			input:   "owner/group/repo@main",
			wantErr: true,
		},
		{
			name:    "invalid format - multiple @ symbols",
			input:   "owner/repo@main@extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRepoRef(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseRepoRef() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseRepoRef() unexpected error: %v", err)
				return
			}

			if got.Owner != tt.wantOwner {
				t.Errorf("ParseRepoRef() Owner = %v, want %v", got.Owner, tt.wantOwner)
			}
			if got.Repo != tt.wantRepo {
				t.Errorf("ParseRepoRef() Repo = %v, want %v", got.Repo, tt.wantRepo)
			}
			if got.Version != tt.wantVersion {
				t.Errorf("ParseRepoRef() Version = %v, want %v", got.Version, tt.wantVersion)
			}
		})
	}
}
