package credential

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSSHCredentail(t *testing.T) {
	cre := NewSSHCredential("id", "username", "passphrase", "privateKey")
	assert.NotNil(t, cre)
	assert.Equal(t, "id", cre.ID)
	assert.Equal(t, "username", cre.Username)
	assert.Equal(t, "passphrase", cre.Passphrase)
	assert.Equal(t, "privateKey", cre.KeySource.PrivateKey)
	assert.Equal(t, DirectSSHCrenditalStaplerClass, cre.KeySource.StaplerClass)
	assert.Equal(t, GLOBALScope, cre.Scope)
	assert.Equal(t, SSHCrenditalStaplerClass, cre.StaplerClass)
}
