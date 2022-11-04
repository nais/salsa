package vcs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestEvenHeadCommit(t *testing.T) {
	metadata := commitMetadata(t)
	context, err := ParseContext(metadata)
	assert.NoError(t, err)
	parsedEvent, err := ParseEvent(context.Event)
	assert.NoError(t, err)
	assert.NotNil(t, parsedEvent)
	assert.Equal(t, "d4cd018b2fe54d8308b78f2bb88db94ac57173ea", parsedEvent.GetHeadCommitId())
	_, err = time.Parse(time.RFC3339, parsedEvent.GetHeadCommitTimestamp())
	assert.NoError(t, err)
	assert.Equal(t, "2022-10-21T11:26:55+02:00", parsedEvent.GetHeadCommitTimestamp())
}

func commitMetadata(t *testing.T) []byte {
	metadata, err := os.ReadFile("testdata/event-head-commit.json")
	assert.NoError(t, err)
	return metadata
}
