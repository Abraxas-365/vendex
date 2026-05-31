package multistore

import (
	"context"
	"strings"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// TenantResolver resolves a tenant ID from a request's Host header.
// Resolution order:
//  1. X-Tenant-ID header (explicit override, useful for dev/testing)
//  2. Subdomain: {slug}.{baseDomain} → lookup tenant by slug
//  3. Custom domain: full host → lookup tenant_domains table
type TenantResolver struct {
	baseDomain string
	repo       TenantDomainRepo
}

// TenantDomainRepo is the minimal interface needed for domain resolution.
type TenantDomainRepo interface {
	GetTenantBySlug(ctx context.Context, slug string) (kernel.TenantID, error)
	GetTenantByDomain(ctx context.Context, domain string) (kernel.TenantID, error)
}

// NewTenantResolver creates a resolver with the given base domain (e.g. "vendex.ai").
func NewTenantResolver(baseDomain string, repo TenantDomainRepo) *TenantResolver {
	return &TenantResolver{baseDomain: baseDomain, repo: repo}
}

// Resolve determines the tenant ID from the host string.
// Returns empty TenantID if resolution fails.
func (r *TenantResolver) Resolve(ctx context.Context, host, headerTenantID string) (kernel.TenantID, error) {
	// 1. Explicit header takes priority (dev/testing)
	if headerTenantID != "" {
		return kernel.TenantID(headerTenantID), nil
	}

	// Strip port from host
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// 2. Check if it's a subdomain of the base domain
	suffix := "." + r.baseDomain
	if strings.HasSuffix(host, suffix) {
		slug := strings.TrimSuffix(host, suffix)
		if slug != "" && !strings.Contains(slug, ".") {
			return r.repo.GetTenantBySlug(ctx, slug)
		}
	}

	// 3. Custom domain lookup
	if host != "localhost" && host != "127.0.0.1" && host != r.baseDomain {
		return r.repo.GetTenantByDomain(ctx, host)
	}

	return "", ErrTenantNotResolved
}
