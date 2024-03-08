package b2c

import (
	"fmt"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/schumann-it/azure-b2c-sdk-for-go/environment"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
)

// Service represents a service that provides operations related to environments and policies.
type Service struct {
	es []environment.Config
	sd string
	td string
}

// NewService creates a new instance of Service by loading the environment configuration from the provided file path
// and initializing the necessary variables.
// It returns a pointer to the Service instance and an error, if any.
func NewService(cp string, sd string, td string) (*Service, error) {
	c, err := environment.NewConfigFromFile(cp)
	if err != nil {
		return nil, err
	}

	asd, err := filepath.Abs(sd)
	if err != nil {
		return nil, err
	}

	tsd, err := filepath.Abs(td)
	if err != nil {
		return nil, err
	}

	return &Service{
		es: *c,
		sd: asd,
		td: tsd,
	}, nil
}

// findConfig searches for the environment configuration with the specified name.
// It returns a pointer to the Config instance and an error, if the environment is not found.
func (s *Service) findConfig(n string) (*environment.Config, error) {
	for _, e := range s.es {
		if e.Name == n {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("environment %s not found", n)
}

// createGraphClient creates a new instance of the Microsoft Graph service client using the provided environment configuration.
// It returns a pointer to the ServiceClient instance and an error if the creation fails.
func (s *Service) createGraphClient(e *environment.Config) (*msgraph.ServiceClient, error) {
	cred, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		return nil, err
	}

	return msgraph.NewClientWithCredential(cred)
}
