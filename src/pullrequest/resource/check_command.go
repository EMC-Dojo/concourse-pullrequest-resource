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

	for _, pull := range pulls {
		versions = append(versions, Version{Ref: *pull.Head.SHA})
	}

	return versions, nil
}
