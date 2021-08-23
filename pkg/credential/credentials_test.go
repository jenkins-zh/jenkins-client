package credential_test

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-client/pkg/credential"
	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("job test", func() {
	var (
		ctrl               *gomock.Controller
		credentialsManager credential.CredentialsManager
		roundTripper       *mhttp.MockRoundTripper
		store              string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		credentialsManager = credential.CredentialsManager{}
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		credentialsManager.RoundTripper = roundTripper
		credentialsManager.URL = "http://localhost"

		store = "system"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("GetList", func() {
		It("should success", func() {
			credential.PrepareForGetCredentialList(roundTripper, credentialsManager.URL, "", "", store)

			list, err := credentialsManager.GetList(store)
			Expect(err).NotTo(HaveOccurred())
			Expect(list).NotTo(BeNil())
			Expect(len(list.Credentials)).To(Equal(1))
		})
	})

	Context("Delete", func() {
		var (
			id = "fake-id"
		)

		It("should success", func() {
			credential.PrepareForDeleteCredential(roundTripper, credentialsManager.URL, "", "", store, id)

			err := credentialsManager.Delete(store, id)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("CreateUsernamePassword", func() {
		It("should success", func() {
			cred := credential.UsernamePasswordCredential{}

			credential.PrepareForCreateUsernamePasswordCredential(roundTripper, credentialsManager.URL,
				"", "", store, cred)

			err := credentialsManager.CreateUsernamePassword(store, cred)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("CreateSecret", func() {
		It("should success", func() {
			cred := credential.StringCredentials{
				Credential: credential.Credential{Scope: "GLOBAL"},
			}

			credential.PrepareForCreateSecretCredential(roundTripper, credentialsManager.URL,
				"", "", store, cred)

			err := credentialsManager.CreateSecret(store, cred)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
