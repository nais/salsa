package intoto

import (
	"encoding/base64"
	"encoding/json"
	"github.com/in-toto/in-toto-golang/in_toto"
)

func ParseEnvelope(dsseEnvelope []byte) (*in_toto.ProvenanceStatement, error) {
	var env = Envelope{}
	err := json.Unmarshal(dsseEnvelope, &env)
	if err != nil {
		return nil, err
	}
	decodedPayload, err := base64.StdEncoding.DecodeString(env.Payload)
	if err != nil {
		return nil, err
	}
	var stat = &in_toto.ProvenanceStatement{}

	err = json.Unmarshal(decodedPayload, &stat)
	if err != nil {
		return nil, err
	}
	return stat, nil
}
