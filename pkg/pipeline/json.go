package pipeline

import (
	"encoding/json"
)

// GitRepo represents a git repository
type GitRepo struct {
	URL    string
	Branch string
}

// FindGit finds the git repositories from a JSON format Jenkinsfile
func FindGit(jenkinsfile string) (repos []GitRepo, err error) {
	var j Jenkinsfile
	if j, err = parseJenkinsfile(jenkinsfile); err == nil {
		for _, stage := range j.Pipeline.Stages {
			targetStep := stage.GetStep("git")
			if targetStep == nil {
				continue
			}

			arg := targetStep.GetArgument("url")
			if arg == nil {
				continue
			}

			branchArg := targetStep.GetArgument("branch")
			repo := GitRepo{
				URL: arg.Value.Value.(string),
			}
			if branchArg != nil {
				repo.Branch = branchArg.Value.Value.(string)
			}
			repos = append(repos, repo)
		}
	}
	return
}

func parseJenkinsfile(jenkinsfile string) (j Jenkinsfile, err error) {
	err = json.Unmarshal([]byte(jenkinsfile), &j)
	return
}

// Jenkinsfile represents a structured Jenkinsfile
//
// We could marshal or unmarshal the JSON format Jenkinsfile
type Jenkinsfile struct {
	Pipeline Pipeline `json:"pipeline"`
}

// Pipeline represents a Pipeline
type Pipeline struct {
	Stages []Stage `json:"stages"`
}

// Stage represents a stage in a Jenkinsfile
type Stage struct {
	Name     string        `json:"name"`
	Branches []StageBranch `json:"branches"`
}

// GetStep finds step by name
func (s Stage) GetStep(name string) (step *Step) {
	for _, branch := range s.Branches {
		for _, s := range branch.Steps {
			if s.Name == name {
				step = &s
				return
			}
		}
	}
	return
}

// StageBranch contains the steps
type StageBranch struct {
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

// Step represents a step of the a Pipeline
type Step struct {
	Name      string         `json:"name"`
	Arguments []StepArgument `json:"arguments"`
}

// GetArgument finds the argument by name
func (s Step) GetArgument(name string) (arg *StepArgument) {
	for _, item := range s.Arguments {
		if item.Key == name {
			arg = &item
			break
		}
	}
	return
}

// StepArgument represents a step argument
type StepArgument struct {
	Key   string            `json:"key"`
	Value StepArgumentValue `json:"value"`
}

// StepArgumentValue represents the value of a step argument
type StepArgumentValue struct {
	IsLiteral bool        `json:"isLiteral"`
	Value     interface{} `json:"value"`
}
