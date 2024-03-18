package b2c

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/schumann-it/azure-b2c-sdk-for-go/environment"
	"github.com/stretchr/testify/assert"
)

var (
	testFixturesDir       = "testdata"
	testFixturesSourceDir = "source"
	testBuildTargetDir    = "/tmp/b2ctests/build"
	testSourceDir         = ""
)

func testHelperSetupService(t *testing.T, env string) *Service {
	_ = os.RemoveAll(testBuildTargetDir)

	r, _ := filepath.Abs(testFixturesDir)
	cp := fmt.Sprintf("%s/%s.yaml", path.Join(testFixturesDir), env)
	testSourceDir = path.Join(r, testFixturesSourceDir)
	s, err := NewServiceFromConfigFile(cp, testSourceDir, testBuildTargetDir)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	return s
}

func testHelperCountFiles(t *testing.T, p string) int {
	c := 0
	_ = filepath.Walk(p, func(_ string, i os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if i.IsDir() {
			return nil
		}
		c++
		return nil
	})

	return c
}

func Test_NewService(t *testing.T) {
	expected := environment.Config{
		Name: "simple",
		Settings: map[string]string{
			"Tenant": "simple.onmicrosoft.com",
		},
	}

	s := testHelperSetupService(t, "simple")
	actual, err := s.FindConfig(expected.Name)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Equal(t, expected, *actual)
}

func Test_NewServiceFailsFonNonExistingConfig(t *testing.T) {
	_, actual := NewServiceFromConfigFile("missing", "source", "build")
	assert.NotNil(t, actual)
}

func Test_NewServiceWithRelativePaths(t *testing.T) {
	cp := fmt.Sprintf("%s/%s.yaml", path.Join(testFixturesDir), "config")
	s, err := NewServiceFromConfigFile(cp, "source", "build")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cwd, _ := os.Getwd()
	assert.Equal(t, path.Join(cwd, "source"), s.sd)
	assert.Equal(t, path.Join(cwd, "build"), s.td)
}

func Test_NewServiceFailsFornNonExistentEnvironment(t *testing.T) {
	s := testHelperSetupService(t, "simple")
	_, actual := s.FindConfig("missing")

	assert.NotNil(t, actual)
}
