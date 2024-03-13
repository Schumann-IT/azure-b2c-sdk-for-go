package msgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// GetApplication retrieves the application with the specified ID from the service.
func (s *ServiceClient) GetApplication(id string) (models.Applicationable, error) {
	return s.gc.Applications().ByApplicationId(id).Get(context.Background(), nil)
}

// PatchApplication updates an application with the specified ID using the provided patch.
// It takes the ID of the application and a map containing the updates to be applied.
// If the patching process encounters an error, it returns an error with a descriptive message.
// If the request to patch the application fails, it returns an error with a descriptive message.
// If the patching process is successful, it returns nil.
func (s *ServiceClient) PatchApplication(id string, patch map[string]interface{}) (models.Applicationable, error) {
	req, err := s.applicationPatchRequest(id, patch)
	if err != nil {
		return nil, fmt.Errorf("failed to patch application %s: %s", id, err)
	}

	err = s.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to patch application %s: %s", id, err)
	}

	return s.GetApplication(id)
}

// applicationPatchRequest sends a PATCH request to update the application with the specified ID in the service.
// It takes the ID of the application to be updated and a map containing the update to be applied to the application as input.
// The update is provided as a map of key-value pairs, where the keys are the fields to be updated and the values are the new values.
// The method returns a pointer to an http.Request object and an error. The request is constructed with the specified ID and update,
// and it is authorized using the ServiceClient's authorization mechanism.
// If an error occurs during the construction or authorization of the request, the method returns nil for the request and the error.
func (s *ServiceClient) applicationPatchRequest(id string, patch map[string]interface{}) (*http.Request, error) {
	b, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	ep := fmt.Sprintf("https://graph.microsoft.com/beta/applications/%s", id)
	req, err := http.NewRequest(http.MethodPatch, ep, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	err = s.Authorize(req)

	return req, err
}
