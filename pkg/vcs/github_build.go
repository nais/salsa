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

func ParseBuild(envs *string) (*CurrentBuildEnvironment, error) {
	current := make(map[string]string)

	decodedEnvsBytes, err := base64.StdEncoding.DecodeString(*envs)
	if err != nil {
		return nil, fmt.Errorf("decoding envs context: %w", err)
	}

	if err := json.Unmarshal(decodedEnvsBytes, &current); err != nil {
		return nil, fmt.Errorf("unmarshal environmental context json: %w", err)
	}

	return &CurrentBuildEnvironment{
		Envs: current,
	}, nil
}

func (in *CurrentBuildEnvironment) FilterEnvs() map[string]string {
	if len(in.Envs) < 1 {
		return map[string]string{}
	}

	for key := range in.Envs {
		in.filterEnvsWithPrefix(key, "INPUT_", "GITHUB_", "RUNNER_", "ACTIONS_")
		in.filterEnv(key, "TOKEN")
	}

	in.removeDuplicateValues()
	return in.Envs
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

func (in *CurrentBuildEnvironment) removeDuplicateValues() {
	var current = make(map[string]struct{})
	for key, v := range in.Envs {
		if _, has := current[v]; has {
			delete(in.Envs, key)
		}
		current[v] = struct{}{}
	}
}
