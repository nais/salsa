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
	ContextTypeKeyGithub ContextType = "GITHUB_ACTIONS"
)

func (in ContextType) String() string {
	return string(in)
}

func ResolveBuildContext(context, ctx, env *string) (ContextEnvironment, error) {
	if len(*context) == 0 {
		return nil, nil
	}

	decodedContext, err := base64.StdEncoding.DecodeString(*context)
	if err != nil {
		return nil, fmt.Errorf("decoding context: %w", err)
	}

	if !isJSON(decodedContext) {
		return nil, nil
	}

	if isGithub() {
		log.Info("prepare Github CI environment...")
		return CreateGithubCIEnvironment(decodedContext, ctx, env)
	}

	return nil, nil
}

func isGithub() bool {
	return os.Getenv(ContextTypeKeyGithub.String()) == "true"
}

func isJSON(str []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(str, &js) == nil
}
