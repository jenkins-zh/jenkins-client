package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jenkins-zh/jenkins-client/pkg/core"

	"github.com/jenkins-zh/jenkins-client/pkg/util"

	"github.com/Pallinder/go-randomdata"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
)

// Client for connect the user
type Client struct {
	core.JenkinsCore
}

// Token is the token of user
type Token struct {
	Status string    `json:"status"`
	Data   TokenData `json:"data"`
}

// TokenData represents the token
type TokenData struct {
	TokenName  string `json:"tokenName"`
	TokenUUID  string `json:"tokenUuid"`
	TokenValue string `json:"tokenValue"`
	UserName   string `json:"userName"`
}

// Get returns a user's detail
func (q *Client) Get() (status *User, err error) {
	api := fmt.Sprintf("/user/%s/api/json", q.UserName)
	err = q.RequestWithData(http.MethodGet, api, nil, nil, 200, &status)
	return
}

// EditDesc update the description of a user
func (q *Client) EditDesc(description string) (err error) {
	formData := url.Values{}
	formData.Add("description", description)
	payload := strings.NewReader(formData.Encode())
	_, err = q.RequestWithoutData(http.MethodPost, fmt.Sprintf("/user/%s/submitDescription", q.UserName),
		map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, payload, 200)
	return
}

// Delete will remove a user from Jenkins
func (q *Client) Delete(username string) (err error) {
	_, err = q.RequestWithoutData(http.MethodPost, fmt.Sprintf("/securityRealm/user/%s/doDelete", username),
		map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, nil, 200)
	return
}

func genSimpleUserAsPayload(username, password string) (payload io.Reader, user *ForCreate) {
	user = &ForCreate{
		User:      User{FullName: username},
		Username:  username,
		Password1: password,
		Password2: password,
		Email:     fmt.Sprintf("%s@%s.com", username, username),
	}

	userData, _ := json.Marshal(user)
	formData := url.Values{
		"json":      {string(userData)},
		"username":  {username},
		"password1": {password},
		"password2": {password},
		"fullname":  {username},
		"email":     {user.Email},
	}
	payload = strings.NewReader(formData.Encode())
	return
}

// Create will create a user in Jenkins
func (q *Client) Create(username, password string) (user *ForCreate, err error) {
	var (
		payload io.Reader
		code    int
	)

	if password == "" {
		password = util.GeneratePassword(8)
	}

	payload, user = genSimpleUserAsPayload(username, password)
	code, err = q.RequestWithoutData(http.MethodPost, "/securityRealm/createAccountByAdmin",
		map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, payload, 200)
	if code == 302 {
		err = nil
	}
	return
}

// CreateWithParams will create a user in Jenkins
func (q *Client) CreateWithParams(data ForCreate) (user *ForCreate, err error) {
	var (
		payload io.Reader
		code    int
	)

	if data.Username == "" {
		err = errors.Join(errors.New("error: username cannot be empty"))
	}
	if data.Password1 == "" {
		err = errors.Join(errors.New("error: password1 cannot be empty"))
	}
	if data.Password2 == "" {
		err = errors.Join(errors.New("error: password2 cannot be empty"))
	}
	if data.Email == "" {
		err = errors.Join(errors.New("error: email cannot be empty"))
	}
	if data.FullName == "" {
		err = errors.Join(errors.New("error: fullname cannot be empty"))
	}
	if err != nil {
		return nil, err
	}

	userData, _ := json.Marshal(data)
	formData := url.Values{
		"json":      {string(userData)},
		"username":  {data.Username},
		"password1": {data.Password1},
		"password2": {data.Password2},
		"fullname":  {data.FullName},
		"email":     {data.Email},
	}
	payload = strings.NewReader(formData.Encode())

	code, err = q.RequestWithoutData(http.MethodPost, "/securityRealm/createAccountByAdmin",
		map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, payload, 200)
	if code == 302 {
		err = nil
	}
	return
}

// CreateToken create a token in Jenkins
func (q *Client) CreateToken(targetUser, newTokenName string) (status *Token, err error) {
	if newTokenName == "" {
		newTokenName = fmt.Sprintf("jcli-%s", randomdata.SillyName())
	}

	if targetUser == "" {
		targetUser = q.UserName
	}

	api := fmt.Sprintf("/user/%s/descriptorByName/jenkins.security.ApiTokenProperty/generateNewToken", targetUser)

	formData := url.Values{}
	formData.Add("newTokenName", newTokenName)
	payload := strings.NewReader(formData.Encode())

	err = q.RequestWithData(http.MethodPost, api,
		map[string]string{httpdownloader.ContentType: httpdownloader.ApplicationForm}, payload, 200, &status)
	return
}

// User for Jenkins
type User struct {
	AbsoluteURL string `json:"absoluteUrl"`
	Description string
	FullName    string `json:"fullname"`
	ID          string
}

// ForCreate is the data for creating a user
type ForCreate struct {
	User      `json:",inline"`
	Username  string `json:"username"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
	Email     string `json:"email"`
}
