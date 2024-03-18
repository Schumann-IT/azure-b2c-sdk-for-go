package b2c

import (
	"os"
	"path"
	"testing"

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
	_ = s.BuildPolicies("test")
	err := s.DeletePolicies()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	_, err = s.gs.GenerateKey("B2C_1A_TokenSigningKeyContainer", "sig", "RSA")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	_, err = s.gs.GenerateKey("B2C_1A_TokenEncryptionKeyContainer", "enc", "RSA")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	err = s.DeployPolicies("test")
	assert.Nil(t, err)

	// cleanup
	_ = s.DeletePolicies()
	_ = s.gs.DeleteKeySet("B2C_1A_TokenSigningKeyContainer")
	_ = s.gs.DeleteKeySet("B2C_1A_TokenEncryptionKeyContainer")
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
