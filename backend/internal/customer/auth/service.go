package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Abraxas-365/vendex/internal/customer"
	"github.com/Abraxas-365/vendex/internal/customer/customersrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost       = 12
	minPasswordLen   = 8
	defaultAccessTTL = 15 * time.Minute
	defaultRefreshTTL = 7 * 24 * time.Hour
	issuer           = "vendex"
)

// CustomerService is the interface for customer business operations needed by auth.
type CustomerService interface {
	Create(ctx context.Context, tenantID kernel.TenantID, in customersrv.CreateInput) (*customer.Customer, error)
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) (*customer.Customer, error)
	GetByEmail(ctx context.Context, tenantID kernel.TenantID, email kernel.Email) (*customer.Customer, error)
	Update(ctx context.Context, c *customer.Customer) error
}

// Service handles customer authentication business logic.
type Service struct {
	credRepo    CredentialsRepository
	customerSvc CustomerService
	jwtSecret   []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

// NewService creates a new customer auth service.
func NewService(
	credRepo CredentialsRepository,
	customerSvc CustomerService,
	jwtSecret string,
) *Service {
	return &Service{
		credRepo:    credRepo,
		customerSvc: customerSvc,
		jwtSecret:   []byte(jwtSecret),
		accessTTL:   defaultAccessTTL,
		refreshTTL:  defaultRefreshTTL,
	}
}

// Register creates a new customer account with password credentials.
func (s *Service) Register(ctx context.Context, tenantID kernel.TenantID, input RegisterInput) (*AuthResponse, error) {
	if input.Email == "" {
		return nil, errx.New("email is required", errx.TypeValidation)
	}
	if len(input.Password) < minPasswordLen {
		return nil, ErrWeakPassword()
	}
	if input.Name == "" {
		return nil, errx.New("name is required", errx.TypeValidation)
	}

	// Check if email already has credentials.
	existing, err := s.credRepo.GetByEmail(ctx, tenantID, input.Email)
	if err != nil && !errx.IsNotFound(err) {
		return nil, errx.Wrap(err, "checking email availability", errx.TypeInternal)
	}
	if existing != nil {
		return nil, ErrEmailTaken()
	}

	// Hash password.
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcryptCost)
	if err != nil {
		return nil, errx.Wrap(err, "hashing password", errx.TypeInternal)
	}

	// Create the customer entity.
	cust, err := s.customerSvc.Create(ctx, tenantID, customersrv.CreateInput{
		Email: input.Email,
		Name:  input.Name,
		Phone: input.Phone,
	})
	if err != nil {
		return nil, err
	}

	// Create credentials.
	now := time.Now()
	creds := &CustomerCredentials{
		ID:           uuid.NewString(),
		CustomerID:   cust.ID,
		TenantID:     tenantID,
		Email:        input.Email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.credRepo.Create(ctx, creds); err != nil {
		return nil, errx.Wrap(err, "saving customer credentials", errx.TypeInternal)
	}

	return s.buildAuthResponse(cust)
}

// Login authenticates a customer and returns tokens.
func (s *Service) Login(ctx context.Context, tenantID kernel.TenantID, input LoginInput) (*AuthResponse, error) {
	if input.Email == "" || input.Password == "" {
		return nil, ErrInvalidCredentials()
	}

	creds, err := s.credRepo.GetByEmail(ctx, tenantID, input.Email)
	if err != nil {
		if errx.IsNotFound(err) {
			return nil, ErrInvalidCredentials()
		}
		return nil, errx.Wrap(err, "fetching customer credentials", errx.TypeInternal)
	}

	if bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(input.Password)) != nil {
		return nil, ErrInvalidCredentials()
	}

	cust, err := s.customerSvc.GetByID(ctx, tenantID, creds.CustomerID)
	if err != nil {
		return nil, errx.Wrap(err, "fetching customer", errx.TypeInternal)
	}

	return s.buildAuthResponse(cust)
}

// RefreshToken validates a refresh token and returns a new token pair.
func (s *Service) RefreshToken(ctx context.Context, refreshTokenStr string) (*AuthResponse, error) {
	claims, err := s.parseRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, ErrTokenInvalid()
	}

	// Re-fetch the customer to ensure they still exist.
	cust, err := s.customerSvc.GetByID(ctx, claims.TenantID, claims.CustomerID)
	if err != nil {
		if errx.IsNotFound(err) {
			return nil, ErrCustomerNotFound()
		}
		return nil, errx.Wrap(err, "fetching customer for refresh", errx.TypeInternal)
	}

	return s.buildAuthResponse(cust)
}

// GetProfile returns the customer profile.
func (s *Service) GetProfile(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (*customer.Customer, error) {
	cust, err := s.customerSvc.GetByID(ctx, tenantID, customerID)
	if err != nil {
		if errx.IsNotFound(err) {
			return nil, ErrCustomerNotFound()
		}
		return nil, err
	}
	return cust, nil
}

// UpdateProfile updates name and phone for a customer.
func (s *Service) UpdateProfile(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, name, phone string) (*customer.Customer, error) {
	cust, err := s.customerSvc.GetByID(ctx, tenantID, customerID)
	if err != nil {
		if errx.IsNotFound(err) {
			return nil, ErrCustomerNotFound()
		}
		return nil, err
	}

	if name != "" {
		cust.Name = name
	}
	if phone != "" {
		cust.Phone = phone
	}

	if err := s.customerSvc.Update(ctx, cust); err != nil {
		return nil, errx.Wrap(err, "updating customer profile", errx.TypeInternal)
	}

	return cust, nil
}

// ChangePassword verifies the old password and sets a new one.
func (s *Service) ChangePassword(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, oldPassword, newPassword string) error {
	creds, err := s.credRepo.GetByCustomerID(ctx, tenantID, customerID)
	if err != nil {
		if errx.IsNotFound(err) {
			return ErrCustomerNotFound()
		}
		return errx.Wrap(err, "fetching customer credentials", errx.TypeInternal)
	}

	if bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(oldPassword)) != nil {
		return ErrInvalidCredentials()
	}

	if len(newPassword) < minPasswordLen {
		return ErrWeakPassword()
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcryptCost)
	if err != nil {
		return errx.Wrap(err, "hashing new password", errx.TypeInternal)
	}

	return s.credRepo.UpdatePassword(ctx, tenantID, customerID, string(hash))
}

// ============================================================================
// Token helpers
// ============================================================================

func (s *Service) generateAccessToken(cust *customer.Customer) (string, error) {
	now := time.Now()
	claims := CustomerClaims{
		CustomerID: cust.ID,
		TenantID:   cust.TenantID,
		Email:      string(cust.Email),
		Name:       cust.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   string(cust.ID),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) generateRefreshToken(cust *customer.Customer) (string, error) {
	now := time.Now()
	// Embed CustomerID and TenantID in the refresh token so we can re-fetch on refresh.
	claims := CustomerClaims{
		CustomerID: cust.ID,
		TenantID:   cust.TenantID,
		Email:      string(cust.Email),
		Name:       cust.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   string(cust.ID),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) parseRefreshToken(tokenStr string) (*CustomerClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomerClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid()
	}
	claims, ok := token.Claims.(*CustomerClaims)
	if !ok {
		return nil, ErrTokenInvalid()
	}
	return claims, nil
}

func (s *Service) buildAuthResponse(cust *customer.Customer) (*AuthResponse, error) {
	accessToken, err := s.generateAccessToken(cust)
	if err != nil {
		return nil, errx.Wrap(err, "generating access token", errx.TypeInternal)
	}
	refreshToken, err := s.generateRefreshToken(cust)
	if err != nil {
		return nil, errx.Wrap(err, "generating refresh token", errx.TypeInternal)
	}
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Customer:     cust,
	}, nil
}
