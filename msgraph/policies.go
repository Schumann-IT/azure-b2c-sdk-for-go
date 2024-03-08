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

func (s *ServiceClient) GetPolicies() ([]string, error) {
	d, err := s.gc.TrustFramework().Policies().Get(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get policies: %s", err)
	}

	var r []string
	for _, p := range d.GetValue() {
		r = append(r, to.String(p.GetId()))
	}

	return r, nil
}

func (s *ServiceClient) DeletePolicies() error {
	ps, err := s.GetPolicies()
	if err != nil {
		return err
	}

	var errs error
	for _, id := range ps {
		err = s.gc.TrustFramework().Policies().ByTrustFrameworkPolicyId(id).Delete(context.Background(), nil)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to delete policy %s: %s", id, err))
			continue
		}
		log.Debugf(fmt.Sprintf("sucessfully deleted policy %s", id))
	}

	return nil
}

func (s *ServiceClient) UploadPolicies(policies []policy.Policy) error {
	var wg sync.WaitGroup
	wg.Add(len(policies))

	errChan := make(chan error, len(policies))
	for _, p := range policies {
		go s.uploadPolicy(p, &wg, errChan)
	}
	wg.Wait()

	var errs error
	for err := range errChan {
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (s *ServiceClient) uploadPolicy(p policy.Policy, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	req, err := s.uploadPolicyRequest(p.Id(), p.Byte())
	if err != nil {
		errChan <- fmt.Errorf("failed to upload policy %s: %s", p.Id(), err)
		return
	}

	err = s.DoRequest(req)
	if err != nil {
		errChan <- fmt.Errorf("failed to upload policy %s: %s", p.Id(), err)
		return
	}

	log.Debugf(fmt.Sprintf("sucessfully uploaded policy %s", p.Id()))

	errChan <- nil
}

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
