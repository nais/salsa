package nodejs

import (
	"encoding/json"
	"fmt"
	"github.com/nais/salsa/pkg/scan"
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
	}
	return npmMetadata
}
