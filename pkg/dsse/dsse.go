package dsse

import v02 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"

type Envelope struct {
	Payload string `json:"payload"`
}

type Attestation struct {
	Predicate v02.ProvenancePredicate `json:"predicate"`
}
