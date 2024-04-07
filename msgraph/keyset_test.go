package msgraph

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/stretchr/testify/assert"
)

func Test_UploadKeySet(t *testing.T) {
	s := testHelperSetupService(t)
	kid := fmt.Sprintf("B2C_1A_%s", acctest.RandStringFromCharSet(10, "ABCDEFGHIJKLMNOPQRSTXYZabcdefghijklmnopqrstxyz"))
	ks, err := s.UploadSecret(kid, "sig", "secret")

	assert.Nil(t, err)
	assert.NotNil(t, ks)

	// cleanup
	_ = s.DeleteKeySet(kid)
	_ = s.DeleteKeySet(fmt.Sprintf("%s.bak", kid))
}
