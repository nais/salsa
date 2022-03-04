package digest

type Digest string

const (
	AlgorithmSHA256 = "sha256"
	AlgorithmSHA1   = "sha1"
)

// func (d Digest) DecodeToString() (string, error) {
// 	// Better to keep it bas64 encoded, can be decoded when verifying is required
// 	decoded, err := base64.StdEncoding.DecodeString(string(d))
// 	if err != nil {
// 		return "", fmt.Errorf("decoding checksum")
// 	}
// 	return fmt.Sprintf("%x", decoded), nil
// }
