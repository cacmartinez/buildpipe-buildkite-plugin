package main

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
)

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// https://github.com/go-yaml/yaml/issues/100
type StringArray []string

func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

type Project struct {
	Label string
	Path  StringArray
	Skip  StringArray
	Env   map[string]string
}

func (p *Project) getMainPath() string {
	return p.Path[0]
}

func (p *Project) checkProjectRules(step map[interface{}]interface{}) bool {
	for _, pattern := range p.Skip {
		label := step["label"].(string)
		if matched, _ := filepath.Match(pattern, label); matched {
			return false
		}
	}
	return true
}

// matchPath checks if the file f matches the path p.
// Taken from https://github.com/chronotc/monorepo-diff-buildkite-plugin/blob/602e650f10d54026fab466521c3202a04ae88afe/pipeline.go#L102
func matchPath(p string, f string) (bool, error) {
	// If the path contains a glob, the `doublestar.Match`
	// method is used to determine the match,
	// otherwise `strings.HasPrefix` is used.
	if strings.Contains(p, "*") {
		match, err := doublestar.Match(p, f)
		if err != nil {
			return false, fmt.Errorf("path matching failed: %v", err)
		}
		if match {
			return true, nil
		}
	}
	if strings.HasPrefix(f, p) {
		return true, nil
	}
	return false, nil
}

func (p *Project) checkAffected(changedFiles []string) bool {
	for _, filePath := range p.Path {
		if filePath == "." {
			return true
		}
		normalizedPath := path.Clean(filePath)
		for _, changedFile := range changedFiles {
			if matched, _ := matchPath(normalizedPath, changedFile); matched {
				return true
			}
		}

		// projectDirs := strings.Split(normalizedPath, "/")
		// for _, changedFile := range changedFiles {
		// 	changedDirs := strings.Split(changedFile, "/")
		// 	if reflect.DeepEqual(changedDirs[:Min(len(projectDirs), len(changedDirs))], projectDirs) {
		// 		return true
		// 	}
		// }
	}
	return false
}
