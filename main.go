package main

import (
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const pluginPrefix = "BUILDKITE_PLUGIN_BUILDPIPE_"

type Config struct {
	Projects []Project         `yaml:"projects"`
	Steps    []interface{}     `yaml:"steps"`
	Env      map[string]string `yaml:"env"`
}

func NewConfig(filename string) *Config {
	config := Config{}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file %s: %s\n", filename, err)
	}

	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalf("Error unmarshalling: %s\n", err)
	}

	return &config
}

func getAffectedProjects(projects []Project, changedFiles []string) []Project {
	affectedProjects := make([]Project, 0)
	for _, project := range projects {
		if project.checkAffected(changedFiles) {
			affectedProjects = append(affectedProjects, project)
		}
	}

	return affectedProjects
}

func projectsFromBuildProjects(buildProjects string, projects []Project) []Project {
	if buildProjects == "*" {
		return projects
	}

	projectNames := strings.Split(buildProjects, ",")

	affectedProjects := make([]Project, 0)
	for _, projectName := range projectNames {
		for _, configProject := range projects {
			if projectName == configProject.Label {
				affectedProjects  = append(affectedProjects, configProject)
			}
		}
	}
	return affectedProjects
}

func main() {
	logLevel := getEnv(pluginPrefix+"LOG_LEVEL", "info")
	ll, err := log.ParseLevel(logLevel)
	if err != nil {
		ll = log.InfoLevel
	}

	log.SetLevel(ll)

	config := NewConfig(os.Getenv(pluginPrefix + "DYNAMIC_PIPELINE"))
	buildProjects := os.Getenv(pluginPrefix+"BUILD_PROJECTS")

	var affectedProjects []Project
	if len(buildProjects) > 0 {
		affectedProjects = projectsFromBuildProjects(buildProjects, config.Projects)
	} else {
		changedFiles := getChangedFiles()
		if len(changedFiles) == 0 {
			log.Info("No files were changed")
			os.Exit(0)
		}

		affectedProjects = getAffectedProjects(config.Projects, changedFiles)
		if len(affectedProjects) == 0 {
			log.Info("No project was affected from git changes")
			os.Exit(0)
		}
	}

	pipeline := generatePipeline(config.Steps, config.Env, affectedProjects)

	uploadPipeline(*pipeline)
}
