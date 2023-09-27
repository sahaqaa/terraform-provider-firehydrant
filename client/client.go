package client

// Actor represents an actor doing things in the FireHydrant API
type Actor struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
}

// PingResponse is the response the ping endpoint gives from FireHydrant
// URL: https://api.firehydrant.io/v1/ping
type PingResponse struct {
	Actor Actor `json:"actor"`
}

// ServiceResponce is the responce for a single service
// URL: https://api.firehydrant.io/v1/services
type ServiceResponce struct {
	ID          string            `json:"id"`
	Slug        string            `json:"slug"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
}
