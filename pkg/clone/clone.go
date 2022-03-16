package clone

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/go-git/go-git/v5"
)

const (
	GithubUrl = "https://github.com"
)

func Repo(owner, repo, path, token string) error {
	auth := transport.AuthMethod(nil)
	if token != "" {
		auth = &http.BasicAuth{
			Username: "github",
			Password: token,
		}
	}
	repoUrl := fmt.Sprintf("%s/%s/%s", GithubUrl, owner, repo)
	log.Printf("cloning repo %s", repoUrl)
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		Auth:     auth,
		URL:      repoUrl,
		Progress: os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("could not clone repo %s, %v", repoUrl, err)
	}
	return nil
}
