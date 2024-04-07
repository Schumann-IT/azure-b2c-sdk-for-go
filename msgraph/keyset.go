package msgraph

import (
	"context"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/trustframework"
)

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
	return s.GraphClient.TrustFramework().KeySets().ByTrustFrameworkKeySetId(id).Get(context.Background(), nil)
}

// GenerateKey generates a new key with the specified settings. It takes in the key set ID, use, and key type as parameters. It creates a request body with the specified use and key.
func (s *ServiceClient) GenerateKey(keySetNameOrId, use, kty string) (models.TrustFrameworkKeySetable, error) {
	ks, err := s.createKeySet(keySetNameOrId)
	if err != nil {
		return nil, err
	}

	r := trustframework.NewKeySetsItemGenerateKeyPostRequestBody()
	r.SetUse(to.StringPtr(use))
	r.SetKty(to.StringPtr(kty))

	key, err := s.GraphClient.TrustFramework().KeySets().ByTrustFrameworkKeySetId(to.String(ks.GetId())).GenerateKey().Post(context.Background(), r, nil)
	if err != nil {
		return nil, err
	}

	ks.SetKeys([]models.TrustFrameworkKeyable{key})

	return ks, err
}

func (s *ServiceClient) UploadSecret(keySetNameOrId, use, secret string) (models.TrustFrameworkKeySetable, error) {
	ks, err := s.createKeySet(keySetNameOrId)
	if err != nil {
		return nil, err
	}

	r := trustframework.NewKeySetsItemUploadSecretPostRequestBody()
	r.SetUse(to.StringPtr(use))
	r.SetK(to.StringPtr(secret))

	key, err := s.GraphClient.TrustFramework().KeySets().ByTrustFrameworkKeySetId(to.String(ks.GetId())).UploadSecret().Post(context.Background(), r, nil)
	if err != nil {
		return nil, err
	}

	ks.SetKeys([]models.TrustFrameworkKeyable{key})

	return ks, err
}

// UploadPkcs12 uploads a PKCS12 certificate to the service for a specific trust framework key set identified by `ksId`. It takes the PKCS12 certificate and password as input.
func (s *ServiceClient) UploadPkcs12(keySetNameOrId, certificate, password string) (models.TrustFrameworkKeySetable, error) {
	ks, err := s.createKeySet(keySetNameOrId)
	if err != nil {
		return nil, err
	}

	b := trustframework.NewKeySetsItemUploadPkcs12PostRequestBody()
	b.SetKey(to.StringPtr(certificate))
	b.SetPassword(to.StringPtr(password))
	key, err := s.GraphClient.TrustFramework().KeySets().ByTrustFrameworkKeySetId(to.String(ks.GetId())).UploadPkcs12().Post(context.Background(), b, nil)
	if err != nil {
		return nil, err
	}

	ks.SetKeys([]models.TrustFrameworkKeyable{key})

	return ks, err
}

// DeleteKeySet deletes a key set with the given ID.
// It sends a DELETE request to the service's TrustFramework KeySets endpoint using the specified ID.
// The context.Background() function is used to create a new background context.
// The function returns an error if the DELETE request fails.
func (s *ServiceClient) DeleteKeySet(id string) error {
	return s.GraphClient.TrustFramework().KeySets().ByTrustFrameworkKeySetId(id).Delete(context.Background(), nil)
}

// createKeySet creates a new key set with the given name. It creates a TrustFrameworkKeySet
// object using the provided name, and then makes a POST request to the TrustFramework KeySets
// endpoint of the ServiceClient with the new key set object. The method returns a
// TrustFrameworkKeySetable object and an error if any.
func (s *ServiceClient) createKeySet(name string) (models.TrustFrameworkKeySetable, error) {
	ks := models.NewTrustFrameworkKeySet()
	ks.SetId(to.StringPtr(name))

	return s.GraphClient.TrustFramework().KeySets().Post(context.Background(), ks, nil)
}
