package k8s

import (
	"io/ioutil"
	"reflect"
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
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &JenkinsConfig{
				Config: tt.fields.Config,
			}
			if err := c.AddPodTemplate(tt.args.podTemplate); (err != nil) != tt.wantErr {
				t.Errorf("AddPodTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, readFileASString(tt.expectResult), c.GetConfigAsString())
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
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &JenkinsConfig{
				Config: tt.fields.Config,
			}
			if err := c.RemovePodTemplate(tt.args.podTemplate); (err != nil) != tt.wantErr {
				t.Errorf("RemovePodTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, readFileASString(tt.expectResult), c.GetConfigAsString())
		})
	}
}

func readPodTemplate(file string) (result *v1.PodTemplate) {
	result = &v1.PodTemplate{}

	data := readFile(file)
	_ = yaml.Unmarshal(data, result)
	return
}

func readFile(file string) (data []byte) {
	data, _ = ioutil.ReadFile(file)
	return
}

func readFileASString(file string) string {
	return string(readFile(file))
}

func TestJenkinsConfig_getJSON(t *testing.T) {
	type fields struct {
		Config []byte
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []byte
		wantErr  bool
	}{{
		name: "normal",
		fields: fields{Config: []byte(`name: name
server: server`)},
		wantData: []byte(`{"name":"name","server":"server"}`),
		wantErr:  false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &JenkinsConfig{
				Config: tt.fields.Config,
			}
			gotData, err := c.getJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("getJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("getJSON() gotData = %v, wantVal %v", gotData, tt.wantData)
			}
		})
	}
}
