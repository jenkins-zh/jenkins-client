package casc

import (
	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"net/http"
)

// Manager is the client of configuration as code
type Manager struct {
	core.JenkinsCore
}

// Export exports the config of configuration-as-code
func (c *Manager) Export() (config string, err error) {
	var (
		data       []byte
		statusCode int
	)

	if statusCode, data, err = c.Request(http.MethodPost, "/configuration-as-code/export",
		nil, nil); err == nil &&
		statusCode != 200 {
		err = c.ErrorHandle(statusCode, data)
	}
	config = string(data)
	return
}

// Schema get the schema of configuration-as-code
func (c *Manager) Schema() (schema string, err error) {
	var (
		data       []byte
		statusCode int
	)

	if statusCode, data, err = c.Request(http.MethodPost, "/configuration-as-code/schema",
		nil, nil); err == nil &&
		statusCode != 200 {
		err = c.ErrorHandle(statusCode, data)
	}
	schema = string(data)
	return
}

// Reload reloads the config of configuration-as-code
func (c *Manager) Reload() (err error) {
	_, err = c.RequestWithoutData(http.MethodPost, "/configuration-as-code/reload",
		nil, nil, 200)
	return
}

// Apply apply the config of configuration-as-code
func (c *Manager) Apply() (err error) {
	_, err = c.RequestWithoutData(http.MethodPost, "/configuration-as-code/apply",
		nil, nil, 200)
	return
}
