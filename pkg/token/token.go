package token

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

const (
	DefaultGithubTokenEnvKey = "GITHUB_TOKEN"
)

func Export(tokenFlag string) error {
	if tokenFlag == "" {
		return nil
	}

	token := strings.Split(strings.TrimSpace(tokenFlag), "=")

	if funky(token) {
		return errors.New("wrong format for --token flag; should be a single value or 'key=value'")
	}

	if len(token) == 1 {
		err := os.Setenv(DefaultGithubTokenEnvKey, token[0])
		if err != nil {
			return fmt.Errorf("exporting varibales%w", err)
		}
		log.Infof("exporting %s", DefaultGithubTokenEnvKey)
		return nil
	}

	key := strings.ToUpper(strings.ReplaceAll(token[0], ".", "_"))
	value := strings.TrimSpace(token[1])
	err := os.Setenv(key, value)

	if err != nil {
		return fmt.Errorf("exporting varibales%w", err)
	}

	log.Infof("exporting %s", key)
	return nil
}

func funky(token []string) bool {
	return len(token) == 1 && token[0] == "" || len(token) == 2 && token[1] == "" || len(token) > 2
}
