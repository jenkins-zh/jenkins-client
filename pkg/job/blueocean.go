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
	// validate option
	api := fmt.Sprintf("%s/%s/%s", organizationAPIPrefix, c.Organization, parsePipelinePath(option.Pipelines))
	if option.Branch != "" {
		api = fmt.Sprintf("%s/branches/%s", api, url.PathEscape(option.Branch))
	}
	api = fmt.Sprintf("%s/runs/", api)
	return api
}

// GetBuildOption contains some options while getting a specific build.
type GetBuildOption struct {
	Pipelines []string
	RunID     string
	Branch    string
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
	api := c.getBuildAPI(BuildOption{
		Pipelines: option.Pipelines,
		Branch:    option.Branch,
	})
	api = api + option.RunID + "/"
	return api
}

// GetNodesOption contains some options while getting nodes detail.
type GetNodesOption struct {
	Pipelines []string
	Branch    string
	RunID     string
	Limit     int
}

// GetNodes gets nodes details
func (c *BlueOceanClient) GetNodes(option GetNodesOption) ([]PipelineRunNode, error) {
	var nodes []PipelineRunNode
	err := c.RequestWithData(http.MethodGet, c.getGetNodesAPI(option), getHeaders(), nil, 200, &nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (c *BlueOceanClient) getGetNodesAPI(option GetNodesOption) string {
	api := c.getGetBuildAPI(GetBuildOption{
		Pipelines: option.Pipelines,
		Branch:    option.Branch,
		RunID:     option.RunID,
	})
	limit := option.Limit
	if limit == 0 {
		// if limit is not set
		limit = 10000
	}
	api = fmt.Sprintf("%snodes/?limit=%d", api, limit)
	return api
}

func getHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}

// PipelineBuild represents a build detail of Pipeline.
type PipelineBuild struct {
	Actions                   []interface{} `json:"actions,omitempty"`
	ArtifactsZipFile          interface{}   `json:"artifactsZipFile,omitempty"`
	CauseOfBlockage           string        `json:"causeOfBlockage,omitempty"`
	Causes                    []interface{} `json:"causes,omitempty"`
	ChangeSet                 []interface{} `json:"changeSet,omitempty"`
	Description               interface{}   `json:"description,omitempty"`
	DurationInMillis          interface{}   `json:"durationInMillis,omitempty"`
	EnQueueTime               Time          `json:"enQueueTime,omitempty"`
	EndTime                   Time          `json:"endTime,omitempty"`
	EstimatedDurationInMillis interface{}   `json:"estimatedDurationInMillis,omitempty"`
	ID                        string        `json:"id,omitempty"`
	Name                      interface{}   `json:"name,omitempty"`
	Organization              string        `json:"organization,omitempty"`
	Pipeline                  string        `json:"pipeline,omitempty"`
	Replayable                bool          `json:"replayable,omitempty"`
	Result                    string        `json:"result,omitempty"`
	RunSummary                interface{}   `json:"runSummary,omitempty"`
	StartTime                 Time          `json:"startTime,omitempty"`
	State                     string        `json:"state,omitempty"`
	Type                      string        `json:"type,omitempty"`
	QueueID                   string        `json:"queueId,omitempty"`
}

type PipelineRunNode struct {
	DisplayDescription string `json:"displayDescription,omitempty"`
	DisplayName        string `json:"displayName,omitempty"`
	DurationInMillis   int    `json:"durationInMillis,omitempty"`
	ID                 string `json:"id,omitempty"`
	Input              *Input `json:"input,omitempty"`
	Result             string `json:"result,omitempty"`
	StartTime          Time   `json:"startTime,omitempty"`
	State              string `json:"state,omitempty"`
	Type               string `json:"type,omitempty"`
	CauseOfBlockage    string `json:"causeOfBlockage,omitempty"`
	Edges              []Edge `json:"edges,omitempty"`
	FirstParent        string `json:"firstParent,omitempty"`
	Restartable        bool   `json:"restartable,omitempty"`
}

type Edge struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type Input struct {
	ID         string                `json:"id,omitempty"`
	Message    string                `json:"message,omitempty"`
	Ok         string                `json:"ok,omitempty"`
	Parameters []ParameterDefinition `json:"parameters,omitempty"`
	Submitter  string                `json:"submitter,omitempty"`
}
