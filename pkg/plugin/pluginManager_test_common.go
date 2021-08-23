package plugin

import (
	"bytes"
	"fmt"
	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
)

// PrepareForOneInstalledPluginWithPluginName only for test
func PrepareForOneInstalledPluginWithPluginName(roundTripper *mhttp.MockRoundTripper, rootURL, pluginName string) (
	request *http.Request, response *http.Response) {
	request, response = PrepareForOneInstalledPluginWithPluginNameAndVer(roundTripper, rootURL, pluginName, "1.0")
	return
}

// PrepareForOneInstalledPluginWithPluginNameAndVer only for test
func PrepareForOneInstalledPluginWithPluginNameAndVer(roundTripper *mhttp.MockRoundTripper, rootURL,
	pluginName, version string) (
	request *http.Request, response *http.Response) {
	request, response = core.PrepareForEmptyInstalledPluginList(roundTripper, rootURL, 1)
	response.Body = ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
			"plugins": [{
				"shortName": "%s",
				"version": "%s",
				"hasUpdate": true,
				"enable": true,
				"active": true
			}]
		}`, pluginName, version)))
	return
}

// PrepareForOneInstalledPlugin only for test
func PrepareForOneInstalledPlugin(roundTripper *mhttp.MockRoundTripper, rootURL string) (
	request *http.Request, response *http.Response) {
	request, response = PrepareForOneInstalledPluginWithPluginName(
		roundTripper, rootURL, "fake")
	return
}
