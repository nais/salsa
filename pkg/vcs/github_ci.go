package vcs

import (
	"fmt"
)

const (
	GithubActionsBuildIdVersion = "v1"
)

type GithubCIEnvironment struct {
	BuildContext     *GithubContext
	Event            *Event
	RunnerContext    *RunnerContext
	BuildEnvironment *CurrentBuildEnvironment
	Actions          *Actions
}

func CreateGithubCIEnvironment(githubContext []byte, runnerContext, envsContext *string) (ContextEnvironment, error) {
	context, err := ParseContext(githubContext)
	if err != nil {
		return nil, fmt.Errorf("parsing context: %w", err)
	}

	runner, err := ParseRunner(runnerContext)
	if err != nil {
		return nil, fmt.Errorf("parsing runner: %w", err)
	}

	event, err := ParseEvent(context.Event)
	if err != nil {
		return nil, fmt.Errorf("parsing event: %w", err)
	}

	// Not required to build a CI environment
	current := &CurrentBuildEnvironment{}
	if envsContext == nil || len(*envsContext) == 0 {
		return BuildEnvironment(context, runner, current, event), nil
	}

	current, err = ParseBuild(envsContext)
	if err != nil {
		return nil, fmt.Errorf("parsing envs: %w", err)
	}

	return BuildEnvironment(context, runner, current, event), nil
}

func BuildEnvironment(context *GithubContext, runner *RunnerContext, current *CurrentBuildEnvironment, event *Event) ContextEnvironment {
	return &GithubCIEnvironment{
		BuildContext:     context,
		Event:            event,
		RunnerContext:    runner,
		BuildEnvironment: current,
		Actions:          BuildId(GithubActionsBuildIdVersion),
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

	// should be filtered to fit the information needed
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

func (in *GithubCIEnvironment) GetHeadCommitTime() string {
	if in.Event == nil {
		return ""
	}

	if in.Event.GetHeadCommitId() == in.Sha() {
		return in.Event.GetHeadCommitTimestamp()
	}

	return ""
}

func (in *GithubCIEnvironment) GetEventMetadata() *Event {
	return in.Event
}
