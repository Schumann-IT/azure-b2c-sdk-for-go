package msgraph

import (
	"context"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// GetTenantInformation retrieves the tenant information for a given tenant ID.
//
// Parameters:
//   - tid (string): The ID of the tenant
//
// Returns:
//   - models.TenantInformationable: The tenant information
//   - error: Any error that occurred during the retrieval
//
// Example usage:
//
//	tid := "abc123"
//	ti, err := s.GetTenantInformation(tid)
//	if err != nil {
//	    log.Errorf("Failed to get tenant information: %v", err)
//	    return
//	}
//	// Use the tenant information
func (s *ServiceClient) GetTenantInformation(tid string) (models.TenantInformationable, error) {
	return s.gc.TenantRelationships().FindTenantInformationByTenantIdWithTenantId(to.StringPtr(tid)).Get(context.Background(), nil)
}
