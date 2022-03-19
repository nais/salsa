package vcs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

type ContextType string

const (
	ContextTypeGithub ContextType = "GITHUB_ACTIONS"
)

func (in ContextType) String() string {
	return string(in)
}

func (in ContextType) Hosted() bool {
	return os.Getenv(in.String()) == "true"
}

func ResolveBuildContext(context, runner, env *string) (ContextEnvironment, error) {
	if !buildContext(context, runner) {
		return nil, nil
	}

	decodedContext, err := base64.StdEncoding.DecodeString(*context)
	if err != nil {
		return nil, fmt.Errorf("decoding context: %w", err)
	}

	if !isJSON(decodedContext) {
		return nil, fmt.Errorf("decoded build context is not in json format")
	}

	if ContextTypeGithub.Hosted() {
		log.Info("prepare Github CI environment...")
		return CreateGithubCIEnvironment(decodedContext, runner, env)
	}

	return nil, fmt.Errorf("build context is not supported")
}

// Required when creating CI Environment, CLI assumed to be run manually without build context
func buildContext(context, runner *string) bool {
	return (context != nil && len(*context) != 0) && (runner != nil && len(*runner) != 0)
}

func isJSON(str []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(str, &js) == nil
}
