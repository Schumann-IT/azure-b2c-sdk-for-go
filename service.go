package b2c

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/schumann-it/azure-b2c-sdk-for-go/environment"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
)

// Service represents a service that provides operations related to environments and policies.
type Service struct {
	es  []environment.Config
	sd  string
	td  string
	gs  *msgraph.ServiceClient
	tid string
}

// NewServiceFromConfigFile creates a new instance of Service by loading the environment configuration from the provided file path
// and initializing the necessary variables.
// It returns a pointer to the Service instance and an error, if any.
func NewServiceFromConfigFile(cp string, sd string, td string) (*Service, error) {
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

func (s *Service) WithEnvironments(environments []environment.Config) {
	s.es = environments
}

func (s *Service) WithSourceDir(dir string) error {
	d, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	s.sd = d

	return nil
}

func (s *Service) WithTargetDir(dir string) error {
	d, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	s.td = d

	return nil
}

func (s *Service) WithTenantId(tid string) error {
	_, err := uuid.Parse(tid)
	if err != nil {
		return err
	}

	s.tid = tid

	return nil
}

// FindConfig searches for the environment configuration with the specified name.
// It returns a pointer to the Config instance and an error, if the environment is not found.
func (s *Service) FindConfig(n string) (*environment.Config, error) {
	for _, e := range s.es {
		if e.Name == n {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("environment %s not found", n)
}

// CreateGraphClientFromEnvironment creates a new instance of the Microsoft Graph service client using environment variable configuration.
// It returns a pointer to the ServiceClient instance and an error if the creation fails.
func (s *Service) CreateGraphClientFromEnvironment() error {
	if s.gs != nil {
		return nil
	}

	var errs error

	tid := os.Getenv("B2C_ARM_TENANT_ID")
	if tid == "" {
		errs = multierror.Append(errs, fmt.Errorf("B2C_ARM_TENANT_ID must be set via environment variable"))
	}
	s.tid = tid

	cid := os.Getenv("B2C_ARM_CLIENT_ID")
	if cid == "" {
		errs = multierror.Append(errs, fmt.Errorf("B2C_ARM_CLIENT_ID must be set via environment variable"))
	}

	cs := os.Getenv("B2C_ARM_CLIENT_SECRET")
	if cs == "" {
		errs = multierror.Append(errs, fmt.Errorf("B2C_ARM_CLIENT_SECRET must be set via environment variable"))
	}

	if errs != nil {
		return errs
	}

	cred, err := azidentity.NewClientSecretCredential(tid, cid, cs, nil)
	if err != nil {
		return err
	}

	c, err := msgraph.NewClientWithCredential(cred)
	if err != nil {
		return err
	}
	s.gs = c

	return nil
}

func (s *Service) GetTenantInformation(tid *string) (models.TenantInformationable, error) {
	if tid == nil {
		tid = &s.tid
	}

	return s.gs.GetTenantInformation(to.String(tid))
}

func (s *Service) CreateGraphClientFromDefaultAzureCredential() error {
	if s.gs != nil {
		return nil
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	c, err := msgraph.NewClientWithCredential(cred)
	if err != nil {
		return err
	}
	s.gs = c

	return nil
}

func (s *Service) GetGraphClient() *msgraph.ServiceClient {
	return s.gs
}
