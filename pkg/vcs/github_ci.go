package vcs

import (
	"fmt"
	"github.com/nais/salsa/pkg/vcs/github"
	"os"
)

const (
	IdentificationVersion = "v1"
)

type GithubCIEnvironment struct {
	GitHubContext          *github.Context
	GithubEvent            *Event
	GithubRunnerContext    *github.RunnerContext
	GithubBuildEnvironment *github.CurrentBuildEnvironment
	GithubStaticBuild      *github.StaticBuild
}

func CreateGithubCIEnvironment(githubContext []byte, runnerContext, envsContext *string) (ContextEnvironment, error) {
	// Required when creating CI CiEnvironment
	if len(githubContext) == 0 || len(*runnerContext) == 0 {
		return nil, nil
	}

	context, err := github.ParseContext(githubContext)
	if err != nil {
		return nil, fmt.Errorf("parsing context: %w", err)
	}

	runner, err := github.ParseRunner(runnerContext)
	if err != nil {
		return nil, fmt.Errorf("parsing runner: %w", err)
	}

	current, err := github.ParseBuild(envsContext)
	if err != nil {
		return nil, fmt.Errorf("parsing envs: %w", err)
	}

	return IntegrationEnvironment(context, runner, current), nil
}

func IntegrationEnvironment(context *github.Context, runner *github.RunnerContext, current *github.CurrentBuildEnvironment) ContextEnvironment {
	return &GithubCIEnvironment{
		GitHubContext: context,
		GithubEvent: &Event{
			Inputs: context.Event,
		},
		GithubRunnerContext:    runner,
		GithubBuildEnvironment: current,
		GithubStaticBuild:      github.Identification(IdentificationVersion),
	}
}

func (in *GithubCIEnvironment) Context() string {
	return in.GitHubContext.Workflow
}

func (in *GithubCIEnvironment) BuildType() string {
	return in.GithubStaticBuild.BuildType
}

func (in *GithubCIEnvironment) RepoUri() string {
	return fmt.Sprintf("%s/%s", in.GitHubContext.ServerUrl, in.GitHubContext.Repository)
}

func (in *GithubCIEnvironment) BuildInvocationId() string {
	return fmt.Sprintf("%s/actions/runs/%s", in.RepoUri(), in.GitHubContext.RunId)
}

func (in *GithubCIEnvironment) Sha() string {
	return in.GitHubContext.SHA
}

func (in *GithubCIEnvironment) BuilderId() string {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return in.RepoUri() + in.GithubStaticBuild.HostedIdSuffix
	}
	return in.RepoUri() + in.GithubStaticBuild.SelfHostedIdSuffix
}

func (in *GithubCIEnvironment) UserDefinedParameters() *Event {
	// Only possible user-defined parameters
	// This is unset/null for all other events.
	if in.GitHubContext.EventName != "workflow_dispatch" {
		return nil
	}

	return in.GithubEvent
}

func (in *GithubCIEnvironment) CurrentFilteredEnvironment() map[string]string {
	if in.GithubBuildEnvironment == nil {
		return map[string]string{}
	}

	return in.GithubBuildEnvironment.FilterEnvs()
}

func (in *GithubCIEnvironment) NonReproducibleMetadata() *Metadata {
	// Other variables that are required to reproduce the build and that cannot be
	// recomputed using existing information.
	//(Documentation would explain how to recompute the rest of the fields.)
	return &Metadata{
		Arch: in.GithubRunnerContext.Arch,
		Env:  in.CurrentFilteredEnvironment(),
		Context: Context{
			Github: Github{
				RunId: in.GitHubContext.RunId,
			},
			Runner: Runner{
				Os:   in.GithubRunnerContext.OS,
				Temp: in.GithubRunnerContext.Temp,
			},
		},
	}
}
