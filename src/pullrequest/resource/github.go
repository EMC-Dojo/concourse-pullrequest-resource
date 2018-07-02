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
	"time"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var githubCheckContext = "concourse/ci"

var downloadPRScriptPath = "/var/download_pr.sh"
var downloadPRScriptBytes = `#!/bin/sh
git -c http.sslVerify=false clone {{.RepoURL}} {{.DestDir}}/
cd {{.DestDir}}
git -c http.sslVerify=false fetch origin pull/{{.PRNumber}}/head:pr
git checkout pr
`

// Pull is
type Pull struct {
	Number          int
	ID              string
	LatestCommitSHA string
	URL             string
	Body            string
	Labels          string
	Title           string
}

// Github is
type Github interface {
	ListPRs() ([]*Pull, error)
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
func (gc *GithubClient) ListPRs() ([]*Pull, error) {
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

	var convertedPulls = []*Pull{}
	for _, pull := range pulls {
		convertedPulls = append(convertedPulls, convertPR(pull))
	}
	return convertedPulls, nil
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

	log.Infof("repo path: %s", destDir)

	cmd := exec.Command("/bin/sh", downloadPRScriptPath)
	if output, err := cmd.Output(); err != nil {
		return fmt.Errorf("executing download script: %s, %+v", string(output), err)
	}

	pull, err := gc.GetPR(prNumber)
	if err != nil {
		return fmt.Errorf("getting pr %d: %+v", prNumber, err)
	}

	return writePullToFile(destDir, pull)
}

// GetPR is
func (gc *GithubClient) GetPR(number int) (*Pull, error) {
	pull, resp, err := gc.client.PullRequests.Get(context.TODO(), gc.owner, gc.repo, number)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return convertPR(pull), nil
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
		State:   &status,
		Context: &githubCheckContext,
		Creator: &github.User{},
	}

	commitHashBytes, err := ioutil.ReadFile(path.Join(sourceDir, repoPath, "pr_last_commit_hash"))
	if err != nil {
		return "", fmt.Errorf("reading pr_last_commit_hash: %+v", err)
	}

	returnedRepoStatus, resp, err := gc.client.Repositories.CreateStatus(context.TODO(), gc.owner, gc.repo, string(commitHashBytes), repoStatus)
	if err != nil {
		return "", fmt.Errorf("creating status: %+v", err)
	}
	if err = resp.Body.Close(); err != nil {
		return "", fmt.Errorf("closing resp body: %+v", err)
	}
	if returnedRepoStatus.GetState() != status {
		return "", errors.New("updating commit status")
	}

	idBytes, err := ioutil.ReadFile(path.Join(sourceDir, repoPath, "pr_id"))
	if err != nil {
		return "", fmt.Errorf("reading pr_id: %+v", err)
	}
	return string(idBytes), nil
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

func writePullToFile(destDir string, pull *Pull) error {
	if err := writeToFile(destDir, "pr_labels", pull.Labels); err != nil {
		return err
	}

	if err := writeToFile(destDir, "pr_number", strconv.Itoa(pull.Number)); err != nil {
		return err
	}

	if err := writeToFile(destDir, "pr_body", pull.Body); err != nil {
		return err
	}

	if err := writeToFile(destDir, "pr_title", pull.Title); err != nil {
		return err
	}

	if err := writeToFile(destDir, "pr_last_commit_hash", pull.LatestCommitSHA); err != nil {
		return err
	}

	if err := writeToFile(destDir, "pr_id", pull.ID); err != nil {
		return err
	}

	return nil
}

func writeToFile(destDir, fileName, content string) error {
	err := ioutil.WriteFile(path.Join(destDir, fileName), []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("writing to %s: %+v", fileName, err)
	}
	return nil
}

func convertPR(pr *github.PullRequest) *Pull {
	var labels string
	for _, label := range pr.Labels {
		labels += label.GetName() + "\n"
	}

	return &Pull{
		Number:          pr.GetNumber(),
		LatestCommitSHA: pr.GetHead().GetSHA(),
		ID:              fmt.Sprintf("%s-%s", pr.GetHead().GetSHA()[0:7], pr.GetUpdatedAt().Format(time.RFC3339)),
		URL:             pr.GetURL(),
		Title:           pr.GetTitle(),
		Body:            pr.GetBody(),
		Labels:          labels,
	}
}
