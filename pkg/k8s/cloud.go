package k8s

import (
	"fmt"
	"strconv"
	"strings"

	unstructured "github.com/linuxsuren/unstructured/pkg"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

// JenkinsConfig represents a Jenkins configuration-as-code object
type JenkinsConfig struct {
	Config []byte
}

// GetConfigAsString returns the config data as string
func (c *JenkinsConfig) GetConfigAsString() string {
	return string(c.Config)
}

// ReplaceOrAddPodTemplate replace the existing PodTemplate, or add it if it's not exist
func (c *JenkinsConfig) ReplaceOrAddPodTemplate(podTemplate *v1.PodTemplate) (err error) {
	if err = c.RemovePodTemplate(podTemplate.Name); err == nil {
		err = c.AddPodTemplate(podTemplate)
	}
	return
}

// AddPodTemplate adds a PodTemplate to the Jenkins cloud config
func (c *JenkinsConfig) AddPodTemplate(podTemplate *v1.PodTemplate) (err error) {
	casc := map[string]interface{}{}
	if err = yaml.Unmarshal(c.Config, &casc); err != nil {
		err = fmt.Errorf("failed to unmarshal YAML to map structure, error: %v", err)
		return
	}

	var templatesObj interface{}
	var ok bool
	if templatesObj, ok, err = unstructured.NestedField(casc, "jenkins", "clouds[0]", "kubernetes", "templates"); !ok {
		err = fmt.Errorf("failed to find jenkins.cloud[0]")
		return
	} else if err != nil {
		return
	}

	var targetPodTemplate JenkinsPodTemplate
	if targetPodTemplate, err = ConvertToJenkinsPodTemplate(podTemplate); err != nil {
		return
	}

	var templates []interface{}
	if templates, ok = templatesObj.([]interface{}); ok {
		templates = append(templates, targetPodTemplate)
	}

	if err = unstructured.SetNestedField(casc, templates, "jenkins", "clouds[0]", "kubernetes", "templates"); err == nil {
		c.Config, err = yaml.Marshal(casc)
	}
	return
}

// ConvertToJenkinsPodTemplate converts a k8s style PodTemplate to a Jenkins style PodTemplate
func ConvertToJenkinsPodTemplate(podTemplate *v1.PodTemplate) (target JenkinsPodTemplate, err error) {
	target.Name = podTemplate.Name
	target.Namespace = podTemplate.Namespace
	target.Label = podTemplate.Name
	target.NodeUsageMode = "EXCLUSIVE"

	// make sure the annotations are not empty
	if podTemplate.Annotations == nil {
		podTemplate.Annotations = map[string]string{}
	}
	annotations := podTemplate.Annotations

	// set the Jenkins agent labels
	if val, ok := annotations["jenkins.agent.labels"]; ok && val != "" {
		labels := strings.Split(val, " ")
		labels = append(labels, podTemplate.Name)
		target.Label = strings.Join(labels, " ")
	}

	// set the template level fields
	target.YAML = annotations["containers.yaml"]
	target.InheritFrom = annotations["inherit.from"]
	if idleMinutes, parseErr := strconv.Atoi(annotations["idleMinutes"]); parseErr == nil {
		target.IdleMinutes = idleMinutes
	}

	// convert the containers
	containers := podTemplate.Template.Spec.Containers
	containersCount := len(containers)
	if containersCount > 0 {
		target.Containers = make([]Container, containersCount)

		for i, container := range containers {
			name := container.Name

			jenkinsAgentContainer := Container{
				Name:       name,
				Image:      container.Image,
				Command:    strings.Join(container.Command, " "),
				Args:       strings.Join(container.Args, " "),
				TtyEnabled: true,
			}

			if mem := container.Resources.Requests.Memory(); mem != nil && !mem.IsZero() {
				jenkinsAgentContainer.ResourceRequestMemory = mem.String()
			}
			if mem := container.Resources.Limits.Memory(); mem != nil && !mem.IsZero() {
				jenkinsAgentContainer.ResourceLimitMemory = mem.String()
			}
			if cpu := container.Resources.Requests.Cpu(); cpu != nil && !cpu.IsZero() {
				jenkinsAgentContainer.ResourceRequestCPU = cpu.String()
			}
			if cpu := container.Resources.Limits.Cpu(); cpu != nil && !cpu.IsZero() {
				jenkinsAgentContainer.ResourceLimitCPU = cpu.String()
			}

			target.Containers[i] = jenkinsAgentContainer
		}

		container := containers[0]
		for _, volMount := range container.VolumeMounts {
			for _, vol := range podTemplate.Template.Spec.Volumes {
				if vol.Name == volMount.Name && vol.HostPath != nil {
					target.Volumes = append(target.Volumes, Volume{
						HostPathVolume{
							HostPath:  vol.HostPath.Path,
							MountPath: volMount.MountPath,
						},
					})
					break
				}
			}
		}
	}
	return
}

// RemovePodTemplate removes a PodTemplate from the Jenkins cloud config
func (c *JenkinsConfig) RemovePodTemplate(name string) (err error) {
	casc := map[string]interface{}{}
	if err = yaml.Unmarshal(c.Config, &casc); err != nil {
		err = fmt.Errorf("failed to unmarshal YAML to map structure, error: %v", err)
		return
	}

	var templatesObj interface{}
	var ok bool
	if templatesObj, ok, err = unstructured.NestedField(casc, "jenkins", "clouds[0]", "kubernetes", "templates"); !ok {
		err = fmt.Errorf("failed to find jenkins.cloud[0]")
		return
	} else if err != nil {
		return
	}

	var templateArray []interface{}
	if templateArray, ok = templatesObj.([]interface{}); ok {
		for i, templateObj := range templateArray {
			var template map[string]interface{}
			if template, ok = templateObj.(map[string]interface{}); ok {
				if template["name"].(string) == name {
					if i == len(template)-1 {
						templateArray = templateArray[0:i]
					} else {
						templateArray = append(templateArray[0:i], templateArray[i+1:]...)
					}
					break
				}
			}
		}
	}

	if err = unstructured.SetNestedField(casc, templateArray, "jenkins", "clouds[0]", "kubernetes", "templates"); err == nil {
		c.Config, err = yaml.Marshal(casc)
	}
	return
}

// CloudAgent represents a Jenkins cloud agent
type CloudAgent struct {
	Kubernetes KubernetesCloud `json:"kubernetes"`
}

// KubernetesCloud represents a Kubernetes connection in Jenkins
type KubernetesCloud struct {
	Name                   string               `json:"name"`
	ServerURL              string               `json:"serverUrl"`
	SkipTLSVerify          bool                 `json:"skipTlsVerify"`
	Namespace              string               `json:"namespace"`
	CredentialsID          string               `json:"credentialsId"`
	JenkinsURL             string               `json:"jenkinsUrl"`
	JenkinsTunnel          string               `json:"jenkinsTunnel"`
	ContainerCapStr        string               `json:"containerCapStr"`
	ConnectTimeout         string               `json:"connectTimeout"`
	ReadTimeout            string               `json:"readTimeout"`
	MaxRequestsPerhHostStr string               `json:"maxRequestsPerhHostStr"`
	Templates              []JenkinsPodTemplate `json:"templates"`
}

// JenkinsPodTemplate represents the PodTemplate that defined in Jenkins
type JenkinsPodTemplate struct {
	Name          string      `json:"name"`
	Namespace     string      `json:"namespace"`
	Label         string      `json:"label"`
	NodeUsageMode string      `json:"nodeUsageMode"`
	IdleMinutes   int         `json:"idleMinutes"`
	Containers    []Container `json:"containers"`
	InheritFrom   string      `json:"inheritFrom,omitempty"`
	Volumes       []Volume    `json:"volumes"`
	// YAML is the YAML format for merging into the whole PodTemplate
	YAML            string          `json:"yaml"`
	WorkspaceVolume WorkspaceVolume `json:"workspaceVolume"`
}

// Volume represents the volume in kubernetes
type Volume struct {
	HostPathVolume HostPathVolume `json:"hostPathVolume"`
}

// Container represents the container that defined in Jenkins
type Container struct {
	Name                  string `json:"name"`
	Image                 string `json:"image"`
	Command               string `json:"command"`
	Args                  string `json:"args"`
	TtyEnabled            bool   `json:"ttyEnabled"`
	Privileged            bool   `json:"privileged"`
	ResourceRequestCPU    string `json:"resourceRequestCpu,omitempty"`
	ResourceLimitCPU      string `json:"resourceLimitCpu,omitempty"`
	ResourceRequestMemory string `json:"resourceRequestMemory,omitempty"`
	ResourceLimitMemory   string `json:"resourceLimitMemory,omitempty"`
}

// WorkspaceVolume is the volume of the Jenkins agent workspace
type WorkspaceVolume struct {
	EmptyDirWorkspaceVolume EmptyDirWorkspaceVolume `json:"emptyDirWorkspaceVolume"`
}

// EmptyDirWorkspaceVolume is an empty dir type of workspace volume
type EmptyDirWorkspaceVolume struct {
	Memory bool `json:"memory"`
}

// HostPathVolume represents a host path volume
type HostPathVolume struct {
	HostPath  string `json:"hostPath"`
	MountPath string `json:"mountPath"`
}
