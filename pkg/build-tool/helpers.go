package build_tool

import (
	"fmt"
	"github.com/briandowns/spinner"
	"time"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func sumSupported(buildFiles ...[]string) (supported []string) {
	for _, supportedBuildFiles := range buildFiles {
		for _, s := range supportedBuildFiles {
			supported = append(supported, s)
		}
	}
	return supported
}

func StartSpinner(project string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 150*time.Millisecond)
	s.Suffix = "\n"
	s.FinalMSG = fmt.Sprintf("provenace for %s finished", project)
	s.Start()
	return s
}
