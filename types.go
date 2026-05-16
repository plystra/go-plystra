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
	UserID        string `json:"user_id"`
	SpaceID       string `json:"space_id"`
	MemberID      string `json:"member_id"`
	UserMemberID  string `json:"user_member_id,omitempty"`
	BindingID     string `json:"binding_id,omitempty"`
	UserEmail     string `json:"user_email,omitempty"`
	UserStatus    string `json:"user_status,omitempty"`
	MemberName    string `json:"member_display_name,omitempty"`
	MemberStatus  string `json:"member_status,omitempty"`
	BindingStatus string `json:"binding_status,omitempty"`
	RelationType  string `json:"relation_type,omitempty"`
	SpaceName     string `json:"space_name,omitempty"`
	SpaceStatus   string `json:"space_status,omitempty"`
}

type AuthzResourceContext struct {
	Type          string `json:"type,omitempty"`
	ID            string `json:"id,omitempty"`
	ExternalID    string `json:"external_id,omitempty"`
	SpaceID       string `json:"space_id,omitempty"`
	GroupID       string `json:"group_id,omitempty"`
	GroupPath     string `json:"group_path,omitempty"`
	OwnerMemberID string `json:"owner_member_id,omitempty"`
	DisplayName   string `json:"display_name,omitempty"`
	Visibility    string `json:"visibility,omitempty"`
	Status        string `json:"status,omitempty"`
	Metadata      Map    `json:"metadata,omitempty"`
}

type AuthzGrantContext struct {
	RoleID               string `json:"role_id,omitempty"`
	RoleKey              string `json:"role_key,omitempty"`
	PermissionID         string `json:"permission_id,omitempty"`
	Resource             string `json:"resource,omitempty"`
	Action               string `json:"action,omitempty"`
	Scope                string `json:"scope"`
	SpaceID              string `json:"space_id,omitempty"`
	ScopeAnchorGroupID   string `json:"scope_anchor_group_id,omitempty"`
	ScopeAnchorGroupPath string `json:"scope_anchor_group_path,omitempty"`
}

type AuthzCheckInput struct {
	Actor        *ActorContext         `json:"actor"`
	ResourceType string                `json:"resource_type,omitempty"`
	ResourceID   string                `json:"resource_id,omitempty"`
	Resource     *AuthzResourceContext `json:"resource,omitempty"`
	Grants       []AuthzGrantContext   `json:"grants,omitempty"`
	Action       string                `json:"action"`
	Explain      bool                  `json:"explain,omitempty"`
}

type SystemService struct{ client *Client }
type AuthzService struct{ client *Client }
type AuditService struct{ client *Client }
type ResourceTypesService struct{ client *Client }
