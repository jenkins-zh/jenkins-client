package artifact

import (
	"io/ioutil"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
)

var _ = Describe("artifacts test", func() {
	var (
		ctrl           *gomock.Controller
		roundTripper   *mhttp.MockRoundTripper
		artifactClient Client

		username string
		password string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		artifactClient = Client{}
		artifactClient.RoundTripper = roundTripper
		artifactClient.URL = "http://localhost"

		username = "admin"
		password = "token"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("List", func() {
		It("should success", func() {
			artifactClient.UserName = username
			artifactClient.Token = password

			jobName := "fakename"
			PrepareGetArtifacts(roundTripper, artifactClient.URL, username, password, jobName, 1)

			artifacts, err := artifactClient.List(jobName, 1)
			Expect(err).To(BeNil())
			Expect(len(artifacts)).To(Equal(1))
		})

		It("should success, with empty artifacts", func() {
			artifactClient.UserName = username
			artifactClient.Token = password

			jobName := "fakename"
			PrepareGetEmptyArtifacts(roundTripper, artifactClient.URL, username, password, jobName, 1)

			artifacts, err := artifactClient.List(jobName, 1)
			Expect(err).To(BeNil())
			Expect(len(artifacts)).To(Equal(0))
		})
	})

	Context("GetArtifactStream", func() {
		It("should success", func() {
			artifactClient.UserName = username
			artifactClient.Token = password

			projectName := "project"
			pipelineName := "pipeline"
			filename := "a.jar"
			PrepareGetArtifact(roundTripper, artifactClient.URL, username, password, projectName, pipelineName, 1, filename)

			body, err := artifactClient.GetArtifact(projectName, pipelineName, 1, filename)
			Expect(err).To(BeNil())
			Expect(func() bool {
				b, err := ioutil.ReadAll(body)
				if err != nil {
					return false
				}
				return len(b) > 0
			}()).To(Equal(true))
		})

		It("should fail", func() {
			artifactClient.UserName = username
			artifactClient.Token = password

			projectName := "project"
			pipelineName := "pipeline"
			filename := "a.jar"
			PrepareGetNoExistsArtifact(roundTripper, artifactClient.URL, username, password, projectName, pipelineName, 1, filename)

			_, err := artifactClient.GetArtifact(projectName, pipelineName, 1, filename)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("generateArtifactURL", func() {
		var (
			projectName   string
			pipelineName  string
			isMultiBranch bool
			branchName    string
			buildID       int
			filename      string
		)
		It("should success with pipeline", func() {
			projectName = "project"
			pipelineName = "pipeline"
			isMultiBranch = false
			branchName = "main"
			buildID = 1
			filename = "a.jar"
			want := "/job/project/job/pipeline/1/artifact/a.jar"
			url := generateArtifactURL(projectName, pipelineName, isMultiBranch, branchName, buildID, filename)
			Expect(url).To(Equal(want))
		})

		It("should success with multi-branch-pipeline", func() {
			projectName = "project"
			pipelineName = "pipeline"
			isMultiBranch = true
			branchName = "main"
			buildID = 1
			filename = "a.jar"
			want := "/job/project/job/pipeline/job/main/1/artifact/a.jar"
			url := generateArtifactURL(projectName, pipelineName, isMultiBranch, branchName, buildID, filename)
			Expect(url).To(Equal(want))
		})
	})

	Context("GetArtifactFromMultiBranchPipeline", func() {
		It("should success with NoScmPipelineType", func() {
			artifactClient.UserName = username
			artifactClient.Token = password

			projectName := "project"
			pipelineName := "pipeline"
			filename := "a.jar"
			branchName := "main"
			PrepareGetArtifact(roundTripper, artifactClient.URL, username, password, projectName, pipelineName, 1, filename)

			body, err := artifactClient.GetArtifactFromMultiBranchPipeline(projectName, pipelineName, false, branchName, 1, filename)
			Expect(err).To(BeNil())
			Expect(func() bool {
				b, err := ioutil.ReadAll(body)
				if err != nil {
					return false
				}
				return len(b) > 0
			}()).To(Equal(true))
		})

		It("should success with MultiBranchPipelineType", func() {
			artifactClient.UserName = username
			artifactClient.Token = password

			projectName := "project"
			pipelineName := "pipeline"
			filename := "a.jar"
			branchName := "main"
			PrepareGetMultiBranchPipelineArtifact(roundTripper, artifactClient.URL, username, password, projectName, pipelineName, 1, filename, branchName)

			body, err := artifactClient.GetArtifactFromMultiBranchPipeline(projectName, pipelineName, true, branchName, 1, filename)
			Expect(err).To(BeNil())
			Expect(func() bool {
				b, err := ioutil.ReadAll(body)
				if err != nil {
					return false
				}
				return len(b) > 0
			}()).To(Equal(true))
		})

		It("should fail with NoScmPipelineType", func() {
			artifactClient.UserName = username
			artifactClient.Token = password

			projectName := "project"
			pipelineName := "pipeline"
			filename := "a.jar"
			PrepareGetNoExistsArtifact(roundTripper, artifactClient.URL, username, password, projectName, pipelineName, 1, filename)

			_, err := artifactClient.GetArtifactFromMultiBranchPipeline(projectName, pipelineName, false, "", 1, filename)
			Expect(err).Should(HaveOccurred())
		})

		It("should fail with MultiBranchPipelineType", func() {
			artifactClient.UserName = username
			artifactClient.Token = password

			projectName := "project"
			pipelineName := "pipeline"
			filename := "a.jar"
			branchName := "main"
			PrepareGetNoExistsMultiBranchPipelineArtifact(roundTripper, artifactClient.URL, username, password, projectName, pipelineName, 1, filename, branchName)

			_, err := artifactClient.GetArtifactFromMultiBranchPipeline(projectName, pipelineName, true, branchName, 1, filename)
			Expect(err).Should(HaveOccurred())
		})
	})
})
