package msgraph

import (
	"context"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/google/uuid"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// CreateServicePrincipal creates a service principal with the given name, appId, and resourceIds.
// It returns the app ID and secret password of the created service principal, or an error if any.
func (s *ServiceClient) CreateServicePrincipal(name string, appId string, resourceIds []string) (*string, *string, error) {
	var ras []models.ResourceAccessable
	for _, rid := range resourceIds {
		ra := models.NewResourceAccess()
		uid := uuid.MustParse(rid)
		ra.SetId(&uid)
		ra.SetTypeEscaped(to.StringPtr("Role"))
		ras = append(ras, ra)
	}
	rra := models.NewRequiredResourceAccess()
	rra.SetResourceAppId(to.StringPtr(appId)) // graph api
	rra.SetResourceAccess(ras)
	rras := []models.RequiredResourceAccessable{
		rra,
	}
	a := models.NewApplication()
	a.SetDisplayName(to.StringPtr(name))
	a.SetRequiredResourceAccess(rras)
	ares, err := s.GraphClient.Applications().Post(context.Background(), a, nil)
	if err != nil {
		return nil, nil, err
	}

	pw, _ := s.AddApplicationPasswordCredentials(to.String(ares.GetId()), "cli")

	spr := models.NewServicePrincipal()
	spr.SetAppId(ares.GetAppId())
	sprres, err := s.GraphClient.ServicePrincipals().Post(context.Background(), spr, nil)
	if err != nil {
		return nil, nil, err
	}

	return sprres.GetAppId(), pw.GetSecretText(), nil
}

// FindServicePrincipal searches for a service principal with the given name.
// It returns true if a service principal with the specified name is found, false otherwise.
func (s *ServiceClient) FindServicePrincipal(name string) bool {
	sps, err := s.GraphClient.ServicePrincipals().Get(context.Background(), nil)
	if err != nil {
		return false
	}

	for _, sp := range sps.GetValue() {
		if to.String(sp.GetDisplayName()) == name {
			return true
		}
	}

	return false
}
