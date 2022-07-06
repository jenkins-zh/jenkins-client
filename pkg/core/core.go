package core

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/jenkins-zh/jenkins-client/pkg/util"
	"go.uber.org/zap"
)

var Logger *zap.Logger

// SetLogger set a global Logger
func SetLogger(zapLogger *zap.Logger) {
	Logger = zapLogger
}

func init() {
	if Logger == nil {
		var err error
		if Logger, err = util.InitLogger("warn"); err != nil {
			panic(err)
		}
	}
}

// Client hold the client of Jenkins core
type Client struct {
	JenkinsCore
}

// Restart will send the restart request
func (q *Client) Restart() (err error) {
	_, err = q.RequestWithoutData(http.MethodPost, "/safeRestart", nil, nil, 503)
	return
}

// RestartDirectly restart Jenkins directly
func (q *Client) RestartDirectly() (err error) {
	_, err = q.RequestWithoutData(http.MethodPost, "/restart", nil, nil, 503)
	return
}

// Shutdown puts Jenkins into the quiet mode, wait for existing builds to be completed, and then shut down Jenkins
func (q *Client) Shutdown(safe bool) (err error) {
	if safe {
		_, err = q.RequestWithoutData(http.MethodPost, "/safeExit", nil, nil, 200)
	} else {
		_, err = q.RequestWithoutData(http.MethodPost, "/exit", nil, nil, 200)
	}
	return
}

// JsonResult represents the JSON result
type JsonResult struct {
	Result string   `json:"result"`
	JSON   string   `json:"json"`
	Errors []string `json:"errors"`
}

// ToJson turns a Jenkinsfile to JSON format
// Read details from https://github.com/jenkinsci/pipeline-model-definition-plugin/blob/master/EXTENDING.md
func (q *Client) ToJson(jenkinsfile string) (result JsonResult, err error) {
	payloadData := url.Values{"jenkinsfile": {jenkinsfile}}
	payload := strings.NewReader(payloadData.Encode())

	err = q.RequestWithData(http.MethodPost, "/pipeline-model-converter/toJson", map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}, payload, http.StatusOK, &result)
	return
}

// JenkinsfileResult represents the Jenkinsfile result
type JenkinsfileResult struct {
	Result      string   `json:"result"`
	Jenkinsfile string   `json:"jenkinsfile"`
	Errors      []string `json:"errors"`
}

// ToJenkinsfile converts a JSON format data to Jenkinsfile
// Read details from https://github.com/jenkinsci/pipeline-model-definition-plugin/blob/master/EXTENDING.md
func (q *Client) ToJenkinsfile(data string) (result JenkinsfileResult, err error) {
	payloadData := url.Values{"json": {data}}
	payload := strings.NewReader(payloadData.Encode())

	err = q.RequestWithData(http.MethodPost, "/pipeline-model-converter/toJenkinsfile", map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}, payload, http.StatusOK, &result)
	return
}

// PrepareShutdown Put Jenkins in a Quiet mode, in preparation for a restart. In that mode Jenkins don’t start any build
func (q *Client) PrepareShutdown(cancel bool) (err error) {
	if cancel {
		_, err = q.RequestWithoutData(http.MethodPost, "/cancelQuietDown", nil, nil, 200)
	} else {
		_, err = q.RequestWithoutData(http.MethodPost, "/quietDown", nil, nil, 200)
	}
	return
}

// JenkinsIdentity belongs to a Jenkins
type JenkinsIdentity struct {
	Fingerprint   string
	PublicKey     string
	SystemMessage string
}

// GetIdentity returns the identity of a Jenkins
func (q *Client) GetIdentity() (identity JenkinsIdentity, err error) {
	err = q.RequestWithData(http.MethodGet, "/instance", nil, nil, 200, &identity)
	return
}
