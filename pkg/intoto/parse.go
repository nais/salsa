package intoto

import (
	"encoding/base64"
	"encoding/json"
)

func ParseEnvelope(dsseEnvelope []byte) (*Statement, error) {
	var env = Envelope{}
	err := json.Unmarshal(dsseEnvelope, &env)
	if err != nil {
		return nil, err
	}
	decoded, err := base64.StdEncoding.DecodeString(env.Payload)
	if err != nil {
		return nil, err
	}
	var att = Attestation{}
	err = json.Unmarshal(decoded, &att)
	if err != nil {
		return nil, err
	}
	var stmt Statement
	err = json.Unmarshal([]byte(att.Predicate.Data), &stmt)
	return &stmt, nil
}
