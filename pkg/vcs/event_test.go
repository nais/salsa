package vcs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEventHeadCommit(t *testing.T) {
	for _, test := range []struct {
		name             string
		workFlowMeatData []byte
		WantTime         string
	}{
		{
			name:             "GitHub Event workflow_run with head_commit",
			workFlowMeatData: commitMetadata(t, "testdata/workflowrun-head-commit.json"),
			WantTime:         "2022-10-21T11:26:55+02:00",
		},
		{
			name:             "GitHub Event pull_request with updated_at",
			workFlowMeatData: commitMetadata(t, "testdata/pull-request-event.json"),
			WantTime:         "2022-11-17T07:46:39Z",
		},
		{
			name:             "GitHub Event workflow_dispatch with head_commit",
			workFlowMeatData: commitMetadata(t, "testdata/github-context.json"),
			WantTime:         "2022-02-14T09:38:16+01:00",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			context, err := ParseContext(test.workFlowMeatData)
			assert.NoError(t, err)
			parsedEvent, err := ParseEvent(context.Event)
			assert.NoError(t, err)
			assert.NotNil(t, parsedEvent)
			assert.Equal(t, test.WantTime, parsedEvent.GetHeadCommitTimestamp())
		})
	}
}

func commitMetadata(t *testing.T, eventFile string) []byte {
	metadata, err := os.ReadFile(eventFile)
	assert.NoError(t, err)
	return metadata
}
