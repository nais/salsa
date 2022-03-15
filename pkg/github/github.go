package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
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

func ParseGithub(github *string, env *Environment) error {
	if len(*github) == 0 {
		return nil
	}

	decodedGithubBytes, err := base64.StdEncoding.DecodeString(*github)
	if err != nil {
		return fmt.Errorf("decoding github context: %w", err)
	}

	if err := json.Unmarshal(decodedGithubBytes, &env.GitHubContext); err != nil {
		if err != nil {
			return fmt.Errorf("unmarshal github context json: %w", err)
		}
	}

	// TODO check if this data is useful, Unmarshal to struct
	if env.GitHubContext.Event != nil {
		env.Event = &Event{
			Inputs: env.GitHubContext.Event,
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

func ParseRunner(runner *string, env *Environment) error {
	if len(*runner) == 0 {
		return nil
	}

	decodedRunnerBytes, err := base64.StdEncoding.DecodeString(*runner)
	if err != nil {
		return fmt.Errorf("decoding runner context: %w", err)
	}

	if err := json.Unmarshal(decodedRunnerBytes, &env.RunnerContext); err != nil {
		return fmt.Errorf("unmarshal runner context json: %w", err)
	}

	return nil
}

type CurrentEnvironment struct {
	Envs map[string]string
}

func ParseEnv(envs *string, env *Environment) error {
	if len(*envs) == 0 {
		return nil
	}

	decodedEnvsBytes, err := base64.StdEncoding.DecodeString(*envs)
	if err != nil {
		return fmt.Errorf("decoding envs context: %w", err)
	}

	env.CurrentEnvironment = &CurrentEnvironment{
		make(map[string]string),
	}

	if err := json.Unmarshal(decodedEnvsBytes, &env.CurrentEnvironment.Envs); err != nil {
		return fmt.Errorf("unmarshal environmental context json: %w", err)
	}

	return nil
}

func (in *CurrentEnvironment) filterEnvs() map[string]string {
	if len(in.Envs) < 1 {
		return map[string]string{}
	}
	for key, val := range in.Envs {
		in.filterEnvsWithPrefix(key, "INPUT_", "GITHUB_", "RUNNER_", "ACTIONS_")
		in.filterEnv(key, "TOKEN")
		in.filterSingleLineEnv(key)
		in.filterEmptyValue(key, val)
	}
	return in.Envs
}

func (in *CurrentEnvironment) filterSingleLineEnv(key string) {
	if !strings.Contains(key, "_") {
		delete(in.Envs, key)
	}
}

func (in *CurrentEnvironment) filterEmptyValue(key, val string) {
	if val == "" {
		delete(in.Envs, key)
	}
}

func (in *CurrentEnvironment) filterEnv(key string, contains ...string) {
	for _, contain := range contains {
		if strings.Contains(key, contain) {
			delete(in.Envs, key)
		}
	}
}

func (in *CurrentEnvironment) filterEnvsWithPrefix(key string, prefixes ...string) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(key, prefix) {
			delete(in.Envs, key)
		}
	}
}
