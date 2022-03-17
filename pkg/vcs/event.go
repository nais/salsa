package vcs

import "encoding/json"

type Event struct {
	Inputs json.RawMessage `json:"inputs"`
}
