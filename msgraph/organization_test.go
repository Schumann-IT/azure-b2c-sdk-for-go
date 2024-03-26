package msgraph

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/corona10/goimagehash"
	"github.com/microsoft/kiota-abstractions-go/store"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/backgroundimage.png
	bgImagePng []byte
	//go:embed testdata/bannerlogo.jpg
	blJpg []byte
	//go:embed testdata/bannerlogo.png
	blPng []byte
	//go:embed testdata/squarelogolight.jpg
	sllJpg []byte
	//go:embed testdata/squarelogodark.jpg
	sldJpg []byte
)

func Test_OrganizationUpdateDisplayName(t *testing.T) {
	if os.Getenv("TEST_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TEST_ACC' set")
		return
	}

	p := []string{
		"de",
		"en",
		"fr",
	}

	s := testHelperSetupService(t)

	o, err := s.OrganizationClient.Get()
	assert.Nil(t, err)
	assert.NotNil(t, o)

	var c []string
	for _, v := range p {
		if v != to.String(o.GetPreferredLanguage()) {
			c = append(c, v)
		}
	}
	o.SetPreferredLanguage(to.StringPtr(c[rand.Intn(len(c))]))

	_, err = s.OrganizationClient.Update(o)
	assert.Nil(t, err)
}

func Test_OrganizationAddBrandingBrandingLocalization(t *testing.T) {
	// We only run acceptance tests if an env var is set because they're
	// slow and generally require some outside configuration.
	if os.Getenv("TEST_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TEST_ACC' set")
		return
	}

	lid := "de-DE"

	expected := map[string]string{
		"backgroundColor":  "#00a075",
		"signInPageText":   lid,
		"usernameHintText": fmt.Sprintf("%s Hint", lid),
	}

	s := testHelperSetupService(t)

	b, err := s.OrganizationClient.GetBrandingLocalization(lid)
	if b == nil {
		b = s.OrganizationClient.NewBrandingLocalization(lid)
	}
	b.SetBackgroundColor(to.StringPtr(expected["backgroundColor"]))
	b.SetSignInPageText(to.StringPtr(expected["signInPageText"]))
	b.SetUsernameHintText(to.StringPtr(expected["usernameHintText"]))

	if err == nil {
		_, err = s.OrganizationClient.UpdateBrandingLocalization(b)
	} else {
		_, err = s.OrganizationClient.CreateBrandingLocalization(b)
	}
	assert.Nil(t, err)

	err = s.OrganizationClient.UploadImage(lid, ImageTypeBannerLogo, blJpg)
	assert.Nil(t, err)
	err = s.OrganizationClient.UploadImage(lid, ImageTypeBackgroundImage, bgImagePng)
	assert.Nil(t, err)
	err = s.OrganizationClient.UploadImage(lid, ImageTypeSquareLogoLight, sllJpg)
	assert.Nil(t, err)
	err = s.OrganizationClient.UploadImage(lid, ImageTypeSquareLogoDark, sldJpg)
	assert.Nil(t, err)

	cb, _ := s.OrganizationClient.GetBrandingLocalization(lid)

	testHelperTestImage(t, blJpg, "jpg", fmt.Sprintf("https://%s/%s", cb.GetCdnList()[0], to.String(cb.GetBannerLogoRelativeUrl())))
	testHelperTestImage(t, bgImagePng, "png", fmt.Sprintf("https://%s/%s", cb.GetCdnList()[0], to.String(cb.GetBackgroundImageRelativeUrl())))
	testHelperTestImage(t, sllJpg, "jpg", fmt.Sprintf("https://%s/%s", cb.GetCdnList()[0], to.String(cb.GetSquareLogoRelativeUrl())))
	testHelperTestImage(t, sldJpg, "jpg", fmt.Sprintf("https://%s/%s", cb.GetCdnList()[0], to.String(cb.GetSquareLogoDarkRelativeUrl())))

	testHelperAssertLocalizationProperties(t, expected, cb.GetBackingStore())
}

func Test_OrganizationUpdateDefaultBrandingBranding(t *testing.T) {
	// We only run acceptance tests if an env var is set because they're
	// slow and generally require some outside configuration.
	if os.Getenv("TEST_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TEST_ACC' set")
		return
	}

	expected := map[string]string{
		"backgroundColor":  "#ffffff",
		"signInPageText":   "Default",
		"usernameHintText": "DefaultHint",
	}

	s := testHelperSetupService(t)

	b := s.OrganizationClient.NewBranding()
	b.SetBackgroundColor(to.StringPtr(expected["backgroundColor"]))
	b.SetSignInPageText(to.StringPtr(expected["signInPageText"]))
	b.SetUsernameHintText(to.StringPtr(expected["usernameHintText"]))
	_, err := s.OrganizationClient.UpdateDefaultBranding(b)
	assert.Nil(t, err)

	err = s.OrganizationClient.UploadDefaultImage(ImageTypeBannerLogo, blPng)
	assert.Nil(t, err)

	cb, _ := s.OrganizationClient.GetDefaultBranding()

	testHelperTestImage(t, blPng, "png", fmt.Sprintf("https://%s/%s", cb.GetCdnList()[0], to.String(cb.GetBannerLogoRelativeUrl())))

	testHelperAssertLocalizationProperties(t, expected, cb.GetBackingStore())
}

func testHelperAssertLocalizationProperties(t *testing.T, expected map[string]string, s store.BackingStore) {
	for k, v := range expected {
		cv, err := s.Get(k)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if cvs, ok := cv.(*string); ok {
			assert.Equal(t, v, to.String(cvs))
		} else {
			t.Fatalf("type assertion failed. should be *string: %v", cv)
		}
	}
}

func testHelperTestImage(t *testing.T, src []byte, srcType string, url string) {
	res, err := http.Get(url)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	defer res.Body.Close()

	var li image.Image
	var ri image.Image
	if srcType == "jpg" {
		li, _ = jpeg.Decode(bytes.NewReader(src))
		ri, _ = jpeg.Decode(res.Body)
	} else {
		li, _ = png.Decode(bytes.NewReader(src))
		ri, _ = png.Decode(res.Body)
	}
	lh, _ := goimagehash.AverageHash(li)
	rh, _ := goimagehash.AverageHash(ri)

	d, err := lh.Distance(rh)
	assert.Nil(t, err)
	assert.Equal(t, 0, d)
}
