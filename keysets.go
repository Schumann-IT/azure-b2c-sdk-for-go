package b2c

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/go-multierror"
	"github.com/schumann-it/azure-b2c-sdk-for-go/environment"
	"github.com/schumann-it/azure-b2c-sdk-for-go/keyset"
)

func (s *Service) SyncKeySets(en string) error {
	e, err := s.findConfig(en)
	if err == nil {
		return err
	}

	c, err := s.createGraphClient(e)
	if err != nil {
		return fmt.Errorf("failed to create graph client: %s", err)
	}

	return c.SyncKeySets(s.createKeySets(e))
}

func (s *Service) DeleteKeySets(en string) error {
	e, err := s.findConfig(en)
	if err == nil {
		return err
	}

	c, err := s.createGraphClient(e)
	if err != nil {
		return fmt.Errorf("failed to create graph client: %s", err)
	}

	var errs error
	for _, set := range s.createKeySets(e) {
		err := c.DeleteKeySet(set)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to delete key set %s: %s", to.String(set.Get().GetId()), err))
		}
	}

	return errs
}

func (s *Service) createKeySets(e *environment.Config) []*keyset.KeySet {
	var ks []*keyset.KeySet
	for _, c := range e.KeySets {
		k := keyset.NewKeySet(c.Name)
		if c.Use != nil {
			k.WithRsaKey(to.String(c.Use))
		}
		if c.CertificateBody != nil && c.Password != nil {
			k.WithCertificate(&keyset.Certificate{
				Key:      to.String(c.CertificateBody),
				Password: to.String(c.Password),
			})
		}
		ks = append(ks, k)
	}
	return ks
}
