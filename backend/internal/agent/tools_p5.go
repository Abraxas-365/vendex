package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Abraxas-365/hada-commerce/internal/abtest"
	"github.com/Abraxas-365/hada-commerce/internal/abtest/abtestsrv"
	"github.com/Abraxas-365/hada-commerce/internal/blog"
	"github.com/Abraxas-365/hada-commerce/internal/blog/blogsrv"
	"github.com/Abraxas-365/hada-commerce/internal/bulkops"
	"github.com/Abraxas-365/hada-commerce/internal/bulkops/bulkopssrv"
	"github.com/Abraxas-365/hada-commerce/internal/collection"
	"github.com/Abraxas-365/hada-commerce/internal/collection/collectionsrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/multistore"
	"github.com/Abraxas-365/hada-commerce/internal/multistore/multistoresrv"
	"github.com/Abraxas-365/hada-commerce/internal/recommendation"
	"github.com/Abraxas-365/hada-commerce/internal/recommendation/recommendationsrv"
)

// Compile-time guards
var (
	_ Tool = (*ListStorefrontsTool)(nil)
	_ Tool = (*CreateStorefrontTool)(nil)
	_ Tool = (*ListBulkOperationsTool)(nil)
	_ Tool = (*CreateBulkOperationTool)(nil)
	_ Tool = (*ProcessBulkOperationTool)(nil)
	_ Tool = (*ListBlogPostsTool)(nil)
	_ Tool = (*CreateBlogPostTool)(nil)
	_ Tool = (*PublishBlogPostTool)(nil)
	_ Tool = (*ListBlogCategoriesTool)(nil)
	_ Tool = (*ListCollectionsTool)(nil)
	_ Tool = (*CreateCollectionTool)(nil)
	_ Tool = (*AddCollectionProductTool)(nil)
	_ Tool = (*ListExperimentsTool)(nil)
	_ Tool = (*CreateExperimentTool)(nil)
	_ Tool = (*GetExperimentResultsTool)(nil)
	_ Tool = (*ListRecommendationRulesTool)(nil)
	_ Tool = (*CreateRecommendationRuleTool)(nil)
	_ Tool = (*GetTrendingProductsTool)(nil)
)

// ─── Multi-Storefront ───

type ListStorefrontsTool struct {
	multistore *multistoresrv.Service
	tenantID   kernel.TenantID
}

func (t *ListStorefrontsTool) Name() string        { return "list_storefronts" }
func (t *ListStorefrontsTool) Description() string  { return "List all storefronts for the tenant" }
func (t *ListStorefrontsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer", "description": "Page number (default 1)"},
			"page_size": map[string]any{"type": "integer", "description": "Items per page (default 20)"},
		},
	}
}
func (t *ListStorefrontsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_storefronts: unmarshal input", errx.TypeValidation)
	}
	if in.Page < 1 { in.Page = 1 }
	if in.PageSize < 1 { in.PageSize = 20 }
	result, err := t.multistore.List(ctx, t.tenantID, in.Page, in.PageSize)
	if err != nil {
		return "", errx.Wrap(err, "list_storefronts", errx.TypeInternal)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d storefronts (page %d/%d):\n", result.Total, result.Page, result.TotalPages)
	for _, s := range result.Items {
		def := ""
		if s.IsDefault { def = " [DEFAULT]" }
		domain := ""
		if s.Domain != nil { domain = *s.Domain }
		fmt.Fprintf(&sb, "- %s (slug: %s, domain: %s, locale: %s, currency: %s, active: %v)%s\n", s.Name, s.Slug, domain, s.DefaultLocale, s.DefaultCurrency, s.IsActive, def)
	}
	return sb.String(), nil
}

type CreateStorefrontTool struct {
	multistore *multistoresrv.Service
	tenantID   kernel.TenantID
}

func (t *CreateStorefrontTool) Name() string        { return "create_storefront" }
func (t *CreateStorefrontTool) Description() string  { return "Create a new storefront" }
func (t *CreateStorefrontTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":     map[string]any{"type": "string"},
			"slug":     map[string]any{"type": "string"},
			"domain":   map[string]any{"type": "string"},
			"locale":   map[string]any{"type": "string"},
			"currency": map[string]any{"type": "string"},
		},
		"required": []string{"name", "slug"},
	}
}
func (t *CreateStorefrontTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in multistore.CreateInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_storefront: unmarshal input", errx.TypeValidation)
	}
	sf, err := t.multistore.Create(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "create_storefront", errx.TypeInternal)
	}
	return fmt.Sprintf("Created storefront %q (ID: %s, slug: %s)", sf.Name, sf.ID, sf.Slug), nil
}

// ─── Bulk Operations ───

type ListBulkOperationsTool struct {
	bulkops  *bulkopssrv.Service
	tenantID kernel.TenantID
}

func (t *ListBulkOperationsTool) Name() string        { return "list_bulk_operations" }
func (t *ListBulkOperationsTool) Description() string  { return "List bulk operations" }
func (t *ListBulkOperationsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer"},
			"page_size": map[string]any{"type": "integer"},
		},
	}
}
func (t *ListBulkOperationsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_bulk_operations: unmarshal input", errx.TypeValidation)
	}
	if in.Page < 1 { in.Page = 1 }
	if in.PageSize < 1 { in.PageSize = 20 }
	result, err := t.bulkops.List(ctx, t.tenantID, in.Page, in.PageSize)
	if err != nil {
		return "", errx.Wrap(err, "list_bulk_operations", errx.TypeInternal)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d bulk operations:\n", result.Total)
	for _, op := range result.Items {
		fmt.Fprintf(&sb, "- %s: type=%s status=%s total=%d success=%d failed=%d\n", op.ID, op.Type, op.Status, op.TotalItems, op.ProcessedItems, op.FailedItems)
	}
	return sb.String(), nil
}

type CreateBulkOperationTool struct {
	bulkops  *bulkopssrv.Service
	tenantID kernel.TenantID
}

func (t *CreateBulkOperationTool) Name() string        { return "create_bulk_operation" }
func (t *CreateBulkOperationTool) Description() string  { return "Create a bulk operation (price_update, status_change, tag_add, tag_remove, delete)" }
func (t *CreateBulkOperationTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"operation_type": map[string]any{"type": "string", "enum": []string{"price_update", "status_change", "tag_add", "tag_remove", "delete"}},
			"resource_type":  map[string]any{"type": "string"},
			"resource_ids":   map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"parameters":     map[string]any{"type": "object"},
		},
		"required": []string{"operation_type", "resource_type", "resource_ids"},
	}
}
func (t *CreateBulkOperationTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in bulkops.CreateInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_bulk_operation: unmarshal input", errx.TypeValidation)
	}
	op, err := t.bulkops.Create(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "create_bulk_operation", errx.TypeInternal)
	}
	return fmt.Sprintf("Created bulk operation %s (type: %s, %d items)", op.ID, op.Type, op.TotalItems), nil
}

type ProcessBulkOperationTool struct {
	bulkops  *bulkopssrv.Service
	tenantID kernel.TenantID
}

func (t *ProcessBulkOperationTool) Name() string        { return "process_bulk_operation" }
func (t *ProcessBulkOperationTool) Description() string  { return "Process a pending bulk operation" }
func (t *ProcessBulkOperationTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id": map[string]any{"type": "string", "description": "Bulk operation ID"},
		},
		"required": []string{"id"},
	}
}
func (t *ProcessBulkOperationTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct{ ID string `json:"id"` }
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "process_bulk_operation: unmarshal input", errx.TypeValidation)
	}
	op, err := t.bulkops.Process(ctx, t.tenantID, kernel.BulkOperationID(in.ID))
	if err != nil {
		return "", errx.Wrap(err, "process_bulk_operation", errx.TypeInternal)
	}
	return fmt.Sprintf("Processed bulk operation %s: %d success, %d failed", op.ID, op.ProcessedItems, op.FailedItems), nil
}

// ─── Blog ───

type ListBlogPostsTool struct {
	blog     *blogsrv.Service
	tenantID kernel.TenantID
}

func (t *ListBlogPostsTool) Name() string        { return "list_blog_posts" }
func (t *ListBlogPostsTool) Description() string  { return "List blog posts with optional status filter" }
func (t *ListBlogPostsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"status":    map[string]any{"type": "string", "enum": []string{"draft", "published", "archived"}},
			"page":      map[string]any{"type": "integer"},
			"page_size": map[string]any{"type": "integer"},
		},
	}
}
func (t *ListBlogPostsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Status   string `json:"status"`
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_blog_posts: unmarshal input", errx.TypeValidation)
	}
	if in.Page < 1 { in.Page = 1 }
	if in.PageSize < 1 { in.PageSize = 20 }
	filter := blog.ListPostsFilter{Status: in.Status, Page: in.Page, PageSize: in.PageSize}
	result, err := t.blog.ListPosts(ctx, t.tenantID, filter)
	if err != nil {
		return "", errx.Wrap(err, "list_blog_posts", errx.TypeInternal)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d blog posts:\n", result.Total)
	for _, p := range result.Items {
		fmt.Fprintf(&sb, "- %s (slug: %s, status: %s, author: %s)\n", p.Title, p.Slug, p.Status, p.AuthorName)
	}
	return sb.String(), nil
}

type CreateBlogPostTool struct {
	blog     *blogsrv.Service
	tenantID kernel.TenantID
}

func (t *CreateBlogPostTool) Name() string        { return "create_blog_post" }
func (t *CreateBlogPostTool) Description() string  { return "Create a new blog post" }
func (t *CreateBlogPostTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title":          map[string]any{"type": "string"},
			"slug":           map[string]any{"type": "string"},
			"content":        map[string]any{"type": "string"},
			"excerpt":        map[string]any{"type": "string"},
			"author":         map[string]any{"type": "string"},
			"featured_image": map[string]any{"type": "string"},
			"tags":           map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		},
		"required": []string{"title", "slug", "content", "author"},
	}
}
func (t *CreateBlogPostTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in blog.CreatePostInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_blog_post: unmarshal input", errx.TypeValidation)
	}
	post, err := t.blog.CreatePost(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "create_blog_post", errx.TypeInternal)
	}
	return fmt.Sprintf("Created blog post %q (ID: %s, slug: %s, status: %s)", post.Title, post.ID, post.Slug, post.Status), nil
}

type PublishBlogPostTool struct {
	blog     *blogsrv.Service
	tenantID kernel.TenantID
}

func (t *PublishBlogPostTool) Name() string        { return "publish_blog_post" }
func (t *PublishBlogPostTool) Description() string  { return "Publish a draft blog post" }
func (t *PublishBlogPostTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id": map[string]any{"type": "string", "description": "Blog post ID"},
		},
		"required": []string{"id"},
	}
}
func (t *PublishBlogPostTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct{ ID string `json:"id"` }
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "publish_blog_post: unmarshal input", errx.TypeValidation)
	}
	post, err := t.blog.PublishPost(ctx, t.tenantID, kernel.BlogPostID(in.ID))
	if err != nil {
		return "", errx.Wrap(err, "publish_blog_post", errx.TypeInternal)
	}
	return fmt.Sprintf("Published blog post %q (ID: %s)", post.Title, post.ID), nil
}

type ListBlogCategoriesTool struct {
	blog     *blogsrv.Service
	tenantID kernel.TenantID
}

func (t *ListBlogCategoriesTool) Name() string        { return "list_blog_categories" }
func (t *ListBlogCategoriesTool) Description() string  { return "List all blog categories" }
func (t *ListBlogCategoriesTool) InputSchema() map[string]any {
	return map[string]any{"type": "object", "properties": map[string]any{}}
}
func (t *ListBlogCategoriesTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	cats, err := t.blog.ListCategories(ctx, t.tenantID)
	if err != nil {
		return "", errx.Wrap(err, "list_blog_categories", errx.TypeInternal)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d blog categories:\n", len(cats))
	for _, c := range cats {
		fmt.Fprintf(&sb, "- %s (slug: %s, ID: %s)\n", c.Name, c.Slug, c.ID)
	}
	return sb.String(), nil
}

// ─── Collections ───

type ListCollectionsTool struct {
	collections *collectionsrv.Service
	tenantID    kernel.TenantID
}

func (t *ListCollectionsTool) Name() string        { return "list_collections" }
func (t *ListCollectionsTool) Description() string  { return "List product collections" }
func (t *ListCollectionsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"active_only": map[string]any{"type": "boolean"},
			"page":        map[string]any{"type": "integer"},
			"page_size":   map[string]any{"type": "integer"},
		},
	}
}
func (t *ListCollectionsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		ActiveOnly bool `json:"active_only"`
		Page       int  `json:"page"`
		PageSize   int  `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_collections: unmarshal input", errx.TypeValidation)
	}
	if in.Page < 1 { in.Page = 1 }
	if in.PageSize < 1 { in.PageSize = 20 }
	result, err := t.collections.List(ctx, t.tenantID, in.ActiveOnly, in.Page, in.PageSize)
	if err != nil {
		return "", errx.Wrap(err, "list_collections", errx.TypeInternal)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d collections:\n", result.Total)
	for _, c := range result.Items {
		fmt.Fprintf(&sb, "- %s (slug: %s, type: %s, active: %v)\n", c.Name, c.Slug, c.Type, c.IsActive)
	}
	return sb.String(), nil
}

type CreateCollectionTool struct {
	collections *collectionsrv.Service
	tenantID    kernel.TenantID
}

func (t *CreateCollectionTool) Name() string        { return "create_collection" }
func (t *CreateCollectionTool) Description() string  { return "Create a new product collection" }
func (t *CreateCollectionTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":        map[string]any{"type": "string"},
			"slug":        map[string]any{"type": "string"},
			"description": map[string]any{"type": "string"},
			"type":        map[string]any{"type": "string", "enum": []string{"manual", "auto"}},
			"image_url":   map[string]any{"type": "string"},
			"sort_order":  map[string]any{"type": "integer"},
		},
		"required": []string{"name", "slug"},
	}
}
func (t *CreateCollectionTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Type        string `json:"type"`
		ImageURL    string `json:"image_url"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_collection: unmarshal input", errx.TypeValidation)
	}
	collType := collection.CollectionType(in.Type)
	if collType == "" {
		collType = "manual"
	}
	result, err := t.collections.Create(ctx, t.tenantID, collection.CreateInput{
		Name:        in.Name,
		Slug:        in.Slug,
		Description: in.Description,
		Type:        collType,
		ImageURL:    in.ImageURL,
		SortOrder:   in.SortOrder,
		IsActive:    true,
	})
	if err != nil {
		return "", errx.Wrap(err, "create_collection", errx.TypeInternal)
	}
	return fmt.Sprintf("Created collection %q (ID: %s, slug: %s, type: %s)", result.Name, result.ID, result.Slug, result.Type), nil
}

type AddCollectionProductTool struct {
	collections *collectionsrv.Service
	tenantID    kernel.TenantID
}

func (t *AddCollectionProductTool) Name() string        { return "add_collection_product" }
func (t *AddCollectionProductTool) Description() string  { return "Add a product to a collection" }
func (t *AddCollectionProductTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"collection_id": map[string]any{"type": "string"},
			"product_id":    map[string]any{"type": "string"},
			"sort_order":    map[string]any{"type": "integer"},
		},
		"required": []string{"collection_id", "product_id"},
	}
}
func (t *AddCollectionProductTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		CollectionID string `json:"collection_id"`
		ProductID    string `json:"product_id"`
		SortOrder    int    `json:"sort_order"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "add_collection_product: unmarshal input", errx.TypeValidation)
	}
	cp, err := t.collections.AddProduct(ctx, t.tenantID, in.CollectionID, in.ProductID, in.SortOrder)
	if err != nil {
		return "", errx.Wrap(err, "add_collection_product", errx.TypeInternal)
	}
	return fmt.Sprintf("Added product %s to collection %s (sort: %d)", cp.ProductID, cp.CollectionID, cp.SortOrder), nil
}

// ─── A/B Testing ───

type ListExperimentsTool struct {
	abtest   *abtestsrv.Service
	tenantID kernel.TenantID
}

func (t *ListExperimentsTool) Name() string        { return "list_experiments" }
func (t *ListExperimentsTool) Description() string  { return "List A/B test experiments" }
func (t *ListExperimentsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"status":    map[string]any{"type": "string", "enum": []string{"draft", "running", "paused", "completed"}},
			"page":      map[string]any{"type": "integer"},
			"page_size": map[string]any{"type": "integer"},
		},
	}
}
func (t *ListExperimentsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Status   string `json:"status"`
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_experiments: unmarshal input", errx.TypeValidation)
	}
	if in.Page < 1 { in.Page = 1 }
	if in.PageSize < 1 { in.PageSize = 20 }
	result, err := t.abtest.ListExperiments(ctx, t.tenantID, in.Status, in.Page, in.PageSize)
	if err != nil {
		return "", errx.Wrap(err, "list_experiments", errx.TypeInternal)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d experiments:\n", result.Total)
	for _, e := range result.Items {
		fmt.Fprintf(&sb, "- %s (status: %s, type: %s, ID: %s)\n", e.Name, e.Status, e.Type, e.ID)
	}
	return sb.String(), nil
}

type CreateExperimentTool struct {
	abtest   *abtestsrv.Service
	tenantID kernel.TenantID
}

func (t *CreateExperimentTool) Name() string        { return "create_experiment" }
func (t *CreateExperimentTool) Description() string  { return "Create an A/B test experiment" }
func (t *CreateExperimentTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":        map[string]any{"type": "string"},
			"description": map[string]any{"type": "string"},
			"type":        map[string]any{"type": "string", "enum": []string{"layout", "pricing", "content", "feature"}},
		},
		"required": []string{"name", "type"},
	}
}
func (t *CreateExperimentTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in abtest.CreateExperimentInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_experiment: unmarshal input", errx.TypeValidation)
	}
	exp, err := t.abtest.CreateExperiment(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "create_experiment", errx.TypeInternal)
	}
	return fmt.Sprintf("Created experiment %q (ID: %s, type: %s, status: %s)", exp.Name, exp.ID, exp.Type, exp.Status), nil
}

type GetExperimentResultsTool struct {
	abtest   *abtestsrv.Service
	tenantID kernel.TenantID
}

func (t *GetExperimentResultsTool) Name() string        { return "get_experiment_results" }
func (t *GetExperimentResultsTool) Description() string  { return "Get results of an A/B test experiment" }
func (t *GetExperimentResultsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id": map[string]any{"type": "string", "description": "Experiment ID"},
		},
		"required": []string{"id"},
	}
}
func (t *GetExperimentResultsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct{ ID string `json:"id"` }
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "get_experiment_results: unmarshal input", errx.TypeValidation)
	}
	results, err := t.abtest.GetResults(ctx, t.tenantID, kernel.ExperimentID(in.ID))
	if err != nil {
		return "", errx.Wrap(err, "get_experiment_results", errx.TypeInternal)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Experiment: %s\nVariant results:\n", results.ExperimentID)
	for _, v := range results.Variants {
		fmt.Fprintf(&sb, "- %s: visitors=%d conversions=%d rate=%.2f%% winner=%v\n",
			v.Name, v.Visitors, v.Conversions, v.ConversionRate*100, v.IsWinner)
	}
	return sb.String(), nil
}

// ─── Recommendations ───

type ListRecommendationRulesTool struct {
	recs     *recommendationsrv.Service
	tenantID kernel.TenantID
}

func (t *ListRecommendationRulesTool) Name() string        { return "list_recommendation_rules" }
func (t *ListRecommendationRulesTool) Description() string  { return "List recommendation rules" }
func (t *ListRecommendationRulesTool) InputSchema() map[string]any {
	return map[string]any{"type": "object", "properties": map[string]any{}}
}
func (t *ListRecommendationRulesTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	rules, err := t.recs.ListRules(ctx, t.tenantID)
	if err != nil {
		return "", errx.Wrap(err, "list_recommendation_rules", errx.TypeInternal)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d recommendation rules:\n", len(rules))
	for _, r := range rules {
		fmt.Fprintf(&sb, "- %s (type: %s, active: %v, ID: %s)\n", r.Name, r.Type, r.IsActive, r.ID)
	}
	return sb.String(), nil
}

type CreateRecommendationRuleTool struct {
	recs     *recommendationsrv.Service
	tenantID kernel.TenantID
}

func (t *CreateRecommendationRuleTool) Name() string        { return "create_recommendation_rule" }
func (t *CreateRecommendationRuleTool) Description() string  { return "Create a recommendation rule" }
func (t *CreateRecommendationRuleTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":              map[string]any{"type": "string"},
			"type":              map[string]any{"type": "string", "enum": []string{"frequently_bought_together", "trending", "recently_viewed", "personalized", "manual"}},
			"source_product_id": map[string]any{"type": "string"},
			"is_active":         map[string]any{"type": "boolean"},
		},
		"required": []string{"name", "type"},
	}
}
func (t *CreateRecommendationRuleTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in recommendation.CreateRuleInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_recommendation_rule: unmarshal input", errx.TypeValidation)
	}
	rule, err := t.recs.CreateRule(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "create_recommendation_rule", errx.TypeInternal)
	}
	return fmt.Sprintf("Created recommendation rule %q (ID: %s, type: %s)", rule.Name, rule.ID, rule.Type), nil
}

type GetTrendingProductsTool struct {
	recs     *recommendationsrv.Service
	tenantID kernel.TenantID
}

func (t *GetTrendingProductsTool) Name() string        { return "get_trending_products" }
func (t *GetTrendingProductsTool) Description() string  { return "Get trending products" }
func (t *GetTrendingProductsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"limit": map[string]any{"type": "integer", "description": "Number of products (default 10)"},
		},
	}
}
func (t *GetTrendingProductsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct{ Limit int `json:"limit"` }
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "get_trending_products: unmarshal input", errx.TypeValidation)
	}
	if in.Limit < 1 { in.Limit = 10 }
	products, err := t.recs.GetTrending(ctx, t.tenantID, in.Limit, 7*24*3600000000000) // 7 days
	if err != nil {
		return "", errx.Wrap(err, "get_trending_products", errx.TypeInternal)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Trending products (%d):\n", len(products))
	for i, p := range products {
		fmt.Fprintf(&sb, "%d. Product %s (score: %.2f)\n", i+1, p.ProductID, p.Score)
	}
	return sb.String(), nil
}
