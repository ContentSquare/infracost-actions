// Package vcsurl builds VCS web URLs from a repository's web URL.
package vcsurl

import (
	"fmt"
	"net/url"
	"strings"
)

// Provider names. Fixed by the v0.1 INFRACOST_VCS_PROVIDER contract, and
// deliberately not the vcs module's package names.
const (
	ProviderGitHub     = "github"
	ProviderGitLab     = "gitlab"
	ProviderAzureRepos = "azure_repos"
	ProviderBitbucket  = "bitbucket"
)

// Providers lists every provider PullRequest can build a URL for.
var Providers = []string{ProviderGitHub, ProviderGitLab, ProviderAzureRepos, ProviderBitbucket}

// ProviderList renders Providers for use in an error message.
func ProviderList() string {
	return strings.Join(Providers, ", ")
}

// PullRequest returns the web URL of a pull request for the given provider.
// repoURL is the repository's web URL; "" and no error when there is no PR.
func PullRequest(provider, repoURL string, number int) (string, error) {
	if repoURL == "" || number <= 0 {
		return "", nil
	}

	if err := checkWebURL(repoURL); err != nil {
		return "", err
	}
	base := strings.TrimSuffix(strings.TrimSuffix(repoURL, "/"), ".git")

	switch provider {
	case ProviderGitHub:
		return fmt.Sprintf("%s/pull/%d", base, number), nil
	case ProviderGitLab:
		// number must be the project-scoped iid, not the global merge request
		// id. A global id is still positive, so it 404s silently.
		return fmt.Sprintf("%s/-/merge_requests/%d", base, number), nil
	case ProviderAzureRepos:
		// Azure PR URLs hang off the repository, not the project, and a project
		// URL is otherwise indistinguishable from a repository one.
		if !strings.Contains(base, "/_git/") {
			return "", fmt.Errorf("repo URL %q must be an Azure Repos repository URL containing /_git/", repoURL)
		}
		return fmt.Sprintf("%s/pullrequest/%d", base, number), nil
	case ProviderBitbucket:
		return fmt.Sprintf("%s/pull-requests/%d", base, number), nil
	default:
		return "", fmt.Errorf("cannot build a pull request URL for VCS provider %q: must be one of %s", provider, ProviderList())
	}
}

// checkWebURL constrains the metadata URL, not the transport: nothing in the
// binary clones, so an SSH remote here only yields an unopenable link.
func checkWebURL(repoURL string) error {
	u, err := url.Parse(repoURL)

	// Checked before any error that echoes repoURL: a credentialed clone URL
	// would otherwise leak its token into the error, metadata and PR key.
	if err == nil && u.User != nil {
		if _, set := u.User.Password(); set {
			return fmt.Errorf("repo URL must not contain a password: pass the repository's web URL, not a credentialed clone URL")
		}
	}

	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("repo URL %q must be an http(s) web URL of the repository, not a clone URL", repoURL)
	}

	// The PR path is appended to the raw URL, so anything after it corrupts.
	if u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("repo URL %q must not contain a query or fragment", repoURL)
	}

	return nil
}
