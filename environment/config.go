package environment

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration structure.
//
// Fields:
// - Name: The name of the configuration.
// - Tenant: The tenant associated with the configuration.
// - IsProduction: Flag indicating if the configuration is for production.
// - Applications: An array of ApplicationConfig representing the applications in the configuration.
// - KeySets: An array of KeySetsConfig representing the key sets in the configuration.
// - Settings: A map of string key-value pairs representing additional optional settings.
type Config struct {
	Name         string              `yaml:"name"`
	Tenant       string              `yaml:"tenant"`
	IsProduction bool                `yaml:"isProduction"`
	Applications []ApplicationConfig `yaml:"applications,omitempty"`
	KeySets      []KeySetsConfig     `yaml:"keySets,omitempty"`
	Settings     map[string]string   `yaml:"settings,omitempty"`
}

// KeySetsConfig represents a key set configuration structure.
// Fields:
// - Name: The name of the key set configuration.
// - Use: The use associated with the key set configuration. It is an optional field.
// - CertificateBody: The certificate body of the key set configuration. It is an optional field.
// - Password: The password of the key set configuration. It is an optional field.
type KeySetsConfig struct {
	Name            string  `yaml:"name"`
	Use             *string `yaml:"use,omitempty"`
	CertificateBody *string `yaml:"cert,omitempty"`
	Password        *string `yaml:"password,omitempty"`
}

// ApplicationConfig represents the configuration for an application.
// Fields:
// - Name: The name of the application.
// - ObjectId: The unique identifier of the application.
// - Patch: A map of string keys and interface{} values representing optional patch data for the application.
type ApplicationConfig struct {
	Name     string                 `yaml:"name"`
	ObjectId string                 `yaml:"objectId"`
	Patch    map[string]interface{} `yaml:"patch,omitempty"`
}

// NewConfig is a function that creates a new Config object by unmarshaling the provided byte slice of YAML data.
// If the unmarshaling process encounters an error, it returns nil and the error. Otherwise, it returns the newly created Config object and nil.
// The Config struct has the following fields:
// - Name: a string field representing the name of the configuration (yaml:"name")
// - Tenant: a string field representing the tenant of the configuration (yaml:"tenant")
// - IsProduction: a boolean field indicating whether the configuration is for production (yaml:"isProduction")
// - Applications: a slice of ApplicationConfig structs representing application configurations (yaml:"applications,omitempty")
// - KeySets: a slice of KeySetsConfig structs representing key set configurations (yaml:"keySets,omitempty")
// - Settings: a map of string key-value pairs representing additional settings (yaml:"settings,omitempty")
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
