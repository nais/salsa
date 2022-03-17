package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type RunnerContext struct {
	Name      string `json:"name"`
	Arch      string `json:"arch"`
	OS        string `json:"os"`
	Temp      string `json:"temp"`
	ToolCache string `json:"tool_cache"`
}

func ParseRunner(runner *string) (*RunnerContext, error) {
	context := RunnerContext{}

	if len(*runner) == 0 {
		return nil, nil
	}

	decodedRunnerBytes, err := base64.StdEncoding.DecodeString(*runner)
	if err != nil {
		return nil, fmt.Errorf("decoding runner context: %w", err)
	}

	if err := json.Unmarshal(decodedRunnerBytes, &context); err != nil {
		return nil, fmt.Errorf("unmarshal runner context json: %w", err)
	}

	return &context, nil
}
