package nodejs

import (
	"encoding/json"
	"fmt"
	"github.com/nais/salsa/pkg/scan"
	"strings"
)

func NpmDeps(packageLockJsonContents string) (*scan.BuildToolMetadata, error) {
	var f interface{}
	err := json.Unmarshal([]byte(packageLockJsonContents), &f)
	if err != nil {
		return nil, fmt.Errorf("unable to parse package-lock.json: %v", err)
	}

	raw := f.(map[string]interface{})
	return transform(raw["dependencies"].(map[string]interface{})), nil
}

func transform(input map[string]interface{}) *scan.BuildToolMetadata {
	npmMetadata := scan.CreateMetadata()
	for key, value := range input {
		dependency := value.(map[string]interface{})
		npmMetadata.Deps[key] = fmt.Sprintf("%s", dependency["version"])
		integrity := fmt.Sprintf("%s", dependency["integrity"])
		shaDig := strings.Split(integrity, "-")
		npmMetadata.Checksums[key] = scan.CheckSum{
			Algorithm: fmt.Sprintf("%s", shaDig[0]),
			Digest:    fmt.Sprintf("%s", shaDig[1]),
		}
	}
	return npmMetadata
}
