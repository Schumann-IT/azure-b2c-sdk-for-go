package b2c

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/stretchr/testify/assert"
)

func Test_BuildPolicies(t *testing.T) {
	s := testHelperSetupService(t, "config")
	err := s.BuildPolicies("test")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expected := testHelperCountFiles(t, testSourceDir)
	actual := testHelperCountFiles(t, path.Join(testBuildTargetDir, "test"))

	assert.Equal(t, expected, actual)
}

func Test_CreateDeployBatch(t *testing.T) {
	s := testHelperSetupService(t, "config")
	_ = s.BuildPolicies("test")

	e, _ := s.FindConfig("test")
	_, err := s.batch(e)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Nil(t, err)
}

func Test_AccDeployPolicies(t *testing.T) {
	// We only run acceptance tests if an env var is set because they're
	// slow and generally require some outside configuration.
	if os.Getenv("TEST_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TEST_ACC' set")
		return
	}

	s := testHelperSetupService(t, "config")
	e, _ := s.FindConfig("test")
	e.Settings["EncKeyID"] = fmt.Sprintf("B2C_1A_%s", acctest.RandStringFromCharSet(10, "ABCDEFGHIJKLMNOPQRSTXYZabcdefghijklmnopqrstxyz"))
	e.Settings["SigKeyID"] = fmt.Sprintf("B2C_1A_%s", acctest.RandStringFromCharSet(10, "ABCDEFGHIJKLMNOPQRSTXYZabcdefghijklmnopqrstxyz"))
	_ = s.BuildPolicies("test")

	err := s.DeletePolicies()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	_, err = s.gs.GenerateKey(e.Settings["EncKeyID"], "sig", "RSA")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	_, err = s.gs.GenerateKey(e.Settings["SigKeyID"], "enc", "RSA")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	err = s.DeployPolicies("test")
	assert.Nil(t, err)

	// cleanup
	_ = s.DeletePolicies()
	_ = s.gs.DeleteKeySet(e.Settings["EncKeyID"])
	_ = s.gs.DeleteKeySet(fmt.Sprintf("%s.bak", e.Settings["EncKeyID"]))
	_ = s.gs.DeleteKeySet(e.Settings["SigKeyID"])
	_ = s.gs.DeleteKeySet(fmt.Sprintf("%s.bak", e.Settings["SigKeyID"]))
}

func Test_AccDeployPoliciesFailsForNonExsistingKeys(t *testing.T) {
	// We only run acceptance tests if an env var is set because they're
	// slow and generally require some outside configuration.
	if os.Getenv("TEST_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TEST_ACC' set")
		return
	}

	s := testHelperSetupService(t, "config")
	_ = s.BuildPolicies("test")
	err := s.DeletePolicies()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	err = s.DeployPolicies("test")
	assert.NotNil(t, err)
}
