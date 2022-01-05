package intoto

type Envelope struct {
	Payload string `json:"payload"`
}

type Attestation struct {
	Predicate Pred `json:"predicate"`
}

type Pred struct {
	Data string `json:"Data"`
}
