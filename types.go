package plystra

import (
	"fmt"
	"net/url"
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

type ActorContext struct {
	UserID       string `json:"user_id"`
	SpaceID      string `json:"space_id"`
	MemberID     string `json:"member_id"`
	UserMemberID string `json:"user_member_id"`
}

type AuthzCheckInput struct {
	Actor        *ActorContext `json:"actor,omitempty"`
	ResourceType string        `json:"resource_type,omitempty"`
	ResourceID   string        `json:"resource_id,omitempty"`
	Resource     *struct {
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
