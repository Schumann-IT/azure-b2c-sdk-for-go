package environment

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration structure.
//
// Fields:
// - Name: The name of the configuration.
// - IsProduction: Flag indicating if the configuration is for production.
// - Settings: A map of string key-value pairs representing additional optional settings.
type Config struct {
	Name     string            `yaml:"name"`
	Settings map[string]string `yaml:"settings,omitempty"`
}

// NewConfig is a function that creates a new Config object by unmarshaling the provided byte slice of YAML data.
// If the unmarshaling process encounters an error, it returns nil and the error. Otherwise, it returns the newly created Config object and nil.
// The Config struct has the following fields:
// - Name: a string field representing the name of the configuration (yaml:"name").
// - Tenant: a string field representing the tenant of the configuration (yaml:"tenant").
// - IsProduction: a boolean field indicating whether the configuration is for production (yaml:"isProduction").
// - Settings: a map of string key-value pairs representing additional settings (yaml:"settings,omitempty").
func NewConfig(b []byte) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewConfigFromFile is a function that creates a new slice of Config objects by unmarshaling the YAML data from the provided file.
// If the file read operation encounters an error, it returns nil and the error. Otherwise, it returns the newly created slice of Config objects and nil.
//
// Example usage:
//
//	c, err := NewConfigFromFile("testdata/config.yaml")
//	if err != nil {
//		log.Fatalf("unexpected error: %s", err)
//	}
//	assert.Equal(t, 3, len(*c))
func NewConfigFromFile(f string) (*[]Config, error) {
	c := &[]Config{}
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
