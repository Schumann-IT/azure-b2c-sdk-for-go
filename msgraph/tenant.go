package msgraph

import (
	"context"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (s *ServiceClient) GetTenantInformation(tid string) (models.TenantInformationable, error) {
	return s.gc.TenantRelationships().FindTenantInformationByTenantIdWithTenantId(to.StringPtr(tid)).Get(context.Background(), nil)
}
