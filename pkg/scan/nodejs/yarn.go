package nodejs

import (
	"github.com/nais/salsa/pkg/scan"
	"regexp"
	"strings"
)

func YarnDeps(yarnLockContents string) *scan.BuildToolMetadata {
	var lastDepName = ""
	var checksums = make(map[string]scan.CheckSum)
	yarnMetadata := scan.CreateMetadata()
	lines := strings.Split(yarnLockContents, "\n")
	startOfEntries := determineStartsOfEntries(lines)
	for _, startOfEntry := range startOfEntries {
		if !isIntegrity(lines[startOfEntry]) {
			depName := trimParseNameDepWithRange(lines[startOfEntry])
			depVersion := parseVersion(lines[startOfEntry+1])
			yarnMetadata.Deps[depName] = depVersion
			lastDepName = depName
		} else {
			checksums[lastDepName] = addYarnShaDigest(lines[startOfEntry], checksums[lastDepName])
		}
	}
	yarnMetadata.Checksums = checksums
	return yarnMetadata
}

func determineStartsOfEntries(yarnLockLines []string) []int {
	var startsOfEntries []int
	previousLineisUseful := false
	for lineNr := range yarnLockLines {
		if isUseful(yarnLockLines[lineNr]) && !previousLineisUseful || isIntegrity(yarnLockLines[lineNr]) {
			startsOfEntries = append(startsOfEntries, lineNr)
			previousLineisUseful = true
		} else if !isUseful(yarnLockLines[lineNr]) {
			previousLineisUseful = false
		}
	}
	return startsOfEntries
}

func isUseful(line string) bool {
	return !(startsWithHash(line) || isEmpty(line))
}

func isEmpty(line string) bool {
	return strings.TrimSpace(line) == ""
}

func startsWithHash(line string) bool {
	return strings.HasPrefix(line, "#")
}

func isIntegrity(line string) bool {
	return strings.Contains(line, "integrity")
}

func trimParseNameDepWithRange(depName string) string {
	if len(strings.Split(depName, ", ")) > 1 {
		depName = parseName(depName)
		return strings.Split(depName, ", ")[1]
	} else {
		return parseName(depName)
	}
}

func parseName(line string) string {
	regex := regexp.MustCompile(`^"?(?P<pkgname>.*)@[^~]?.*$`)
	matches := regex.FindStringSubmatch(line)
	pkgnameIndex := regex.SubexpIndex("pkgname")
	return matches[pkgnameIndex]
}

func parseVersion(line string) string {
	regex := regexp.MustCompile(`.*"(?P<pkgversion>.*)"$`)
	matches := regex.FindStringSubmatch(line)
	pkgversionIndex := regex.SubexpIndex("pkgversion")
	return matches[pkgversionIndex]
}

func addYarnShaDigest(line string, s scan.CheckSum) scan.CheckSum {
	pattern := regexp.MustCompile(`(?P<name>integrity.*)`)
	integrity := pattern.FindString(line)
	trimPrefixIntegrity := strings.TrimPrefix(integrity, "integrity ")
	algoDigest := strings.Split(trimPrefixIntegrity, "-")
	s.Algorithm = strings.TrimSpace(algoDigest[0])
	// Better to keep the digest base64 encoded when signing envelope
	// decodedDigest, _ := base64.StdEncoding.DecodeString(algoDigest[1])
	// s.Digest = fmt.Sprintf("%x", decodedDigest)
	s.Digest = strings.TrimSpace(algoDigest[1])
	return s
}
