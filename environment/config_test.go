package environment

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_NewConfigFromFile(t *testing.T) {
	c, err := NewConfigFromFile("testdata/config.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Equal(t, 3, len(*c))
}

func Test_NewConfig(t *testing.T) {
	s := "name: test"
	c, err := NewConfig([]byte(s))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Equal(t, "test", c.Name)
}

func Test_InvalidConfig(t *testing.T) {
	s := "name:\n- test\n"
	_, err := NewConfig([]byte(s))
	if err == nil {
		t.Fatalf("expected yaml error")
	}
	var yerr *yaml.TypeError
	errors.As(err, &yerr)
	assert.Equal(t, 1, len(yerr.Errors))
}
