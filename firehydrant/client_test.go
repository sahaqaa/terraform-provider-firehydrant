package firehydrant

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	pingResponseJSON    = `{"response":"pong","actor":{"id":"3dcdaf4c-54a1-4688-bd38-2839c14ee529","name":"Oleksandr","email":"saha338-firehydrant@yahoo.com","type":"firehydrant_user"},"organization":{"name":"Example.com","id":"03d01b60-7166-4d5c-ad2a-56e316968efa"}}`
	serviceResponceJSON = `{"id":"9538c479-51d6-410e-985b-81325c04ffd4","name":"qqq","description":"","slug":"qqq","created_at":"2023-10-03T06:29:03.321Z","updated_at":"2023-10-03T06:29:03.321Z","labels":{}}`
)

func TestClientInitialization(t *testing.T) {
	var requestPathReceived, token string

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPathReceived = req.URL.Path
		// t.Logf("Headers! %+v", req.Header)
		token = req.Header.Get("Authorization")

		w.Write([]byte(pingResponseJSON))
	})

	ts := httptest.NewServer(h)

	defer ts.Close()

	testToken := "testing-123"
	c, err := NewRestClient(testToken, WithBaseURL(ts.URL))

	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	res, err := c.Ping()
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}

	actorID := res.Actor.ID
	actorEmail := res.Actor.Email

	if expected := "/ping"; expected != requestPathReceived {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathReceived)
	}

	if expected := "Bearer " + testToken; expected != token {
		t.Fatalf("Expected %s, Got: %s for bearer token", expected, token)
	}

	if expected := "3dcdaf4c-54a1-4688-bd38-2839c14ee529"; expected != actorID {
		t.Fatalf("Expected %s, Got: %s for actor ID", expected, actorID)
	}

	if expected := "saha338-firehydrant@yahoo.com"; expected != actorEmail {
		t.Fatalf("Expected %s, Got: %s for actor email", expected, actorEmail)
	}
}

func TestGetService(t *testing.T) {
	var requestPathReceived string

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPathReceived = req.URL.Path
		// t.Logf("Headers! %+v", req.Header)
		//token = req.Header.Get("Authorization")

		w.Write([]byte(serviceResponceJSON))
	})

	ts := httptest.NewServer(h)

	defer ts.Close()

	testToken := "testing-123"
	c, err := NewRestClient(testToken, WithBaseURL(ts.URL))

	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	testServiceID := "test-service-id"
	res, err := c.GetService(testServiceID)
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}

	serviceID := res.ID
	serviceName := res.Name

	if expected := "/services/" + testServiceID; expected != requestPathReceived {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathReceived)
	}

	if expected := "9538c479-51d6-410e-985b-81325c04ffd4"; expected != serviceID {
		t.Fatalf("Expected %s, Got: %s for service ID", expected, serviceID)
	}

	if expected := "qqq"; expected != serviceName {
		t.Fatalf("Expected %s, Got: %s for service name", expected, serviceName)
	}
}

func TestCreateService(t *testing.T) {
	var requestPathReceived string

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPathReceived = req.URL.Path

		w.Write([]byte(serviceResponceJSON))
	})

	ts := httptest.NewServer(h)

	defer ts.Close()

	testToken := "testing-123"
	c, err := NewRestClient(testToken, WithBaseURL(ts.URL))

	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	testServiceID := "test-service-id"
	res, err := c.GetService(testServiceID)
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}

	serviceID := res.ID
	serviceName := res.Name

	if expected := "/services/" + testServiceID; expected != requestPathReceived {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathReceived)
	}

	if expected := "9538c479-51d6-410e-985b-81325c04ffd4"; expected != serviceID {
		t.Fatalf("Expected %s, Got: %s for service ID", expected, serviceID)
	}

	if expected := "qqq"; expected != serviceName {
		t.Fatalf("Expected %s, Got: %s for service name", expected, serviceName)
	}
}
