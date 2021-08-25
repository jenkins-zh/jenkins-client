[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/jenkins-zh/jenkins-client)
[![](https://goreportcard.com/badge/jenkins-zh/jenkins-client)](https://goreportcard.com/report/jenkins-zh/jenkins-client)
[![codecov](https://codecov.io/gh/jenkins-zh/jenkins-client/branch/main/graph/badge.svg?token=8N1vvFPxPm)](https://codecov.io/gh/jenkins-zh/jenkins-client)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/2b413ddf89d64212b8ddf3db23cb3c47)](https://www.codacy.com/gh/jenkins-zh/jenkins-client/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=jenkins-zh/jenkins-client&amp;utm_campaign=Badge_Grade)
[![Contributors](https://img.shields.io/github/contributors/jenkins-zh/jenkins-client.svg)](https://github.com/jenkins-zh/jenkins-client/graphs/contributors)

# jenkins-client Document

## How to get it

go.mod

```
require github.com/jenkins-zh/jenkins-client
```

## Configuration

Examples of jcli configuration

```
- name: dev
  url: http://192.168.1.10
  username: admin
  token: 11132c9ae4b20edbe56ac3e09cb5a3c8c2
  proxy: http://192.168.10.10:47586
  proxyAuth: username:password
```

So when you use jcli-client, you need to configure and read configuration files like this

```
type Config struct {
	URL      string `yaml:"url"`
	UserName string `yaml:"username"`
	Token    string `yaml:"token"`
}

func GetJenkinsCore() (core client.JenkinsCore, e error) {
	// read configuration files
	jenkinsConfigPath := "./jenkins.yml"
	yamlFile, e := ioutil.ReadFile(jenkinsConfigPath)
	if e != nil {
		return
	}
	// yaml parse
	var config Config
	e = yaml.Unmarshal(yamlFile, &config)
	if e != nil {
		return
	}
	// capsulate JenkinsCore
	core = client.JenkinsCore{
		URL:      config.URL,
		UserName: config.UserName,
		Token:    config.Token,
	}
	crumbIssuer, e := core.GetCrumb()
	if e != nil {
		return
	} else if crumbIssuer != nil {
		core.JenkinsCrumb = *crumbIssuer
	}
	return
}
```

## Job API

More used API

```
// GetBuild get build information of a job
func (q *JobClient) GetBuild(jobName string, id int) (job *JobBuild, err error){...}

// Build trigger a job
func (q *JobClient) Build(jobName string) (err error) {...}

// BuildWithParams build a job which has params
func (q *JobClient) BuildWithParams(jobName string, parameters []ParameterDefinition) (err error) {...}

// GetHistory returns the build history of a job
func (q *JobClient) GetHistory(name string) (builds []*JobBuild, err error) {...}

// Log get the log of a job
func (q *JobClient) Log(jobName string, history int, start int64) (jobLog JobLog, err error){...}
```

## User API

More used API

```
// Create will create a user in Jenkins
func (q *UserClient) Create(username, password string) (user *UserForCreate, err error) {...}

// Delete will remove a user from Jenkins
func (q *UserClient) Delete(username string) (err error) {...}
```

## Examples

```
type JobBuildOptions struct {
	Env       string 
	JobName   string 
	BranchTag string 
}


func BuildJob(jobBuild JobBuildOptions) (e error) {
	core, e := GetJenkinsCore()
	jobClient = client.JobClient{core}
	if e != nil {
		return
	}
	param1 := client.ParameterDefinition{
		Description:           "pre and prd please use tag",
		Name:                  "BRANCH_TAG",
		Type:                  "Branch or Tag",
		Value:                 jobBuild.BranchTag,
		DefaultParameterValue: client.DefaultParameterValue{Value: "origin/master"},
	}
	param2 := client.ParameterDefinition{
		Description:           "choice",
		Name:                  "ENV",
		Type:                  "Choice Parameter",
		Value:                 jobBuild.Env,
		DefaultParameterValue: client.DefaultParameterValue{Value: "qa"},
	}
	param3 := client.ParameterDefinition{
		Description:           "Deploy Type",
		Name:                  "DEPLOY_TYPE",
		Type:                  "Choice Parameter",
		Value:                 "publish",
		DefaultParameterValue: client.DefaultParameterValue{Value: "publish"},
	}
	params := []client.ParameterDefinition{param1, param2, param3}
	e = jobClient.BuildWithParams(jobBuild.JobName, params)
	return
}
```

## Plugin Dependencies

Pay attention to the part of the Jenkins CLI, they need some special plugins. Please read through the following table：

| API | Plugin |
|---|---|
| Search Job | [pipeline-restful-api](https://github.com/jenkinsci/pipeline-restful-api-plugin) |
