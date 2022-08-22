package credential

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSecretCredential(t *testing.T) {
	cre := NewSecretTextCredential("id", "secret")
	assert.NotNil(t, cre)
	assert.Equal(t, "secret", cre.Secret)
}
