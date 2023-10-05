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
	ac := m.(firehydrant.Client)
	serviceID := d.Get("id").(string)

	r, err := ac.GetService(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
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

	return ds
}

func resourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantService,
		UpdateContext: updateResourceFireHydrantService, // not sure if this is working lol
		ReadContext:   readResourceFireHydrantService,   // need to fix later
		DeleteContext: deleteResourceFireHydrantService, // need to fix later
		Schema: map[string]*schema.Schema{
			// "id": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			"name": {
				Type:     schema.TypeString,
				Required: true,
				//ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew: true,
			},
		},
	}
}

func createResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceName := d.Get("name").(string)
	serviceDescription := d.Get("description").(string)

	r := firehydrant.CreateServiceRequest{
		Name:        serviceName,
		Description: serviceDescription,
	}

	newService, err := ac.CreateService(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newService.ID)
	d.Set("name", newService.Name)
	d.Set("description", newService.Description)

	var ds diag.Diagnostics

	return ds
}

func updateResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceID := d.Id()
	serviceName := d.Get("name").(string)
	serviceDescription := d.Get("description").(string)

	r := firehydrant.UpdateServiceRequest{
		Name:        serviceName,
		Description: serviceDescription,
	}

	_, err := ac.UpdateService(ctx, serviceID, r)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func readResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceID := d.Id()

	r, err := ac.GetService(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	svc := map[string]string{
		"name":        r.Name,
		"description": r.Description,
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return ds
}

func deleteResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceID := d.Id()

	err := ac.DeleteService(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}
