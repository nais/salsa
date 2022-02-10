package vcs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Environment struct {
	GitHubContext `json:"github"`
	AnyEvent      `json:"event"`
}

func CreateCIEnvironment(inputContext *string) (*Environment, error) {
	env := Environment{}
	if len(*inputContext) > 0 {
		if err := json.Unmarshal([]byte(*inputContext), &env.GitHubContext); err != nil {
			if err != nil {
				return nil, err
			}
		}
		if err := json.Unmarshal(env.GitHubContext.Event, &env.AnyEvent); err != nil {
			if err != nil {
				return nil, err
			}
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
