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
