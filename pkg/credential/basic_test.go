package credential

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUsernamePasswordCredential(t *testing.T) {
	cre := NewUsernamePasswordCredential("id", "username", "password")
	assert.NotNil(t, cre)
	assert.Equal(t, "id", cre.ID)
	assert.Equal(t, "username", cre.Username)
	assert.Equal(t, "password", cre.Password)
}
