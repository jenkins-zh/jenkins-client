package util

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Test open browser", func() {
	It("should success", func() {
		err := Open("fake://url", "", FakeExecCommandSuccess)
		Expect(err).To(BeNil())
	})
})

// TestShellProcessSuccess only for test
func TestShellProcessSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}
	//os.Exit(0)
}

func TestFakeLookPath(t *testing.T) {
	val, err := FakeLookPath("fake")
	assert.Nil(t, err)
	assert.Equal(t, "fake", val)

	err = FakeSystemCallExecSuccess("", nil, nil)
	assert.Nil(t, err)
}
