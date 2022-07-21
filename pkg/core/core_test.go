package core

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("core test", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		coreClient   Client

		username string
		password string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		coreClient = Client{}
		coreClient.RoundTripper = roundTripper
		coreClient.URL = "http://localhost"

		username = "admin"
		password = "token"

		coreClient.UserName = username
		coreClient.Token = password
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Get data", func() {
		It("should success", func() {
			PrepareRestart(roundTripper, coreClient.URL, username, password, 503)

			err := coreClient.Restart()
			Expect(err).To(BeNil())
		})

		It("should error, 400", func() {
			PrepareRestart(roundTripper, coreClient.URL, username, password, 400)

			err := coreClient.Restart()
			Expect(err).To(HaveOccurred())
		})

		It("should success", func() {
			PrepareRestartDirectly(roundTripper, coreClient.URL, username, password, 503)

			err := coreClient.RestartDirectly()
			Expect(err).To(BeNil())
		})

		It("GetIdentity", func() {
			PrepareForGetIdentity(roundTripper, coreClient.URL, username, password)

			identity, err := coreClient.GetIdentity()
			Expect(err).NotTo(HaveOccurred())
			Expect(identity).To(Equal(JenkinsIdentity{
				Fingerprint:   "fingerprint",
				PublicKey:     "publicKey",
				SystemMessage: "systemMessage",
			}))
		})
	})

	Context("shutdown", func() {
		var (
			err  error
			safe bool
		)

		JustBeforeEach(func() {
			PrepareForShutdown(roundTripper, coreClient.URL, username, password, safe)
			err = coreClient.Shutdown(safe)
		})

		Context("shutdown safely", func() {
			BeforeEach(func() {
				safe = true
			})
			It("should success", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("shutdown not safely", func() {
			BeforeEach(func() {
				safe = false
			})
			It("should success", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Context("prepare shutdown", func() {
		var (
			err    error
			cancel bool
		)

		JustBeforeEach(func() {
			PrepareForCancelShutdown(roundTripper, coreClient.URL, username, password, cancel)
			err = coreClient.PrepareShutdown(cancel)
		})

		Context("cancelQuietDown", func() {
			BeforeEach(func() {
				cancel = true
			})
			It("should success", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("quietDown", func() {
			BeforeEach(func() {
				cancel = false
			})
			It("should success", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Context("toJson", func() {
		var (
			result GenericResult
			err    error
		)
		JustBeforeEach(func() {
			PrepareForToJSON(roundTripper, coreClient.URL, username, password)
			result, err = coreClient.ToJSON("jenkinsfile")
		})
		It("normal", func() {
			Expect(err).To(BeNil())
			Expect(result.GetResult()).To(Equal(`{"a":"b"}`))
		})
	})

	Context("toJenkinsfile", func() {
		var (
			result GenericResult
			err    error
		)
		JustBeforeEach(func() {
			PrepareForToJenkinsfile(roundTripper, coreClient.URL, username, password)
			result, err = coreClient.ToJenkinsfile("json")
		})
		It("normal", func() {
			Expect(err).To(BeNil())
			Expect(result.GetResult()).To(Equal("jenkinsfile"))
		})
	})

	Context("GetLabels", func() {
		var (
			labels *LabelsResponse
			err    error
		)
		JustBeforeEach(func() {
			PrepareForToGetLabels(roundTripper, coreClient.URL, username, password)
			labels, err = coreClient.GetLabels()
		})
		It("normal", func() {
			Expect(err).To(BeNil())
			Expect(labels).To(Equal(&LabelsResponse{
				Status: "ok",
				Data: []AgentLabel{{
					CloudsCount:                    0,
					Description:                    "",
					HasMoreThanOneJob:              false,
					JobsCount:                      0,
					JobsWithLabelDefaultValue:      []string{},
					JobsWithLabelDefaultValueCount: 0,
					Label:                          "java",
					LabelURL:                       "label/java/",
					NodesCount:                     1,
					PluginActiveForLabel:           false,
					TriggeredJobs:                  []string{},
					TriggeredJobsCount:             0,
				}},
			}))
		})
	})
})

func TestLabelsResponse_GetLabels(t *testing.T) {
	type fields struct {
		Status string
		Data   []AgentLabel
	}
	tests := []struct {
		name       string
		fields     fields
		wantLabels []string
	}{{
		name: "normal case",
		fields: fields{
			Status: "",
			Data: []AgentLabel{{
				Label: "good",
			}, {
				Label: "bad",
			}},
		},
		wantLabels: []string{"good", "bad"},
	}, {
		name:       "no data field",
		fields:     fields{},
		wantLabels: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LabelsResponse{
				Status: tt.fields.Status,
				Data:   tt.fields.Data,
			}
			if gotLabels := l.GetLabels(); !reflect.DeepEqual(gotLabels, tt.wantLabels) {
				t.Errorf("GetLabels() = %v, want %v", gotLabels, tt.wantLabels)
			}
		})
	}
}
