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
	BuildContext     *github.Context
	Event            *Event
	RunnerContext    *github.RunnerContext
	BuildEnvironment *github.CurrentBuildEnvironment
	StaticBuild      *github.StaticBuild
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
		BuildContext: context,
		Event: &Event{
			Inputs: context.Event,
		},
		RunnerContext:    runner,
		BuildEnvironment: current,
		StaticBuild:      github.Identification(IdentificationVersion),
	}
}

func (in *GithubCIEnvironment) Context() string {
	return in.BuildContext.Workflow
}

func (in *GithubCIEnvironment) BuildType() string {
	return in.StaticBuild.BuildType
}

func (in *GithubCIEnvironment) RepoUri() string {
	return fmt.Sprintf("%s/%s", in.BuildContext.ServerUrl, in.BuildContext.Repository)
}

func (in *GithubCIEnvironment) BuildInvocationId() string {
	return fmt.Sprintf("%s/actions/runs/%s", in.RepoUri(), in.BuildContext.RunId)
}

func (in *GithubCIEnvironment) Sha() string {
	return in.BuildContext.SHA
}

func (in *GithubCIEnvironment) BuilderId() string {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return in.RepoUri() + in.StaticBuild.HostedIdSuffix
	}
	return in.RepoUri() + in.StaticBuild.SelfHostedIdSuffix
}

func (in *GithubCIEnvironment) UserDefinedParameters() *Event {
	// Only possible user-defined parameters
	// This is unset/null for all other events.
	if in.BuildContext.EventName != "workflow_dispatch" {
		return nil
	}

	return in.Event
}

func (in *GithubCIEnvironment) CurrentFilteredEnvironment() map[string]string {
	if in.BuildEnvironment == nil {
		return map[string]string{}
	}

	return in.BuildEnvironment.FilterEnvs()
}

func (in *GithubCIEnvironment) NonReproducibleMetadata() *Metadata {
	// Other variables that are required to reproduce the build and that cannot be
	// recomputed using existing information.
	//(Documentation would explain how to recompute the rest of the fields.)
	return &Metadata{
		Arch: in.RunnerContext.Arch,
		Env:  in.CurrentFilteredEnvironment(),
		Context: Context{
			Github: Github{
				RunId: in.BuildContext.RunId,
			},
			Runner: Runner{
				Os:   in.RunnerContext.OS,
				Temp: in.RunnerContext.Temp,
			},
		},
	}
}
