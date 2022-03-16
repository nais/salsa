package config

import (
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/vcs"
	"github.com/spf13/cobra"
)

type ScanConfiguration struct {
	WorkDir            string
	RepoName           string
	Dependencies       *build.ArtifactDependencies
	ContextEnvironment vcs.ContextEnvironment
	Cmd                *cobra.Command
}
