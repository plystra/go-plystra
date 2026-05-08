# Plystra Go SDK

Official Go SDK for Plystra Core v1.0.

Repository/module name: `plystra/go-plystra`

## Install

```bash
go get github.com/plystra/go-plystra
```

## Usage

```go
import plystra "github.com/plystra/go-plystra"

client := plystra.NewClient("http://localhost:8080")
_, _ = client.Auth.Login(ctx, "alice@example.com", "plystra-demo")
decision, err := client.Authz.Check(ctx, plystra.AuthzCheckInput{
    Actor: plystra.ActorContext{
        UserID: "user_alice",
        SpaceID: "space_acme",
        MemberID: "member_finance_reviewer",
        UserMemberID: "um_alice_finance_reviewer",
    },
    ResourceType: "invoice",
    ResourceID: "invoice_001",
    Action: "approve",
})
```

Non-public endpoints require a Bearer session whose user has an active admin grant.
