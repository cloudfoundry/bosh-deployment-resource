package tools

import (
	"path/filepath"
	"fmt"
)

func UnfurlGlobs(globs ...string) ([]string, error) {
	allPaths := []string{}
	uniquePaths := map[string]bool{}
	for _, glob := range globs {
		nextPaths, err := filepath.Glob(glob)
		if err != nil {
			return []string{}, fmt.Errorf("%s is not a valid file glob", glob)
		}

		for _, nextPath := range nextPaths {
			uniquePaths[nextPath] = true
		}
	}
	for key := range uniquePaths {
		allPaths = append(allPaths, key)
	}
	return allPaths, nil
}


