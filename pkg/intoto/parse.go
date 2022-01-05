package intoto

import (
	"encoding/base64"
	"encoding/json"
	"github.com/in-toto/in-toto-golang/in_toto"
)

func ParseEnvelope(dsseEnvelope []byte) (*in_toto.Statement, error) {
	var env = Envelope{}
	err := json.Unmarshal(dsseEnvelope, &env)
	if err != nil {
		return nil, err
	}
	decodedPayload, err := base64.StdEncoding.DecodeString(env.Payload)
	if err != nil {
		return nil, err
	}
	var att = Attestation{}
	err = json.Unmarshal(decodedPayload, &att)
	if err != nil {
		return nil, err
	}
	var stmt in_toto.Statement
	err = json.Unmarshal([]byte(att.Predicate.Data), &stmt)
	return &stmt, nil
}
