package artifact

import (
	"fmt"
	"io"
	"net/http"

	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"github.com/jenkins-zh/jenkins-client/pkg/job"
)

// Artifact represents the artifacts from Jenkins build
type Artifact struct {
	ID   string
	Name string
	Path string
	URL  string
	Size int64
}

// Client is client for getting the artifacts
type Client struct {
	core.JenkinsCore
}

// List get the list of artifacts from a build
func (q *Client) List(jobName string, buildID int) (artifacts []Artifact, err error) {
	path := job.ParseJobPath(jobName)
	var api string
	if buildID < 1 {
		api = fmt.Sprintf("%s/lastBuild/wfapi/artifacts", path)
	} else {
		api = fmt.Sprintf("%s/%d/wfapi/artifacts", path, buildID)
	}
	err = q.RequestWithData(http.MethodGet, api, nil, nil, 200, &artifacts)
	return
}

// GetArtifact download artifact using stream
func (q *Client) GetArtifact(projectName, pipelineName string, buildId int, filename string) (io.ReadCloser, error) {
	artifactUrl := fmt.Sprintf("/job/%s/job/%s/%d/artifact/%s", projectName, pipelineName, buildId, filename)
	resp, err := q.RequestWithResponse(http.MethodGet, artifactUrl, nil, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get artifact. the HTTP status code is %d", resp.StatusCode)
	}

	return resp.Body, nil
}
