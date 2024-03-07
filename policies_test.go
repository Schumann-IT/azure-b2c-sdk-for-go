package b2c

import (
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

	e, _ := s.findConfig("test")
	_, err := s.batch(e)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Nil(t, err)
}
