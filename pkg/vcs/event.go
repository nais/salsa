package vcs

import (
	"encoding/json"
	"time"
)

type Event struct {
	EventMetadata *EventMetadata `json:"event"`
}

type EventMetadata struct {
	HeadCommit  *HeadCommit  `json:"head_commit"`
	PullRequest *PullRequest `json:"pull_request"`
	WorkFlowRun *WorkFlow    `json:"workflow_run"`
}

type HeadCommit struct {
	Timestamp string `json:"timestamp"`
}

type PullRequest struct {
	UpdatedAt string `json:"updated_at"`
}

type WorkFlow struct {
	HeadCommit *HeadCommit `json:"head_commit"`
}

func ParseEvent(inputs []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(inputs, &event.EventMetadata)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (in *Event) GetHeadCommitTimestamp() string {
	if in.EventMetadata.HeadCommit != nil {
		return in.EventMetadata.HeadCommit.Timestamp
	}

	if in.EventMetadata.WorkFlowRun != nil {
		return in.EventMetadata.WorkFlowRun.HeadCommit.Timestamp
	}

	if in.EventMetadata.PullRequest != nil {
		return in.EventMetadata.PullRequest.UpdatedAt
	}

	return time.Now().UTC().Round(time.Second).Format(time.RFC3339)
}
