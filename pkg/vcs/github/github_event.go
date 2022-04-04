package github

import (
	"encoding/json"
	"fmt"
	"time"
)

type EventMetadata struct {
	HeadCommit HeadCommit `json:"head_commit"`
}

type HeadCommit struct {
	Timestamp string `json:"timestamp"`
}

func ParseEventMetaData(githubContext *Context) (*EventMetadata, error) {
	eventData := EventMetadata{}

	if err := json.Unmarshal(githubContext.Event, &eventData); err != nil {
		if err != nil {
			return nil, fmt.Errorf("unmarshal github event meatdata: %w", err)
		}
	}

	return &eventData, nil
}

func (in EventMetadata) BuildStartedOn() (time.Time, error) {
	t, err := time.Parse(time.RFC3339, in.HeadCommit.Timestamp)

	if err != nil {
		return time.Time{}, fmt.Errorf("event metdata timestamp")
	}
	return t, nil
}
