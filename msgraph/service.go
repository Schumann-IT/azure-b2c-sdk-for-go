package msgraph

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type ServiceClient struct {
	gc *msgraph.GraphServiceClient
	ac azcore.TokenCredential
	s  []string
	t  *azcore.AccessToken
}

var scopes = []string{"https://graph.microsoft.com/.default"}

func NewClientWithCredential(cred azcore.TokenCredential) (*ServiceClient, error) {
	g, err := msgraph.NewGraphServiceClientWithCredentials(cred, scopes)
	if err != nil {
		return nil, err
	}

	return &ServiceClient{
		s:  scopes,
		ac: cred,
		gc: g,
	}, nil
}

func (s *ServiceClient) Token() (*azcore.AccessToken, error) {
	if s.t == nil {
		t, err := s.ac.GetToken(context.Background(), policy.TokenRequestOptions{
			Scopes: s.s,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get token: %s", err.Error())
		}

		s.t = &t
	}

	return s.t, nil
}

func (s *ServiceClient) Authorize(req *http.Request) error {
	t, err := s.Token()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Token))

	return nil
}

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
