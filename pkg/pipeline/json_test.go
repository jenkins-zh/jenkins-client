package pipeline

import (
	"reflect"
	"testing"

	"github.com/jenkins-zh/jenkins-client/pkg/util"
)

func TestFindGit(t *testing.T) {
	type args struct {
		jenkinsfile string
	}
	tests := []struct {
		name      string
		args      args
		wantRepos []GitRepo
	}{{
		name: "jenkinsfile contains a git step",
		args: args{jenkinsfile: util.ReadFileASString("testdata/jenkinsfile.json")},
		wantRepos: []GitRepo{{
			URL:    "https://github.com/kubesphere/ks-devops/",
			Branch: "master",
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepos, err := FindGit(tt.args.jenkinsfile)
			if !reflect.DeepEqual(gotRepos, tt.wantRepos) {
				t.Errorf("FindGit() = %v, want %v", gotRepos, tt.wantRepos)
			}
			if err != nil {
				t.Errorf("shoud not have error, but got %v", err)
			}
		})
	}
}

func TestGitRepos_GetURLs(t *testing.T) {
	tests := []struct {
		name string
		g    GitRepos
		want string
	}{{
		name: "normal",
		g: []GitRepo{{
			URL: "http://git.com",
		}, {
			URL: "http://fake.com",
		}},
		want: "http://git.com,http://fake.com",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.GetURLs(); got != tt.want {
				t.Errorf("GetURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRepos_GetBranchesAsJSONString(t *testing.T) {
	tests := []struct {
		name string
		g    GitRepos
		want string
	}{{
		name: "normal",
		g: []GitRepo{{
			Branch: "master",
		}, {
			Branch: "feat-x",
		}},
		want: `["master","feat-x"]`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.GetBranchesAsJSONString(); got != tt.want {
				t.Errorf("GetBranchesAsJSONString() = %v, want %v", got, tt.want)
			}
		})
	}
}
