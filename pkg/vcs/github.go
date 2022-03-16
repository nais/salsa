package vcs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	GitHubHostedIdSuffix = "/Attestations/GitHubHostedActions@v1"
	SelfHostedIdSuffix   = "/Attestations/SelfHostedActions@v1"
	// BuildType Describes what the invocations buildConfig and materials was created
	BuildType = "https://github.com/Attestations/GitHubActionsWorkflow@v1"
	// AdHocBuildType no entry point, and the commands were run in an ad-hoc fashion
	AdHocBuildType = "https://github.com/nais/salsa/ManuallyRunCommands@v1"
	DefaultBuildId = "https://github.com/nais/salsa"
)

type GitHubContext struct {
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

type Event struct {
	Inputs json.RawMessage `json:"inputs"`
}

func ParseContext(github *string) (*GitHubContext, error) {
	context := GitHubContext{}
	if len(*github) == 0 {
		return nil, nil
	}

	decodedGithubBytes, err := base64.StdEncoding.DecodeString(*github)
	if err != nil {
		return nil, fmt.Errorf("decoding github context: %w", err)
	}

	if err := json.Unmarshal(decodedGithubBytes, &context); err != nil {
		if err != nil {
			return nil, fmt.Errorf("unmarshal github context json: %w", err)
		}
	}

	// Ensure we dont misuse token.
	context.Token = ""
	return &context, nil
}
