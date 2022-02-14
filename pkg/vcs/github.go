package vcs

import (
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

func (in *GitHubContext) Parse(github *string, env *Environment) error {
	if len(*github) == 0 {
		return nil
	}

	if err := json.Unmarshal([]byte(*github), &env.GitHubContext); err != nil {
		if err != nil {
			return fmt.Errorf("unmarshal github context json: %w", err)
		}
	}
	if env.GitHubContext.Event != nil {
		if err := json.Unmarshal(env.GitHubContext.Event, &env.Event); err != nil {
			if err != nil {
				return fmt.Errorf("unmarshal github event json: %w", err)
			}
		}
	}

	// Ensure we dont misuse token.
	env.GitHubContext.Token = ""
	return nil
}

type RunnerContext struct {
	Name      string `json:"name"`
	Arch      string `json:"arch"`
	OS        string `json:"os"`
	Temp      string `json:"temp"`
	ToolCache string `json:"tool_cache"`
}

func (in *RunnerContext) Parse(runner *string, env *Environment) error {
	if len(*runner) == 0 {
		return nil
	}

	if err := json.Unmarshal([]byte(*runner), &env.RunnerContext); err != nil {
		return fmt.Errorf("unmarshal runner context json: %w", err)
	}

	return nil
}

type CurrentEnvironment struct {
	Envs map[string]string
}

func (in *CurrentEnvironment) Parse(envs *string, env *Environment) error {
	if len(*envs) == 0 {
		return nil
	}

	env.CurrentEnvironment = &CurrentEnvironment{
		make(map[string]string),
	}

	if err := json.Unmarshal([]byte(*envs), &env.CurrentEnvironment.Envs); err != nil {
		return fmt.Errorf("unmarshal environmental context json: %w", err)
	}

	return nil
}
