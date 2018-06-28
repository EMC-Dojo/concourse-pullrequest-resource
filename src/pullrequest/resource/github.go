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
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var downloadPRScriptPath = "/var/download_pr.sh"

var downloadPRScriptBytes = `#!/bin/sh
git -c http.sslVerify=false clone {{.RepoURL}} {{.DestDir}}/
cd {{.DestDir}}
git -c http.sslVerify=false fetch origin pull/{{.PRNumber}}/head:pr
git checkout pr
`

// Github is
type Github interface {
	ListPRs() ([]*github.PullRequest, error)
	DownloadPR(string, int) error
	UpdatePR(string, string, string) (string, error)
}

// GithubClient is
type GithubClient struct {
	client *github.Client
	owner  string
	repo   string
	token  string
	ctx    context.Context
}

// NewGithubClient is
func NewGithubClient(source Source) (*GithubClient, error) {
	var httpClient = &http.Client{}
	var ctx = context.Background()
	var err error

	if source.Insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	}

	if source.AccessToken != "" {
		httpClient, err = oauthClient(ctx, source)
		if err != nil {
			return nil, fmt.Errorf("constructing oauth2 client: %+v", err)
		}
	}

	var client *github.Client
	if source.APIURL == "" {
		client = github.NewClient(httpClient)
	} else {
		client, err = github.NewEnterpriseClient(source.APIURL, source.APIURL+"/upload", httpClient)
		if err != nil {
			return nil, fmt.Errorf("construting enterprise oauth2 client: %+v", err)
		}
	}

	return &GithubClient{
		client: client,
		owner:  source.Owner,
		repo:   source.Repo,
		token:  source.AccessToken,
	}, nil
}

// ListPRs is
func (gc *GithubClient) ListPRs() ([]*github.PullRequest, error) {
	options := &github.PullRequestListOptions{
		Sort:      "updated",
		Direction: "asc",
	}

	pulls, resp, err := gc.client.PullRequests.List(context.TODO(), gc.owner, gc.repo, options)
	if err != nil {
		return nil, fmt.Errorf("listing pr: %+v", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return pulls, nil
}

type pullFetcher struct {
	RepoURL  string
	DestDir  string
	PRNumber int
}

// DownloadPR is
func (gc *GithubClient) DownloadPR(destDir string, prNumber int) error {
	repo, resp, err := gc.client.Repositories.Get(context.TODO(), gc.owner, gc.repo)
	if err != nil {
		return fmt.Errorf("getting repos: %+v", err)
	}

	if err = resp.Body.Close(); err != nil {
		return fmt.Errorf("closing resp body: %+v", err)
	}

	file, err := os.OpenFile(downloadPRScriptPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("opening download script: %+v", err)
	}

	tmpl := template.Must(template.New("").Parse(downloadPRScriptBytes))
	pf := pullFetcher{
		RepoURL:  buildURLWithToken(repo.GetHTMLURL(), gc.token),
		DestDir:  destDir,
		PRNumber: prNumber,
	}

	if err = tmpl.Execute(file, pf); err != nil {
		return fmt.Errorf("executing template: %+v", err)
	}

	log.Infof("download script path: %s", downloadPRScriptPath)
	log.Infof("repo path: %s", destDir)

	cmd := exec.Command("/bin/sh", downloadPRScriptPath)
	if output, err := cmd.Output(); err != nil {
		return fmt.Errorf("executing download script: %s, %+v", string(output), err)
	}

	log.Infof("listing comments for PR: %d", prNumber)
	pull, resp, err := gc.client.PullRequests.Get(context.TODO(), gc.owner, gc.repo, prNumber)
	if err != nil {
		return fmt.Errorf("listing comments: %+v", err)
	}

	if err = resp.Body.Close(); err != nil {
		return fmt.Errorf("closing resp body: %+v", err)
	}

	var labels string
	for _, label := range pull.Labels {
		labels += " " + label.GetName()
	}

	log.Infof("fetched %d labels: %s", len(pull.Labels), labels)

	err = ioutil.WriteFile(path.Join(destDir, "pr_labels"), []byte(labels), 0644)
	if err != nil {
		return fmt.Errorf("writing to pr_labels: %+v", err)
	}

	err = ioutil.WriteFile(path.Join(destDir, "pr_number"), []byte(strconv.Itoa(prNumber)), 0644)
	if err != nil {
		return fmt.Errorf("writing to pr_number: %+v", err)
	}

	err = ioutil.WriteFile(path.Join(destDir, "pr_comment"), []byte(*pull.Body), 0644)
	if err != nil {
		return fmt.Errorf("writing to pr_comment: %+v", err)
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
func (gc *GithubClient) UpdatePR(sourceDir, status, repoPath string) (string, error) {
	switch status {
	case
		"error",
		"failure",
		"pending",
		"success":
		break
	default:
		return "", fmt.Errorf("%s is not a valid status", status)
	}
	repoStatus := &github.RepoStatus{
		State: &status,
	}

	repoDir := path.Join(sourceDir, repoPath)
	log.Infof("repo dir: %s", repoDir)

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path.Join(repoDir)
	output, err := cmd.Output()
	log.Infof("git rev-parse output: %+v", string(output))
	if err != nil {
		return "", fmt.Errorf("getting pr commit hash: %+v", err)
	}
	ref := strings.TrimRight(string(output), "\n\r")

	returnedRepoStatus, resp, err := gc.client.Repositories.CreateStatus(context.TODO(), gc.owner, gc.repo, string(ref), repoStatus)
	if err != nil {
		return "", fmt.Errorf("creating status: %+v", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("closing resp body: %+v", err)
	}

	if returnedRepoStatus.GetState() != status {
		return "", errors.New("updating commit status")
	}
	return ref, nil
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
