package vcs

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func CloneRepo(repoUrl string, path string) error {
	fmt.Printf("cloning repo %s", repoUrl)
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: repoUrl,
	})

	if err != nil {
		return fmt.Errorf("could not clone repo %s, got error %v", repoUrl, err)
	}
	return nil
}
