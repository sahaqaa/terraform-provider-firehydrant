package firehydrant

import (
	"context"
	"fmt"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
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

// APIClient is the client that accesses all of the api.firehydrant.io resources
type APIClient struct {
	baseURL string
	token   string

	client *sling.Sling
}

const (
	// DefaultBaseURL is the URL that is used to make requests to the FireHydrant API
	DefaultBaseURL = "https://api.firehydrant.io/v1/"
)

// Client is the client that makes requests to FireHydrant
type Client interface {
	Ping(ctx context.Context) (*PingResponse, error)
	GetService(ctx context.Context, id string) (*ServiceResponce, error)
	CreateService(ctx context.Context, req CreateServiceRequest) (*ServiceResponce, error)
	UpdateService(ctx context.Context, serviceID string, req UpdateServiceRequest) (*ServiceResponce, error)
	DeleteService(ctx context.Context, serviceID string) error
}

// OptFunc is a function that sets a setting on a client
type OptFunc func(c *APIClient) error

// WithBaseURL modifies the base URL for all requests
func WithBaseURL(baseURL string) OptFunc {
	return func(c *APIClient) error {
		c.baseURL = baseURL
		return nil
	}
}

// NewRestClient initializes a new API client for FireHydrant
func NewRestClient(token string, opts ...OptFunc) (*APIClient, error) {
	c := &APIClient{
		baseURL: DefaultBaseURL,
		token:   token,
	}

	for _, f := range opts {
		if err := f(c); err != nil {
			return nil, err
		}
	}

	c.client = sling.New().Base(c.baseURL).
		Set("User-Agent", fmt.Sprintf("%s (%s)", UserAgentPrefix, Version)).
		Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return c, nil
}

// Ping hits and verifies the HTTP of FireHydrant
// TODO: Check failure case
func (c *APIClient) Ping(ctx context.Context) (*PingResponse, error) {
	res := &PingResponse{}

	if _, err := c.client.Get("ping").Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not ping")
	}

	return res, nil
}

// GetService retrieves a service from the FireHydrant API
// TODO: Check failure case
func (c *APIClient) GetService(ctx context.Context, id string) (*ServiceResponce, error) {
	res := &ServiceResponce{}

	if _, err := c.client.Get("services/"+id).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not get service")
	}

	return res, nil
}

// CreateService creates a brand spanking new service in FireHydrant
// TODO: Check failure case
func (c *APIClient) CreateService(ctx context.Context, createReq CreateServiceRequest) (*ServiceResponce, error) {
	res := &ServiceResponce{}

	if _, err := c.client.Post("services").BodyJSON(&createReq).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not create service")
	}

	return res, nil
}

// UpdateService updates a old spanking service in FireHydrant
// TODO: Check failure case
func (c *APIClient) UpdateService(ctx context.Context, serviceID string, updateReq UpdateServiceRequest) (*ServiceResponce, error) {
	res := &ServiceResponce{}

	if _, err := c.client.Patch("services/"+serviceID).BodyJSON(&updateReq).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not update service")
	}

	return res, nil
}

// DeleteService deletes a service in FireHydrant
// TODO: Check failure case
func (c *APIClient) DeleteService(ctx context.Context, serviceID string) error {
	if _, err := c.client.Delete("services/"+serviceID).Receive(nil, nil); err != nil {
		return errors.Wrap(err, "could not delete service")
	}

	return nil
}
