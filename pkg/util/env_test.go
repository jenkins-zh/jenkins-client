package util

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	val := GetEnvOrDefault(fmt.Sprintf("%d", time.Now().Nanosecond()), "fake")
	assert.Equal(t, "fake", val)

	env := os.Environ()
	if len(env) > 0 {
		val = GetEnvOrDefault(strings.Split(env[0], "=")[0], "fake")
		assert.Equal(t, strings.Split(env[0], "=")[1], val)
	}
}
