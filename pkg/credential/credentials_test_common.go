package credential

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/jenkins-zh/jenkins-client/pkg/core"

	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"

	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
)

// PrepareForGetCredentialList only for test
func PrepareForGetCredentialList(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, store string) {
	api := fmt.Sprintf("%s/credentials/store/%s/domain/_/api/json?pretty=true&depth=1", rootURL, store)
	request, _ := http.NewRequest(http.MethodGet, api, nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(PrepareForCredentialListJSON())),
	}
	roundTripper.EXPECT().
		RoundTrip(core.NewRequestMatcher(request)).Return(response, nil)
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
}

// PrepareForDeleteCredential only for test
func PrepareForDeleteCredential(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, store, id string) {
	api := fmt.Sprintf("%s/credentials/store/%s/domain/_/credential/%s/doDelete", rootURL, store, id)
	request, _ := http.NewRequest(http.MethodPost, api, nil)
	core.PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareForDeleteCredentialInFolder only for test
func PrepareForDeleteCredentialInFolder(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, folder, id string) {
	api := fmt.Sprintf("%s/job/%s/credentials/store/folder/domain/_/credential/%s/doDelete", rootURL, folder, id)
	request, _ := http.NewRequest(http.MethodPost, api, nil)
	core.PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareForCreateCredentialInFolder only for test
func PrepareForCreateCredentialInFolder(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, folder string, payload io.Reader) {
	api := fmt.Sprintf("%s/job/%s/credentials/store/folder/domain/_/createCredentials", rootURL, folder)
	request, _ := http.NewRequest(http.MethodPost, api, payload)
	request.Header.Add(httpdownloader.ContentType, httpdownloader.ApplicationForm)
	core.PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareForUpdateCredentialInFolder only for test
func PrepareForUpdateCredentialInFolder(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, folder, id string, payload io.Reader) {
	api := fmt.Sprintf("%s/job/%s/credentials/store/folder/domain/_/credential/%s/updateSubmit", rootURL, folder, id)
	request, _ := http.NewRequest(http.MethodPost, api, payload)
	request.Header.Add(httpdownloader.ContentType, httpdownloader.ApplicationForm)
	core.PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareForGetCredentialInFolder only for test
func PrepareForGetCredentialInFolder(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, folder, id string, payload io.Reader) {
	api := fmt.Sprintf("%s/job/%s/credentials/store/folder/domain/_/credential/%s", rootURL, folder, id)
	request, _ := http.NewRequest(http.MethodGet, api, payload)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body:       ioutil.NopCloser(bytes.NewBufferString(PrepareForCredentialListJSON())),
	}
	roundTripper.EXPECT().
		RoundTrip(core.NewRequestMatcher(request)).Return(response, nil)
}

// PrepareForCreateCredential only for test
func PrepareForCreateCredential(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, store, credential string) {
	api := fmt.Sprintf("%s/credentials/store/%s/domain/_/createCredentials", rootURL, store)

	formData := url.Values{}
	formData.Add("json", fmt.Sprintf(`{"credentials": %s}`, credential))
	payload := strings.NewReader(formData.Encode())

	request, _ := http.NewRequest(http.MethodPost, api, payload)
	request.Header.Add(httpdownloader.ContentType, httpdownloader.ApplicationForm)
	core.PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareForCreateUsernamePasswordCredential only for test
func PrepareForCreateUsernamePasswordCredential(roundTripper *mhttp.MockRoundTripper, rootURL, user, password,
	store string, cred UsernamePasswordCredential) {
	cred.Class = "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"
	if payload, err := json.Marshal(cred); err == nil {
		PrepareForCreateCredential(roundTripper, rootURL, user, password, store, string(payload))
	}
}

// PrepareForCreateSecretCredential only for test
func PrepareForCreateSecretCredential(roundTripper *mhttp.MockRoundTripper, rootURL, user, password,
	store string, cred StringCredentials) {
	cred.Class = "org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl"
	cred.Scope = "GLOBAL"
	if payload, err := json.Marshal(cred); err == nil {
		PrepareForCreateCredential(roundTripper, rootURL, user, password, store, string(payload))
	}
}

// PrepareForCredentialListJSON only for test
func PrepareForCredentialListJSON() string {
	return `{
  "_class" : "com.cloudbees.plugins.credentials.CredentialsStoreAction$DomainWrapper",
  "credentials" : [
    {
      "description" : "",
      "displayName" : "displayName",
      "fingerprint" : {
      },
      "fullName" : "system/_/19c27487-acca-4a39-9889-9ddd500388f3",
      "id" : "19c27487-acca-4a39-9889-9ddd500388f3",
      "typeName" : "Username with password"
    }
  ],
  "description" : "Credentials that should be available irrespective of domain specification to requirements matching.",
  "displayName" : "全局凭据 (unrestricted)",
  "fullDisplayName" : "系统 » 全局凭据 (unrestricted)",
  "fullName" : "system/_",
  "global" : true,
  "urlName" : "_"
}`
}
