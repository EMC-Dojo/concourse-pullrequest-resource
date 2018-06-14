package resource

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Github is
type Github interface {
	ListPRs() ([]*github.PullRequest, error)
}

// GithubClient is
type GithubClient struct {
	client     *github.Client
	owner      string
	repository string
}

// NewGithubClient is
func NewGithubClient(source Source) (*GithubClient, error) {
	var httpClient = &http.Client{}
	var ctx = context.TODO()

	if source.Insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	}

	if source.AccessToken != "" {
		var err error
		httpClient, err = oauthClient(ctx, source)
		if err != nil {
			return nil, err
		}
	}

	client := github.NewClient(httpClient)

	return &GithubClient{
		client,
		source.Owner,
		source.Repo,
	}, nil
}

// ListPRs is
func (g *GithubClient) ListPRs(opts *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	pulls, resp, err := g.client.PullRequests.List(context.TODO(), g.owner, g.repository, nil)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return pulls, nil
}

func oauthClient(ctx context.Context, source Source) (*http.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: source.AccessToken,
	})

	oauthClient := oauth2.NewClient(ctx, ts)

	githubHTTPClient := &http.Client{
		Transport: oauthClient.Transport,
	}
	return githubHTTPClient, nil
}
