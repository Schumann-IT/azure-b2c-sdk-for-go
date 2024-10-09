package msgraph

import (
	_ "embed"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func testMustTokenCredential(auth string) azcore.TokenCredential {
	var tc azcore.TokenCredential
	var err error

	if auth == "env" {
		tc, err = azidentity.NewEnvironmentCredential(nil)
	} else {
		tc, err = azidentity.NewDefaultAzureCredential(nil)
	}

	if err != nil {
		log.Fatalf("unexpected error: %s", err.Error())
	}

	return tc
}

func testHelperSetupService(t *testing.T, tc azcore.TokenCredential) *ServiceClient {
	s, err := NewClientWithCredential(tc)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	s.CreateOrganizationClient(os.Getenv("AZURE_TENANT_ID"))

	return s
}
