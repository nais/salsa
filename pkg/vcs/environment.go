package vcs

import (
	"encoding/json"
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
