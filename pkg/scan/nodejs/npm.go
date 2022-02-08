package nodejs

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nais/salsa/pkg/scan"
)

func NpmDeps(packageLockJsonContents string) ([]scan.Dependency, error) {
	var f interface{}
	err := json.Unmarshal([]byte(packageLockJsonContents), &f)
	if err != nil {
		return nil, fmt.Errorf("unable to parse package-lock.json: %v", err)
	}
	raw := f.(map[string]interface{})
	return transform(raw["dependencies"].(map[string]interface{})), nil
}

func transform(input map[string]interface{}) []scan.Dependency {
	deps := make([]scan.Dependency, 0)
	for key, value := range input {
		dependency := value.(map[string]interface{})
		integrity := fmt.Sprintf("%s", dependency["integrity"])
		shaDig := strings.Split(integrity, "-")
		deps = append(deps, scan.Dependency{
			Coordinates: key,
			Version:     fmt.Sprintf("%s", dependency["version"]),
			CheckSum: scan.CheckSum{
				Algorithm: fmt.Sprintf("%s", shaDig[0]),
				Digest:    fmt.Sprintf("%s", shaDig[1]),
			},
		})
	}
	return deps
}
