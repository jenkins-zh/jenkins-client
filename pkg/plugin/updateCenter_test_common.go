package plugin

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-client/pkg/core"
	"net/http"
	"net/url"
	"strings"

	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
)

// PrepareForSetMirrorCertificate only for test
func PrepareForSetMirrorCertificate(roundTripper *mhttp.MockRoundTripper, rootURL, user, password string, enable bool) {
	api := "/update-center-mirror/use"
	if !enable {
		api = "/update-center-mirror/remove"
	}

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", rootURL, api), nil)
	request.Header.Add(httpdownloader.ContentType, httpdownloader.ApplicationForm)
	core.PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}

// PrepareForChangeUpdateCenterSite only for test
func PrepareForChangeUpdateCenterSite(roundTripper *mhttp.MockRoundTripper, rootURL, user, password, name, updateCenterURL string) {
	formData := url.Values{}
	formData.Add("site", updateCenterURL)
	payload := strings.NewReader(formData.Encode())

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/pluginManager/siteConfigure", rootURL), payload)
	request.Header.Add(httpdownloader.ContentType, httpdownloader.ApplicationForm)
	core.PrepareCommonPost(request, "", roundTripper, user, password, rootURL)
}
