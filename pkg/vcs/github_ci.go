package vcs

import (
	"fmt"
	"os"
)

type GithubCIEnvironment struct {
	GitHubContext    *GitHubContext
	Event            *Event
	RunnerContext    *RunnerContext
	BuildEnvironment BuildEnvironment
}

func integrationEnvironment(context *GitHubContext, runner *RunnerContext, current BuildEnvironment) ContextEnvironment {
	return &GithubCIEnvironment{
		GitHubContext: context,
		Event: &Event{
			Inputs: context.Event,
		},
		RunnerContext:    runner,
		BuildEnvironment: current,
	}
}

func CreateGithubCIEnvironment(githubContext, runnerContext, envsContext *string) (ContextEnvironment, error) {
	// Required when creating CI CiEnvironment
	if len(*githubContext) == 0 || len(*runnerContext) == 0 {
		return nil, nil
	}

	context, err := ParseContext(githubContext)
	if err != nil {
		return nil, fmt.Errorf("parsing context: %w", err)
	}

	runner, err := ParseRunner(runnerContext)
	if err != nil {
		return nil, fmt.Errorf("parsing runner: %w", err)
	}

	current, err := ParseBuild(envsContext)
	if err != nil {
		return nil, fmt.Errorf("parsing envs: %w", err)
	}

	return integrationEnvironment(context, runner, current), nil
}

func (in *GithubCIEnvironment) Context() string {
	return in.GitHubContext.Workflow
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
		return in.RepoUri() + GitHubHostedIdSuffix
	}
	return in.RepoUri() + GitHubHostedIdSuffix
}

func (in *GithubCIEnvironment) UserDefinedParameters() *Event {
	// Only possible user-defined parameters
	// This is unset/null for all other events.
	if in.GitHubContext.EventName != "workflow_dispatch" {
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
				RunId: in.GitHubContext.RunId,
			},
			Runner: Runner{
				Os:   in.RunnerContext.OS,
				Temp: in.RunnerContext.Temp,
			},
		},
	}
}
