package msgraph

import (
	_ "embed"
	"os"
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
	s.CreateOrganizationClient(os.Getenv("AZURE_TENANT_ID"))

	return s
}
