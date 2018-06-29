package fake

import "pullrequest/resource"

// FGithub is
type FGithub struct {
	Pulls []*resource.Pull

	UpdatePRResult string
	UpdatePRError  error
}

// ListPRs is
func (fg *FGithub) ListPRs() ([]*resource.Pull, error) {
	return fg.Pulls, nil
}

// DownloadPR is
func (fg *FGithub) DownloadPR(destDir string, prNumber int) error {
	return nil
}

// UpdatePR is
func (fg *FGithub) UpdatePR(sourceDir, status, path string) (string, error) {
	return fg.UpdatePRResult, fg.UpdatePRError
}
