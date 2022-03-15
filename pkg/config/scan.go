package config

import (
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/github"
	"github.com/spf13/cobra"
)

type ScanConfiguration struct {
	WorkDir       string
	RepoName      string
	Dependencies  *build.ArtifactDependencies
	CiEnvironment *github.Environment
	Cmd           *cobra.Command
}
