package vcs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Environment struct {
	GitHubContext `json:"github"`
	Event         `json:"event"`
	RunnerContext `json:"runner_context"`
}

func CreateCIEnvironment(githubContext, runnerContext *string) (*Environment, error) {
	env := Environment{}

	if githubContext != nil && len(*githubContext) > 0 {
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

	if runnerContext != nil && len(*runnerContext) > 0 {
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

func (in *Environment) EventBytes() []byte {
	output, _ := in.Event.Inputs.MarshalJSON()
	return output
}
