package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("common test", func() {
	var (
		ctrl         *gomock.Controller
		jenkinsCore  JenkinsCore
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		jenkinsCore = JenkinsCore{}
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jenkinsCore.RoundTripper = roundTripper
		jenkinsCore.URL = "http://localhost"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Request", func() {
		var (
			method  string
			api     string
			headers map[string]string
			payload io.Reader
		)

		BeforeEach(func() {
			method = http.MethodGet
			api = "/fake"
		})

		It("normal case for get request", func() {
			request, _ := http.NewRequest(method, fmt.Sprintf("%s%s", jenkinsCore.URL, api), payload)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Header:     http.Header{},
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(request)).Return(response, nil)

			statusCode, data, err := jenkinsCore.Request(method, api, headers, payload)
			Expect(err).To(BeNil())
			Expect(statusCode).To(Equal(200))
			Expect(string(data)).To(Equal(""))
		})

		It("normal case for post request", func() {
			method = http.MethodPost
			request, _ := http.NewRequest(method, fmt.Sprintf("%s%s", jenkinsCore.URL, api), payload)
			request.Header.Add("CrumbRequestField", "Crumb")
			request.Header.Add("Fake", "fake")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(request)).Return(response, nil)

			requestCrumb, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", jenkinsCore.URL, "/crumbIssuer/api/json"), payload)
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(requestCrumb)).Return(responseCrumb, nil)

			headers = make(map[string]string, 1)
			headers["fake"] = "fake"
			statusCode, data, err := jenkinsCore.Request(method, api, headers, payload)
			Expect(err).To(BeNil())
			Expect(statusCode).To(Equal(200))
			Expect(string(data)).To(Equal(""))
		})
	})

	Context("GetCrumb", func() {
		It("without crumb setting", func() {
			requestCrumb, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", jenkinsCore.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 404,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(requestCrumb)).Return(responseCrumb, nil)

			_, err := jenkinsCore.GetCrumb()
			Expect(err).NotTo(HaveOccurred())
		})

		It("with crumb setting", func() {
			PrepareForGetIssuer(roundTripper, jenkinsCore.URL, "", "")

			crumb, err := jenkinsCore.GetCrumb()
			Expect(err).To(BeNil())
			Expect(crumb).NotTo(BeNil())
			Expect(crumb.CrumbRequestField).To(Equal("CrumbRequestField"))
			Expect(crumb.Crumb).To(Equal("Crumb"))
		})

		It("with error from server", func() {
			//requestCrumb, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", jenkinsCore.URL, "/crumbIssuer/api/json"), nil)
			//responseCrumb := &http.Response{
			//	StatusCode: 500,
			//	Proto:      "HTTP/1.1",
			//	Request:    requestCrumb,
			//	Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			//}
			//roundTripper.EXPECT().
			//	RoundTrip(NewRequestMatcher(requestCrumb)).Return(responseCrumb, nil)
			PrepareForGetIssuerWith500(roundTripper, jenkinsCore.URL, "", "")

			_, err := jenkinsCore.GetCrumb()
			Expect(err).To(HaveOccurred())
		})

		It("with Language", func() {
			request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", jenkinsCore.URL, "/view/all/itemCategories?depth=3"), nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`number name                       type
0      构建一个自由风格的软件项目 Standalone Projects
1      构建一个maven项目          Standalone Projects
2      流水线                     Standalone Projects
3      构建一个多配置项目         Standalone Projects
0      Bitbucket Team/Project     Nested Projects
1      文件夹                     Nested Projects
2      GitHub 组织                Nested Projects
3      多分支流水线               Nested Projects
`)),
			}
			request.Header.Set("Accept-Language", "zh-CN")
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(request)).Return(response, nil)

			SetLanguage("zh-CN")
			statusCode, data, err := jenkinsCore.Request(http.MethodGet, "/view/all/itemCategories?depth=3", nil, nil)
			SetLanguage("")
			Expect(err).To(BeNil())
			Expect(statusCode).To(Equal(200))
			Expect(string(data)).To(Equal(`number name                       type
0      构建一个自由风格的软件项目 Standalone Projects
1      构建一个maven项目          Standalone Projects
2      流水线                     Standalone Projects
3      构建一个多配置项目         Standalone Projects
0      Bitbucket Team/Project     Nested Projects
1      文件夹                     Nested Projects
2      GitHub 组织                Nested Projects
3      多分支流水线               Nested Projects
`))
		})

		It("with 404 error from server", func() {
			err := jenkinsCore.ErrorHandle(404, []byte{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("not found resources"))
		})

		It("with 403 error from server", func() {
			err := jenkinsCore.ErrorHandle(403, []byte{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("the current user has not permission, code 403"))
		})

		It("with CrumbHandle error from server", func() {
			requestCrumb, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", jenkinsCore.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 500,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(requestCrumb)).Return(responseCrumb, nil)
			err := jenkinsCore.CrumbHandle(requestCrumb)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unexpected status code: 500"))
		})

		It("handle a request contains crumb in it", func() {
			requestCrumb, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", jenkinsCore.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
 "crumb": "3c4525418803ddec7003a6a03995ba94dc151c9686825e139a032b3142249942",
 "crumbRequestField": "Jenkins-Crumb"
}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(NewRequestMatcher(requestCrumb)).Return(responseCrumb, nil)

			fakePostRequest, _ := http.NewRequest(http.MethodGet, "fake", nil)
			fakePostRequest.Header.Add("Jenkins-Crumb", "fake")

			err := jenkinsCore.CrumbHandle(fakePostRequest)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(fakePostRequest.Header.Values("Jenkins-Crumb"))).To(Equal(1))
			Expect(fakePostRequest.Header.Values("Jenkins-Crumb")[0]).To(Equal("3c4525418803ddec7003a6a03995ba94dc151c9686825e139a032b3142249942"))
		})

		It("test GetClient", func() {
			jenkinsCore.RoundTripper = nil
			jenkinsCore.Proxy = "kljasdsll"
			jenkinsCore.ProxyAuth = "kljaslkdjkslad"
			jclient := jenkinsCore.GetClient()
			Expect(jclient).NotTo(BeNil())
		})
	})
})

func TestRemoveSliceItem(t *testing.T) {
	tests := []struct {
		name   string
		items  []int
		target int
		expect []int
	}{{
		name:   "empty slice",
		items:  []int{},
		target: 12,
		expect: []int{},
	}, {
		name:   "one item, not found the target",
		items:  []int{12},
		target: 13,
		expect: []int{12},
	}, {
		name:   "one item, match with the target",
		items:  []int{12},
		target: 12,
		expect: []int{},
	}, {
		name:   "two items, not found the target",
		items:  []int{12, 13},
		target: 14,
		expect: []int{12, 13},
	}, {
		name:   "two items, match with the target",
		items:  []int{12, 13},
		target: 12,
		expect: []int{13},
	}, {
		name:   "more items, not found the target",
		items:  []int{12, 13, 14, 15},
		target: 10,
		expect: []int{12, 13, 14, 15},
	}, {
		name:   "more items, match with the target",
		items:  []int{12, 13, 14, 15},
		target: 14,
		expect: []int{12, 13, 15},
	}}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeSliceItem(tt.items, tt.target)
			assert.ElementsMatch(t, tt.expect, result, "failed in case [%s]-[%d]", tt.name, i)
		})
	}
}
