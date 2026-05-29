package kernel

// AuthContext carries authenticated identity for a request.
// It is injected by auth middleware and read by Fiber handlers via c.Locals("auth").
type AuthContext struct {
	TenantID TenantID
	UserID   UserID
}
