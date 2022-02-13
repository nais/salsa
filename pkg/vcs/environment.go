package vcs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Environment struct {
	GitHubContext `json:"github"`
	*Event        `json:"event,omitempty"`
	RunnerContext `json:"context"`
}

func CreateCIEnvironment(githubContext, runnerContext *string) (*Environment, error) {
	if githubContext == nil || runnerContext == nil {
		return nil, nil
	}

	env := Environment{}

	if len(*githubContext) > 0 {
		if err := json.Unmarshal([]byte(*githubContext), &env.GitHubContext); err != nil {
			if err != nil {
				return nil, fmt.Errorf("unmarshal github context json: %w", err)
			}
		}
		if env.GitHubContext.Event != nil {
			if err := json.Unmarshal(env.GitHubContext.Event, &env.Event); err != nil {
				if err != nil {
					return nil, fmt.Errorf("unmarshal github event json: %w", err)
				}
			}
		}
	}

	// Ensure we dont misuse token.
	env.Token = ""

	if len(*runnerContext) > 0 {
		if err := json.Unmarshal([]byte(*runnerContext), &env.RunnerContext); err != nil {
			return nil, fmt.Errorf("unmarshal runner context json: %w", err)
		}
	}

	return &env, nil
}

func (in *Environment) RepoUri() string {
	return fmt.Sprintf("%s/%s", in.GitHubContext.ServerUrl, in.GitHubContext.Repository)
}

func (in *Environment) BuildInvocationId() string {
	return fmt.Sprintf("%s/actions/runs/%s", in.RepoUri(), in.GitHubContext.RunId)
}

func (in *Environment) GithubSha() string {
	return in.GitHubContext.SHA
}

func (in *Environment) BuilderId() string {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return in.RepoUri() + GitHubHostedIdSuffix
	}
	return in.RepoUri() + GitHubHostedIdSuffix
}

func (in *Environment) EventInputs() json.RawMessage {
	if in.EventName != "workflow_dispatch" {
		return nil
	}
	return in.Event.Inputs
}

func (in *Environment) FilteredEnvironment() *Environment {
	// Should also contain environment variables.
	// These are always set because it is not possible
	// to know whether they were referenced or not.
	return &Environment{
		GitHubContext: GitHubContext{
			RunId: in.RunId,
		},
		Event:         nil,
		RunnerContext: in.RunnerContext,
	}
}
