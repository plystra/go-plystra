package plystra

import (
	"context"
	"fmt"
)

func (s SystemService) Version(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/version", nil)
}
func (s SystemService) Health(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/health", nil)
}
func (s SystemService) Ready(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/ready", nil)
}

func (s AuthService) Login(ctx context.Context, email, password string) (Map, error) {
	out, err := s.client.postMap(ctx, "/api/v1/auth/login", Map{"email": email, "password": password})
	if err == nil {
		s.client.storeTokens(out)
	}
	return out, err
}

func (s AuthService) Refresh(ctx context.Context, refreshToken string) (Map, error) {
	if refreshToken == "" {
		refreshToken = s.client.RefreshToken
	}
	if refreshToken == "" {
		return nil, fmt.Errorf("plystra: refresh token is required")
	}
	out, err := s.client.postMap(ctx, "/api/v1/auth/refresh", Map{"refresh_token": refreshToken})
	if err == nil {
		s.client.storeTokens(out)
	}
	return out, err
}

func (s AuthService) Logout(ctx context.Context, refreshToken string) (Map, error) {
	if refreshToken == "" {
		refreshToken = s.client.RefreshToken
	}
	out, err := s.client.postMap(ctx, "/api/v1/auth/logout", Map{"refresh_token": refreshToken})
	if err == nil {
		s.client.SetTokens("", "")
	}
	return out, err
}

func (s ActorService) Context(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/actor/context", nil)
}

func (s ActorService) SwitchMember(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/actor/switch-member", input)
}

func (s AdminService) Me(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/admin/me", nil)
}

func (s AdminService) ListGrants(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/admin/grants", query)
}

func (s AdminService) GetGrant(ctx context.Context, grantID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/admin/grants/"+esc(grantID), nil)
}

func (s AdminService) CreateGrant(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/admin/grants", input)
}

func (s AdminService) RevokeGrant(ctx context.Context, grantID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/admin/grants/"+esc(grantID)+"/revoke", input)
}

func (s APIKeysService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/api-keys", query)
}

func (s APIKeysService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/api-keys", input)
}

func (s APIKeysService) Get(ctx context.Context, apiKeyID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/api-keys/"+esc(apiKeyID), nil)
}

func (s APIKeysService) Revoke(ctx context.Context, apiKeyID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/api-keys/"+esc(apiKeyID)+"/revoke", input)
}

func (s AuthzService) Check(ctx context.Context, input AuthzCheckInput) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/authz/check", input)
}

func (s AuthzService) Explain(ctx context.Context, input AuthzCheckInput) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/authz/explain", input)
}

func (s AuditService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/audit/logs", query)
}

func (s AuditService) Get(ctx context.Context, auditLogID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/audit/logs/"+esc(auditLogID), nil)
}
