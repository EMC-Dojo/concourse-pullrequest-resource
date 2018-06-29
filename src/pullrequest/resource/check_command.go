package resource

// CheckCommand is
type CheckCommand struct {
	github Github
}

// NewCheckCommand is
func NewCheckCommand(g Github) *CheckCommand {
	return &CheckCommand{g}
}

// Run is
func (cc *CheckCommand) Run(request CheckRequest) ([]Version, error) {
	versions := []Version{}

	pulls, err := cc.github.ListPRs()
	if err != nil {
		return versions, err
	}

	if len(pulls) == 0 {
		return versions, nil
	}

	for i := len(pulls) - 1; i >= 0; i-- {
		versions = append([]Version{Version{Ref: pulls[i].SHA}}, versions...)
		if request.Version.Ref == pulls[i].SHA {
			break
		}
	}
	return versions, nil
}
