package vcs

import (
	"encoding/json"
)

type Event struct {
	EventMetadata *EventMetadata `json:"event"`
}

type EventMetadata struct {
	HeadCommit *HeadCommit `json:"head_commit"`
}

type HeadCommit struct {
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

func ParseEvent(inputs []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(inputs, &event.EventMetadata)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (in *Event) GetHeadCommitId() string {
	return in.EventMetadata.HeadCommit.Id
}

func (in *Event) GetHeadCommitTimestamp() string {
	return in.EventMetadata.HeadCommit.Timestamp
}
