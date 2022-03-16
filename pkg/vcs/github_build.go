package vcs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type CurrentBuildEnvironment struct {
	Envs map[string]string
}

func build(envs map[string]string) BuildEnvironment {
	return &CurrentBuildEnvironment{
		Envs: envs,
	}
}

func ParseBuild(envs *string) (BuildEnvironment, error) {
	current := make(map[string]string)
	if len(*envs) == 0 {
		return nil, nil
	}

	decodedEnvsBytes, err := base64.StdEncoding.DecodeString(*envs)
	if err != nil {
		return nil, fmt.Errorf("decoding envs context: %w", err)
	}

	if err := json.Unmarshal(decodedEnvsBytes, &current); err != nil {
		return nil, fmt.Errorf("unmarshal environmental context json: %w", err)
	}

	return build(current), nil
}

func (in *CurrentBuildEnvironment) FilterEnvs() map[string]string {
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

func (in *CurrentBuildEnvironment) filterSingleLineEnv(key string) {
	if !strings.Contains(key, "_") {
		delete(in.Envs, key)
	}
}

func (in *CurrentBuildEnvironment) filterEmptyValue(key, val string) {
	if val == "" {
		delete(in.Envs, key)
	}
}

func (in *CurrentBuildEnvironment) filterEnv(key string, contains ...string) {
	for _, contain := range contains {
		if strings.Contains(key, contain) {
			delete(in.Envs, key)
		}
	}
}

func (in *CurrentBuildEnvironment) filterEnvsWithPrefix(key string, prefixes ...string) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(key, prefix) {
			delete(in.Envs, key)
		}
	}
}

func (in *CurrentBuildEnvironment) GetEnvs() map[string]string {
	return in.Envs
}
