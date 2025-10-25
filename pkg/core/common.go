package core

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jenkins-zh/jenkins-client/pkg/util"
	"go.uber.org/zap"
	"moul.io/http2curl"

	ext "github.com/linuxsuren/cobra-extension/version"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg/net"
)

// language is for global Accept Language
var language string

// SetLanguage set the language
func SetLanguage(lan string) {
	language = lan
}

// JenkinsCore core information of Jenkins
type JenkinsCore struct {
	JenkinsCrumb
	Timeout            time.Duration
	URL                string
	InsecureSkipVerify bool
	UserName           string
	Token              string
	Proxy              string
	ProxyAuth          string

	Debug        bool
	Output       io.Writer
	RoundTripper http.RoundTripper

	Cookies []*http.Cookie
}

// JenkinsCrumb crumb for Jenkins
type JenkinsCrumb struct {
	CrumbRequestField string
	Crumb             string
}

// GetClient get the default http Jenkins client
func (j *JenkinsCore) GetClient() (client *http.Client) {
	var roundTripper http.RoundTripper
	if j.RoundTripper != nil {
		roundTripper = j.RoundTripper
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: j.InsecureSkipVerify},
		}
		if err := httpdownloader.SetProxy(j.Proxy, j.ProxyAuth, tr); err != nil {
			log.Fatal(err)
		}
		roundTripper = tr
	}

	// make sure have a default timeout here
	if j.Timeout <= 0 {
		j.Timeout = 15
	}

	client = &http.Client{
		Transport: roundTripper,
		Timeout:   j.Timeout * time.Second,
	}
	return
}

// ProxyHandle takes care of the proxy setting
func (j *JenkinsCore) ProxyHandle(request *http.Request) {
	if j.ProxyAuth != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(j.ProxyAuth))
		Logger.Debug("setting proxy for HTTP request", zap.String("header", basicAuth))
		request.Header.Add("Proxy-Authorization", basicAuth)
	}
}

// AuthHandle takes care of the auth
func (j *JenkinsCore) AuthHandle(request *http.Request) (err error) {
	if j.UserName != "" && j.Token != "" {
		request.SetBasicAuth(j.UserName, j.Token)
	}

	// not add the User-Agent for tests
	if j.RoundTripper == nil {
		request.Header.Set("User-Agent", ext.GetCombinedVersion())
	}

	j.ProxyHandle(request)

	// all post request to Jenkins must be has the crumb
	if request.Method == http.MethodPost {
		err = j.CrumbHandle(request)
	}
	return
}

// CrumbHandle handle crum with http request
func (j *JenkinsCore) CrumbHandle(request *http.Request) error {
	if c, err := j.GetCrumb(); err == nil && c != nil {
		// cannot get the crumb could be a normal situation
		j.CrumbRequestField = c.CrumbRequestField
		j.Crumb = c.Crumb
		request.Header.Set(j.CrumbRequestField, j.Crumb)
	} else {
		return err
	}

	return nil
}

// GetCrumb get the crumb from Jenkins
func (j *JenkinsCore) GetCrumb() (crumbIssuer *JenkinsCrumb, err error) {
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = j.Request(http.MethodGet, "/crumbIssuer/api/json", nil, nil); err == nil {
		if statusCode == 200 {
			err = json.Unmarshal(data, &crumbIssuer)
		} else if statusCode == 404 {
			// return 404 if Jenkins does no have crumb
			//err = fmt.Errorf("crumb is disabled")
		} else {
			err = fmt.Errorf("unexpected status code: %d", statusCode)
		}
	}
	return
}

// RequestWithData requests the api and parse the data into an interface
func (j *JenkinsCore) RequestWithData(method, api string, headers map[string]string,
	payload io.Reader, successCode int, obj interface{}) (err error) {
	var (
		statusCode int
		data       []byte
	)

	if statusCode, data, err = j.Request(method, api, headers, payload); err == nil {
		if statusCode == successCode {
			err = json.Unmarshal(data, obj)
		} else {
			err = j.ErrorHandle(statusCode, data)
		}
	}
	return
}

// RequestWithoutData requests the api without handling data
func (j *JenkinsCore) RequestWithoutData(method, api string, headers map[string]string,
	payload io.Reader, successCode int) (statusCode int, err error) {
	var (
		data []byte
	)

	if statusCode, data, err = j.Request(method, api, headers, payload); err == nil &&
		statusCode != successCode {
		err = j.ErrorHandle(statusCode, data)
	}
	return
}

// RequestBuilder is a helper for the HTTP request
type RequestBuilder struct {
	client      *JenkinsCore
	acceptCodes []int
	method      string
	api         string
	headers     map[string]string
	payload     io.Reader

	responseCode int
	data         []byte
}

// NewRequest creates a HTTP request builder instance
func NewRequest(api string, j *JenkinsCore) *RequestBuilder {
	return &RequestBuilder{
		api:         api,
		method:      http.MethodGet,
		headers:     map[string]string{},
		acceptCodes: []int{http.StatusOK},
		client:      j,
	}
}

// AcceptStatusCode accept status code
func (r *RequestBuilder) AcceptStatusCode(code int) *RequestBuilder {
	r.acceptCodes = append(r.acceptCodes, code)
	return r
}

// RejectStatusCode reject specific status code
func (r *RequestBuilder) RejectStatusCode(code int) *RequestBuilder {
	r.acceptCodes = removeSliceItem(r.acceptCodes, code)
	return r
}

func removeSliceItem(items []int, target int) []int {
	for i := range items {
		if items[i] == target {
			count := len(items)
			if count >= 2 {
				items[i] = items[len(items)-1]
				items = items[:len(items)-1]
			} else {
				items = []int{}
			}
			break
		}
	}
	return items
}

// WithMethod sets the HTTP request method
func (r *RequestBuilder) WithMethod(method string) *RequestBuilder {
	r.method = method
	return r
}

// WithPostMethod sets the request method be POSt
func (r *RequestBuilder) WithPostMethod() *RequestBuilder {
	return r.WithMethod(http.MethodPost)
}

// WithPayload sets the payload for the request
func (r *RequestBuilder) WithPayload(payload io.Reader) *RequestBuilder {
	r.payload = payload
	return r
}

// WithValues sets the payload with values format
func (r *RequestBuilder) WithValues(values url.Values) *RequestBuilder {
	return r.WithPayload(strings.NewReader(values.Encode()))
}

// AddHeader adds a header
func (r *RequestBuilder) AddHeader(key, val string) *RequestBuilder {
	r.headers[key] = val
	return r
}

// AsFormRequest makes this request as a form request
func (r *RequestBuilder) AsFormRequest() *RequestBuilder {
	return r.AddHeader("Content-Type", "application/x-www-form-urlencoded")
}

// AsPostFormRequest make this request as a POST form request
func (r *RequestBuilder) AsPostFormRequest() *RequestBuilder {
	return r.AsFormRequest().WithMethod(http.MethodPost)
}

// GetData returns the response data
func (r *RequestBuilder) GetData() []byte {
	return r.data
}

// GetObject parses the data to an interface
func (r *RequestBuilder) GetObject(obj interface{}) error {
	return json.Unmarshal(r.GetData(), obj)
}

// Do runs the HTTP request
func (r *RequestBuilder) Do() (err error) {
	if r.responseCode, r.data, err = r.client.Request(r.method, r.api, r.headers, r.payload); err == nil {
		found := false
		for _, code := range r.acceptCodes {
			if code == r.responseCode {
				found = true
				break
			}
		}
		if !found {
			err = r.client.ErrorHandle(r.responseCode, r.data)
		}
	}
	return
}

// ErrorHandle handles the error cases
func (j *JenkinsCore) ErrorHandle(statusCode int, data []byte) (err error) {
	if statusCode >= 400 && statusCode < 500 {
		err = j.PermissionError(statusCode)
	} else {
		err = fmt.Errorf("unexpected status code: %d", statusCode)
	}

	Logger.Debug("get response", zap.String("data", string(data)))
	return
}

// PermissionError handles the no permission
func (j *JenkinsCore) PermissionError(statusCode int) (err error) {
	switch statusCode {
	case 400:
		err = fmt.Errorf("bad request, code %d", statusCode)
	case 404:
		err = fmt.Errorf("not found resources")
	default:
		err = fmt.Errorf("the current user has not permission, code %d", statusCode)
	}
	return
}

// RequestWithResponseHeader make a common request
func (j *JenkinsCore) RequestWithResponseHeader(method, api string, headers map[string]string, payload io.Reader, obj interface{}) (
	response *http.Response, err error) {
	response, err = j.RequestWithResponse(method, api, headers, payload)

	if err == nil && obj != nil && response.StatusCode == 200 {
		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			err = json.Unmarshal(data, obj)
		}
	}

	return
}

// RequestWithResponse make a common request
func (j *JenkinsCore) RequestWithResponse(method, api string, headers map[string]string, payload io.Reader) (
	response *http.Response, err error) {
	var (
		req *http.Request
	)

	if req, err = http.NewRequest(method, fmt.Sprintf("%s%s", j.URL, api), payload); err != nil {
		return
	}
	if err = j.AuthHandle(req); err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := j.GetClient()

	if curlCmd, curlErr := http2curl.GetCurlCommand(req); curlErr == nil {
		Logger.Debug("HTTP request as curl", zap.String("cmd", curlCmd.String()))
	}
	return client.Do(req)
}

// Request make a common request
func (j *JenkinsCore) Request(method, api string, headers map[string]string, payload io.Reader) (
	statusCode int, data []byte, err error) {
	var (
		req        *http.Request
		response   *http.Response
		requestURL string
	)

	if requestURL, err = util.URLJoinAsString(j.URL, api); err != nil {
		err = fmt.Errorf("cannot parse the URL of Jenkins, error is %v", err)
		return
	}

	Logger.Debug("send HTTP request", zap.String("URL", requestURL), zap.String("method", method))
	if req, err = http.NewRequest(method, requestURL, payload); err != nil {
		return
	}
	if language != "" {
		req.Header.Set("Accept-Language", language)
	}
	if err = j.AuthHandle(req); err != nil {
		return
	}

	for _, c := range j.Cookies {
		req.AddCookie(c)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if curlCmd, curlErr := http2curl.GetCurlCommand(req); curlErr == nil {
		Logger.Debug("HTTP request as curl", zap.String("cmd", curlCmd.String()))
	}

	client := j.GetClient()
	if response, err = client.Do(req); err == nil {
		if len(response.Cookies()) > 0 {
			j.Cookies = response.Cookies()
		}
		statusCode = response.StatusCode
		data, err = ioutil.ReadAll(response.Body)
	}
	return
}
