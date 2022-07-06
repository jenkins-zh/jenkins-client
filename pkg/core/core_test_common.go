package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
)

// PrepareRestart only for test
func PrepareRestart(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string, statusCode int) {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/safeRestart", rootURL), nil)
	response := PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	response.StatusCode = statusCode
}

// PrepareRestartDirectly only for test
func PrepareRestartDirectly(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string, statusCode int) {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/restart", rootURL), nil)
	response := PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
	response.StatusCode = statusCode
}

// PrepareForShutdown only for test
func PrepareForShutdown(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string, safe bool) {
	var request *http.Request
	if safe {
		request, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/safeExit", rootURL), nil)
	} else {
		request, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/exit", rootURL), nil)
	}
	PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareForCancelShutdown only for test
func PrepareForCancelShutdown(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string, cancel bool) {
	var request *http.Request
	if cancel {
		request, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/cancelQuietDown", rootURL), nil)
	} else {
		request, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/quietDown", rootURL), nil)
	}
	PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareForToJSON only for test
func PrepareForToJSON(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string) {
	payloadData := url.Values{"jenkinsfile": {"jenkinsfile"}}
	payload := strings.NewReader(payloadData.Encode())

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/pipeline-model-converter/toJson", rootURL), payload)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	PrepareCommonPost(request, `{"result":"success","json":"json"}`, roundTripper, user, password, rootURL)
}

// PrepareForToJenkinsfile only for test
func PrepareForToJenkinsfile(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string) {
	payloadData := url.Values{"json": {"json"}}
	payload := strings.NewReader(payloadData.Encode())

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/pipeline-model-converter/toJenkinsfile", rootURL), payload)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	PrepareCommonPost(request, `{"result":"success","jenkinsfile":"jenkinsfile"}`, roundTripper, user, password, rootURL)
}

// PrepareForGetIdentity only for test
func PrepareForGetIdentity(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string) {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/instance", rootURL), nil)
	response := &http.Response{
		StatusCode: 200,
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
{"fingerprint":"fingerprint","publicKey":"publicKey","systemMessage":"systemMessage"}`)),
	}
	roundTripper.EXPECT().
		RoundTrip(NewRequestMatcher(request)).Return(response, nil)

	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
}
