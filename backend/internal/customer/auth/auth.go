package auth

import (
	"net/http"
	"time"

	"github.com/Abraxas-365/vendex/internal/customer"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/golang-jwt/jwt/v5"
)

// CustomerCredentials holds the hashed password credentials for a customer.
type CustomerCredentials struct {
	ID           string            `json:"id" db:"id"`
	CustomerID   kernel.CustomerID `json:"customer_id" db:"customer_id"`
	TenantID     kernel.TenantID   `json:"tenant_id" db:"tenant_id"`
	Email        string            `json:"email" db:"email"`
	PasswordHash string            `json:"-" db:"password_hash"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
}

// RegisterInput holds the data needed to register a customer.
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone,omitempty"`
}

// LoginInput holds the data needed to log in a customer.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is returned on successful authentication.
type AuthResponse struct {
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	Customer     *customer.Customer `json:"customer"`
}

// CustomerClaims are the JWT claims for customer tokens.
type CustomerClaims struct {
	CustomerID kernel.CustomerID `json:"customer_id"`
	TenantID   kernel.TenantID   `json:"tenant_id"`
	Email      string            `json:"email"`
	Name       string            `json:"name"`
	jwt.RegisteredClaims
}

// CustomerAuthContext is stored in Fiber locals for authenticated customers.
type CustomerAuthContext struct {
	CustomerID kernel.CustomerID
	TenantID   kernel.TenantID
	Email      string
	Name       string
}

// ============================================================================
// Errors
// ============================================================================

var ErrRegistry = errx.NewRegistry("CUSTOMER_AUTH")

var (
	CodeInvalidCredentials = ErrRegistry.Register("INVALID_CREDENTIALS", errx.TypeAuthorization, http.StatusUnauthorized, "invalid email or password")
	CodeEmailTaken         = ErrRegistry.Register("EMAIL_TAKEN", errx.TypeConflict, http.StatusConflict, "email is already registered")
	CodeWeakPassword       = ErrRegistry.Register("WEAK_PASSWORD", errx.TypeValidation, http.StatusBadRequest, "password must be at least 8 characters")
	CodeCustomerNotFound   = ErrRegistry.Register("CUSTOMER_NOT_FOUND", errx.TypeNotFound, http.StatusNotFound, "customer not found")
	CodeTokenInvalid       = ErrRegistry.Register("TOKEN_INVALID", errx.TypeAuthorization, http.StatusUnauthorized, "invalid or expired token")
	CodeTokenRequired      = ErrRegistry.Register("TOKEN_REQUIRED", errx.TypeAuthorization, http.StatusUnauthorized, "authentication token required")
)

func ErrInvalidCredentials() *errx.Error  { return ErrRegistry.New(CodeInvalidCredentials) }
func ErrEmailTaken() *errx.Error          { return ErrRegistry.New(CodeEmailTaken) }
func ErrWeakPassword() *errx.Error        { return ErrRegistry.New(CodeWeakPassword) }
func ErrCustomerNotFound() *errx.Error    { return ErrRegistry.New(CodeCustomerNotFound) }
func ErrTokenInvalid() *errx.Error        { return ErrRegistry.New(CodeTokenInvalid) }
func ErrTokenRequired() *errx.Error       { return ErrRegistry.New(CodeTokenRequired) }
