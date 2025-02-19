// Apply provides a CLI to deploy a project config to GCP.
//
// Usage:
//   $ bazel run :apply -- --project_yaml_path=${PROJECT_YAML_PATH?} --project=${PROJECT_ID?}
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"flag"
	
	"github.com/GoogleCloudPlatform/healthcare/deploy/apply"
	"github.com/GoogleCloudPlatform/healthcare/deploy/config"
	"github.com/ghodss/yaml"
)

var (
	projectYAMLPath = flag.String("project_yaml_path", "", "Path to project yaml file")
	projectID       = flag.String("project", "", "Project within the project yaml file to deploy config resources for")
)

func main() {
	flag.Parse()

	if *projectYAMLPath == "" {
		log.Fatal("--project_yaml_path must be set")
	}
	if *projectID == "" {
		log.Fatal("--project must be set")
	}

	// TODO: handle split yaml configs
	b, err := ioutil.ReadFile(*projectYAMLPath)
	if err != nil {
		log.Fatalf("failed to read input projects yaml file at path %q: %v", *projectYAMLPath, err)
	}

	conf := new(config.Config)
	if err := yaml.Unmarshal(b, conf); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
	if err := conf.Init(); err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	proj, err := findProject(*projectID, conf)
	if err != nil {
		log.Fatal(err)
	}

	if err := apply.Apply(conf, proj); err != nil {
		log.Fatalf("failed to deploy %q resources: %v", *projectID, err)
	}

	log.Println("Config deployed successfully")
}

func findProject(id string, c *config.Config) (*config.Project, error) {
	for _, p := range c.AllProjects() {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("failed to find project %q", id)
}
