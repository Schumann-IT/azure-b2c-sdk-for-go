package msgraph

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/go-multierror"
	"github.com/schumann-it/azure-b2c-sdk-for-go/policy"
)

// GetPolicies retrieves the list of policies from the ServiceClient.
// It returns a slice of strings, where each string represents a policy ID.
// If an error occurs while retrieving the policies, it returns nil and an error.
func (s *ServiceClient) GetPolicies() ([]string, error) {
	d, err := s.GraphClient.TrustFramework().Policies().Get(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get policies: %w", err)
	}

	var r []string
	for _, p := range d.GetValue() {
		r = append(r, to.String(p.GetId()))
	}

	return r, nil
}

// DeletePolicies deletes all policies from the ServiceClient.
// It first calls the GetPolicies method to retrieve the list of policies.
// If an error occurs while retrieving the policies, it returns the error.
// It then iterates over each policy ID and deletes the policy using the corresponding API call.
// If an error occurs while deleting a policy, it adds the error to a multierror and continues to the next policy.
// After deleting all policies, it returns nil if there were no errors.
func (s *ServiceClient) DeletePolicies() error {
	ps, err := s.GetPolicies()
	if err != nil {
		return err
	}

	var errs error
	for _, id := range ps {
		err = s.GraphClient.TrustFramework().Policies().ByTrustFrameworkPolicyId(id).Delete(context.Background(), nil)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to delete policy %s: %w", id, err))
			continue
		}
		log.Debugf(fmt.Sprintf("successfully deleted policy %s", id))
	}

	return nil
}

// UploadPolicies takes a slice of policy.Policy objects and uploads them to the ServiceClient.
// It performs the upload concurrently by launching a goroutine for each policy in the slice.
// The method waits for all the goroutines to finish using a sync.WaitGroup.
// If any errors occur during the upload, they are collected and returned as a single error using multierror.Append.
// The function returns nil if all the policies were uploaded successfully, or the aggregated error if any upload failed.
// Please note that it is important to handle and propagate the returned error appropriately in the calling code.
func (s *ServiceClient) UploadPolicies(policies []policy.Policy) error {
	var wg sync.WaitGroup
	wg.Add(len(policies))

	res := make(chan error, len(policies))
	for _, p := range policies {
		go s.uploadPolicy(p, &wg, res)
	}
	wg.Wait()
	close(res)

	var errs error
	for err := range res {
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

// uploadPolicy uploads a policy to the ServiceClient.
// It takes a policy.Policy object and a wait group (wg) as input parameters.
// The error channel (errChan) is used to report any errors that occur during the execution of the method.
// The method defers the completion of the wait group using wg.Done().
// It then creates an uploadPolicyRequest using the policy ID and the policy data (byte array).
// If there is an error creating the upload request, it sends an error message to errChan and returns.
// Next, it calls the DoRequest method of the ServiceClient to upload the policy.
// If there is an error during the request, it sends an error message to errChan and returns.
// Finally, it logs a debug message indicating the successful upload of the policy.
// It then sends a nil value to errChan to indicate the successful completion of the method.
func (s *ServiceClient) uploadPolicy(p policy.Policy, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	req, err := s.uploadPolicyRequest(p.Id(), p.Byte())
	if err != nil {
		errChan <- fmt.Errorf("failed to upload policy %s: %w", p.Id(), err)
		return
	}

	err = s.DoRequest(req)
	if err != nil {
		errChan <- fmt.Errorf("failed to upload policy %s: %w", p.Id(), err)
		return
	}

	log.Debugf(fmt.Sprintf("successfully uploaded policy %s", p.Id()))

	errChan <- nil
}

// uploadPolicyRequest creates an HTTP request to upload a policy with the specified ID and body.
// The ID is used to construct the endpoint URL and the body is used as the request body.
// It returns a pointer to the created request and an error if any.
// The request has the HTTP method set to PUT and the endpoint URL is constructed using the ID.
// The request body is set to the provided body parameter.
// The Content-Type header is set to "application/xml; charset=utf-8".
// The request is authorized using the ServiceClient's Authorize method.
func (s *ServiceClient) uploadPolicyRequest(id string, body []byte) (*http.Request, error) {
	ep := fmt.Sprintf("https://graph.microsoft.com/beta/trustFramework/policies/%s/$value", id)
	req, err := http.NewRequest(http.MethodPut, ep, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	err = s.Authorize(req)

	return req, err
}
