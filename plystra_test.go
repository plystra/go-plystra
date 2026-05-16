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

func TestClientPublicSystemRoutesSkipAPIKey(t *testing.T) {
	var seen [][2]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		seen = append(seen, [2]string{r.URL.Path, r.Header.Get("X-Plystra-API-Key")})
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"status": "ok"}})
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("ply_kernel_secret"))
	if _, err := client.System.Health(context.Background()); err != nil {
		t.Fatalf("health: %v", err)
	}
	if _, err := client.System.Ready(context.Background()); err != nil {
		t.Fatalf("ready: %v", err)
	}
	if _, err := client.System.Version(context.Background()); err != nil {
		t.Fatalf("version: %v", err)
	}

	want := [][2]string{
		{"/api/v1/health", ""},
		{"/api/v1/ready", ""},
		{"/api/v1/version", ""},
	}
	if len(seen) != len(want) {
		t.Fatalf("seen = %#v", seen)
	}
	for i := range want {
		if seen[i] != want[i] {
			t.Fatalf("seen[%d] = %#v, want %#v", i, seen[i], want[i])
		}
	}
}

func TestClientProtectedKernelRoutesUseAPIKey(t *testing.T) {
	var seen [][3]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		seen = append(seen, [3]string{r.Method, r.URL.Path, r.Header.Get("X-Plystra-API-Key")})
		switch r.URL.Path {
		case "/api/v1/capabilities":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{{"name": "authorization.resource"}}})
		case "/api/v1/resource-types":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{{"key": "invoice"}}})
		case "/api/v1/audit/logs":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{}})
		default:
			t.Fatalf("unexpected route: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("ply_kernel_secret"))
	if _, err := client.System.Capabilities(context.Background()); err != nil {
		t.Fatalf("capabilities: %v", err)
	}
	if _, err := client.ResourceTypes.List(context.Background(), nil); err != nil {
		t.Fatalf("resource types: %v", err)
	}
	if _, err := client.Audit.List(context.Background(), url.Values{"limit": []string{"5"}}); err != nil {
		t.Fatalf("audit: %v", err)
	}

	want := [][3]string{
		{"GET", "/api/v1/capabilities", "ply_kernel_secret"},
		{"GET", "/api/v1/resource-types", "ply_kernel_secret"},
		{"GET", "/api/v1/audit/logs", "ply_kernel_secret"},
	}
	if len(seen) != len(want) {
		t.Fatalf("seen = %#v", seen)
	}
	for i := range want {
		if seen[i] != want[i] {
			t.Fatalf("seen[%d] = %#v, want %#v", i, seen[i], want[i])
		}
	}
}

func TestClientSendsAuthzContextModePayload(t *testing.T) {
	var seenAPIKey string
	var seenBody AuthzCheckInput
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		seenAPIKey = r.Header.Get("X-Plystra-API-Key")
		if r.URL.Path != "/api/v1/authz/check" {
			t.Fatalf("unexpected route: %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&seenBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{
			"decision":  "allow",
			"allow":     true,
			"trace_id":  "trc_test",
			"deny_code": nil,
		}})
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("ply_kernel_secret"))
	decision, err := client.Authz.Check(context.Background(), AuthzCheckInput{
		Actor: &ActorContext{
			UserID:     "user_external_alice",
			MemberID:   "member_finance_reviewer",
			BindingID:  "binding_external_alice_finance",
			SpaceID:    "space_acme",
			UserEmail:  "alice@example.com",
			MemberName: "Finance Reviewer",
		},
		Resource: &AuthzResourceContext{
			Type:          "invoice",
			ExternalID:    "invoice_001",
			SpaceID:       "space_acme",
			GroupPath:     "finance.apac",
			OwnerMemberID: "member_invoice_creator",
			Metadata:      Map{"amount": float64(1250)},
		},
		Grants: []AuthzGrantContext{{
			RoleKey:              "finance_approver",
			Resource:             "invoice",
			Action:               "approve",
			Scope:                "group_tree",
			SpaceID:              "space_acme",
			ScopeAnchorGroupPath: "finance",
		}},
		Action: "approve",
	})
	if err != nil {
		t.Fatalf("authz check: %v", err)
	}
	if decision["decision"] != "allow" {
		t.Fatalf("decision = %#v", decision)
	}
	if seenAPIKey != "ply_kernel_secret" {
		t.Fatalf("X-Plystra-API-Key = %q", seenAPIKey)
	}
	if seenBody.Actor.BindingID != "binding_external_alice_finance" {
		t.Fatalf("binding id = %q", seenBody.Actor.BindingID)
	}
	if seenBody.Resource.ExternalID != "invoice_001" || seenBody.Resource.GroupPath != "finance.apac" {
		t.Fatalf("resource = %#v", seenBody.Resource)
	}
	if len(seenBody.Grants) != 1 || seenBody.Grants[0].ScopeAnchorGroupPath != "finance" {
		t.Fatalf("grants = %#v", seenBody.Grants)
	}
}

func TestClientSendsRequestIDFromContext(t *testing.T) {
	var seenRequestID string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		seenRequestID = r.Header.Get("X-Request-ID")
		if r.URL.Path != "/api/v1/resource-types" {
			t.Fatalf("unexpected route: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{}})
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("ply_kernel_secret"))
	ctx := WithRequestID(context.Background(), "req_go_test")
	if _, err := client.ResourceTypes.List(ctx, nil); err != nil {
		t.Fatalf("resource types: %v", err)
	}
	if seenRequestID != "req_go_test" {
		t.Fatalf("X-Request-ID = %q", seenRequestID)
	}
}

func TestClientListUsesQueryAndAPIErrorCarriesDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/resource-types":
			if got := r.URL.Query().Get("limit"); got != "10" {
				t.Fatalf("limit = %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{{"key": "invoice"}}})
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

	client := NewClient(server.URL, WithAPIKey("ply_kernel_secret"))
	query := url.Values{"limit": []string{"10"}}
	rows, err := client.ResourceTypes.List(context.Background(), query)
	if err != nil {
		t.Fatalf("resource types list: %v", err)
	}
	if len(rows) != 1 || rows[0]["key"] != "invoice" {
		t.Fatalf("rows = %#v", rows)
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
