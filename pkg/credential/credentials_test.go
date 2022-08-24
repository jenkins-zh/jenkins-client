package credential_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-client/pkg/credential"
	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
	"github.com/jenkins-zh/jenkins-client/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
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
			id     = "fake-id"
			folder = "fake-folder"
		)

		It("should success", func() {
			credential.PrepareForDeleteCredential(roundTripper, credentialsManager.URL, "", "", store, id)

			err := credentialsManager.Delete(store, id)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should success", func() {
			credential.PrepareForDeleteCredentialInFolder(roundTripper, credentialsManager.URL, "", "", folder, id)

			err := credentialsManager.DeleteInFolder(folder, id)
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

func TestCreateInFolder(t *testing.T) {
	tests := []struct {
		name    string
		folder  string
		cre     interface{}
		prepare func(*testing.T, interface{}, *credential.CredentialsManager)
		wantErr bool
	}{{
		name:   "create username-password credential",
		folder: "fake",
		cre: &credential.UsernamePasswordCredential{
			Credential: credential.Credential{ID: "id"},
			Username:   "username",
			Password:   "password",
		},
		prepare: func(t *testing.T, obj interface{}, manager *credential.CredentialsManager) {
			ctrl := gomock.NewController(t)
			roundTripper := mhttp.NewMockRoundTripper(ctrl)

			formData := url.Values{}
			formData.Add("json", fmt.Sprintf(`{"credentials": %s}`, util.TOJSON(obj)))
			payload := strings.NewReader(formData.Encode())

			manager.URL = "http://localhost"
			manager.RoundTripper = roundTripper
			credential.PrepareForCreateCredentialInFolder(roundTripper, manager.URL,
				"", "", "fake", payload)
		},
		wantErr: false,
	}}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &credential.CredentialsManager{}
			if tt.prepare != nil {
				tt.prepare(t, tt.cre, client)
			}

			err := client.CreateInFolder(tt.folder, tt.cre)
			if tt.wantErr {
				assert.NotNil(t, err, "should have error in case [%s]-[%d]", tt.name, i)
			} else {
				assert.Nil(t, err, "should not have error in case [%s]-[%d]", tt.name, i)
			}
		})
	}
}

func TestUpdateInFolder(t *testing.T) {
	tests := []struct {
		name    string
		folder  string
		id      string
		cre     interface{}
		prepare func(*testing.T, interface{}, *credential.CredentialsManager)
		wantErr bool
	}{{
		name:   "update username-password credential",
		folder: "fake",
		id:     "id",
		cre: &credential.UsernamePasswordCredential{
			Credential: credential.Credential{ID: "id"},
			Username:   "username",
			Password:   "password",
		},
		prepare: func(t *testing.T, obj interface{}, manager *credential.CredentialsManager) {
			ctrl := gomock.NewController(t)
			roundTripper := mhttp.NewMockRoundTripper(ctrl)

			formData := url.Values{}
			formData.Add("json", fmt.Sprintf(`{"credentials": %s}`, util.TOJSON(obj)))
			payload := strings.NewReader(formData.Encode())

			manager.URL = "http://localhost"
			manager.RoundTripper = roundTripper
			credential.PrepareForUpdateCredentialInFolder(roundTripper, manager.URL,
				"", "", "fake", "id", payload)
		},
		wantErr: false,
	}}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &credential.CredentialsManager{}
			if tt.prepare != nil {
				tt.prepare(t, tt.cre, client)
			}

			err := client.UpdateInFolder(tt.folder, tt.id, tt.cre)
			if tt.wantErr {
				assert.NotNil(t, err, "should have error in case [%s]-[%d]", tt.name, i)
			} else {
				assert.Nil(t, err, "should not have error in case [%s]-[%d]", tt.name, i)
			}
		})
	}
}

func TestGetInFolder(t *testing.T) {
	tests := []struct {
		name    string
		folder  string
		id      string
		cre     interface{}
		prepare func(*testing.T, interface{}, *credential.CredentialsManager)
		wantErr bool
	}{{
		name:   "get username-password credential",
		folder: "fake",
		id:     "id",
		cre: &credential.UsernamePasswordCredential{
			Credential: credential.Credential{ID: "id"},
			Username:   "username",
			Password:   "password",
		},
		prepare: func(t *testing.T, obj interface{}, manager *credential.CredentialsManager) {
			ctrl := gomock.NewController(t)
			roundTripper := mhttp.NewMockRoundTripper(ctrl)

			formData := url.Values{}
			formData.Add("json", fmt.Sprintf(`{"credentials": %s}`, util.TOJSON(obj)))
			payload := strings.NewReader(formData.Encode())

			manager.URL = "http://localhost"
			manager.RoundTripper = roundTripper
			credential.PrepareForGetCredentialInFolder(roundTripper, manager.URL,
				"", "", "fake", "id", payload)
		},
		wantErr: false,
	}}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &credential.CredentialsManager{}
			if tt.prepare != nil {
				tt.prepare(t, tt.cre, client)
			}

			obj, err := client.GetInFolder(tt.folder, tt.id)
			assert.NotNil(t, obj)
			if tt.wantErr {
				assert.NotNil(t, err, "should have error in case [%s]-[%d]", tt.name, i)
			} else {
				assert.Nil(t, err, "should not have error in case [%s]-[%d]", tt.name, i)
			}
		})
	}
}
