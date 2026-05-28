package plystra

import "context"

func (s AuthService) Register(ctx context.Context, input AuthRegisterInput) (Map, error) {
	return s.client.postPublicMap(ctx, "/api/v1/auth/register", input)
}

func (s AuthService) Login(ctx context.Context, input AuthLoginInput) (Map, error) {
	return s.client.postPublicMap(ctx, "/api/v1/auth/login", input)
}

func (s AuthService) Refresh(ctx context.Context, refreshToken string) (Map, error) {
	return s.client.postPublicMap(ctx, "/api/v1/auth/refresh", AuthRefreshInput{RefreshToken: refreshToken})
}

func (s AuthService) Logout(ctx context.Context, refreshToken string) (Map, error) {
	return s.client.postPublicMap(ctx, "/api/v1/auth/logout", AuthLogoutInput{RefreshToken: refreshToken})
}

func (s ActorService) Context(ctx context.Context) (Map, error) {
	return s.client.getBearerMap(ctx, "/api/v1/actor/context", nil)
}

func (s ActorService) SwitchMember(ctx context.Context, input SwitchMemberInput) (Map, error) {
	return s.client.postBearerMap(ctx, "/api/v1/actor/switch-member", input)
}

func (s SystemService) Version(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/version", nil, true)
}

func (s SystemService) Health(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/health", nil, true)
}

func (s SystemService) Ready(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/ready", nil, true)
}

func (s SystemService) Capabilities(ctx context.Context) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/capabilities", nil)
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
	return s.client.getMap(ctx, "/api/v1/audit/logs/"+esc(auditLogID), nil, false)
}
