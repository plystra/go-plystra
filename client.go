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

type contextKey string

const requestIDContextKey contextKey = "plystra_request_id"

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

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	value, _ := ctx.Value(requestIDContextKey).(string)
	return value
}

func NewClient(baseURL string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL:        strings.TrimRight(baseURL, "/"),
		HTTPClient:     &http.Client{Timeout: 10 * time.Second},
		DefaultHeaders: http.Header{},
		UserAgent:      "go-plystra/1.0.0-rc6",
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
	if requestID := RequestIDFromContext(ctx); requestID != "" && req.Header.Get("X-Request-ID") == "" {
		req.Header.Set("X-Request-ID", requestID)
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
