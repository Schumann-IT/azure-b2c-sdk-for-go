package msgraph

import (
	_ "embed"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func testHelperSetupService(t *testing.T) *ServiceClient {
	tc, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	s, err := NewClientWithCredential(tc)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	s.CreateOrganizationClient("09cd16c8-453f-4a03-b72d-342409d41ed5")

	return s
}
