package provider

import (
	"context"
	"fmt"

	"github.com/sahaqaa/terraform-provider-firehydrant/client"

	"github.com/dghubble/sling"
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
		DataSourcesMap: map[string]*schema.Resource{
			"firehydrant_service": dataSourceService(),
		},
		ConfigureContextFunc: setupFireHydrantContext,
	}
}

func setupFireHydrantContext(ctx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := rd.Get(apiKeyName).(string)
	firehydrantBaseURL := rd.Get(firehydrantBaseURLName).(string)

	apiClient := sling.New().Base(firehydrantBaseURL).
		Set("User-Agent", fmt.Sprintf("%s (%s)", UserAgentPrefix, Version)).
		Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	var pr client.PingResponse

	resp, err := apiClient.Get("ping").Receive(&pr, nil)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// return nil, diag.FromErr(fmt.Errorf(resp.Request.URL.String()))

	if resp.StatusCode != 200 {
		return nil, diag.FromErr(fmt.Errorf("Invalid response code from FireHydrant API: %d", resp.StatusCode))
	}

	return apiClient, nil
}

func dataSourceService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantService,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*sling.Sling)
	serviceID := d.Get("id").(string)

	var r client.ServiceResponce
	_, err := apiClient.Get("services/").Get(serviceID).Receive(&r, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	svc := map[string]string{
		"name":        r.Name,
		"description": r.Description,
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(r.ID)

	return nil
}
