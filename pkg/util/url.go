package util

import (
	"net/url"
	"path"
	"strings"
)

// URLJoin is a util function to join host URL and API URL
func URLJoin(host, api string) (targetURL *url.URL, err error) {
	if targetURL, err = url.Parse(host); err == nil {
		pathURL, _ := url.Parse(path.Join(targetURL.Path, api))
		// handle API like "path/" correctly
		if strings.HasSuffix(api, "/") {
			pathURL.Path = pathURL.Path + "/"
		}
		targetURL = targetURL.ResolveReference(pathURL)
	}
	return
}

//URLJoinAsString  is an util function to join host URL and API URL
func URLJoinAsString(host, api string) (targetURLStr string, err error) {
	var targetURL *url.URL
	if targetURL, err = URLJoin(host, api); err == nil {
		targetURLStr = targetURL.String()
	}
	return
}
