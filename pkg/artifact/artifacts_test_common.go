package artifact

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"github.com/jenkins-zh/jenkins-client/pkg/job"

	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
)

// PrepareGetArtifacts only for test
func PrepareGetArtifacts(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd,
	jobName string, buildID int) (response *http.Response) {
	path := job.ParseJobPath(jobName)
	var api string
	if buildID <= 0 {
		api = fmt.Sprintf("%s/lastBuild/wfapi/artifacts", path)
	} else {
		api = fmt.Sprintf("%s/%d/wfapi/artifacts", path, buildID)
	}
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", rootURL, api), nil)
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`[{"id":"n1","name":"a.log","path":"a.log","url":"/job/pipeline/1/artifact/a.log","size":0}]`)),
	}
	roundTripper.EXPECT().
		RoundTrip(core.NewRequestMatcher(request)).Return(response, nil)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
	return
}

// PrepareGetEmptyArtifacts only for test
func PrepareGetEmptyArtifacts(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd,
	jobName string, buildID int) (response *http.Response) {
	response = PrepareGetArtifacts(roundTripper, rootURL, user, passwd, jobName, buildID)
	response.Body = ioutil.NopCloser(bytes.NewBufferString(`[]`))
	return
}

// PrepareGetArtifact only for test
func PrepareGetArtifact(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd, projectName, pipelineName string, buildId int, filename string) (response *http.Response) {
	var api = fmt.Sprintf("/job/%s/job/%s/%d/artifact/%s", projectName, pipelineName, buildId, filename)
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", rootURL, api), nil)
	response = &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`the is test file`)),
	}
	roundTripper.EXPECT().
		RoundTrip(core.NewRequestMatcher(request)).Return(response, nil)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
	return
}

// PrepareGetNoExistsArtifact only for test
func PrepareGetNoExistsArtifact(roundTripper *mhttp.MockRoundTripper, rootURL, user, passwd, projectName, pipelineName string, buildId int, filename string) (response *http.Response) {
	var api = fmt.Sprintf("/job/%s/job/%s/%d/artifact/%s", projectName, pipelineName, buildId, filename)
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", rootURL, api), nil)
	response = &http.Response{
		StatusCode: 404,
		Request:    request,
		Body:       nil,
	}
	roundTripper.EXPECT().
		RoundTrip(core.NewRequestMatcher(request)).Return(response, nil)

	if user != "" && passwd != "" {
		request.SetBasicAuth(user, passwd)
	}
	return
}
