package job

import (
	"encoding/json"
	"fmt"
	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	searchAPIPrefix       = "/blue/rest/search"
	organizationAPIPrefix = "/blue/rest/organizations"
)

// BlueOceanClient is client for operating pipelines via BlueOcean RESTful API.
type BlueOceanClient struct {
	core.JenkinsCore
	Organization string
}

// Parameter contains name and value of an option.
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Search searches jobs via the BlueOcean API
func (c *BlueOceanClient) Search(name string, start, limit int) (items []JenkinsItem, err error) {
	api := fmt.Sprintf("%s/?q=pipeline:*%s*;type:pipeline;organization:%s;excludedFromFlattening=jenkins.branch.MultiBranchProject,com.cloudbees.hudson.plugins.folder.AbstractFolder&filter=no-folders&start=%d&limit=%d",
		searchAPIPrefix, name, c.Organization, start, limit)
	err = c.RequestWithData(http.MethodGet, api,
		nil, nil, 200, &items)
	return
}

// BuildOption contains some options of Build method.
type BuildOption struct {
	Pipelines  []string
	Parameters []Parameter
	Branch     string
}

// GetBuildOption contains some options while getting a specific build.
type GetBuildOption struct {
	Pipelines []string
	RunID     string
	Branch    string
}

// Build builds a pipeline for specific organization and pipelines.
func (c *BlueOceanClient) Build(option BuildOption) (*PipelineBuild, error) {
	var pb PipelineBuild
	var payloadReader io.Reader
	if len(option.Parameters) > 0 {
		// ignore this error due to never happened
		payloadBytes, _ := json.Marshal(map[string][]Parameter{
			"parameters": option.Parameters,
		})
		payloadReader = strings.NewReader(string(payloadBytes))
	}
	err := c.RequestWithData(http.MethodPost, c.getBuildAPI(option), getHeaders(), payloadReader, 200, &pb)
	if err != nil {
		return nil, err
	}
	return &pb, nil
}

func (c *BlueOceanClient) getBuildAPI(option BuildOption) string {
	api := fmt.Sprintf("%s/%s/%s", organizationAPIPrefix, c.Organization, parsePipelinePath(option.Pipelines))
	if option.Branch != "" {
		api = fmt.Sprintf("%s/branches/%s", api, url.PathEscape(option.Branch))
	}
	api = fmt.Sprintf("%s/runs/", api)
	return api
}

// GetBuild gets build result for specific organization, run ID and pipelines.
func (c *BlueOceanClient) GetBuild(option GetBuildOption) (*PipelineBuild, error) {
	var pb PipelineBuild
	err := c.RequestWithData(http.MethodGet, c.getGetBuildAPI(option), getHeaders(), nil, 200, &pb)
	if err != nil {
		return nil, err
	}
	return &pb, nil
}

func (c *BlueOceanClient) getGetBuildAPI(option GetBuildOption) string {
	api := fmt.Sprintf("%s/%s/%s", organizationAPIPrefix, c.Organization, parsePipelinePath(option.Pipelines))
	if option.Branch != "" {
		api = fmt.Sprintf("%s/branches/%s", api, url.PathEscape(option.Branch))
	}
	api = fmt.Sprintf("%s/runs/%s/", api, option.RunID)
	return api
}

func getHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}

// PipelineBuild represents a build detail of Pipeline.
type PipelineBuild struct {
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
