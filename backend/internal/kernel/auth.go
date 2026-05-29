package kernel

// AuthContext carries identity information extracted from a JWT token.
// It is stored in Fiber's Locals under the "auth" key by the auth middleware.
type AuthContext struct {
	UserID   UserID
	TenantID TenantID
	Email    string
	Role     string
}
