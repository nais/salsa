package vcs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestEvenCommit(t *testing.T) {
	metadata := commitMetadata(t)
	event := NewEvent(metadata)
	parsedEvent, err := event.ParseEvent()
	assert.NoError(t, err)
	assert.NotNil(t, parsedEvent)
	assert.Equal(t, "d4cd018b2fe54d8308b78f2bb88db94ac57173ea", parsedEvent.GetHeadCommitId())
	_, err = time.Parse(time.RFC3339, parsedEvent.GetHeadCommitTimestamp())
	assert.NoError(t, err)
	assert.Equal(t, "2022-10-21T11:26:55+02:00", parsedEvent.GetHeadCommitTimestamp())
}

func commitMetadata(t *testing.T) []byte {
	metadata, err := os.ReadFile("testdata/event-commit.json")
	assert.NoError(t, err)
	return metadata
}
