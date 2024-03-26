package msgraph

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	sdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// ServiceClient is a client for interacting with a service.
type ServiceClient struct {
	GraphClient        *sdk.GraphServiceClient
	OrganizationClient *OrganizationClient
	ac                 azcore.TokenCredential
	s                  []string
	t                  *azcore.AccessToken
}

// scopes is a variable of type []string that represents the different scopes for Microsoft Graph API requests.
var scopes = []string{"https://graph.microsoft.com/.default"}

// NewClientWithCredential creates a new ServiceClient with the provided TokenCredential.
// It returns a pointer to the newly created ServiceClient and any error encountered.
//
// The ServiceClient uses the provided TokenCredential for authentication and authorization.
// The TokenCredential must implement the azcore.TokenCredential interface.
//
// Example usage:
//
//	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, authorities...)
//	if err != nil {
//	  log.Fatal(err)
//	}
//	client, err := NewClientWithCredential(cred)
//	if err != nil {
//	  log.Fatal(err)
//	}
func NewClientWithCredential(cred azcore.TokenCredential) (*ServiceClient, error) {
	g, err := sdk.NewGraphServiceClientWithCredentials(cred, scopes)
	if err != nil {
		return nil, err
	}

	return &ServiceClient{
		s:           scopes,
		ac:          cred,
		GraphClient: g,
	}, nil
}

func (s *ServiceClient) CreateOrganizationClient(tid string) {
	s.OrganizationClient = &OrganizationClient{
		s:  s,
		id: tid,
	}
}

// Token retrieves an access token from the ServiceClient's TokenCredential.
// It returns a pointer to the AccessToken and any error encountered.
//
// If the ServiceClient's AccessToken is not already set, it calls the GetToken method on the TokenCredential,
// passing the ServiceClient's scopes as the TokenRequestOptions.
// If an error occurs while getting the token, it returns a nil AccessToken and an error message indicating the failure.
//
// Example usage:
//
//	token, err := client.Token()
//	if err != nil {
//	  log.Fatal(err)
//	}
//	fmt.Println("Access Token:", token.Token)
//
// Returns:
//   - AccessToken: The access token retrieved from the TokenCredential
//   - error: Any error encountered while getting the token
func (s *ServiceClient) Token() (*azcore.AccessToken, error) {
	if s.t == nil {
		t, err := s.ac.GetToken(context.Background(), policy.TokenRequestOptions{
			Scopes: s.s,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get token: %w", err)
		}

		s.t = &t
	}

	return s.t, nil
}

// Authorize sets the Authorization header in the given *http.Request
// with the access token retrieved from the ServiceClient's TokenCredential.
// It returns an error if there is an error retrieving the token.
//
// Example usage:
//
//	req, _ := http.NewRequest("GET", "https://api.example.com", nil)
//	err := client.Authorize(req)
//	if err != nil {
//	  log.Fatal(err)
//	}
//
// Returns:
//   - error: Any error encountered while getting the token
func (s *ServiceClient) Authorize(req *http.Request) error {
	t, err := s.Token()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Token))

	return nil
}

// DoRequest sends an HTTP request and checks the response status code.
// It takes an *http.Request as input and returns any error encountered.
//
// It creates an http.Client and defers the closure of idle connections.
// It then sends the request using the client's Do method and assigns the response to resp.
// If an error occurs while sending the request, it is returned.
//
// The function defers the closure of the response body to ensure it is closed when the function returns.
// It reads the entire response body using io.ReadAll and assigns it to the body variable.
// If an error occurs while reading the response body, it is returned.
//
// If the response status code is greater than or equal to 400, it returns an error message indicating the request failure,
// using the body contents as the error message.
//
// Example usage:
//
//	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
//	if err != nil {
//	  log.Fatal(err)
//	}
//	err = client.DoRequest(req)
//	if err != nil {
//	  log.Fatal(err)
//	}
//
// Returns:
//   - error: Any error encountered while sending the request or reading the response body
func (s *ServiceClient) DoRequest(req *http.Request) error {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed: %s", string(body))
	}

	return nil

}
