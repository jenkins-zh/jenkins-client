package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "fake")
	assert.Nil(t, err)
	assert.NotNil(t, f)
	defer func() {
		_ = os.RemoveAll(f.Name())
	}()

	err = os.WriteFile(f.Name(), []byte("hello"), 0600)
	assert.Nil(t, err)

	assert.Equal(t, "hello", string(ReadFile(f.Name())))
	assert.Equal(t, "hello", ReadFileASString(f.Name()))
}
