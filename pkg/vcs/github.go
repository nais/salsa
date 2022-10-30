package vcs

import (
	"encoding/json"
	"fmt"
)

type GithubContext struct {
	Action          string          `json:"action"`
	Actor           string          `json:"actor"`
	Event           json.RawMessage `json:"event"`
	EventName       string          `json:"event_name"`
	EventPath       string          `json:"event_path"`
	Job             string          `json:"job"`
	Ref             string          `json:"ref"`
	Repository      string          `json:"repository"`
	RepositoryOwner string          `json:"repository_owner"`
	RunId           string          `json:"run_id"`
	RunNumber       string          `json:"run_number"`
	ServerUrl       string          `json:"server_url"`
	SHA             string          `json:"sha"`
	Token           string          `json:"token,omitempty"`
	Workflow        string          `json:"workflow"`
	Workspace       string          `json:"workspace"`
}

func ParseContext(github []byte) (*GithubContext, error) {
	context := GithubContext{}
	if len(github) == 0 {
		return nil, nil
	}

	if err := json.Unmarshal(github, &context); err != nil {
		if err != nil {
			return nil, fmt.Errorf("unmarshal github context json: %w", err)
		}
	}

	// Ensure we dont misuse token.
	context.Token = ""
	return &context, nil
}

type Actions struct {
	HostedIdSuffix     string
	SelfHostedIdSuffix string
	BuildType          string
}

// BuildId
// The GitHub Actions team has not yet reviewed or approved this design,
// and it is not yet implemented. Details are subject to change!
func BuildId(version string) *Actions {
	return &Actions{
		HostedIdSuffix: fmt.Sprintf("/Attestations/GitHubHostedActions@%s", version),
		// Self-hosted runner: Not yet supported.
		SelfHostedIdSuffix: fmt.Sprintf("/Attestations/SelfHostedActions@%s", version),
		// BuildType Describes what the invocations buildConfig and materials was created
		BuildType: fmt.Sprintf("https://github.com/Attestations/GitHubActionsWorkflow@%s", version),
	}
}
