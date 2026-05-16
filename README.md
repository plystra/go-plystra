# Plystra Go SDK

Official Go SDK for Plystra Core v1.0.

Repository/module name: `plystra/go-plystra`

## Install

```bash
go get github.com/plystra/go-plystra
```

## Usage

For production services, prefer a scoped API key or a server-side access token issued by your frontend/gateway session. Keep password login for admin tools and bootstrap flows.

```go
client := plystra.NewClient(
	"https://plystra.internal",
	plystra.WithAPIKey(os.Getenv("PLYSTRA_API_KEY")),
)
```

Attach an application request id to every call through the Go context:

```go
ctx = plystra.WithRequestID(ctx, "req_01HY...")
decision, err := client.Authz.Check(ctx, plystra.AuthzCheckInput{
	ResourceType: "invoice",
	ResourceID:   "invoice_001",
	Action:       "approve",
})
```

`Authz.Check` may omit `Actor` when using an access token; Core uses the token's active actor. API key calls must pass `Actor` explicitly.

Context Mode lets a trusted backend protect one existing business action without syncing users, organizations, or resources first:

```go
decision, err := client.Authz.Check(ctx, plystra.AuthzCheckInput{
	Actor: &plystra.ActorContext{
		UserID:    "user_external_alice",
		MemberID:  "member_finance_reviewer",
		BindingID: "binding_external_alice_finance",
		SpaceID:   "space_acme",
	},
	Resource: &plystra.AuthzResourceContext{
		Type:       "invoice",
		ExternalID: "invoice_001",
		SpaceID:    "space_acme",
		GroupPath:  "finance.apac",
	},
	Grants: []plystra.AuthzGrantContext{{
		RoleKey:              "finance_approver",
		Resource:             "invoice",
		Action:               "approve",
		Scope:                "group_tree",
		SpaceID:              "space_acme",
		ScopeAnchorGroupPath: "finance",
	}},
	Action: "approve",
})
```

Inline actor, resource, and grant context is API-key-only. Build it from trusted server-side session and database state, never directly from browser input.

```go
import plystra "github.com/plystra/go-plystra"

client := plystra.NewClient("http://localhost:8080")
_, _ = client.Auth.Login(ctx, "alice@example.com", "plystra-demo")
_, _ = client.Auth.Refresh(ctx, "") // Uses the stored refresh token and persists the rotated token pair.
decision, err := client.Authz.Check(ctx, plystra.AuthzCheckInput{
	Actor: &plystra.ActorContext{
		UserID:       "user_alice",
		SpaceID:      "space_acme",
		MemberID:     "member_finance_reviewer",
		UserMemberID: "um_alice_finance_reviewer",
	},
	ResourceType: "invoice",
	ResourceID:   "invoice_001",
	Action:       "approve",
})
```

Non-public endpoints require either a Bearer session whose user has an active admin grant or a scoped API key with matching permission keys.

`Data`, `Plugins`, and `Templates` clients wrap preview Core metadata surfaces. Data Console is disabled by default, plugin routes do not represent a stable plugin runtime, and template routes do not represent a stable template ecosystem.

Core rotates refresh tokens. Keep `client.Tokens()` in your server-side encrypted session store after `Login` and `Refresh`; pass the stored values back with `WithAccessToken` and `WithRefreshToken` when creating a client for the next request.

Registration is disabled by default in Core. When an operator enables it and provides a registration token, `Auth.Register` creates the user/session and stores the returned token pair:

```go
_, err := client.Auth.Register(ctx, plystra.Map{
	"email":              "founder@example.com",
	"password":           "long-enough-password",
	"space_name":         "Founder Space",
	"registration_token": registrationToken,
})
```
