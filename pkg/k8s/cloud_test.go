package k8s

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func TestJenkinsConfig_AddPodTemplate(t *testing.T) {
	type fields struct {
		Config []byte
	}
	type args struct {
		podTemplate *v1.PodTemplate
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectResult string
		wantErr      bool
	}{{
		name:         "normal",
		fields:       fields{Config: readFile("testdata/casc.yaml")},
		args:         args{podTemplate: readPodTemplate("testdata/k8s-podtemplate.yaml")},
		expectResult: "testdata/result-casc.yaml",
		wantErr:      false,
	}, {
		name:    "casc is not a valid YAMl",
		fields:  fields{Config: []byte(`fake`)},
		args:    args{podTemplate: readPodTemplate("testdata/k8s-podtemplate.yaml")},
		wantErr: true,
	}, {
		name:    "casc has not the expect structure",
		fields:  fields{Config: []byte(`name: rick`)},
		args:    args{podTemplate: readPodTemplate("testdata/k8s-podtemplate.yaml")},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &JenkinsConfig{
				Config: tt.fields.Config,
			}
			if err := c.AddPodTemplate(tt.args.podTemplate); (err != nil) != tt.wantErr {
				t.Errorf("AddPodTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.expectResult != "" {
				assert.Equal(t, readFileASString(tt.expectResult), c.GetConfigAsString())
			}
		})
	}
}

func TestJenkinsConfig_RemovePodTemplate(t *testing.T) {
	type fields struct {
		Config []byte
	}
	type args struct {
		podTemplate string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		expectResult string
	}{{
		name:         "normal",
		fields:       fields{Config: readFile("testdata/result-casc.yaml")},
		args:         args{podTemplate: "base"},
		expectResult: "testdata/casc.yaml",
		wantErr:      false,
	}, {
		name:    "casc is invalid",
		fields:  fields{Config: []byte("fake")},
		args:    args{podTemplate: "base"},
		wantErr: true,
	}, {
		name:    "casc has an unexpected structure",
		fields:  fields{Config: []byte(`name: rick`)},
		args:    args{podTemplate: "base"},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &JenkinsConfig{
				Config: tt.fields.Config,
			}
			if err := c.RemovePodTemplate(tt.args.podTemplate); (err != nil) != tt.wantErr {
				t.Errorf("RemovePodTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.expectResult != "" {
				assert.Equal(t, readFileASString(tt.expectResult), c.GetConfigAsString())
			}
		})
	}
}

func readPodTemplate(file string) (result *v1.PodTemplate) {
	result = &v1.PodTemplate{}

	data := readFile(file)
	_ = yaml.Unmarshal(data, result)
	return
}

func readJenkinsPodTemplate(file string) (result JenkinsPodTemplate) {
	data := readFile(file)
	_ = yaml.Unmarshal(data, &result)
	return
}

func readFile(file string) (data []byte) {
	data, _ = ioutil.ReadFile(file)
	return
}

func readFileASString(file string) string {
	return string(readFile(file))
}

func TestJenkinsConfig_ReplaceOrAddPodTemplate(t *testing.T) {
	type fields struct {
		Config []byte
	}
	type args struct {
		podTemplate *v1.PodTemplate
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectResult string
		wantErr      assert.ErrorAssertionFunc
	}{{
		name:         "normal",
		fields:       fields{Config: readFile("testdata/result-casc.yaml")},
		args:         args{podTemplate: readPodTemplate("testdata/k8s-podtemplate.yaml")},
		expectResult: "testdata/result-casc.yaml",
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.Nil(t, err)
			return true
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &JenkinsConfig{
				Config: tt.fields.Config,
			}
			tt.wantErr(t, c.ReplaceOrAddPodTemplate(tt.args.podTemplate), fmt.Sprintf("ReplaceOrAddPodTemplate(%v)", tt.args.podTemplate))
			assert.Equal(t, readFileASString(tt.expectResult), c.GetConfigAsString())
		})
	}
}

func TestConvertToJenkinsPodTemplate(t *testing.T) {
	type args struct {
		podTemplate *v1.PodTemplate
	}
	tests := []struct {
		name       string
		args       args
		wantTarget JenkinsPodTemplate
		wantErr    assert.ErrorAssertionFunc
	}{{
		name: "annotation is nil",
		args: args{
			podTemplate: readPodTemplate("testdata/k8s-podtemplate-without-annotations.yaml"),
		},
		wantTarget: readJenkinsPodTemplate("testdata/jenkins-pod-template.yaml"),
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.Nil(t, err)
			return true
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTarget, err := ConvertToJenkinsPodTemplate(tt.args.podTemplate)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvertToJenkinsPodTemplate(%v)", tt.args.podTemplate)) {
				return
			}
			assert.Equalf(t, tt.wantTarget, gotTarget, "ConvertToJenkinsPodTemplate(%v)", tt.args.podTemplate)
		})
	}
}
