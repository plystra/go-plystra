# Plystra Go SDK

Official Go SDK for the Plystra Kernel Phase 1 API.

Repository/module name: `plystra/go-plystra`

## Install

```bash
go get github.com/plystra/go-plystra
```

## Usage

Phase 1 is Context Mode: your trusted backend keeps its existing users, organizations, roles, and business rows, then calls Plystra to protect one action.

```go
client := plystra.NewClient(
	"https://plystra.internal",
	plystra.WithAPIKey(os.Getenv("PLYSTRA_API_KEY")),
)

decision, err := client.Authz.Check(ctx, plystra.AuthzCheckInput{
	Actor: &plystra.ActorContext{
		UserID:    "user_external_alice",
		MemberID:  "member_finance_reviewer",
		BindingID: "binding_external_alice_finance",
		SpaceID:   "space_acme",
	},
	Resource: &plystra.AuthzResourceContext{
		Type:          "invoice",
		ExternalID:    "invoice_001",
		SpaceID:       "space_acme",
		GroupPath:     "finance.apac",
		OwnerMemberID: "member_invoice_creator",
	},
	Grants: []plystra.AuthzGrantContext{{
		RoleKey:              "finance_approver",
		Resource:             "invoice",
		Action:               "approve",
		Scope:                "group_tree",
		SpaceID:              "space_acme",
		ScopeAnchorGroupPath: "finance",
	}},
	Action:  "approve",
	Explain: true,
})
```

Inline actor, resource, and grant context is trusted server-side input. Build it from your authenticated session and database state, never directly from browser-submitted JSON.

## Kernel Surfaces

- `System.Health`, `System.Ready`, `System.Version`
- `System.Capabilities`
- `ResourceTypes.List`
- `Authz.Check` and `Authz.Explain`
- `Audit.List` and `Audit.Get`
- `Client.Request` for low-level calls

Attach a correlation id to a group of calls:

```go
ctx = plystra.WithRequestID(ctx, "req_01HY...")
_, err := client.Authz.Explain(ctx, contextModeRequest)
```

Protected routes require a scoped server API key. Public health, readiness, and version checks do not send the key.
