package storefrontsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/storefront"
)

// Service implements the storefront business logic.
type Service struct {
	pages    storefront.PageRepository
	versions storefront.PageVersionRepository
}

// New creates a new storefront Service.
func New(pages storefront.PageRepository, versions storefront.PageVersionRepository) *Service {
	return &Service{pages: pages, versions: versions}
}

// CreatePageInput holds the data needed to create a new page.
type CreatePageInput struct {
	TenantID  kernel.TenantID
	Slug      string
	Title     string
	HTML      string
	CSS       string
	Meta      storefront.PageMeta
	CreatedBy string
	// ByAgent — when true the page starts as pending_review instead of draft.
	ByAgent bool
}

// CreatePage persists a new page. Agent-created pages always land in pending_review.
func (s *Service) CreatePage(ctx context.Context, input CreatePageInput) (*storefront.Page, error) {
	status := storefront.PageStatusDraft
	if input.ByAgent {
		status = storefront.PageStatusPendingReview
	}

	now := time.Now().UTC()
	page := &storefront.Page{
		ID:        kernel.PageID(newID()),
		TenantID:  input.TenantID,
		Slug:      input.Slug,
		Title:     input.Title,
		HTML:      input.HTML,
		CSS:       input.CSS,
		Meta:      input.Meta,
		Status:    status,
		Version:   1,
		CreatedBy: input.CreatedBy,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.pages.Create(ctx, page); err != nil {
		return nil, errx.Wrap(err, "create page", errx.TypeInternal)
	}

	// Snapshot version 1.
	v := &storefront.PageVersion{
		ID:        kernel.PageVersionID(newID()),
		PageID:    page.ID,
		TenantID:  page.TenantID,
		Version:   1,
		HTML:      page.HTML,
		CSS:       page.CSS,
		EditedBy:  input.CreatedBy,
		Comment:   "initial version",
		CreatedAt: now,
	}
	if err := s.versions.Create(ctx, v); err != nil {
		return nil, errx.Wrap(err, "create initial version", errx.TypeInternal)
	}

	return page, nil
}

// UpdatePageInput holds the fields that can change on an edit.
type UpdatePageInput struct {
	TenantID kernel.TenantID
	ID       kernel.PageID
	Title    *string
	HTML     *string
	CSS      *string
	Meta     *storefront.PageMeta
	EditedBy string
	Comment  string
}

// UpdatePage applies a content edit, bumps the version counter, and saves a snapshot.
func (s *Service) UpdatePage(ctx context.Context, input UpdatePageInput) (*storefront.Page, error) {
	page, err := s.pages.GetByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}
	if !page.CanBeEdited() {
		return nil, storefront.ErrPageArchived
	}

	if input.Title != nil {
		page.Title = *input.Title
	}
	if input.HTML != nil {
		page.HTML = *input.HTML
	}
	if input.CSS != nil {
		page.CSS = *input.CSS
	}
	if input.Meta != nil {
		page.Meta = *input.Meta
	}

	page.Version++
	page.UpdatedAt = time.Now().UTC()

	if err := s.pages.Update(ctx, page); err != nil {
		return nil, errx.Wrap(err, "update page", errx.TypeInternal)
	}

	// Snapshot the new version.
	v := &storefront.PageVersion{
		ID:        kernel.PageVersionID(newID()),
		PageID:    page.ID,
		TenantID:  page.TenantID,
		Version:   page.Version,
		HTML:      page.HTML,
		CSS:       page.CSS,
		EditedBy:  input.EditedBy,
		Comment:   input.Comment,
		CreatedAt: page.UpdatedAt,
	}
	if err := s.versions.Create(ctx, v); err != nil {
		return nil, errx.Wrap(err, "create version snapshot", errx.TypeInternal)
	}

	return page, nil
}

// Publish transitions a page to published status.
func (s *Service) Publish(ctx context.Context, tenantID kernel.TenantID, id kernel.PageID) (*storefront.Page, error) {
	page, err := s.pages.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if !page.CanBePublished() {
		return nil, storefront.ErrInvalidStatus
	}

	now := time.Now().UTC()
	page.Status = storefront.PageStatusPublished
	page.PublishedAt = &now
	page.UpdatedAt = now

	if err := s.pages.Update(ctx, page); err != nil {
		return nil, errx.Wrap(err, "publish page", errx.TypeInternal)
	}
	return page, nil
}

// Unpublish transitions a published page back to draft.
func (s *Service) Unpublish(ctx context.Context, tenantID kernel.TenantID, id kernel.PageID) (*storefront.Page, error) {
	page, err := s.pages.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if page.Status != storefront.PageStatusPublished {
		return nil, storefront.ErrInvalidStatus
	}

	page.Status = storefront.PageStatusDraft
	page.UpdatedAt = time.Now().UTC()

	if err := s.pages.Update(ctx, page); err != nil {
		return nil, errx.Wrap(err, "unpublish page", errx.TypeInternal)
	}
	return page, nil
}

// Archive permanently archives a page. Archived pages cannot be edited or published.
func (s *Service) Archive(ctx context.Context, tenantID kernel.TenantID, id kernel.PageID) (*storefront.Page, error) {
	page, err := s.pages.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if page.Status == storefront.PageStatusArchived {
		return nil, storefront.ErrInvalidStatus
	}

	page.Status = storefront.PageStatusArchived
	page.UpdatedAt = time.Now().UTC()

	if err := s.pages.Update(ctx, page); err != nil {
		return nil, errx.Wrap(err, "archive page", errx.TypeInternal)
	}
	return page, nil
}

// ReviseFromFeedback applies reviewer feedback, bumps the version, and puts the page
// back into pending_review so it can be approved again.
func (s *Service) ReviseFromFeedback(ctx context.Context, input UpdatePageInput) (*storefront.Page, error) {
	page, err := s.UpdatePage(ctx, input)
	if err != nil {
		return nil, err
	}

	// After agent revises based on feedback, return to pending_review.
	if page.Status == storefront.PageStatusDraft {
		page.Status = storefront.PageStatusPendingReview
		page.UpdatedAt = time.Now().UTC()
		if err := s.pages.Update(ctx, page); err != nil {
			return nil, errx.Wrap(err, "set pending_review after revision", errx.TypeInternal)
		}
	}
	return page, nil
}

// GetPublishedPage returns a published page by slug — used for public rendering.
func (s *Service) GetPublishedPage(ctx context.Context, tenantID kernel.TenantID, slug string) (*storefront.Page, error) {
	return s.pages.GetPublished(ctx, tenantID, slug)
}

// GetPage returns any page by ID for admin use.
func (s *Service) GetPage(ctx context.Context, tenantID kernel.TenantID, id kernel.PageID) (*storefront.Page, error) {
	return s.pages.GetByID(ctx, tenantID, id)
}

// ListPages returns pages with optional status filter.
func (s *Service) ListPages(ctx context.Context, tenantID kernel.TenantID, status *storefront.PageStatus, p kernel.PaginationOptions) (kernel.Paginated[storefront.Page], error) {
	if status != nil {
		return s.pages.ListByStatus(ctx, tenantID, *status, p)
	}
	return s.pages.List(ctx, tenantID, p)
}

// ListVersions returns all version snapshots for a page.
func (s *Service) ListVersions(ctx context.Context, tenantID kernel.TenantID, pageID kernel.PageID) ([]storefront.PageVersion, error) {
	return s.versions.ListByPage(ctx, tenantID, pageID)
}

// GetVersion retrieves a specific version snapshot.
func (s *Service) GetVersion(ctx context.Context, tenantID kernel.TenantID, pageID kernel.PageID, version int) (*storefront.PageVersion, error) {
	return s.versions.GetByVersion(ctx, tenantID, pageID, version)
}

// ResolveTemplateTags parses and resolves dynamic template tags embedded in page HTML.
// Tag syntax example: {{products "featured" limit=8}}
// TODO: implement tag parsing and delegation to product/promo/catalog resolvers.
func (s *Service) ResolveTemplateTags(ctx context.Context, tenantID kernel.TenantID, html string) (string, []storefront.TemplateTag, error) {
	// TODO: walk html, extract {{...}} tags, build TemplateTag structs, call resolvers.
	return html, nil, nil
}

// newID generates a new UUID-like unique string identifier.
func newID() string {
	return generateUUID()
}
