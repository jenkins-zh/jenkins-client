package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	fake := &struct {
		Name string
	}{
		Name: "fake",
	}

	assert.Equal(t, `{"Name":"fake"}`, TOJSON(fake))
}
