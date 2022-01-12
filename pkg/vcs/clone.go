package vcs

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
)

func CloneRepo(repoUrl string, path string) error {
	log.Printf("cloning repo %s", repoUrl)
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: repoUrl,
	})

	if err != nil {
		return fmt.Errorf("could not clone repo %s, %v", repoUrl, err)
	}
	return nil
}
