package msgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (s *ServiceClient) GetApplication(id string) (models.Applicationable, error) {
	return s.gc.Applications().ByApplicationId(id).Get(context.Background(), nil)
}

func (s *ServiceClient) PatchApplication(id string, patch map[string]interface{}) error {
	req, err := s.applicationPatchRequest(id, patch)
	if err != nil {
		return fmt.Errorf("failed to patch application %s: %s", id, err)
	}

	err = s.DoRequest(req)
	if err != nil {
		return fmt.Errorf("failed to patch application %s: %s", id, err)
	}

	return nil
}

func (s *ServiceClient) applicationPatchRequest(id string, patch map[string]interface{}) (*http.Request, error) {
	b, err := json.Marshal(patch)
	ep := fmt.Sprintf("https://graph.microsoft.com/beta/applications/%s", id)
	req, err := http.NewRequest(http.MethodPatch, ep, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	err = s.Authorize(req)

	return req, err
}
