package msgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (s *ServiceClient) GetApplication(id string) (models.Applicationable, error) {
	return s.gc.Applications().ByApplicationId(id).Get(context.Background(), nil)
}

func (s *ServiceClient) PatchApplication(id string, patch map[string]interface{}) error {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	b, err := json.Marshal(patch)
	ep := fmt.Sprintf("https://graph.microsoft.com/beta/applications/%s", id)
	req, err := http.NewRequest(http.MethodPatch, ep, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to patch application %s: %s", id, err)
	}

	t, err := s.Token()
	if err != nil {
		return fmt.Errorf("failed to patch application %s: %s", id, err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Token))
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("failed to patch application %s: %s", id, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to patch application %s: %s", id, err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to patch application %s: %s", id, string(body))
	}

	return nil
}
