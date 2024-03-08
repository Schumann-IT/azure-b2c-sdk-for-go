package msgraph

import (
	"context"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/trustframework"
	"github.com/schumann-it/azure-b2c-sdk-for-go/keyset"
)

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

func (s *ServiceClient) DeleteKeySet(id string) error {
	return s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(id).Delete(context.Background(), nil)
}

func (s *ServiceClient) GetKeySet(id string) (models.TrustFrameworkKeySetable, error) {
	return s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(id).Get(context.Background(), nil)
}

func (s *ServiceClient) CreateKeySet(name string) (models.TrustFrameworkKeySetable, error) {
	ks := models.NewTrustFrameworkKeySet()
	ks.SetId(to.StringPtr(name))

	return s.gc.TrustFramework().KeySets().Post(context.Background(), ks, nil)
}

func (s *ServiceClient) CreateKey(ks *keyset.KeySet) (models.TrustFrameworkKeySetable, error) {
	r, err := s.gc.TrustFramework().KeySets().Post(context.Background(), ks.Get(), nil)
	if err != nil {
		return nil, err
	}

	if to.String(r.GetId()) != to.String(ks.Get().GetId()) {
		log.Warningf("id changed while creating %s: new id is %s", to.String(ks.Get().GetId()), to.String(r.GetId()))
		ks.Get().SetId(r.GetId())
	}

	var key models.TrustFrameworkKeyable
	if ks.Key() != nil {
		key, err = s.GenerateKey(to.String(ks.Get().GetId()), to.String(ks.Key().GetUse()), to.String(ks.Key().GetKty()))
		if err != nil {
			return nil, err
		}
	}

	if ks.Certificate() != nil {
		key, err = s.UploadPkcs12(to.String(ks.Get().GetId()), ks.Certificate().Key, ks.Certificate().Password)
		if err != nil {
			return nil, err
		}
	}

	r.SetKeys([]models.TrustFrameworkKeyable{key})

	return r, err
}

func (s *ServiceClient) getKeySets() (models.TrustFrameworkKeySetCollectionResponseable, error) {
	return s.gc.TrustFramework().KeySets().Get(context.Background(), nil)
}

func (s *ServiceClient) GenerateKey(ksId, use, kty string) (models.TrustFrameworkKeyable, error) {
	r := trustframework.NewKeySetsItemGenerateKeyPostRequestBody()
	r.SetUse(to.StringPtr(use))
	r.SetKty(to.StringPtr(kty))

	key, err := s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(ksId).GenerateKey().Post(context.Background(), r, nil)

	return key, err
}

func (s *ServiceClient) UploadPkcs12(ksId, certificate, password string) (models.TrustFrameworkKeyable, error) {
	b := trustframework.NewKeySetsItemUploadPkcs12PostRequestBody()
	b.SetKey(to.StringPtr(certificate))
	b.SetPassword(to.StringPtr(password))
	key, err := s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(ksId).UploadPkcs12().Post(context.Background(), b, nil)

	return key, err
}

func (s *ServiceClient) keySetExists(c models.TrustFrameworkKeySetCollectionResponseable, ks *keyset.KeySet) bool {
	for _, ek := range c.GetValue() {
		if to.String(ks.Get().GetId()) == to.String(ek.GetId()) {
			return true
		}
	}

	return false
}
