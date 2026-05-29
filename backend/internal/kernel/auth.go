package kernel

// AuthContext holds authenticated user and tenant information extracted from
// JWT claims or the auth middleware. It is stored in Fiber's c.Locals("auth").
type AuthContext struct {
	TenantID TenantID
	UserID   UserID
	Email    string
	Role     string
}
