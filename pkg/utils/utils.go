package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

func ProvenanceFile(repoName string) string {
	return fmt.Sprintf("%s.provenance", repoName)
}

func StartSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 150*time.Millisecond)
	s.FinalMSG = message
	s.Start()
	return s
}

func DecodeDigest(base64Encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return "", fmt.Errorf("decoding base64 encoded checksum")
	}
	return fmt.Sprintf("%x", decoded), nil
}
