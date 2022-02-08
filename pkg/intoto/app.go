package intoto

import (
	"os"
	"time"

	v02 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/scan"
	"github.com/nais/salsa/pkg/vcs"
)

type App struct {
	Name              string
	BuilderId         string
	BuildType         string
	Dependencies      *scan.ArtifactDependencies
	BuildStartedOn    time.Time
	BuildInvocationId string
	BuilderRepoDigest v02.ProvenanceMaterial
}

func CreateApp(name string, deps *scan.ArtifactDependencies) App {
	return App{
		Name:           name,
		BuildType:      "todoType",
		Dependencies:   deps,
		BuildStartedOn: time.Now().UTC(),
	}
}

func (a App) With(context *vcs.AnyContext) App {
	if context == nil {
		// Required
		a.BuilderId = "default"
		return a
	}

	repoURI := "https://github.com/" + context.GitHubContext.Repository
	a.BuildInvocationId = repoURI + "/actions/runs/" + context.GitHubContext.RunId

	a.BuilderRepoDigest = v02.ProvenanceMaterial{
		URI:    "git+" + repoURI,
		Digest: v02.DigestSet{"sha1": context.GitHubContext.SHA},
	}

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		a.BuilderId = repoURI + vcs.GitHubHostedIdSuffix
	} else {
		a.BuilderId = repoURI + vcs.GitHubHostedIdSuffix
	}

	return a
}
