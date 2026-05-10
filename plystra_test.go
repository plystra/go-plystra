package plystra

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClientLoginSetsBearerAndRoutesModules(t *testing.T) {
	var seenAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/auth/login":
			if r.Method != http.MethodPost {
				t.Fatalf("login method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"access_token": "token-123", "refresh_token": "refresh-123"}})
		case "/api/v1/auth/refresh":
			var body map[string]string
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("refresh body: %v", err)
			}
			if body["refresh_token"] != "refresh-123" {
				t.Fatalf("refresh token = %q", body["refresh_token"])
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"access_token": "token-456", "refresh_token": "refresh-456"}})
		case "/api/v1/admin/me":
			seenAuth = r.Header.Get("Authorization")
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"user_id": "user_alice"}})
		case "/api/v1/resources/invoice/invoice_001":
			seenAuth = r.Header.Get("Authorization")
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": "invoice_001"}})
		default:
			t.Fatalf("unexpected route: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)
	if _, err := client.Auth.Login(context.Background(), "alice@example.com", "plystra-demo"); err != nil {
		t.Fatalf("login: %v", err)
	}
	if client.AccessToken != "token-123" {
		t.Fatalf("access token = %q", client.AccessToken)
	}
	if client.RefreshToken != "refresh-123" {
		t.Fatalf("refresh token = %q", client.RefreshToken)
	}
	if _, err := client.Auth.Refresh(context.Background(), ""); err != nil {
		t.Fatalf("refresh: %v", err)
	}
	accessToken, refreshToken := client.Tokens()
	if accessToken != "token-456" || refreshToken != "refresh-456" {
		t.Fatalf("tokens = %q %q", accessToken, refreshToken)
	}
	if _, err := client.Admin.Me(context.Background()); err != nil {
		t.Fatalf("admin me: %v", err)
	}
	if seenAuth != "Bearer token-456" {
		t.Fatalf("Authorization = %q", seenAuth)
	}
	if _, err := client.Resources.Get(context.Background(), "invoice", "invoice_001"); err != nil {
		t.Fatalf("resource get: %v", err)
	}
}

func TestClientListUsesQueryAndAPIErrorCarriesDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/users":
			if got := r.URL.Query().Get("limit"); got != "10" {
				t.Fatalf("limit = %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{{"id": "user_alice"}}})
		case "/api/v1/audit/logs":
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{
				"code":       "ADMIN_PERMISSION_REQUIRED",
				"message":    "permission required",
				"request_id": "req_test",
			}})
		default:
			t.Fatalf("unexpected route: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)
	query := url.Values{"limit": []string{"10"}}
	users, err := client.Users.List(context.Background(), query)
	if err != nil {
		t.Fatalf("users list: %v", err)
	}
	if len(users) != 1 || users[0]["id"] != "user_alice" {
		t.Fatalf("users = %#v", users)
	}
	_, err = client.Audit.List(context.Background(), nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %T %[1]v, want APIError", err)
	}
	if apiErr.StatusCode != http.StatusForbidden || apiErr.Code != "ADMIN_PERMISSION_REQUIRED" || apiErr.RequestID != "req_test" {
		t.Fatalf("api error = %#v", apiErr)
	}
}

func TestClientWrapsInvalidJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("upstream unavailable"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.System.Health(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %T %[1]v, want APIError", err)
	}
	if apiErr.Code != "INVALID_JSON_RESPONSE" || apiErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("api error = %#v", apiErr)
	}
}
