package casc

import (
	"fmt"
	"net/url"

	"github.com/jenkins-zh/jenkins-client/pkg/core"
)

// Manager is the client of configuration as code
type Manager struct {
	core.JenkinsCore
}

// Export exports the config of configuration-as-code
func (c *Manager) Export() (config string, err error) {
	request := core.NewRequest("/configuration-as-code/export", &c.JenkinsCore)
	request.WithPostMethod()
	if err = request.Do(); err == nil {
		config = string(request.GetData())
	}
	return
}

// Schema get the schema of configuration-as-code
func (c *Manager) Schema() (schema string, err error) {
	request := core.NewRequest("/configuration-as-code/schema", &c.JenkinsCore)
	request.WithPostMethod()
	if err = request.Do(); err == nil {
		schema = string(request.GetData())
	}
	return
}

// Reload reloads the config of configuration-as-code
func (c *Manager) Reload() (err error) {
	request := core.NewRequest("/configuration-as-code/reload", &c.JenkinsCore)
	err = request.WithPostMethod().Do()
	return
}

// Replace replaces the new source
func (c *Manager) Replace(source string) (err error) {
	formValue := make(url.Values)
	formValue.Set("json", fmt.Sprintf(`{"newSource": "%s"}`, source))
	formValue.Set("_.newSource", source)

	// Jenkins does not have a standard API. This is a form submit, so the expected code is not 200
	request := core.NewRequest("/configuration-as-code/replace", &c.JenkinsCore)
	request.WithPostMethod().AsFormRequest().WithValues(formValue).AcceptStatusCode(302)
	err = request.Do()
	if urlErr, ok := err.(*url.Error); ok && urlErr.Err.Error() == "302 response missing Location header" {
		err = nil
	}
	return
}

// CheckNewSource checks the new source of CasC
func (c *Manager) CheckNewSource(source string) (err error) {
	formValue := make(url.Values)
	formValue.Set("newSource", source)

	request := core.NewRequest("/configuration-as-code/checkNewSource", &c.JenkinsCore)
	request.WithPostMethod().AsFormRequest().WithValues(formValue)
	err = request.Do()
	return
}

// Apply applies the config of configuration-as-code
func (c *Manager) Apply() (err error) {
	request := core.NewRequest("/configuration-as-code/apply", &c.JenkinsCore)
	request.WithPostMethod()
	err = request.Do()
	return
}
