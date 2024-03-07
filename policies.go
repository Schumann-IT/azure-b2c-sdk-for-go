package b2c

import (
	"fmt"
	"path"

	"github.com/hashicorp/go-multierror"
	"github.com/schumann-it/azure-b2c-sdk-for-go/environment"
	"github.com/schumann-it/azure-b2c-sdk-for-go/policy"
)

func (s *Service) BuildPolicies(en string) error {
	e, err := s.findConfig(en)
	if err != nil {
		return err
	}

	b := policy.NewBuilder()
	err = b.Read(s.sd)
	if err != nil {
		return fmt.Errorf("failed to build %s: read from %s failed: %s", en, s.sd, err)
	}
	err = b.Process(e.Settings)
	if err != nil {
		return fmt.Errorf("failed to process %s: %s", en, err)
	}
	err = b.Write(path.Join(s.td, e.Name))
	if err != nil {
		return fmt.Errorf("failed to write to %s/%s: %s", s.td, e.Name, err)
	}

	return nil
}

func (s *Service) ListPolicies(en string) error {
	e, err := s.findConfig(en)
	if err != nil {
		return err
	}

	c, err := s.createGraphClient(e)
	if err != nil {
		return fmt.Errorf("failed to create graph client: %s", err)
	}

	ps, err := c.GetPolicies()

	for _, p := range ps {
		log.Infof("found policy %s", p)
	}

	return nil
}

func (s *Service) DeletePolicies(en string) error {
	e, err := s.findConfig(en)
	if err != nil {
		return err
	}

	c, err := s.createGraphClient(e)
	if err != nil {
		return fmt.Errorf("failed to create graph client: %s", err)
	}

	return c.DeletePolicies()
}

func (s *Service) DeployPolicies(en string) error {
	e, err := s.findConfig(en)
	if err != nil {
		return err
	}

	c, err := s.createGraphClient(e)
	if err != nil {
		return fmt.Errorf("failed to create graph client: %s", err)
	}

	bs, err := s.batch(e)
	if err != nil {
		return err
	}

	var errs error
	for i, b := range bs {
		err = c.UploadPolicies(b)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to upload batch %d from %s: %s", i, en, err))
		}
	}

	return errs
}

func (s *Service) batch(e *environment.Config) ([][]policy.Policy, error) {
	t := &policy.Tree{}

	td := path.Join(s.td, e.Name)
	err := t.Read(td)
	if err != nil {
		return nil, fmt.Errorf("failed to read from %s, did you run build?: %s", td, err)
	}

	return t.Batches(), nil
}
