package provider

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sahaqaa/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	apiKeyName             = "api_key"
	firehydrantBaseURLName = "firehydrant_base_url"

	apiClientUserAgent = "github.com/sahaqaa/terraform-provider-firehydrant"
)

const (
	// MajorVersion is the major version
	MajorVersion = 0
	// MinorVersion is the minor version
	MinorVersion = 1
	// PatchVersion is the patch version
	PatchVersion = 0

	// UserAgentPrefix is the prefix of User-Agent header that all Terraform REST calls perform
	UserAgentPrefix = "firehydrant-terraform-provider"
)

// Version is the semver of this provider
var Version = fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)

// Provider returns a Terraform provider for the FireHydrant API
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			apiKeyName: {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("FIREHYDRANT_API_KEY", nil),
			},
			firehydrantBaseURLName: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FIREHYDRANT_BASE_URL", "https://api.firehydrant.io/v1/"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"firehydrant_service": resourceService(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"firehydrant_service": dataSourceService(),
		},
		ConfigureContextFunc: setupFireHydrantContext,
	}
}

func setupFireHydrantContext(ctx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := rd.Get(apiKeyName).(string)
	firehydrantBaseURL := rd.Get(firehydrantBaseURLName).(string)

	// apiClient := sling.New().Base(firehydrantBaseURL).
	// 	Set("User-Agent", fmt.Sprintf("%s (%s)", UserAgentPrefix, Version)).
	// 	Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	ac, err := firehydrant.NewRestClient(apiKey, firehydrant.WithBaseURL(firehydrantBaseURL))
	if err != nil {
		return nil, diag.FromErr(errors.Wrap(err, "cound not initialize API client"))
	}

	_, err = ac.Ping(ctx)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// return nil, diag.FromErr(fmt.Errorf(resp.Request.URL.String()))

	// if resp.StatusCode != 200 {
	// 	return nil, diag.FromErr(fmt.Errorf("Invalid response code from FireHydrant API: %d", resp.StatusCode))
	// }

	return ac, nil
}
