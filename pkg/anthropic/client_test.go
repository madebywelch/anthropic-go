package anthropic

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client, err := NewClient(apiKey)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if client.apiKey != apiKey {
		t.Errorf("expected apiKey %q, but got %q", apiKey, client.apiKey)
	}
}

func TestDoRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.baseURL = server.URL

	request, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	response, err := client.doRequest(request)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, response.StatusCode)
	}
}
