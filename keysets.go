package b2c

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/go-multierror"
	"github.com/schumann-it/azure-b2c-sdk-for-go/environment"
	"github.com/schumann-it/azure-b2c-sdk-for-go/keyset"
)

// SyncKeySets syncs the key sets for a given environment.
//
// It finds the configuration for the environment by calling the `findConfig` method.
// If the configuration is not found, an error is returned.
//
// It then creates a graph client by calling the `createGraphClient` method using the found configuration.
// If there is an error creating the client, an error is returned wrapping the original error.
//
// Finally, it calls the `SyncKeySets` method on the graph client, passing in the key sets created by calling the `createKeySets` method using the found configuration.
// If there is an error syncing the key sets, an error is returned.
//
// Example:
//
//	service := &Service{}
//	err := service.SyncKeySets("env")
//
//	if err != nil {
//	    log.Fatal(err)
//	}
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

// DeleteKeySets deletes the key sets for a given environment.
//
// It finds the configuration for the environment by calling the `findConfig` method.
// If the configuration is not found, an error is returned.
//
// It then creates a graph client by calling the `createGraphClient` method using the found configuration.
// If there is an error creating the client, an error is returned wrapping the original error.
//
// It iterates over the key sets created by calling the `createKeySets` method using the found configuration,
// and calls the `DeleteKeySet` method on the graph client for each key set.
// If there is an error deleting a key set, an error is appended to the `errs` error variable.
//
// The final `errs` error variable is returned.
//
// Example:
//
//	service := &Service{}
//	err := service.DeleteKeySets("env")
//
//	if err != nil {
//	    log.Fatal(err)
//	}
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
		err := c.DeleteKeySet(to.String(set.Get().GetId()))
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to delete key set %s: %s", to.String(set.Get().GetId()), err))
		}
	}

	return errs
}

// createKeySets creates an array of key sets based on the provided environment config.
// It iterates over each KeySetsConfig in the environment config and creates a new KeySet using the name specified.
// If the Use field is not nil, it sets the RSA key with the specified value.
// If both the CertificateBody and Password fields are not nil, it sets the Certificate with the specified key and password.
// The created KeySet is then appended to the array.
// Example:
//
//	service := &Service{}
//	keySets := service.createKeySets(environmentConfig)
//	for _, ks := range keySets {
//	    // Do something with ks
//	}
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
