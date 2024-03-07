package b2c

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/stretchr/testify/assert"
)

func Test_CreateKeySets(t *testing.T) {
	expected := []string{"B2C_1A_TokenSigningKeyContainer", "B2C_1A_TokenEncryptionKeyContainer", "B2C_1A_SamlIdpCert"}

	s := testHelperSetupService(t, "config")
	e, _ := s.findConfig("test")
	ks := s.createKeySets(e)

	var actual []string
	for _, k := range ks {
		actual = append(actual, to.String(k.Get().GetId()))
	}

	assert.Equal(t, expected, actual)
}
