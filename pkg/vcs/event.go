package vcs

import (
	"encoding/json"
)

type EventInput struct {
	Inputs json.RawMessage `json:"inputs"`
}

type Event struct {
	Event eventMetadata `json:"event"`
}

type eventMetadata struct {
	HeadCommit headCommit `json:"head_commit"`
}

type headCommit struct {
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

func NewEvent(metadata []byte) *EventInput {
	return &EventInput{
		Inputs: metadata,
	}
}

func (in *EventInput) ParseEvent() (*Event, error) {
	var event Event
	err := json.Unmarshal(in.Inputs, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (in *Event) GetHeadCommitId() string {
	return in.Event.HeadCommit.Id
}

func (in *Event) GetHeadCommitTimestamp() string {
	return in.Event.HeadCommit.Timestamp
}
