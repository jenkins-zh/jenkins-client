package credential

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateKubeConfigCredential(t *testing.T) {
	cre := NewKubeConfigCredential("id", "kubeconfig")
	assert.NotNil(t, cre)
	assert.Equal(t, "id", cre.ID)
	assert.Equal(t, "kubeconfig", cre.KubeconfigSource.Content)
}
