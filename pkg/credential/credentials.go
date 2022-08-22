package credential

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"github.com/jenkins-zh/jenkins-client/pkg/util"

	"go.uber.org/zap"
)

// CredentialsManager hold the info of credentials client
type CredentialsManager struct {
	core.JenkinsCore
}

// GetList returns the credential list
func (c *CredentialsManager) GetList(store string) (credentialList List, err error) {
	api := fmt.Sprintf("/credentials/store/%s/domain/_/api/json?pretty=true&depth=1", store)
	request := core.NewRequest(api, &c.JenkinsCore)
	if err = request.Do(); err == nil {
		err = request.GetObject(&credentialList)
	}
	return
}

// Delete deletes a credential by id from a store
func (c *CredentialsManager) Delete(store, id string) (err error) {
	api := fmt.Sprintf("/credentials/store/%s/domain/_/credential/%s/doDelete", store, id)
	request := core.NewRequest(api, &c.JenkinsCore)
	err = request.WithPostMethod().Do()
	return
}

// DeleteInFolder deletes a credential by id from a folder
func (c *CredentialsManager) DeleteInFolder(folder, id string) (err error) {
	api := fmt.Sprintf("/job/%s/credentials/store/folder/domain/_/credential/%s/doDelete", folder, id)
	request := core.NewRequest(api, &c.JenkinsCore)
	err = request.WithPostMethod().Do()
	return
}

// Create create a credential in Jenkins
func (c *CredentialsManager) Create(store, credential string) (err error) {
	api := fmt.Sprintf("/credentials/store/%s/domain/_/createCredentials", store)
	core.Logger.Debug("create credential", zap.String("api", api), zap.String("payload", credential))

	formData := url.Values{}
	formData.Add("json", fmt.Sprintf(`{"credentials": %s}`, credential))

	request := core.NewRequest(api, &c.JenkinsCore)
	request.AsPostFormRequest().WithValues(formData)
	err = request.Do()
	return
}

// CreateInFolder creates a credential in a folder
func (c *CredentialsManager) CreateInFolder(folder string, cre interface{}) (err error) {
	api := fmt.Sprintf("/job/%s/credentials/store/folder/domain/_/createCredentials", folder)

	formData := url.Values{}
	formData.Add("json", fmt.Sprintf(`{"credentials": %s}`, util.TOJSON(cre)))

	request := core.NewRequest(api, &c.JenkinsCore)
	request.AsPostFormRequest().WithValues(formData)
	err = request.Do()
	return
}

// UpdateInFolder updates a credential in a folder
func (c *CredentialsManager) UpdateInFolder(folder, id string, cre interface{}) (err error) {
	api := fmt.Sprintf("/job/%s/credentials/store/folder/domain/_/credential/%s/updateSubmit", folder, id)

	formData := url.Values{}
	formData.Add("json", fmt.Sprintf(`{"credentials": %s}`, util.TOJSON(cre)))

	request := core.NewRequest(api, &c.JenkinsCore)
	request.AsPostFormRequest().WithValues(formData).AcceptStatusCode(http.StatusNotFound)
	err = request.Do()
	return
}

// GetInFolder gets a credential in a folder
func (c *CredentialsManager) GetInFolder(folder, id string) (cre Credential, err error) {
	api := fmt.Sprintf("/job/%s/credentials/store/folder/domain/_/credential/%s", folder, id)

	request := core.NewRequest(api, &c.JenkinsCore)
	if err = request.WithValues(url.Values{"depth": {"2"}}).Do(); err == nil {
		err = request.GetObject(&cre)
	}
	return
}

// CreateUsernamePassword create username and password credential in Jenkins
func (c *CredentialsManager) CreateUsernamePassword(store string, cred UsernamePasswordCredential) (err error) {
	var payload []byte
	cred.Class = "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"
	if payload, err = json.Marshal(cred); err == nil {
		err = c.Create(store, string(payload))
	}
	return
}

// CreateSecret create token credential in Jenkins
func (c *CredentialsManager) CreateSecret(store string, cred StringCredentials) (err error) {
	var payload []byte
	cred.Class = "org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl"
	if payload, err = json.Marshal(cred); err == nil {
		err = c.Create(store, string(payload))
	}
	return
}

// Credential of Jenkins
type Credential struct {
	Description  string `json:"description"`
	DisplayName  string
	Fingerprint  interface{}
	FullName     string
	ID           string `json:"id"`
	TypeName     string
	Class        string `json:"$class"`
	StaplerClass string `json:"stapler-class"`
	Scope        string `json:"scope"`
}

// UsernamePasswordCredential hold the username and password
type UsernamePasswordCredential struct {
	Credential `json:",inline"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

// SSHCredential represents a SSH type of credential
type SSHCredential struct {
	Credential `json:",inline"`
	Username   string           `json:"username"`
	Passphrase string           `json:"passphrase"`
	KeySource  PrivateKeySource `json:"privateKeySource"`
}

// PrivateKeySource represents SSH private key
type PrivateKeySource struct {
	StaplerClass string `json:"stapler-class"`
	PrivateKey   string `json:"privateKey"`
}

// StringCredentials hold a token
type StringCredentials struct {
	Credential `json:",inline"`
	Secret     string `json:"secret"`
}

// KubeConfigCredential represents a KubeConfig credentail
type KubeConfigCredential struct {
	Credential       `json:",inline"`
	KubeconfigSource KubeconfigSource `json:"kubeconfigSource"`
}

// KubeconfigSource represents the KubeConfig content
type KubeconfigSource struct {
	StaplerClass string `json:"stapler-class"`
	Content      string `json:"content"`
}

// List contains many credentials
type List struct {
	Description     string
	DisplayName     string
	FullDisplayName string
	FullName        string
	Global          bool
	URLName         string
	Credentials     []Credential
}

const (
	// SSHCrenditalStaplerClass is the Jenkins class
	SSHCrenditalStaplerClass = "com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey"
	// DirectSSHCrenditalStaplerClass is the Jenkins class
	DirectSSHCrenditalStaplerClass = "com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey$DirectEntryPrivateKeySource"
	// UsernamePassswordCredentialStaplerClass is the Jenkins class
	UsernamePassswordCredentialStaplerClass = "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"
	// SecretTextCredentialStaplerClass is the Jenkins class
	SecretTextCredentialStaplerClass = "org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl"
	// KubeconfigCredentialStaplerClass is the Jenkins class
	KubeconfigCredentialStaplerClass = "com.microsoft.jenkins.kubernetes.credentials.KubeconfigCredentials"
	// DirectKubeconfigCredentialStaperClass is the Jenkins class
	DirectKubeconfigCredentialStaperClass = "com.microsoft.jenkins.kubernetes.credentials.KubeconfigCredentials$DirectEntryKubeconfigSource"
	// GLOBALScope is the Jenkins class
	GLOBALScope = "GLOBAL"
)
