package msgraph

import (
	"context"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/trustframework"
	"github.com/schumann-it/azure-b2c-sdk-for-go/keyset"
)

// SyncKeySets synchronizes the given key sets with the service's key sets. It checks if each key set in the input exists in the service's key sets. If a key set does not exist, it creates
func (s *ServiceClient) SyncKeySets(ks []*keyset.KeySet) error {
	r, err := s.getKeySets()
	if err != nil {
		return err
	}

	for _, rk := range ks {
		if !s.keySetExists(r, rk) {
			_, err = s.CreateKeySet(to.String(rk.Get().GetId()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteKeySet deletes a key set with the given ID.
// It sends a DELETE request to the service's TrustFramework KeySets endpoint using the specified ID.
// The context.Background() function is used to create a new background context.
// The function returns an error if the DELETE request fails.
func (s *ServiceClient) DeleteKeySet(id string) error {
	return s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(id).Delete(context.Background(), nil)
}

// GetKeySet retrieves a trust framework key set from the service by its ID.
// It returns the trust framework key set object and an error, if any.
// Example usage:
//
//	keySet, err := service.GetKeySet("keySetId")
//	if err != nil {
//		// handle error
//	}
//	// use keySet
//
// Parameters:
//
//	id - the ID of the trust framework key set to retrieve
//
// Returns:
//
//	The trust framework key set object and an error, if any.
func (s *ServiceClient) GetKeySet(id string) (models.TrustFrameworkKeySetable, error) {
	return s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(id).Get(context.Background(), nil)
}

// CreateKeySet creates a new key set with the given name. It creates a TrustFrameworkKeySet
// object using the provided name, and then makes a POST request to the TrustFramework KeySets
// endpoint of the ServiceClient with the new key set object. The method returns a
// TrustFrameworkKeySetable object and an error if any.
func (s *ServiceClient) CreateKeySet(name string) (models.TrustFrameworkKeySetable, error) {
	ks := models.NewTrustFrameworkKeySet()
	ks.SetId(to.StringPtr(name))

	return s.gc.TrustFramework().KeySets().Post(context.Background(), ks, nil)
}

// getKeySets retrieves the collection of key sets from the service.
func (s *ServiceClient) getKeySets() (models.TrustFrameworkKeySetCollectionResponseable, error) {
	return s.gc.TrustFramework().KeySets().Get(context.Background(), nil)
}

// GenerateKey generates a new key with the specified settings. It takes in the key set ID, use, and key type as parameters. It creates a request body with the specified use and key
func (s *ServiceClient) GenerateKey(keySetNameOrId, use, kty string) (models.TrustFrameworkKeySetable, error) {
	ks, err := s.CreateKeySet(keySetNameOrId)
	if err != nil {
		return nil, err
	}

	r := trustframework.NewKeySetsItemGenerateKeyPostRequestBody()
	r.SetUse(to.StringPtr(use))
	r.SetKty(to.StringPtr(kty))

	key, err := s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(to.String(ks.GetId())).GenerateKey().Post(context.Background(), r, nil)
	if err != nil {
		return nil, err
	}

	ks.SetKeys([]models.TrustFrameworkKeyable{key})

	return ks, err
}

// UploadPkcs12 uploads a PKCS12 certificate to the service for a specific trust framework key set identified by `ksId`. It takes the PKCS12 certificate and password as input. It creates
func (s *ServiceClient) UploadPkcs12(keySetNameOrId, certificate, password string) (models.TrustFrameworkKeySetable, error) {
	ks, err := s.CreateKeySet(keySetNameOrId)
	if err != nil {
		return nil, err
	}

	b := trustframework.NewKeySetsItemUploadPkcs12PostRequestBody()
	b.SetKey(to.StringPtr(certificate))
	b.SetPassword(to.StringPtr(password))
	key, err := s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(to.String(ks.GetId())).UploadPkcs12().Post(context.Background(), b, nil)
	if err != nil {
		return nil, err
	}

	ks.SetKeys([]models.TrustFrameworkKeyable{key})

	return ks, err
}

// keySetExists checks if the given key set exists in the TrustFrameworkKeySetCollectionResponseable.
// It iterates through each key set in the collection and compares the ID of the given key set with the ID of each key set in the collection.
// If a key set with the same ID is found, it returns true. Otherwise, it returns false.
func (s *ServiceClient) keySetExists(c models.TrustFrameworkKeySetCollectionResponseable, ks *keyset.KeySet) bool {
	for _, ek := range c.GetValue() {
		if to.String(ks.Get().GetId()) == to.String(ek.GetId()) {
			return true
		}
	}

	return false
}
