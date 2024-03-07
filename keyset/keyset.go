package keyset

import (
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

type KeySet struct {
	ks *models.TrustFrameworkKeySet
	k  *models.TrustFrameworkKey
	c  *Certificate
}

func NewKeySet(id string) *KeySet {
	ks := models.NewTrustFrameworkKeySet()
	ks.SetId(to.StringPtr(id))

	return &KeySet{
		ks: ks,
	}
}

func (ks *KeySet) WithCertificate(cert *Certificate) {
	ks.c = cert
}

func (ks *KeySet) WithRsaKey(use string) {
	k := models.NewTrustFrameworkKey()
	k.SetUse(to.StringPtr(use))
	k.SetKty(to.StringPtr("RSA"))

	ks.k = k
}

func (ks *KeySet) Get() *models.TrustFrameworkKeySet {
	return ks.ks
}

func (ks *KeySet) Certificate() *Certificate {
	return ks.c
}

func (ks *KeySet) Key() *models.TrustFrameworkKey {
	return ks.k
}
