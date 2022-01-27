package main

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
	log "github.com/sirupsen/logrus"
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
	Label       string
	Path        StringArray
	ExcludePath StringArray `yaml:"exclude_path"`
	Skip        StringArray
	Env         map[string]string
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
func matchPath(p string, f string) bool {
	// If the path contains a glob, the `doublestar.Match`
	// method is used to determine the match,
	// otherwise `strings.HasPrefix` is used.
	if strings.Contains(p, "*") {
		match, err := doublestar.Match(p, f)
		if err != nil {
			log.Errorf("path matching failed: %v", err)
			return false
		}
		if match {
			return true
		}
	}
	if strings.HasPrefix(f, p) {
		return true
	}
	return false
}

func (p *Project) checkAffected(changedFiles []string) bool {
	filteredChangedFiles := p.filterExcludedFiles(changedFiles)

	for _, filePath := range p.Path {
		if filePath == "." {
			return true
		}
		normalizedPath := path.Clean(filePath)

		for _, changedFile := range filteredChangedFiles {
			if matched := matchPath(normalizedPath, changedFile); matched {
				return true
			}
		}
	}
	return false
}

func (p *Project) filterExcludedFiles(changedFiles []string) []string {
	result := make([]string, 0, len(changedFiles))

	if len(p.ExcludePath) <= 0 {
		return changedFiles
	}

CHANGED_FILES_LOOP:
	for _, changedFile := range changedFiles {
		for _, excludedPath := range p.ExcludePath {
			if matched := matchPath(excludedPath, changedFile); matched {
				continue CHANGED_FILES_LOOP
			}
		}
		result = append(result, changedFile)
	}

	return result
}
