package job

import "fmt"

// BlueItemRun contains basic metadata of a build.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/a7cbc946b73d89daf9dfd91cd713cc7ab64a2d95/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueItemRun.java
type BlueItemRun struct {
	ArtifactsZipFile          *BlueArtifact        `json:"artifactsZipFile,omitempty"`
	CauseOfBlockage           string               `json:"causeOfBlockage,omitempty"`
	Causes                    []Cause              `json:"causes,omitempty"`
	ChangeSet                 []BlueChangeSetEntry `json:"changeSet,omitempty"`
	Description               string               `json:"description,omitempty"`
	DurationInMillis          *int64               `json:"durationInMillis,omitempty"`
	EnQueueTime               Time                 `json:"enQueueTime,omitempty"`
	EndTime                   Time                 `json:"endTime,omitempty"`
	StartTime                 Time                 `json:"startTime,omitempty"`
	EstimatedDurationInMillis *int64               `json:"estimatedDurationInMillis,omitempty"`
	ID                        string               `json:"id,omitempty"`
	Name                      string               `json:"name,omitempty"`
	Organization              string               `json:"organization,omitempty"`
	Pipeline                  string               `json:"pipeline,omitempty"`
	Replayable                bool                 `json:"replayable,omitempty"`
	Result                    string               `json:"result,omitempty"`
	RunSummary                string               `json:"runSummary,omitempty"`
	State                     string               `json:"state,omitempty"`
	Type                      string               `json:"type,omitempty"`
}

// BlueArtifact contains fields that artifact owns.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueArtifact.java#L9
type BlueArtifact struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Path         string `json:"path,omitempty"`
	URL          string `json:"url,omitempty"`
	Size         int64  `json:"size,omitempty"`
	Downloadable bool   `json:"downloadable,omitempty"`
}

type BlueChangeSetEntry struct {
	CommitID      string    `json:"commitId,omitempty"`
	Author        *BlueUser `json:"author,omitempty"`
	Timestamp     Time      `json:"timestamp,omitempty"`
	Message       string    `json:"msg,omitempty"`
	AffectedPaths []string  `json:"affectedPaths,omitempty"`
	URL           string    `json:"url,omitempty"`
	CheckoutCount int       `json:"checkoutCount,omitempty"`
	Issues        []Issue   `json:"issues,omitempty"`
}

type BlueUser struct {
	ID         string              `json:"id,omitempty"`
	FullName   string              `json:"fullName,omitempty"`
	Email      string              `json:"email,omitempty"`
	Avatar     string              `json:"avatar,omitempty"`
	Permission *BlueUserPermission `json:"permission,omitempty"`
}

// Permissions are a map (string, bool) to represent permission of some operations.
type Permissions map[string]bool

// BlueUserPermission contains administrator, pipeline permission and credential permission.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueUserPermission.java#L14
type BlueUserPermission struct {
	Administrator         bool         `json:"administrator,omitempty"`
	PipelinePermissions   *Permissions `json:"pipeline,omitempty"`
	CredentialPermissions *Permissions `json:"credential,omitempty"`
}

// Cause is a map structure and contains shortDescription field at least.
// Reference:
// - https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueRun.java#L257
// - https://github.com/jenkinsci/jenkins/blob/f8b43c12fb6f07c1118d63f1435e6564723a3253/core/src/main/java/hudson/model/Cause.java#L401
type Cause map[string]interface{}

// GetShortDescription gets short description of current cause.
func (cause Cause) GetShortDescription() string {
	if shortDescription, ok := cause["shortDescription"]; ok && shortDescription != nil {
		return fmt.Sprint(shortDescription)
	}
	return ""
}

// PipelineRun represents a build detail of Pipeline.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/a7cbc946b73d89daf9dfd91cd713cc7ab64a2d95/blueocean-pipeline-api-impl/src/main/java/io/jenkins/blueocean/rest/impl/pipeline/PipelineRunImpl.java
type PipelineRun struct {
	BlueItemRun
	QueueID     string       `json:"queueId,omitempty"`
	CommitID    string       `json:"commitId,omitempty"`
	CommitURL   string       `json:"commitUrl,omitempty"`
	PullRequest *PullRequest `json:"pullRequest,omitempty"`
	Branch      *Branch      `json:"branch,omitempty"`
}

// PipelineRunSummary is summary of a PipelineRun.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest-impl/src/main/java/io/jenkins/blueocean/service/embedded/rest/AbstractBlueRunSummary.java#L18
type PipelineRunSummary struct {
	BlueItemRun
}

// PullRequest contains metadata of pull request.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/a7cbc946b73d89daf9dfd91cd713cc7ab64a2d95/blueocean-pipeline-api-impl/src/main/java/io/jenkins/blueocean/rest/impl/pipeline/BranchImpl.java#L143
type PullRequest struct {
	ID     string `json:"id,omitempty"`
	Author string `json:"author,omitempty"`
	Title  string `json:"title,omitempty"`
	URL    string `json:"url,omitempty"`
}

// Branch contains metadata of branch.
type Branch struct {
	URL       string  `json:"url,omitempty"`
	IsPrimary bool    `json:"isPrimary,omitempty"`
	Issues    []Issue `json:"issues,omitempty"`
}

// Issue holds issue ID and URL.
type Issue struct {
	ID  string `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}

// Node represents a node detail of a PipelineRun.
type Node struct {
	DisplayDescription string `json:"displayDescription,omitempty"`
	DisplayName        string `json:"displayName,omitempty"`
	DurationInMillis   int    `json:"durationInMillis,omitempty"`
	ID                 string `json:"id,omitempty"`
	Input              *Input `json:"input,omitempty"`
	Result             string `json:"result,omitempty"`
	StartTime          Time   `json:"startTime,omitempty"`
	State              string `json:"state,omitempty"`
	Type               string `json:"type,omitempty"`
	CauseOfBlockage    string `json:"causeOfBlockage,omitempty"`
	Edges              []Edge `json:"edges,omitempty"`
	FirstParent        string `json:"firstParent,omitempty"`
	Restartable        bool   `json:"restartable,omitempty"`
}

// Edge represents edge of SimplePipeline flow graph.
type Edge struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

// Input contains input step data.
type Input struct {
	ID         string                `json:"id,omitempty"`
	Message    string                `json:"message,omitempty"`
	Ok         string                `json:"ok,omitempty"`
	Parameters []ParameterDefinition `json:"parameters,omitempty"`
	Submitter  string                `json:"submitter,omitempty"`
}

// Step represents a step of Pipeline.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BluePipelineStep.java
type Step struct {
	ID                 string `json:"id,omitempty"`
	DisplayName        string `json:"displayName,omitempty"`
	DisplayDescription string `json:"displayDescription,omitempty"`
	Type               string `json:"type,omitempty"`
	Result             string `json:"result,omitempty"`
	State              string `json:"state,omitempty"`
	StartTime          Time   `json:"startTime,omitempty"`
	DurationInMillis   int64  `json:"durationInMillis,omitempty"`
	Input              *Input `json:"input,omitempty"`
}

// BlueRunnableItem contains some fields runnable PipelineRun owns only.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueRunnableItem.java
type BlueRunnableItem struct {
	WeatherScore              int                   `json:"weatherScore,omitempty"`
	LatestRun                 *PipelineRunSummary   `json:"latestRun,omitempty"`
	EstimatedDurationInMillis int64                 `json:"estimatedDurationInMillis,omitempty"`
	Permissions               *Permissions          `json:"permissions,omitempty"`
	Parameters                []ParameterDefinition `json:"parameters,omitempty"`
}

// BluePipelineItem only contains minimal fields PipelineRun owns.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BluePipelineItem.java
type BluePipelineItem struct {
	Organization    string `json:"organization,omitempty"`
	Name            string `json:"name,omitempty"`
	Disabled        bool   `json:"disabled,omitempty"`
	DisplayName     string `json:"displayName,omitempty"`
	FullName        string `json:"fullName,omitempty"`
	FullDisplayName string `json:"fullDisplayName,omitempty"`
}

// BlueContainerItem only contains folders information of a PipelineRun.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueContainerItem.java
type BlueContainerItem struct {
	NumberOfPipelines   int      `json:"numberOfPipelines,omitempty"`
	NumberOfFolders     int      `json:"numberOfFolders,omitempty"`
	PipelineFolderNames []string `json:"pipelineFolderNames,omitempty"`
}

// BlueMultiBranchItem contains some fields multi-branch PipelineRun owns only.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueMultiBranchItem.java
type BlueMultiBranchItem struct {
	TotalNumberOfBranches          int      `json:"totalNumberOfBranches,omitempty"`
	NumberOfFailingBranches        int      `json:"numberOfFailingBranches,omitempty"`
	NumberOfSuccessfulBranches     int      `json:"numberOfSuccessfulBranches,omitempty"`
	TotalNumberOfPullRequests      int      `json:"totalNumberOfPullRequests,omitempty"`
	NumberOfFailingPullRequests    int      `json:"numberOfFailingPullRequests,omitempty"`
	NumberOfSuccessfulPullRequests int      `json:"numberOfSuccessfulPullRequests,omitempty"`
	BranchNames                    []string `json:"branchNames,omitempty"`
}

// BlueMultiBranchPipeline contains all fields mult-branch PipelineRun owns.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-pipeline-api-impl/src/main/java/io/jenkins/blueocean/rest/impl/pipeline/MultiBranchPipelineImpl.java
type BlueMultiBranchPipeline struct {
	BlueRunnableItem
	BluePipelineItem
	BlueContainerItem
	BlueMultiBranchItem
	SCMSource  *SCMSource `json:"scmSource,omitempty"`
	ScriptPath string     `json:"scriptPath,omitempty"`
}

// Pipeline represents a Jenkins BlueOcean Pipeline data
type Pipeline struct {
	BlueMultiBranchPipeline
}

// SCMSource provides metadata about the backing SCM for a BluePipeline.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/868c0ea4354f19e8d509deacc94325f97151aec0/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueScmSource.java
type SCMSource struct {
	ID     string `json:"id,omitempty"`
	APIUrl string `json:"apiUrl,omitempty"`
}

// PipelineBranch is like Pipeline but contains some additional data, such as branch and pull request.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-pipeline-api-impl/src/main/java/io/jenkins/blueocean/rest/impl/pipeline/BranchImpl.java#L37
type PipelineBranch struct {
	BlueRunnableItem
	BluePipelineItem
	Branch      *Branch      `json:"branch,omitempty"`
	PullRequest *PullRequest `json:"pullRequest,omitempty"`
}
