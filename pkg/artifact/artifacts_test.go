package artifact

import (
	"io/ioutil"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

			projectName := "fakename"
			pipelineName := "fakename"
			filename := "fakename"
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

			projectName := "fakename"
			pipelineName := "fakename"
			filename := "fakename"
			PrepareGetNoExistsArtifact(roundTripper, artifactClient.URL, username, password, projectName, pipelineName, 1, filename)

			_, err := artifactClient.GetArtifact(projectName, pipelineName, 1, filename)
			Expect(err).Should(HaveOccurred())
		})
	})
})
