package job

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

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

func (c *BlueOceanClient) getPipelineAPI(folders ...string) (api string) {
	api = fmt.Sprintf("%s/%s/pipelines", organizationAPIPrefix, c.Organization)
	for _, folder := range folders {
		api = fmt.Sprintf("%s/%s/pipelines/", api, folder)
	}
	return
}

// GetPipeline obtains Pipeline metadata with Pipeline name and folders.
func (c *BlueOceanClient) GetPipeline(pipelineName string, folders ...string) (*Pipeline, error) {
	api := c.getGetPipelineAPI(pipelineName, folders...)
	pipeline := &Pipeline{}
	if err := c.RequestWithData(http.MethodGet, api, nil, nil, 200, pipeline); err != nil {
		return nil, err
	}
	return pipeline, nil
}

func (c *BlueOceanClient) getGetPipelineAPI(pipelineName string, folders ...string) string {
	api := fmt.Sprintf("%s/%s", organizationAPIPrefix, c.Organization)
	folders = append(folders, pipelineName)
	for _, folder := range folders {
		api = fmt.Sprintf("%s/pipelines/%s", api, folder)
	}
	return api
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
	if err := c.RequestWithData(http.MethodPost, c.getReplayAPI(&option), getHeaders(), nil, 200, pipelineRun); err != nil {
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

// GetStepsOption holds options for getting steps data.
type GetStepsOption struct {
	Folders      []string
	PipelineName string
	Branch       string
	RunID        string
	NodeID       string
}

// GetSteps returns all steps of the given Pipeline.
// Reference: https://github.com/jenkinsci/blueocean-plugin/tree/master/blueocean-rest#get-pipeline-steps
func (c *BlueOceanClient) GetSteps(option GetStepsOption) ([]Step, error) {
	api := c.getGetStepsAPI(&option)
	steps := make([]Step, 0)
	if err := c.RequestWithData(http.MethodGet, api, nil, nil, 200, &steps); err != nil {
		return nil, err
	}
	return steps, nil
}

func (c *BlueOceanClient) getGetStepsAPI(option *GetStepsOption) string {
	// api := c.getGetPipelineAPI(pipelineName, folders...)
	api := c.getGetPipelineAPI(option.PipelineName, option.Folders...)
	if option.Branch != "" {
		api = api + "/branches/" + option.Branch
	}
	api = api + "/runs/" + option.RunID
	if option.NodeID != "" {
		api = api + "/nodes/" + option.NodeID
	}
	api = api + "/steps/"
	return api
}

// Filter is Pipeline job filter.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/a7cbc946b73d89daf9dfd91cd713cc7ab64a2d95/blueocean-pipeline-api-impl/src/main/java/io/jenkins/blueocean/rest/impl/pipeline/PipelineJobFilters.java
type Filter string

const (
	// FolderJobFilter will filter out non folder Pipelines.
	FolderJobFilter Filter = "no-folders"
	// OriginFilter will filter out branches that are not pull requests.
	OriginFilter Filter = "origin"
	// PullRequestFilter will filter out branches that are pull requests.
	PullRequestFilter Filter = "pull-requests"
)

// GetBranchesOption contains some options for getting Pipeline branches.
type GetBranchesOption struct {
	Folders      []string
	PipelineName string
	Filter       Filter
	Start        int
	Limit        int
}

// GetBranches gets branches of a Pipeline.
func (c *BlueOceanClient) GetBranches(option GetBranchesOption) ([]PipelineBranch, error) {
	api := c.getGetBranchesAPI(&option)
	branches := []PipelineBranch{}
	if err := c.RequestWithData(http.MethodGet, api, nil, nil, http.StatusOK, &branches); err != nil {
		return nil, err
	}
	return branches, nil
}

func (c *BlueOceanClient) getGetBranchesAPI(option *GetBranchesOption) string {
	api := c.getGetPipelineAPI(option.PipelineName, option.Folders...)
	api = api + "/branches/"
	apiURL := &url.URL{
		Path: api,
	}
	query := apiURL.Query()
	if option.Filter != "" {
		query.Add("filter", string(option.Filter))
	}
	if option.Start > 0 {
		query.Add("start", strconv.Itoa(option.Start))
	}
	if option.Limit > 0 {
		query.Add("limit", strconv.Itoa(option.Limit))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}
