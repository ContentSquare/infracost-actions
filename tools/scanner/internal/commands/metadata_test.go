package commands

import (
	"os"
	"testing"

	"github.com/infracost/actions/tools/scanner/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveVCSProvider(t *testing.T) {
	tests := []struct {
		name          string
		provider      string
		githubActions string
		want          string
		wantErr       string
	}{
		{
			name:     "github",
			provider: "github",
			want:     "github",
		},
		{
			name:     "gitlab",
			provider: "gitlab",
			want:     "gitlab",
		},
		{
			name:     "azure repos",
			provider: "azure_repos",
			want:     "azure_repos",
		},
		{
			name:     "bitbucket",
			provider: "bitbucket",
			want:     "bitbucket",
		},
		{
			name:          "unset falls back to github on github actions",
			githubActions: "true",
			want:          "github",
		},
		{
			name:    "unset elsewhere",
			wantErr: "cannot determine the VCS provider: set INFRACOST_VCS_PROVIDER to one of github, gitlab, azure_repos, bitbucket",
		},
		{
			name:     "typo",
			provider: "githbu",
			wantErr:  `unsupported VCS provider "githbu": set INFRACOST_VCS_PROVIDER to one of github, gitlab, azure_repos, bitbucket`,
		},
		{
			// The vcs module's package name, and the dashboard's ci platform
			// namespaces — all near misses worth rejecting explicitly.
			name:     "vcs module package name",
			provider: "azure",
			wantErr:  `unsupported VCS provider "azure": set INFRACOST_VCS_PROVIDER to one of github, gitlab, azure_repos, bitbucket`,
		},
		{
			name:     "ci platform name",
			provider: "gitlab_ci",
			wantErr:  `unsupported VCS provider "gitlab_ci": set INFRACOST_VCS_PROVIDER to one of github, gitlab, azure_repos, bitbucket`,
		},
		{
			name:     "wrong case",
			provider: "GitHub",
			wantErr:  `unsupported VCS provider "GitHub": set INFRACOST_VCS_PROVIDER to one of github, gitlab, azure_repos, bitbucket`,
		},
		{
			name:          "an explicit typo is not rescued by the fallback",
			provider:      "githbu",
			githubActions: "true",
			wantErr:       `unsupported VCS provider "githbu": set INFRACOST_VCS_PROVIDER to one of github, gitlab, azure_repos, bitbucket`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// make test-unit runs on GitHub Actions, where this is already set.
			t.Setenv("GITHUB_ACTIONS", tt.githubActions)
			if tt.githubActions == "" {
				require.NoError(t, os.Unsetenv("GITHUB_ACTIONS"))
			}

			got, err := resolveVCSProvider(&config.Config{VCSProvider: tt.provider})

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				assert.Empty(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
