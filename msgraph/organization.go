package msgraph

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest/to"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/organization"
)

type OrganizationClient struct {
	s  *ServiceClient
	id string
}

type ImageType int

const (
	ImageTypeBackgroundImage ImageType = iota
	ImageTypeBannerLogo
	ImageTypeSquareLogoLight
	ImageTypeSquareLogoDark
)

func (it *ImageType) update(req *organization.ItemBrandingLocalizationsOrganizationalBrandingLocalizationItemRequestBuilder, body []byte) ([]byte, error) {
	ct := to.StringPtr(http.DetectContentType(body))
	switch *it {
	case ImageTypeBackgroundImage:
		return req.BackgroundImage().Put(context.Background(), body, ct, nil)
	case ImageTypeBannerLogo:
		return req.BannerLogo().Put(context.Background(), body, ct, nil)
	case ImageTypeSquareLogoLight:
		return req.SquareLogo().Put(context.Background(), body, ct, nil)
	case ImageTypeSquareLogoDark:
		return req.SquareLogoDark().Put(context.Background(), body, ct, nil)
	default:
		return nil, fmt.Errorf("image type %d not found", *it)
	}
}

func (c *OrganizationClient) GetAll() (models.OrganizationCollectionResponseable, error) {
	return c.s.GraphClient.Organization().Get(context.Background(), nil)
}

func (c *OrganizationClient) Get() (models.Organizationable, error) {
	return c.s.GraphClient.Organization().ByOrganizationId(c.id).Get(context.Background(), nil)
}

func (c *OrganizationClient) NewBranding() models.OrganizationalBrandingable {
	return models.NewOrganizationalBranding()
}

func (c *OrganizationClient) GetDefaultBranding() (models.OrganizationalBrandingable, error) {
	h := abstractions.NewRequestHeaders()
	h.Add("Accept-Language", "0")
	cfg := &organization.ItemBrandingRequestBuilderGetRequestConfiguration{
		Headers: h,
	}

	b, err := c.s.GraphClient.Organization().ByOrganizationId(c.id).Branding().Get(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	b.GetBackingStore().SetInitializationCompleted(true)

	return b, nil
}

func (c *OrganizationClient) UpdateDefaultBranding(b models.OrganizationalBrandingable) (models.OrganizationalBrandingable, error) {
	h := abstractions.NewRequestHeaders()
	h.Add("Accept-Language", "0")
	cfg := &organization.ItemBrandingRequestBuilderPatchRequestConfiguration{
		Headers: h,
	}

	return c.s.GraphClient.Organization().ByOrganizationId(c.id).Branding().Patch(context.Background(), b, cfg)
}

func (c *OrganizationClient) NewBrandingLocalization(lid string) models.OrganizationalBrandingLocalizationable {
	b := models.NewOrganizationalBrandingLocalization()
	b.SetId(to.StringPtr(lid))

	return b
}

func (c *OrganizationClient) GetBrandingLocalization(lid string) (models.OrganizationalBrandingLocalizationable, error) {
	b, err := c.s.GraphClient.Organization().ByOrganizationId(c.id).Branding().Localizations().ByOrganizationalBrandingLocalizationId(lid).Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	b.GetBackingStore().SetInitializationCompleted(true)

	return b, nil
}

func (c *OrganizationClient) UpdateBrandingLocalization(b models.OrganizationalBrandingLocalizationable) (models.OrganizationalBrandingLocalizationable, error) {
	return c.s.GraphClient.Organization().ByOrganizationId(c.id).Branding().Localizations().ByOrganizationalBrandingLocalizationId(to.String(b.GetId())).Patch(context.Background(), b, nil)
}

func (c *OrganizationClient) CreateBrandingLocalization(b models.OrganizationalBrandingLocalizationable) (models.OrganizationalBrandingLocalizationable, error) {
	return c.s.GraphClient.Organization().ByOrganizationId(c.id).Branding().Localizations().Post(context.Background(), b, nil)
}

func (c *OrganizationClient) DeleteBrandingLocalization(lid string) error {
	b := c.NewBrandingLocalization(lid)
	b.SetBackgroundColor(nil)
	b.SetSignInPageText(nil)
	b.SetUsernameHintText(nil)

	_, err := c.UpdateBrandingLocalization(b)
	if err != nil {
		return err
	}

	return c.s.GraphClient.Organization().ByOrganizationId(c.id).Branding().Localizations().ByOrganizationalBrandingLocalizationId(lid).Delete(context.Background(), nil)
}

func (c *OrganizationClient) DeleteDefaultBranding() error {
	return c.DeleteBrandingLocalization("0")
}

func (c *OrganizationClient) UploadImage(lid string, it ImageType, body []byte) error {
	req := c.s.GraphClient.Organization().ByOrganizationId(c.id).Branding().Localizations().ByOrganizationalBrandingLocalizationId(lid)
	_, err := it.update(req, body)

	return err
}

func (c *OrganizationClient) UploadDefaultImage(it ImageType, body []byte) error {
	return c.UploadImage("0", it, body)
}

func (c *OrganizationClient) GetInformation() (models.TenantInformationable, error) {
	return c.s.GraphClient.TenantRelationships().FindTenantInformationByTenantIdWithTenantId(&c.id).Get(context.Background(), nil)
}
