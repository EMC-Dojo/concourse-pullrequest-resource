package resource

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Github is
type Github interface {
	ListPRs() ([]*github.PullRequest, error)
	DownloadPR(string, int) error
	UpdatePR(string, string) error
}

// GithubClient is
type GithubClient struct {
	client *github.Client
	owner  string
	repo   string
	token  string
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

	log.Infof("AccessToken: %s", source.AccessToken)

	if source.AccessToken != "" {
		var err error
		httpClient, err = oauthClient(ctx, source)
		if err != nil {
			return nil, fmt.Errorf("constructing oauth2 client: %+v", err)
		}
	}

	client := github.NewClient(httpClient)

	return &GithubClient{
		client,
		source.Owner,
		source.Repo,
		source.AccessToken,
	}, nil
}

// ListPRs is
func (gc *GithubClient) ListPRs() ([]*github.PullRequest, error) {
	pulls, resp, err := gc.client.PullRequests.List(context.TODO(), gc.owner, gc.repo, nil)
	if err != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Infof("resp: %+v", string(bodyBytes))
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return pulls, nil
}

type pullFetcher struct {
	RepoURL  string
	RepoDir  string
	PRNumber int
}

// DownloadPR is
func (gc *GithubClient) DownloadPR(sourceDir string, prNumber int) error {
	repo, resp, err := gc.client.Repositories.Get(context.TODO(), gc.owner, gc.repo)
	if err != nil {
		return err
	}

	if err = resp.Body.Close(); err != nil {
		return fmt.Errorf("closing resp body: %+v", err)
	}

	file, err := os.OpenFile(path.Join(sourceDir, "fetch.sh"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("opening file: %+v", err)
	}

	tmpl := template.Must(template.New("").Parse(`#!/bin/bash
git clone {{.RepoURL}} {{.RepoDir}}
pushd {{.RepoDir}}
  git fetch origin pull/{{.PRNumber}}/head:pr
  git checkout pr
popd
rm -- "$0"
`))

	pf := pullFetcher{
		RepoURL:  buildURLWithToken(repo.GetHTMLURL(), gc.token),
		RepoDir:  path.Join(sourceDir, repo.GetName()),
		PRNumber: prNumber,
	}

	if err = tmpl.Execute(file, pf); err != nil {
		return fmt.Errorf("executing template: %+v", err)
	}

	cmd := exec.Command("./fetch.sh")
	cmd.Dir = sourceDir
	if _, err = cmd.Output(); err != nil {
		return err
	}

	return nil
}

// GetPR is
func (gc *GithubClient) GetPR(number int) (*github.PullRequest, error) {
	pull, resp, err := gc.client.PullRequests.Get(context.TODO(), gc.owner, gc.repo, number)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return pull, nil
}

// UpdatePR is
func (gc *GithubClient) UpdatePR(sourceDir, status string) error {
	switch status {
	case
		"error",
		"failure",
		"pending",
		"success":
		break
	default:
		return fmt.Errorf("%s is not a valid status", status)
	}

	repoStatus := &github.RepoStatus{
		State: &status,
	}

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path.Join(sourceDir, gc.repo)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	ref := strings.TrimRight(string(output), "\n\r")
	fmt.Println("ref", ref, "end")

	returnedRepoStatus, resp, err := gc.client.Repositories.CreateStatus(context.TODO(), gc.owner, gc.repo, string(ref), repoStatus)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	if returnedRepoStatus.GetState() != status {
		return errors.New("updating commit status")
	}

	return nil
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

func buildURLWithToken(url, token string) string {
	return fmt.Sprintf("https://%s@%s", token, url[8:])
}
