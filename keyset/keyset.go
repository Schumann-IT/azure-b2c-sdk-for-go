package keyset

import (
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// KeySet represents a set of keys for a trust framework.
type KeySet struct {
	ks *models.TrustFrameworkKeySet
	k  *models.TrustFrameworkKey
	c  *Certificate
}

// NewKeySet creates a new KeySet with the given ID.
func NewKeySet(id string) *KeySet {
	ks := models.NewTrustFrameworkKeySet()
	ks.SetId(to.StringPtr(id))

	return &KeySet{
		ks: ks,
	}
}

// WithCertificate sets the certificate for the KeySet.
// It takes a pointer to a Certificate struct as an argument.
// Example usage:
//
//	ks := keyset.NewKeySet("example")
//	cert := &keyset.Certificate{
//	    Key:      "certificate key",
//	    Password: "certificate password",
//	}
//	ks.WithCertificate(cert)
func (ks *KeySet) WithCertificate(cert *Certificate) {
	ks.c = cert
}

// WithRsaKey sets the RSA key for the KeySet.
// It takes a string argument 'use' which represents the use of the key.
// Example usage:
//
//	ks := keyset.NewKeySet("example")
//	ks.WithRsaKey("signing")
func (ks *KeySet) WithRsaKey(use string) {
	k := models.NewTrustFrameworkKey()
	k.SetUse(to.StringPtr(use))
	k.SetKty(to.StringPtr("RSA"))

	ks.k = k
}

// Get returns the TrustFrameworkKeySet stored in the KeySet.
// Example usage:
//
//	ks := keyset.NewKeySet("example")
//	tfs := ks.Get()
func (ks *KeySet) Get() *models.TrustFrameworkKeySet {
	return ks.ks
}

// Certificate returns the certificate associated with the KeySet.
// It does not take any arguments.
// It returns a pointer to a Certificate struct.
// Example usage:
//
//	ks := keyset.NewKeySet("example")
//	cert := ks.Certificate()
//	fmt.Println(cert.Key)      // Output: certificate key
//	fmt.Println(cert.Password) // Output: certificate password
func (ks *KeySet) Certificate() *Certificate {
	return ks.c
}

// Key returns the TrustFrameworkKey stored in the KeySet.
// It returns a pointer to a TrustFrameworkKey struct.
// Example usage:
//
//	ks := keyset.NewKeySet("example")
//	key := ks.Key()
//	fmt.Println(key)
func (ks *KeySet) Key() *models.TrustFrameworkKey {
	return ks.k
}
