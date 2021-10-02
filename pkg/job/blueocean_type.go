package job

// PipelineRun represents a build detail of Pipeline.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/a7cbc946b73d89daf9dfd91cd713cc7ab64a2d95/blueocean-pipeline-api-impl/src/main/java/io/jenkins/blueocean/rest/impl/pipeline/PipelineRunImpl.java
type PipelineRun struct {
	ArtifactsZipFile          interface{}   `json:"artifactsZipFile,omitempty"`
	CauseOfBlockage           string        `json:"causeOfBlockage,omitempty"`
	Causes                    []interface{} `json:"causes,omitempty"`
	ChangeSet                 []interface{} `json:"changeSet,omitempty"`
	Description               string        `json:"description,omitempty"`
	DurationInMillis          *int64        `json:"durationInMillis,omitempty"`
	EnQueueTime               Time          `json:"enQueueTime,omitempty"`
	EndTime                   Time          `json:"endTime,omitempty"`
	EstimatedDurationInMillis *int64        `json:"estimatedDurationInMillis,omitempty"`
	ID                        string        `json:"id,omitempty"`
	Name                      string        `json:"name,omitempty"`
	Organization              string        `json:"organization,omitempty"`
	Pipeline                  string        `json:"pipeline,omitempty"`
	Replayable                bool          `json:"replayable,omitempty"`
	Result                    string        `json:"result,omitempty"`
	RunSummary                string        `json:"runSummary,omitempty"`
	StartTime                 Time          `json:"startTime,omitempty"`
	State                     string        `json:"state,omitempty"`
	Type                      string        `json:"type,omitempty"`
	QueueID                   string        `json:"queueId,omitempty"`
	CommitID                  string        `json:"commitId,omitempty"`
	CommitURL                 string        `json:"commitUrl,omitempty"`
	PullRequest               interface{}   `json:"pullRequest,omitempty"`
	Branch                    interface{}   `json:"branch,omitempty"`
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

// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueRunnableItem.java
type blueRunnableItem struct {
	WeatherScore              int                   `json:"weatherScore,omitempty"`
	LatestRun                 *PipelineRun          `json:"latestRun,omitempty"`
	EstimatedDurationInMillis int64                 `json:"estimatedDurationInMillis,omitempty"`
	Permissions               map[string]bool       `json:"permissions,omitempty"`
	Parameters                []ParameterDefinition `json:"parameters,omitempty"`
}

// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BluePipelineItem.java
type bluePipelineItem struct {
	Organization    string `json:"organization,omitempty"`
	Name            string `json:"name,omitempty"`
	Disabled        bool   `json:"disabled,omitempty"`
	DisplayName     string `json:"displayName,omitempty"`
	FullName        string `json:"fullName,omitempty"`
	FullDisplayName string `json:"fullDisplayName,omitempty"`
}

// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueContainerItem.java
type blueContainerItem struct {
	NumberOfPipelines   int      `json:"numberOfPipelines,omitempty"`
	NumberOfFolders     int      `json:"numberOfFolders,omitempty"`
	PipelineFolderNames []string `json:"pipelineFolderNames,omitempty"`
}

// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueMultiBranchItem.java
type blueMultiBranchItem struct {
	TotalNumberOfBranches          int      `json:"totalNumberOfBranches,omitempty"`
	NumberOfFailingBranches        int      `json:"numberOfFailingBranches,omitempty"`
	NumberOfSuccessfulBranches     int      `json:"numberOfSuccessfulBranches,omitempty"`
	TotalNumberOfPullRequests      int      `json:"totalNumberOfPullRequests,omitempty"`
	NumberOfFailingPullRequests    int      `json:"numberOfFailingPullRequests,omitempty"`
	NumberOfSuccessfulPullRequests int      `json:"numberOfSuccessfulPullRequests,omitempty"`
	BranchNames                    []string `json:"branchNames,omitempty"`
}

// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/6b27be3724c892427b732f30575fdcc2977cfaef/blueocean-pipeline-api-impl/src/main/java/io/jenkins/blueocean/rest/impl/pipeline/MultiBranchPipelineImpl.java
type blueMultiBranchPipeline struct {
	blueRunnableItem
	bluePipelineItem
	blueContainerItem
	blueMultiBranchItem
	SCMSource  *SCMSource `json:"scmSource,omitempty"`
	ScriptPath string     `json:"scriptPath,omitempty"`
}

// Pipeline represents a Jenkins BlueOcean Pipeline data
type Pipeline struct {
	blueMultiBranchPipeline
}

// SCMSource provides metadata about the backing SCM for a BluePipeline.
// Reference: https://github.com/jenkinsci/blueocean-plugin/blob/868c0ea4354f19e8d509deacc94325f97151aec0/blueocean-rest/src/main/java/io/jenkins/blueocean/rest/model/BlueScmSource.java
type SCMSource struct {
	ID     string `json:"id,omitempty"`
	APIUrl string `json:"apiUrl,omitempty"`
}
