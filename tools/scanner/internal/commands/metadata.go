package commands

import (
	"fmt"
	"os"
	"slices"

	"github.com/infracost/actions/tools/scanner/internal/config"
	"github.com/infracost/actions/tools/scanner/internal/vcsurl"
)

// resolveVCSProvider falls back to github when GITHUB_ACTIONS is set, which
// implies GitHub including Enterprise Server. Needed because pinned action
// versions run new binaries with old YAML that sets no provider.
func resolveVCSProvider(cfg *config.Config) (string, error) {
	provider := cfg.VCSProvider
	if provider == "" {
		if os.Getenv("GITHUB_ACTIONS") == "" {
			return "", fmt.Errorf("cannot determine the VCS provider: set INFRACOST_VCS_PROVIDER to one of %s", vcsurl.ProviderList())
		}
		return vcsurl.ProviderGitHub, nil
	}

	// Rejected here rather than at the URL builder so a typo fails before any
	// scan runs, for all three commands.
	if !slices.Contains(vcsurl.Providers, provider) {
		return "", fmt.Errorf("unsupported VCS provider %q: set INFRACOST_VCS_PROVIDER to one of %s", provider, vcsurl.ProviderList())
	}

	return provider, nil
}
