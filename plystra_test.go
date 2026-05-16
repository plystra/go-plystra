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

func TestClientRegisterStoresTokens(t *testing.T) {
	var seenBody Map
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/api/v1/auth/register" {
			t.Fatalf("unexpected route: %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&seenBody); err != nil {
			t.Fatalf("register body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"access_token": "register-token", "refresh_token": "register-refresh"}})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	if _, err := client.Auth.Register(context.Background(), Map{
		"email":              "founder@example.com",
		"password":           "long-enough-password",
		"space_name":         "Founder Space",
		"registration_token": "registration-token",
	}); err != nil {
		t.Fatalf("register: %v", err)
	}
	accessToken, refreshToken := client.Tokens()
	if accessToken != "register-token" || refreshToken != "register-refresh" {
		t.Fatalf("tokens = %q %q", accessToken, refreshToken)
	}
	if seenBody["space_name"] != "Founder Space" {
		t.Fatalf("register body = %#v", seenBody)
	}
}

func TestClientSendsAPIKeyAndRoutesAPIKeysModule(t *testing.T) {
	var seenAPIKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		seenAPIKey = r.Header.Get("X-Plystra-API-Key")
		switch r.URL.Path {
		case "/api/v1/authz/check":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"decision": "allow"}})
		case "/api/v1/api-keys":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{{"id": "ak_test"}}})
		default:
			t.Fatalf("unexpected route: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("ply_ak_test.secret"))
	decision, err := client.Authz.Check(context.Background(), AuthzCheckInput{Action: "approve"})
	if err != nil {
		t.Fatalf("authz check: %v", err)
	}
	if decision["decision"] != "allow" {
		t.Fatalf("decision = %#v", decision)
	}
	keys, err := client.APIKeys.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("api keys list: %v", err)
	}
	if len(keys) != 1 || keys[0]["id"] != "ak_test" {
		t.Fatalf("api keys = %#v", keys)
	}
	if seenAPIKey != "ply_ak_test.secret" {
		t.Fatalf("X-Plystra-API-Key = %q", seenAPIKey)
	}
}

func TestClientPrefersAPIKeyOverBearerToken(t *testing.T) {
	var seenAPIKey string
	var seenAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		seenAPIKey = r.Header.Get("X-Plystra-API-Key")
		seenAuth = r.Header.Get("Authorization")
		if r.URL.Path != "/api/v1/authz/check" {
			t.Fatalf("unexpected route: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"decision": "allow"}})
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAccessToken("access-session"), WithAPIKey("ply_ak_test.secret"))
	decision, err := client.Authz.Check(context.Background(), AuthzCheckInput{Action: "approve"})
	if err != nil {
		t.Fatalf("authz check: %v", err)
	}
	if decision["decision"] != "allow" {
		t.Fatalf("decision = %#v", decision)
	}
	if seenAPIKey != "ply_ak_test.secret" {
		t.Fatalf("X-Plystra-API-Key = %q", seenAPIKey)
	}
	if seenAuth != "" {
		t.Fatalf("Authorization = %q, want empty", seenAuth)
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

	client := NewClient(server.URL, WithAPIKey("ply_ak_test.secret"))
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
	if seenAPIKey != "ply_ak_test.secret" {
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
		switch r.URL.Path {
		case "/api/v1/health":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"status": "ok"}})
		default:
			t.Fatalf("unexpected route: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := WithRequestID(context.Background(), "req_go_test")
	if _, err := client.System.Health(ctx); err != nil {
		t.Fatalf("health: %v", err)
	}
	if seenRequestID != "req_go_test" {
		t.Fatalf("X-Request-ID = %q", seenRequestID)
	}
}

func TestClientUsesCanonicalCoreRoutes(t *testing.T) {
	seen := map[string]bool{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		seen[r.Method+" "+r.URL.Path] = true
		switch r.URL.Path {
		case "/api/v1/spaces/space_acme/member-roles":
			if r.Method == http.MethodGet {
				_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{}})
			} else {
				_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"ok": true}})
			}
		case "/api/v1/plugins/com.example.payments/resources",
			"/api/v1/plugins/com.example.payments/permissions",
			"/api/v1/plugins/com.example.payments/audit-events",
			"/api/v1/plugins/com.example.payments/admin-menus",
			"/api/v1/plugins/com.example.payments/settings":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{}})
		case "/api/v1/spaces/space_acme/member-roles/mr_finance_reviewer_approver",
			"/api/v1/spaces/space_acme/member-roles/mr_finance_reviewer_approver/revoke":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"ok": true}})
		default:
			t.Fatalf("unexpected route: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := context.Background()
	if _, err := client.Spaces.MemberRoleGrants(ctx, "space_acme", nil); err != nil {
		t.Fatalf("space member roles: %v", err)
	}
	if _, err := client.MemberRoles.List(ctx, "space_acme", nil); err != nil {
		t.Fatalf("member roles list: %v", err)
	}
	if _, err := client.MemberRoles.Create(ctx, "space_acme", Map{"member_id": "member_finance_reviewer", "role_id": "role_finance_approver"}); err != nil {
		t.Fatalf("member roles create: %v", err)
	}
	if _, err := client.MemberRoles.Get(ctx, "space_acme", "mr_finance_reviewer_approver"); err != nil {
		t.Fatalf("member roles get: %v", err)
	}
	if _, err := client.MemberRoles.Revoke(ctx, "space_acme", "mr_finance_reviewer_approver", nil); err != nil {
		t.Fatalf("member roles revoke: %v", err)
	}
	if _, err := client.Plugins.Resources(ctx, "com.example.payments"); err != nil {
		t.Fatalf("plugin resources: %v", err)
	}
	if _, err := client.Plugins.Permissions(ctx, "com.example.payments"); err != nil {
		t.Fatalf("plugin permissions: %v", err)
	}
	if _, err := client.Plugins.AuditEvents(ctx, "com.example.payments"); err != nil {
		t.Fatalf("plugin audit events: %v", err)
	}
	if _, err := client.Plugins.AdminMenu(ctx, "com.example.payments"); err != nil {
		t.Fatalf("plugin admin menu: %v", err)
	}
	if _, err := client.Plugins.Settings(ctx, "com.example.payments", nil); err != nil {
		t.Fatalf("plugin settings: %v", err)
	}

	for _, want := range []string{
		"GET /api/v1/spaces/space_acme/member-roles",
		"POST /api/v1/spaces/space_acme/member-roles",
		"GET /api/v1/spaces/space_acme/member-roles/mr_finance_reviewer_approver",
		"POST /api/v1/spaces/space_acme/member-roles/mr_finance_reviewer_approver/revoke",
		"GET /api/v1/plugins/com.example.payments/admin-menus",
	} {
		if !seen[want] {
			t.Fatalf("missing route %s; saw %#v", want, seen)
		}
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
