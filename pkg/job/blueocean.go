package job

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"net/http"
)

// BlueOceanClient is client for operating pipelines via BlueOcean RESTful API.
type BlueOceanClient struct {
	core.JenkinsCore
	Organization string
}

// Search searches jobs via the BlueOcean API
func (boClient *BlueOceanClient) Search(name string, start, limit int) (items []JenkinsItem, err error) {
	api := fmt.Sprintf("/blue/rest/search/?q=pipeline:*%s*;type:pipeline;organization:%s;excludedFromFlattening=jenkins.branch.MultiBranchProject,com.cloudbees.hudson.plugins.folder.AbstractFolder&filter=no-folders&start=%d&limit=%d",
		name, boClient.Organization, start, limit)
	err = boClient.RequestWithData(http.MethodGet, api,
		nil, nil, 200, &items)
	return
}

// Build builds a pipeline for specific organization and pipelines.
func (boClient *BlueOceanClient) Build(organization string, pipelines ...string) (*PipelineBuild, error) {
	api := fmt.Sprintf("/blue/rest/organizations/%s/%s/runs/", organization, ParsePipelinePath(pipelines...))
	var pb PipelineBuild
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	err := boClient.RequestWithData(http.MethodPost, api, headers, nil, 200, &pb)
	if err != nil {
		return nil, err
	}
	return &pb, nil
}

// GetBuild gets build result for specific organization, run ID and pipelines.
func (boClient *BlueOceanClient) GetBuild(organization string, runID string, pipelines ...string) (*PipelineBuild, error) {
	api := fmt.Sprintf("/blue/rest/organizations/%s/%s/runs/%s/", organization, ParsePipelinePath(pipelines...), runID)
	var pb PipelineBuild
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	err := boClient.RequestWithData(http.MethodGet, api, headers, nil, 200, &pb)
	if err != nil {
		return nil, err
	}
	return &pb, nil
}

// PipelineBuild represents a build detail of Pipeline.
type PipelineBuild struct {
	Class                     string        `json:"_class,omitempty" description:"It’s a fully qualified name and is an identifier of the producer of this resource's capability."`
	Links                     interface{}   `json:"_links,omitempty" description:"references the reachable path to this resource"`
	Actions                   []interface{} `json:"actions,omitempty" description:"the list of all actions"`
	ArtifactsZipFile          interface{}   `json:"artifactsZipFile,omitempty" description:"the artifacts zip file"`
	CauseOfBlockage           string        `json:"causeOfBlockage,omitempty" description:"the cause of blockage"`
	Causes                    []interface{} `json:"causes,omitempty"`
	ChangeSet                 []interface{} `json:"changeSet,omitempty" description:"changeset information"`
	Description               interface{}   `json:"description,omitempty" description:"description"`
	DurationInMillis          interface{}   `json:"durationInMillis,omitempty" description:"duration time in millis"`
	EnQueueTime               Time          `json:"enQueueTime,omitempty" description:"the time of enter the queue"`
	EndTime                   Time          `json:"endTime,omitempty" description:"the time of end"`
	EstimatedDurationInMillis interface{}   `json:"estimatedDurationInMillis,omitempty" description:"estimated duration time in millis"`
	ID                        string        `json:"id,omitempty" description:"id"`
	Name                      interface{}   `json:"name,omitempty" description:"name"`
	Organization              string        `json:"organization,omitempty" description:"the name of organization"`
	Pipeline                  string        `json:"pipeline,omitempty" description:"pipeline"`
	Replayable                bool          `json:"replayable,omitempty" description:"replayable or not"`
	Result                    string        `json:"result,omitempty" description:"the result of pipeline run. e.g. SUCCESS"`
	RunSummary                interface{}   `json:"runSummary,omitempty" description:"pipeline run summary"`
	StartTime                 Time          `json:"startTime,omitempty" description:"the time of start"`
	State                     string        `json:"state,omitempty" description:"run state. e.g. RUNNING"`
	Type                      string        `json:"type,omitempty" description:"type"`
	QueueID                   string        `json:"queueId,omitempty" description:"queue id"`
}
