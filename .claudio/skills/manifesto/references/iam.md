# IAM Packages

All IAM packages live under `internal/iam/`. Each is a self-contained bounded context.

---

## Auth (JWT)

Import: `github.com/Abraxas-365/manifesto/internal/iam/auth`

```go
// Generate token
token, err := auth.GenerateToken(ctx, claims, auth.WithExpiry(24*time.Hour))

// Validate token
claims, err := auth.ValidateToken(ctx, tokenString)

// Middleware (HTTP)
router.Use(auth.Middleware(authClient))
```

Claims carry `UserID` (`kernel.UserID`) and `TenantID` (`kernel.TenantID`).

---

## Tenant

Import: `github.com/Abraxas-365/manifesto/internal/iam/tenant`

```go
tenant, err := tenantSvc.Create(ctx, input)
tenant, err := tenantSvc.GetByID(ctx, tenantID)
err         := tenantSvc.Delete(ctx, tenantID)
```

`kernel.TenantID` is the canonical type for all tenant references across the system.

---

## User

Import: `github.com/Abraxas-365/manifesto/internal/iam/user`

```go
user, err := userSvc.Create(ctx, tenantID, input)
user, err := userSvc.GetByEmail(ctx, tenantID, email)
err       := userSvc.UpdatePassword(ctx, userID, newHash)
```

Always scope user queries by `TenantID`.

---

## API Keys

Import: `github.com/Abraxas-365/manifesto/internal/iam/apikey`

```go
key, err  := apikeySvc.Create(ctx, tenantID, userID, label)
valid, err := apikeySvc.Validate(ctx, rawKey)
err        = apikeySvc.Revoke(ctx, keyID)
```

---

## OTP

Import: `github.com/Abraxas-365/manifesto/internal/iam/otp`

```go
err        := otpSvc.Send(ctx, userID, channel) // channel: "email" | "sms"
valid, err := otpSvc.Verify(ctx, userID, code)
```

---

## Invitations

Import: `github.com/Abraxas-365/manifesto/internal/iam/invitation`

```go
inv, err := invitationSvc.Create(ctx, tenantID, email, role)
err       = invitationSvc.Accept(ctx, token, newUserInput)
err       = invitationSvc.Revoke(ctx, invitationID)
```

---

## Kernel value objects

Import: `github.com/Abraxas-365/hada-commerce/internal/kernel`

```go
tenantID := kernel.TenantID("t_abc123")
userID   := kernel.UserID("u_xyz789")
email, err := kernel.NewEmail("user@example.com") // validates format
```

Use these types throughout the codebase — never raw `string` for IDs or emails.
