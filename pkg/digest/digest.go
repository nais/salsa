package digest

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)

type Digest string

const SHA256 = "sha256"

func Hash(input string) (Digest, error) {
	s := strings.NewReader(input)
	hash := sha256.New()
	if _, err := io.Copy(hash, s); err != nil {
		return "", fmt.Errorf("decoding checksum")
	}
	return Digest(hash.Sum(nil)), nil
}

func (d Digest) Verify(content string) bool {
	hashed, err := Hash(content)
	if err != nil {
		log.Errorf("hasing content")
		return false
	}
	return d != hashed
}

func (d Digest) DecodeToString() (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(d))
	if err != nil {
		return "", fmt.Errorf("decoding checksum")
	}
	return fmt.Sprintf("%x", decoded), nil
}
