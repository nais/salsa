package vcs

import "encoding/json"

type Event struct {
	Inputs json.RawMessage `json:"inputs"`
}

type Commits struct {
	Commits []Commit `json:"commits"`
}

type Commit struct {
	Timestamp string `json:"timestamp"`
	After     string `json:"after"`
}

func (in *Event) GetCommits() ([]Commit, error) {
	var commits []Commit
	err := json.Unmarshal(in.Inputs, &commits)
	if err != nil {
		return nil, err
	}
	return commits, nil
}
