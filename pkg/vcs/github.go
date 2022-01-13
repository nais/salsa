package vcs

import (
	"encoding/json"
)

const (
	GitHubHostedIdSuffix = "/Attestations/GitHubHostedActions@v1"
	SelfHostedIdSuffix   = "/Attestations/SelfHostedActions@v1"
)

type GitHubContext struct {
	Action          string          `json:"action"`
	Actor           string          `json:"actor"`
	BaseRef         string          `json:"base_ref"`
	Event           json.RawMessage `json:"event"`
	EventName       string          `json:"event_name"`
	EventPath       string          `json:"event_path"`
	Job             string          `json:"job"`
	Ref             string          `json:"ref"`
	Repository      string          `json:"repository"`
	RepositoryOwner string          `json:"repository_owner"`
	RunId           string          `json:"run_id"`
	RunNumber       string          `json:"run_number"`
	SHA             string          `json:"sha"`
	Token           string          `json:"token,omitempty"`
	Workflow        string          `json:"workflow"`
	Workspace       string          `json:"workspace"`
}

type AnyEvent struct {
	Inputs json.RawMessage `json:"inputs"`
}

type AnyContext struct {
	GitHubContext `json:"github"`
	AnyEvent      `json:"event"`
}

func CreateCIContext(inputContext *string) (*AnyContext, error) {
	context := AnyContext{}
	if len(*inputContext) > 0 {
		if err := json.Unmarshal([]byte(*inputContext), &context.GitHubContext); err != nil {
			if err != nil {
				return nil, err
			}
		}
		if err := json.Unmarshal(context.GitHubContext.Event, &context.AnyEvent); err != nil {
			if err != nil {
				return nil, err
			}
		}
	}
	return &context, nil
}
