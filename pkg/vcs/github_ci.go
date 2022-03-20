package vcs

import (
	"fmt"
	"github.com/nais/salsa/pkg/vcs/github"
)

const (
	GithubActionsBuildIdVersion = "v1"
)

type GithubCIEnvironment struct {
	BuildContext     *github.Context
	Event            *Event
	RunnerContext    *github.RunnerContext
	BuildEnvironment *github.CurrentBuildEnvironment
	Actions          *github.Actions
}

func CreateGithubCIEnvironment(githubContext []byte, runnerContext, envsContext *string) (ContextEnvironment, error) {
	context, err := github.ParseContext(githubContext)
	if err != nil {
		return nil, fmt.Errorf("parsing context: %w", err)
	}

	runner, err := github.ParseRunner(runnerContext)
	if err != nil {
		return nil, fmt.Errorf("parsing runner: %w", err)
	}

	// Not required to build a CI environment
	current := &github.CurrentBuildEnvironment{}
	if envsContext == nil || len(*envsContext) == 0 {
		return BuildEnvironment(context, runner, current), nil
	}

	current, err = github.ParseBuild(envsContext)
	if err != nil {
		return nil, fmt.Errorf("parsing envs: %w", err)
	}

	return BuildEnvironment(context, runner, current), nil
}

func BuildEnvironment(context *github.Context, runner *github.RunnerContext, current *github.CurrentBuildEnvironment) ContextEnvironment {
	return &GithubCIEnvironment{
		BuildContext: context,
		Event: &Event{
			Inputs: context.Event,
		},
		RunnerContext:    runner,
		BuildEnvironment: current,
		Actions:          github.BuildId(GithubActionsBuildIdVersion),
	}
}

func (in *GithubCIEnvironment) Context() string {
	return in.BuildContext.Workflow
}

func (in *GithubCIEnvironment) BuildType() string {
	return in.Actions.BuildType
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
	if ContextTypeGithub.Hosted() {
		return in.RepoUri() + in.Actions.HostedIdSuffix
	}
	return in.RepoUri() + in.Actions.SelfHostedIdSuffix
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
