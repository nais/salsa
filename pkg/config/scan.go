package config

import (
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/vcs"
	"github.com/spf13/cobra"
	"time"
)

type ScanConfiguration struct {
	BuildStartedOn     time.Time
	Cmd                *cobra.Command
	ContextEnvironment vcs.ContextEnvironment
	Dependencies       *build.ArtifactDependencies
	WorkDir            string
	RepoName           string
}
