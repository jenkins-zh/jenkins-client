package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipeline test via BlueOcean RESTful API", func() {
	var (
		ctrl         *gomock.Controller
		c            BlueOceanClient
		roundTripper *mhttp.MockRoundTripper
	)

	const organization = "jenkins"

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		c = BlueOceanClient{}
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		c.RoundTripper = roundTripper
		c.URL = "http://localhost"
		c.Organization = organization
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Search", func() {
		given := func(name string, givenResponseJSON string) {
			api := fmt.Sprintf("/blue/rest/search/?q=pipeline:*%s*;type:pipeline;organization:%s;excludedFromFlattening=jenkins.branch.MultiBranchProject,com.cloudbees.hudson.plugins.folder.AbstractFolder&filter=no-folders&start=%d&limit=%d",
				name, organization, 0, 50)
			request, _ := http.NewRequest(http.MethodGet, api, nil)
			response := &http.Response{
				StatusCode: 200,
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(givenResponseJSON)),
			}
			roundTripper.EXPECT().
				RoundTrip(core.NewRequestMatcher(request)).Return(response, nil)
		}

		It("Basic case with one result item", func() {
			name := "fake"
			given(name, `
[
  {
    "name": "fake",
    "displayName": "fake",
    "description": null,
    "type": "WorkflowJob",
    "shortURL": "job/fake/",
    "url": "job/fake/"
  }
]`)

			result, err := c.Search(name, 0, 50)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(len(result)).To(Equal(1))
			Expect(result[0].Name).To(Equal("fake"))
		})

		It("Basic case without any result items", func() {
			name := "fake"
			given(name, `[]`)

			result, err := c.Search(name, 0, 50)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(len(result)).To(Equal(0))
		})
	})

	Context("Build", func() {
		given := func(pipelineName string,
			branch string,
			statusCode int,
			requestCustomizer func(request *http.Request),
			responseCustomizer func(response *http.Response)) {
			api := c.URL + "/blue/rest/organizations/" + organization + "/pipelines/" + pipelineName
			if branch != "" {
				api = api + "/branches/" + branch
			}
			api = api + "/runs/"

			request, _ := http.NewRequest(http.MethodPost, api, nil)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Add("CrumbRequestField", "Crumb")

			if requestCustomizer != nil {
				requestCustomizer(request)
			}

			response := &http.Response{
				StatusCode: statusCode,
				Proto:      "HTTP/1.1",
				Request:    request,
			}
			if responseCustomizer != nil {
				responseCustomizer(response)
			}

			roundTripper.EXPECT().RoundTrip(core.NewRequestMatcher(request).WithQuery().WithBody()).Return(response, nil)

			requestCrumb, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", c.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}`)),
			}
			roundTripper.EXPECT().RoundTrip(core.NewRequestMatcher(requestCrumb).WithQuery().WithBody()).Return(responseCrumb, nil)
		}

		It("Trigger a simple Pipeline via Blue Ocean REST API", func() {
			const pipelineName = "fakePipeline"
			given(pipelineName, "", 200, nil, func(response *http.Response) {
				response.Body = io.NopCloser(strings.NewReader(`
{
  "expectedBuildNumber" : 1,
  "id" : "3",
  "enQueueTime": null
}`))
			})
			pipelineBuild, err := c.Build(BuildOption{
				Pipelines: []string{pipelineName},
			})
			Expect(err).To(BeNil())
			Expect(pipelineBuild).NotTo(BeNil())
			Expect(pipelineBuild.EnQueueTime.IsZero()).Should(BeTrue())
			Expect(pipelineBuild.ID).Should(Equal("3"))
		})

		It("Trigger a simple Pipeline with parameters", func() {
			const pipelineName = "fakePipeline"
			var parameters = []Parameter{{
				Name:  "this_is_a_name",
				Value: "this_is_a_value",
			}}
			paramBytes, err := json.Marshal(map[string][]Parameter{
				"parameters": parameters,
			})
			Expect(err).To(BeNil())

			given(pipelineName, "", 200, func(request *http.Request) {
				request.Body = io.NopCloser(strings.NewReader(string(paramBytes)))
			}, func(response *http.Response) {
				response.Body = io.NopCloser(strings.NewReader(`
{
 "expectedBuildNumber" : 1,
 "id" : "3",
 "enQueueTime": null
}`))
			})

			pipelineBuild, err := c.Build(BuildOption{
				Pipelines:  []string{pipelineName},
				Parameters: parameters,
			})

			Expect(err).To(BeNil())
			Expect(pipelineBuild).NotTo(BeNil())
			Expect(pipelineBuild.EnQueueTime.IsZero()).Should(BeTrue())
			Expect(pipelineBuild.ID).Should(Equal("3"))
		})

		It("Trigger a simple Pipeline via Blue Ocean REST API with an error", func() {
			const pipelineName = "fakePipeline"
			given(pipelineName, "", 400, nil, func(response *http.Response) {
				response.Body = io.NopCloser(strings.NewReader(`
{
    "message": "parameters.name is required element",
    "code": 400,
    "errors": []
}`))
			})

			_, err := c.Build(BuildOption{
				Pipelines: []string{pipelineName},
			})
			Expect(err).NotTo(BeNil())
		})

		It("Trigger a multi branch Pipeline", func() {
			const (
				pipelineName = "fakePipeline"
				branch       = "feature-1"
			)
			given(pipelineName, "feature-1", 200, nil, func(response *http.Response) {
				response.Body = io.NopCloser(strings.NewReader(`
{
  "expectedBuildNumber" : 1,
  "id" : "3",
  "enQueueTime": null
}`))
			})

			pb, err := c.Build(BuildOption{
				Pipelines: []string{pipelineName},
				Branch:    branch,
			})
			Expect(err).To(BeNil())
			Expect(pb).NotTo(BeNil())
			Expect(pb.ID).Should(Equal("3"))
		})
	})

	Context("GetBuild", func() {
		given := func(pipelineName, runID string, branch string, statusCode int, givenResponseJSON string) {
			api := fmt.Sprintf("%s/blue/rest/organizations/%s/pipelines/%s", c.URL, organization, pipelineName)
			if branch != "" {
				api = fmt.Sprintf("%s/branches/%s", api, branch)
			}
			api = fmt.Sprintf("%s/runs/%s/", api, runID)
			request, _ := http.NewRequest(http.MethodGet, api, nil)
			request.Header.Set("Content-Type", "application/json")
			response := &http.Response{
				StatusCode: statusCode,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       io.NopCloser(bytes.NewBufferString(givenResponseJSON)),
			}
			roundTripper.EXPECT().RoundTrip(core.NewRequestMatcher(request)).Return(response, nil)
		}
		It("Get specific Pipeline run", func() {
			const (
				pipelineName = "fakePipeline"
				runID        = "1"
			)
			given(pipelineName, runID, "", 200, `
{
  "enQueueTime": "2021-08-25T07:29:13.483+0000",
  "endTime": null,
  "estimatedDurationInMillis": 35661,
  "id": "5",
  "result": "UNKNOWN",
  "startTime": "2021-08-25T07:29:13.499+0000",
  "state": "RUNNING"
}`)
			pipelineBuild, err := c.GetBuild(GetBuildOption{
				RunID:     runID,
				Pipelines: []string{pipelineName},
			})
			Expect(err).Should(BeNil())
			Expect(pipelineBuild).ShouldNot(BeNil())
		})
		It("Get specific Pipeline run with an error", func() {
			const (
				pipelineName = "fakePipeline"
				runID        = "1"
			)
			given(pipelineName, runID, "", 500, `
{
  "message" : "Failed to create Git pipeline: demo",
  "code" : 400,
  "errors" : [ {
    "message" : "demo already exists",
    "code" : "ALREADY_EXISTS",
    "field" : "name"
  } ]
}`)
			_, err := c.GetBuild(GetBuildOption{
				RunID:     runID,
				Pipelines: []string{pipelineName},
			})
			Expect(err).ShouldNot(BeNil())
		})
		It("Get multi branch Pipeline run", func() {
			const (
				pipelineName = "fakePipeline"
				runID        = "1"
				branch       = "main"
			)
			given(pipelineName, runID, branch, 200, `
{
  "enQueueTime": "2021-08-25T07:29:13.483+0000",
  "endTime": null,
  "estimatedDurationInMillis": 35661,
  "id": "5",
  "result": "UNKNOWN",
  "startTime": "2021-08-25T07:29:13.499+0000",
  "state": "RUNNING"
}`)
			pipelineBuild, err := c.GetBuild(GetBuildOption{
				RunID:     runID,
				Pipelines: []string{pipelineName},
				Branch:    branch,
			})
			Expect(err).Should(BeNil())
			Expect(pipelineBuild).ShouldNot(BeNil())
		})

	})
})
