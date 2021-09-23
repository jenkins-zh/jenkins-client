package job

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"strings"

	"github.com/jenkins-zh/jenkins-client/pkg/core"
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

// GetPipelines returns the Pipeline list which comes from the possible nest folders
func (c *BlueOceanClient) GetPipelines(folders ...string) (pipelines []Pipeline, err error) {
	api := c.getPipelineAPI(folders...)
	err = c.RequestWithData(http.MethodGet, api,
		nil, nil, 200, &pipelines)
	return
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
func (c *BlueOceanClient) Build(option BuildOption) (*PipelineRun, error) {
	var pr PipelineRun
	var payloadReader io.Reader
	// we allow developers to pass an empty parameters, but nil parameters
	if option.Parameters != nil {
		// ignore this error due to never happened
		payloadBytes, _ := json.Marshal(map[string][]Parameter{
			"parameters": option.Parameters,
		})
		payloadReader = strings.NewReader(string(payloadBytes))
	}
	err := c.RequestWithData(http.MethodPost, c.getBuildAPI(option), getHeaders(), payloadReader, 200, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
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
func (c *BlueOceanClient) GetBuild(option GetBuildOption) (*PipelineRun, error) {
	var pr PipelineRun
	err := c.RequestWithData(http.MethodGet, c.getGetBuildAPI(option), getHeaders(), nil, 200, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (c *BlueOceanClient) getPipelineAPI(folders ...string) (api string) {
	api = fmt.Sprintf("%s/%s/pipelines", organizationAPIPrefix, c.Organization)
	for _, folder := range folders {
		api = fmt.Sprintf("%s/%s/pipelines/", api, folder)
	}
	return
}

// GetPipelineRuns returns a PipelineRun which in the possible nest folders
func (c *BlueOceanClient) GetPipelineRuns(pipeline string, folders ...string) (runs []PipelineRun, err error) {
	api := c.getPipelineAPI(folders...)
	api = fmt.Sprintf("%s/%s/runs/", api, pipeline)
	err = c.RequestWithData(http.MethodGet, api,
		nil, nil, 200, &runs)
	return
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
func (c *BlueOceanClient) GetNodes(option GetNodesOption) ([]Node, error) {
	var nodes []Node
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

// ReplayOption contains some options while replaying a PipelineRun.
type ReplayOption struct {
	Folders []string
	Branch  string
	RunID   string
}

// Replay will queue up a replay of the pipeline run with the same commit id as the run used.
// Reference: https://github.com/jenkinsci/blueocean-plugin/tree/master/blueocean-rest#replay-a-pipeline-build
func (c *BlueOceanClient) Replay(option ReplayOption) (*PipelineRun, error) {
	pipelineRun := &PipelineRun{}
	if err := c.RequestWithData(http.MethodPost, c.getReplayAPI(&option), nil, nil, 200, pipelineRun); err != nil {
		return nil, err
	}
	return pipelineRun, nil
}

func (c *BlueOceanClient) getReplayAPI(option *ReplayOption) string {
	api := c.getGetBuildAPI(GetBuildOption{
		Pipelines: option.Folders,
		Branch:    option.Branch,
		RunID:     option.RunID,
	})
	api = api + "replay/"
	return api
}

func getHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}

// PipelineRun represents a build detail of Pipeline.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/a7cbc946b73d89daf9dfd91cd713cc7ab64a2d95/blueocean-pipeline-api-impl/src/main/java/io/jenkins/blueocean/rest/impl/pipeline/PipelineRunImpl.java
type PipelineRun struct {
	ArtifactsZipFile          interface{}   `json:"artifactsZipFile,omitempty"`
	CauseOfBlockage           string        `json:"causeOfBlockage,omitempty"`
	Causes                    []interface{} `json:"causes,omitempty"`
	ChangeSet                 []interface{} `json:"changeSet,omitempty"`
	Description               string        `json:"description,omitempty"`
	DurationInMillis          *int64        `json:"durationInMillis,omitempty"`
	EnQueueTime               Time          `json:"enQueueTime,omitempty"`
	EndTime                   Time          `json:"endTime,omitempty"`
	EstimatedDurationInMillis *int64        `json:"estimatedDurationInMillis,omitempty"`
	ID                        string        `json:"id,omitempty"`
	Name                      string        `json:"name,omitempty"`
	Organization              string        `json:"organization,omitempty"`
	Pipeline                  string        `json:"pipeline,omitempty"`
	Replayable                bool          `json:"replayable,omitempty"`
	Result                    string        `json:"result,omitempty"`
	RunSummary                string        `json:"runSummary,omitempty"`
	StartTime                 Time          `json:"startTime,omitempty"`
	State                     string        `json:"state,omitempty"`
	Type                      string        `json:"type,omitempty"`
	QueueID                   string        `json:"queueId,omitempty"`
	CommitID                  string        `json:"commitId,omitempty"`
	CommitURL                 string        `json:"commitUrl,omitempty"`
	PullRequest               interface{}   `json:"pullRequest,omitempty"`
	Branch                    interface{}   `json:"branch,omitempty"`
}

// Node represents a node detail of a PipelineRun.
type Node struct {
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

// Edge represents edge of SimplePipeline flow graph.
type Edge struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

// Input contains input step data.
type Input struct {
	ID         string                `json:"id,omitempty"`
	Message    string                `json:"message,omitempty"`
	Ok         string                `json:"ok,omitempty"`
	Parameters []ParameterDefinition `json:"parameters,omitempty"`
	Submitter  string                `json:"submitter,omitempty"`
}

// Pipeline represents a Jenkins BlueOcean Pipeline data
type Pipeline struct {
	Name         string
	Disabled     bool
	DisplayName  string
	WeatherScore int
}
