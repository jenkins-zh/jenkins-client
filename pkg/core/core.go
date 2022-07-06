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

type Result struct {
	Status string        `json:"status"`
	Data   GenericResult `json:"data"`
}

type GenericResult interface {
	GetResult() string
	GetErrors() []string
	GetStatus() string
}

// JSONResult represents the JSON result
type JSONResult struct {
	Result string   `json:"result"`
	JSON   string   `json:"json"`
	Errors []string `json:"errors"`
}

func (r JSONResult) GetResult() string {
	return r.JSON
}

func (r JSONResult) GetErrors() []string {
	return r.Errors
}

func (r JSONResult) GetStatus() string {
	return r.Result
}

// ToJSON turns a Jenkinsfile to JSON format
// Read details from https://github.com/jenkinsci/pipeline-model-definition-plugin/blob/master/EXTENDING.md
func (q *Client) ToJSON(jenkinsfile string) (result GenericResult, err error) {
	payloadData := url.Values{"jenkinsfile": {jenkinsfile}}
	payload := strings.NewReader(payloadData.Encode())

	genericResult := &Result{
		Data: &JSONResult{},
	}
	err = q.RequestWithData(http.MethodPost, "/pipeline-model-converter/toJson", map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}, payload, http.StatusOK, genericResult)
	if err == nil {
		result = genericResult.Data
	}
	return
}

// JenkinsfileResult represents the Jenkinsfile result
type JenkinsfileResult struct {
	Result      string   `json:"result"`
	Jenkinsfile string   `json:"jenkinsfile"`
	Errors      []string `json:"errors"`
}

func (r JenkinsfileResult) GetResult() string {
	return r.Jenkinsfile
}

func (r JenkinsfileResult) GetErrors() []string {
	return r.Errors
}

func (r JenkinsfileResult) GetStatus() string {
	return r.Result
}

// ToJenkinsfile converts a JSON format data to Jenkinsfile
// Read details from https://github.com/jenkinsci/pipeline-model-definition-plugin/blob/master/EXTENDING.md
func (q *Client) ToJenkinsfile(data string) (result GenericResult, err error) {
	payloadData := url.Values{"json": {data}}
	payload := strings.NewReader(payloadData.Encode())

	genericResult := &Result{
		Data: &JenkinsfileResult{},
	}
	err = q.RequestWithData(http.MethodPost, "/pipeline-model-converter/toJenkinsfile", map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}, payload, http.StatusOK, genericResult)
	if err == nil {
		result = genericResult.Data
	}
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
