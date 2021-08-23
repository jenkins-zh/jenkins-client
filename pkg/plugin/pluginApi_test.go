package plugin

import (
	"os"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("plugin api test", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		pluginAPI    API
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		pluginAPI = API{
			RoundTripper: roundTripper,
			ShowProgress: false,
			UseMirror:    false,
			SkipOptional: true,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("ShowTrend", func() {
		It("basic case", func() {
			keyword := "fake"

			PrepareShowTrend(roundTripper, keyword)

			trend, err := pluginAPI.ShowTrend(keyword)
			Expect(err).To(BeNil())
			Expect(trend).NotTo(Equal(""))
		})
	})

	Context("DownloadPlugins", func() {
		var (
			names []string
		)

		BeforeEach(func() {
			names = []string{}
		})

		It("empty name list", func() {
			// ignore error
			_ = pluginAPI.DownloadPlugins(names)
		})

		It("one plugin name", func() {
			names = []string{"fake"}

			PrepareOnePluginInfo(roundTripper, "fake")
			PrepareDownloadPlugin(roundTripper)
			// ignore error
			_ = pluginAPI.DownloadPlugins(names)

			_, err := os.Stat("fake.hpi")
			defer func() {
				// ignore this error
				_ = os.Remove("fake.hpi")
			}()
			Expect(err).To(BeNil())
		})

		It("use mirror", func() {
			pluginAPI.UseMirror = true
			pluginAPI.MirrorURL = "http://updates.jenkins-ci.org/download/"
			names = []string{"fake"}

			PrepareOnePluginInfo(roundTripper, "fake")
			PrepareDownloadPlugin(roundTripper)
			// ignore error
			_ = pluginAPI.DownloadPlugins(names)

			_, err := os.Stat("fake.hpi")
			defer func() {
				// ignore error
				_ = os.Remove("fake.hpi")
			}()
			Expect(err).To(BeNil())
		})

		It("with dependency which is not optional", func() {
			names = []string{"fake"}

			PrepareOnePluginWithDep(roundTripper, "fake")
			PrepareDownloadPlugin(roundTripper)
			PrepareDownloadPlugin(roundTripper)
			pluginAPI.SkipDependency = false
			// ignore error
			_ = pluginAPI.DownloadPlugins(names)

			var err error
			_, err = os.Stat("fake.hpi")
			Expect(err).To(BeNil())
			_, err = os.Stat("fake-1.hpi")
			Expect(err).To(BeNil())

			defer func() {
				// ignore error
				_ = os.Remove("fake.hpi")
			}()
			defer func() {
				// ignore error
				_ = os.Remove("fake-1.hpi")
			}()
		})

		It("with dependency which is optional", func() {
			names = []string{"fake"}

			PrepareOnePluginWithOptionalDep(roundTripper, "fake")
			PrepareDownloadPlugin(roundTripper)
			pluginAPI.SkipDependency = false
			pluginAPI.SkipOptional = true
			_ = pluginAPI.DownloadPlugins(names)

			var err error
			_, err = os.Stat("fake.hpi")
			Expect(err).To(BeNil())

			defer func() {
				// ignore error
				_ = os.Remove("fake.hpi")
			}()
		})

		It("skip dependency", func() {
			names = []string{"fake"}

			PrepareOnePluginWithOptionalDep(roundTripper, "fake")
			PrepareDownloadPlugin(roundTripper)
			pluginAPI.SkipDependency = true
			pluginAPI.SkipOptional = true
			_ = pluginAPI.DownloadPlugins(names)

			var err error
			_, err = os.Stat("fake.hpi")
			Expect(err).To(BeNil())

			defer func() {
				// ignore error
				_ = os.Remove("fake.hpi")
			}()
		})

		It("batch search plugins", func() {
			names = []string{"fake"}

			PrepareShowPlugins(roundTripper, "fake")

			plugins, _ := pluginAPI.BatchSearchPlugins("fake")
			Expect(len(plugins)).To(Equal(0))
		})
	})
})
