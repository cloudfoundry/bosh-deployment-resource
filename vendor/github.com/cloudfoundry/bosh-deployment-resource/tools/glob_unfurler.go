package tools

import (
	"fmt"
	"path"
	"path/filepath"
)

func UnfurlGlobs(baseDirectory string, globs []string) ([]string, error) {
	allPaths := []string{}
	uniquePaths := map[string]bool{}
	for _, glob := range globs {
		nextPaths, err := filepath.Glob(path.Join(baseDirectory, glob))
		if err != nil {
			return []string{}, fmt.Errorf("%s is not a valid file glob", glob)
		}

		if len(nextPaths) == 0 {
			return []string{}, fmt.Errorf("%s does not match any files", glob)
		}

		for _, nextPath := range nextPaths {
			if !uniquePaths[nextPath] {
				allPaths = append(allPaths, nextPath)
				uniquePaths[nextPath] = true
			}
		}
	}

	return allPaths, nil
}
