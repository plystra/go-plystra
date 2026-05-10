package plystra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Map map[string]any
type Query = url.Values

type Pagination struct {
	Limit   int     `json:"limit"`
	Cursor  *string `json:"cursor"`
	HasMore bool    `json:"has_more"`
}

type Envelope[T any] struct {
	Data       T           `json:"data"`
	Error      *ErrorBody  `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	RequestID  string      `json:"request_id,omitempty"`
	Meta       Map         `json:"meta,omitempty"`
}

type ErrorBody struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
	DenyCode   string `json:"deny_code,omitempty"`
	TraceID    string `json:"trace_id,omitempty"`
	AuditLogID string `json:"audit_log_id,omitempty"`
}

type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Details    any
	RequestID  string
	TraceID    string
	AuditLogID string
}

func (e APIError) Error() string {
	return fmt.Sprintf("plystra: %s: %s", e.Code, e.Message)
}

type Client struct {
	BaseURL        string
	AccessToken    string
	RefreshToken   string
	APIKey         string
	HTTPClient     *http.Client
	DefaultHeaders http.Header
	UserAgent      string

	System          SystemService
	Authz           AuthzService
	Audit           AuditService
	Auth            AuthService
	Actor           ActorService
	Admin           AdminService
	APIKeys         APIKeysService
	Users           UsersService
	Spaces          SpacesService
	Groups          GroupsService
	Members         MembersService
	UserMembers     UserMembersService
	Roles           RolesService
	MemberRoles     MemberRolesService
	Permissions     PermissionsService
	RolePermissions RolePermissionsService
	ResourceTypes   ResourceTypesService
	Resources       ResourcesService
	Data            DataService
	Plugins         PluginsService
	Templates       TemplatesService
}

type ClientOption func(*Client)

func WithAccessToken(token string) ClientOption {
	return func(c *Client) { c.AccessToken = token }
}

func WithRefreshToken(token string) ClientOption {
	return func(c *Client) { c.RefreshToken = token }
}

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) { c.APIKey = apiKey }
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.HTTPClient = httpClient
		}
	}
}

func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		if c.DefaultHeaders == nil {
			c.DefaultHeaders = http.Header{}
		}
		c.DefaultHeaders.Set(key, value)
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) { c.UserAgent = userAgent }
}

func NewClient(baseURL string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL:        strings.TrimRight(baseURL, "/"),
		HTTPClient:     &http.Client{Timeout: 10 * time.Second},
		DefaultHeaders: http.Header{},
		UserAgent:      "go-plystra/1.0.0",
	}
	for _, opt := range opts {
		opt(c)
	}
	c.System = SystemService{client: c}
	c.Authz = AuthzService{client: c}
	c.Audit = AuditService{client: c}
	c.Auth = AuthService{client: c}
	c.Actor = ActorService{client: c}
	c.Admin = AdminService{client: c}
	c.APIKeys = APIKeysService{client: c}
	c.Users = UsersService{client: c}
	c.Spaces = SpacesService{client: c}
	c.Groups = GroupsService{client: c}
	c.Members = MembersService{client: c}
	c.UserMembers = UserMembersService{client: c}
	c.Roles = RolesService{client: c}
	c.MemberRoles = MemberRolesService{client: c}
	c.Permissions = PermissionsService{client: c}
	c.RolePermissions = RolePermissionsService{client: c}
	c.ResourceTypes = ResourceTypesService{client: c}
	c.Resources = ResourcesService{client: c}
	c.Data = DataService{client: c}
	c.Plugins = PluginsService{client: c}
	c.Templates = TemplatesService{client: c}
	return c
}

func (c *Client) SetAccessToken(token string) {
	c.AccessToken = token
}

func (c *Client) SetRefreshToken(token string) {
	c.RefreshToken = token
}

func (c *Client) SetAPIKey(apiKey string) {
	c.APIKey = apiKey
}

func (c *Client) SetTokens(accessToken, refreshToken string) {
	c.AccessToken = accessToken
	c.RefreshToken = refreshToken
}

func (c *Client) Tokens() (accessToken, refreshToken string) {
	return c.AccessToken, c.RefreshToken
}

type ActorContext struct {
	UserID       string `json:"user_id"`
	SpaceID      string `json:"space_id"`
	MemberID     string `json:"member_id"`
	UserMemberID string `json:"user_member_id"`
}

type AuthzCheckInput struct {
	Actor             *ActorContext `json:"actor,omitempty"`
	ActorUserID       string        `json:"actor_user_id,omitempty"`
	ActorMemberID     string        `json:"actor_member_id,omitempty"`
	ActorUserMemberID string        `json:"actor_user_member_id,omitempty"`
	SpaceID           string        `json:"space_id,omitempty"`
	ResourceType      string        `json:"resource_type,omitempty"`
	ResourceID        string        `json:"resource_id,omitempty"`
	Resource          *struct {
		Type string `json:"type,omitempty"`
		ID   string `json:"id,omitempty"`
	} `json:"resource,omitempty"`
	Action  string `json:"action"`
	Explain bool   `json:"explain,omitempty"`
}

type SystemService struct{ client *Client }
type AuthzService struct{ client *Client }
type AuditService struct{ client *Client }
type AuthService struct{ client *Client }
type ActorService struct{ client *Client }
type AdminService struct{ client *Client }
type APIKeysService struct{ client *Client }
type UsersService struct{ client *Client }
type SpacesService struct{ client *Client }
type GroupsService struct{ client *Client }
type MembersService struct{ client *Client }
type UserMembersService struct{ client *Client }
type RolesService struct{ client *Client }
type MemberRolesService struct{ client *Client }
type PermissionsService struct{ client *Client }
type RolePermissionsService struct{ client *Client }
type ResourceTypesService struct{ client *Client }
type ResourcesService struct{ client *Client }
type DataService struct{ client *Client }
type PluginsService struct{ client *Client }
type TemplatesService struct{ client *Client }

func (s SystemService) Version(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/system/version", nil)
}
func (s SystemService) Health(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/system/health", nil)
}
func (s SystemService) Ready(ctx context.Context) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/system/ready", nil)
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

func (s UsersService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/users", query)
}
func (s UsersService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/users", input)
}
func (s UsersService) Get(ctx context.Context, userID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/users/"+esc(userID), nil)
}
func (s UsersService) Update(ctx context.Context, userID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/users/"+esc(userID), input)
}
func (s UsersService) Disable(ctx context.Context, userID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/users/"+esc(userID)+"/disable", input)
}
func (s UsersService) Restore(ctx context.Context, userID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/users/"+esc(userID)+"/restore", input)
}

func (s SpacesService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces", query)
}
func (s SpacesService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces", input)
}
func (s SpacesService) Get(ctx context.Context, spaceID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID), nil)
}
func (s SpacesService) Update(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID), input)
}
func (s SpacesService) Disable(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/disable", input)
}
func (s SpacesService) Restore(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/restore", input)
}
func (s SpacesService) Groups(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups", query)
}
func (s SpacesService) GroupTree(ctx context.Context, spaceID string) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups/tree", nil)
}
func (s SpacesService) Members(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/members", query)
}
func (s SpacesService) UserMembers(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/user-members", query)
}
func (s SpacesService) Roles(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/roles", query)
}
func (s SpacesService) MemberRoleGrants(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/member-role-grants", query)
}
func (s SpacesService) Resources(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/resources", query)
}
func (s SpacesService) AuditLogs(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/audit-logs", query)
}

func (s GroupsService) Create(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups", input)
}
func (s GroupsService) Get(ctx context.Context, groupID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/groups/"+esc(groupID), nil)
}
func (s GroupsService) GetInSpace(ctx context.Context, spaceID, groupID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups/"+esc(groupID), nil)
}
func (s GroupsService) Update(ctx context.Context, spaceID, groupID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups/"+esc(groupID), input)
}
func (s GroupsService) Disable(ctx context.Context, spaceID, groupID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups/"+esc(groupID)+"/disable", input)
}

func (s MembersService) Create(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/members", input)
}
func (s MembersService) Get(ctx context.Context, memberID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/members/"+esc(memberID), nil)
}
func (s MembersService) GetInSpace(ctx context.Context, spaceID, memberID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/members/"+esc(memberID), nil)
}
func (s MembersService) Update(ctx context.Context, spaceID, memberID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/members/"+esc(memberID), input)
}
func (s MembersService) Disable(ctx context.Context, spaceID, memberID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/members/"+esc(memberID)+"/disable", input)
}

func (s UserMembersService) Create(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/user-members", input)
}
func (s UserMembersService) Get(ctx context.Context, userMemberID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/user-members/"+esc(userMemberID), nil)
}
func (s UserMembersService) GetInSpace(ctx context.Context, spaceID, userMemberID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/user-members/"+esc(userMemberID), nil)
}
func (s UserMembersService) Update(ctx context.Context, spaceID, userMemberID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/user-members/"+esc(userMemberID), input)
}
func (s UserMembersService) Revoke(ctx context.Context, spaceID, userMemberID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/user-members/"+esc(userMemberID)+"/revoke", input)
}

func (s RolesService) Create(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/roles", input)
}
func (s RolesService) Get(ctx context.Context, roleID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/roles/"+esc(roleID), nil)
}
func (s RolesService) GetInSpace(ctx context.Context, spaceID, roleID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/roles/"+esc(roleID), nil)
}
func (s RolesService) Update(ctx context.Context, spaceID, roleID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/roles/"+esc(roleID), input)
}
func (s RolesService) Disable(ctx context.Context, spaceID, roleID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/roles/"+esc(roleID)+"/disable", input)
}

func (s MemberRolesService) List(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/member-role-grants", query)
}
func (s MemberRolesService) Create(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/member-role-grants", input)
}
func (s MemberRolesService) Get(ctx context.Context, spaceID, memberRoleID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/member-role-grants/"+esc(memberRoleID), nil)
}
func (s MemberRolesService) Revoke(ctx context.Context, spaceID, memberRoleID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/member-role-grants/"+esc(memberRoleID)+"/revoke", input)
}

func (s PermissionsService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/permissions", query)
}
func (s PermissionsService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/permissions", input)
}
func (s PermissionsService) Get(ctx context.Context, permissionID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/permissions/"+esc(permissionID), nil)
}
func (s PermissionsService) Update(ctx context.Context, permissionID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/permissions/"+esc(permissionID), input)
}
func (s PermissionsService) Disable(ctx context.Context, permissionID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/permissions/"+esc(permissionID)+"/disable", input)
}

func (s RolePermissionsService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/role-permissions", query)
}
func (s RolePermissionsService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/role-permissions", input)
}
func (s RolePermissionsService) Get(ctx context.Context, rolePermissionID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/role-permissions/"+esc(rolePermissionID), nil)
}
func (s RolePermissionsService) Revoke(ctx context.Context, rolePermissionID string, input Map) (Map, error) {
	return s.client.deleteMap(ctx, "/api/v1/role-permissions/"+esc(rolePermissionID), input)
}

func (s ResourceTypesService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/resource-types", query)
}
func (s ResourceTypesService) Upsert(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/resource-types", input)
}
func (s ResourceTypesService) Get(ctx context.Context, key string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/resource-types/"+esc(key), nil)
}
func (s ResourceTypesService) Actions(ctx context.Context, key string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/resource-types/"+esc(key)+"/actions", query)
}
func (s ResourceTypesService) UpsertAction(ctx context.Context, key string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/resource-types/"+esc(key)+"/actions", input)
}
func (s ResourceTypesService) Mapping(ctx context.Context, key string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/resource-types/"+esc(key)+"/mapping", nil)
}
func (s ResourceTypesService) UpsertMapping(ctx context.Context, key string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/resource-types/"+esc(key)+"/mapping", input)
}

func (s ResourcesService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/resources", query)
}
func (s ResourcesService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/resources", input)
}
func (s ResourcesService) CreateInSpace(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/resources", input)
}
func (s ResourcesService) Get(ctx context.Context, resourceType, resourceID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/resources/"+esc(resourceType)+"/"+esc(resourceID), nil)
}
func (s ResourcesService) GetInSpace(ctx context.Context, spaceID, resourceID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/resources/"+esc(resourceID), nil)
}
func (s ResourcesService) Update(ctx context.Context, spaceID, resourceID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/resources/"+esc(resourceID), input)
}
func (s ResourcesService) Archive(ctx context.Context, spaceID, resourceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/resources/"+esc(resourceID)+"/archive", input)
}

func (s DataService) Tables(ctx context.Context) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/data/tables", nil)
}
func (s DataService) ListRows(ctx context.Context, resourceType string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/data/rows/"+esc(resourceType), query)
}
func (s DataService) GetRow(ctx context.Context, resourceType, resourceID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/data/rows/"+esc(resourceType)+"/"+esc(resourceID), nil)
}
func (s DataService) CreateRow(ctx context.Context, resourceType string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/data/rows/"+esc(resourceType), input)
}
func (s DataService) UpdateRow(ctx context.Context, resourceType, resourceID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/data/rows/"+esc(resourceType)+"/"+esc(resourceID), input)
}
func (s DataService) DeleteRow(ctx context.Context, resourceType, resourceID string, input Map) (Map, error) {
	return s.client.deleteMap(ctx, "/api/v1/data/rows/"+esc(resourceType)+"/"+esc(resourceID), input)
}

func (s PluginsService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins", query)
}
func (s PluginsService) Get(ctx context.Context, pluginID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/plugins/"+esc(pluginID), nil)
}
func (s PluginsService) ValidateManifest(ctx context.Context, manifest any) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/plugins/validate-manifest", manifest)
}
func (s PluginsService) Install(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/plugins/install", input)
}
func (s PluginsService) Enable(ctx context.Context, pluginID string) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/plugins/"+esc(pluginID)+"/enable", Map{})
}
func (s PluginsService) Disable(ctx context.Context, pluginID string) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/plugins/"+esc(pluginID)+"/disable", Map{})
}
func (s PluginsService) Uninstall(ctx context.Context, pluginID string) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/plugins/"+esc(pluginID)+"/uninstall", Map{})
}
func (s PluginsService) Resources(ctx context.Context, pluginID string) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins/"+esc(pluginID)+"/resources", nil)
}
func (s PluginsService) Permissions(ctx context.Context, pluginID string) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins/"+esc(pluginID)+"/permissions", nil)
}
func (s PluginsService) AuditEvents(ctx context.Context, pluginID string) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins/"+esc(pluginID)+"/audit-events", nil)
}
func (s PluginsService) AdminMenu(ctx context.Context, pluginID string) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins/"+esc(pluginID)+"/admin-menu", nil)
}
func (s PluginsService) Settings(ctx context.Context, pluginID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins/"+esc(pluginID)+"/settings", query)
}
func (s PluginsService) UpdateSettings(ctx context.Context, pluginID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/plugins/"+esc(pluginID)+"/settings", input)
}

func (s TemplatesService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/templates", query)
}
func (s TemplatesService) Get(ctx context.Context, templateID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/templates/"+esc(templateID), nil)
}
func (s TemplatesService) PreviewInstall(ctx context.Context, templateID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/templates/"+esc(templateID)+"/preview-install", input)
}
func (s TemplatesService) Install(ctx context.Context, templateID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/templates/"+esc(templateID)+"/install", input)
}

func (c *Client) getMap(ctx context.Context, path string, query Query) (Map, error) {
	var out Map
	err := c.do(ctx, http.MethodGet, withQuery(path, query), nil, &out)
	return out, err
}

func (c *Client) getList(ctx context.Context, path string, query Query) ([]Map, error) {
	var out []Map
	err := c.do(ctx, http.MethodGet, withQuery(path, query), nil, &out)
	return out, err
}

func (c *Client) postMap(ctx context.Context, path string, body any) (Map, error) {
	var out Map
	err := c.do(ctx, http.MethodPost, path, body, &out)
	return out, err
}

func (c *Client) patchMap(ctx context.Context, path string, body any) (Map, error) {
	var out Map
	err := c.do(ctx, http.MethodPatch, path, body, &out)
	return out, err
}

func (c *Client) deleteMap(ctx context.Context, path string, body any) (Map, error) {
	var out Map
	err := c.do(ctx, http.MethodDelete, path, body, &out)
	return out, err
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, reader)
	if err != nil {
		return err
	}
	for key, values := range c.DefaultHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if c.APIKey != "" {
		req.Header.Set("X-Plystra-API-Key", c.APIKey)
	}
	if c.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	if len(raw) == 0 && resp.StatusCode == http.StatusNoContent {
		return nil
	}

	var envelope struct {
		Data      json.RawMessage `json:"data"`
		Error     *ErrorBody      `json:"error"`
		RequestID string          `json:"request_id"`
		Meta      Map             `json:"meta"`
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&envelope); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       "INVALID_JSON_RESPONSE",
			Message:    fmt.Sprintf("Plystra returned non-JSON response with HTTP %d", resp.StatusCode),
			Details:    string(raw),
		}
	}
	if envelope.Error != nil || resp.StatusCode >= 400 {
		if envelope.Error == nil {
			return &APIError{StatusCode: resp.StatusCode, Code: "HTTP_ERROR", Message: http.StatusText(resp.StatusCode), RequestID: envelope.RequestID}
		}
		requestID := envelope.Error.RequestID
		if requestID == "" {
			requestID = envelope.RequestID
		}
		if requestID == "" {
			if value, ok := envelope.Meta["request_id"].(string); ok {
				requestID = value
			}
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       envelope.Error.Code,
			Message:    envelope.Error.Message,
			Details:    envelope.Error.Details,
			RequestID:  requestID,
			TraceID:    envelope.Error.TraceID,
			AuditLogID: envelope.Error.AuditLogID,
		}
	}
	if out == nil || len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return nil
	}
	return json.Unmarshal(envelope.Data, out)
}

func (c *Client) storeTokens(data Map) {
	if token, ok := data["access_token"].(string); ok {
		c.AccessToken = token
	}
	if token, ok := data["refresh_token"].(string); ok {
		c.RefreshToken = token
	}
}

func withQuery(path string, query Query) string {
	if len(query) == 0 {
		return path
	}
	return path + "?" + query.Encode()
}

func esc(value string) string {
	return url.PathEscape(value)
}
