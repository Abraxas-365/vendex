package agent

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Abraxas-365/vendex/internal/abtest"
	"github.com/Abraxas-365/vendex/internal/abtest/abtestsrv"
	"github.com/Abraxas-365/vendex/internal/blog"
	"github.com/Abraxas-365/vendex/internal/blog/blogsrv"
	"github.com/Abraxas-365/vendex/internal/bulkops"
	"github.com/Abraxas-365/vendex/internal/bulkops/bulkopssrv"
	"github.com/Abraxas-365/vendex/internal/collection"
	"github.com/Abraxas-365/vendex/internal/collection/collectionsrv"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/multistore"
	"github.com/Abraxas-365/vendex/internal/multistore/multistoresrv"
	"github.com/Abraxas-365/vendex/internal/recommendation"
	"github.com/Abraxas-365/vendex/internal/recommendation/recommendationsrv"
)

// ── Multi-store stubs ──

type stubMultistoreRepo struct{}

func (s *stubMultistoreRepo) Create(_ context.Context, sf *multistore.Storefront) error {
	sf.ID = "sf-1"
	return nil
}
func (s *stubMultistoreRepo) GetByID(context.Context, kernel.TenantID, kernel.StorefrontEntryID) (*multistore.Storefront, error) {
	panic("unused")
}
func (s *stubMultistoreRepo) GetBySlug(_ context.Context, _ kernel.TenantID, slug string) (*multistore.Storefront, error) {
	return nil, multistore.ErrNotFound
}
func (s *stubMultistoreRepo) GetByDomain(context.Context, string) (*multistore.Storefront, error) {
	return nil, multistore.ErrNotFound
}
func (s *stubMultistoreRepo) List(_ context.Context, _ kernel.TenantID, _, _ int) (kernel.Paginated[multistore.Storefront], error) {
	return kernel.Paginated[multistore.Storefront]{
		Items: []multistore.Storefront{{ID: "sf-1", Name: "Default", Slug: "default", IsDefault: true}},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (s *stubMultistoreRepo) Update(context.Context, *multistore.Storefront) error { panic("unused") }
func (s *stubMultistoreRepo) Delete(context.Context, kernel.TenantID, kernel.StorefrontEntryID) error {
	panic("unused")
}
func (s *stubMultistoreRepo) SetDefault(context.Context, kernel.TenantID, kernel.StorefrontEntryID) error {
	panic("unused")
}
func (s *stubMultistoreRepo) ClearDefault(context.Context, kernel.TenantID) error { return nil }
func (s *stubMultistoreRepo) AddCatalog(context.Context, *multistore.StorefrontCatalog) error {
	panic("unused")
}
func (s *stubMultistoreRepo) RemoveCatalog(context.Context, kernel.TenantID, kernel.StorefrontEntryID, string) error {
	panic("unused")
}
func (s *stubMultistoreRepo) ListCatalogs(context.Context, kernel.TenantID, kernel.StorefrontEntryID) ([]multistore.StorefrontCatalog, error) {
	panic("unused")
}

// ── Bulk ops stubs ──

type stubBulkOpsRepo struct{}

func (s *stubBulkOpsRepo) Create(_ context.Context, op *bulkops.BulkOperation, _ []bulkops.BulkOperationItem) error {
	op.ID = "bop-1"
	return nil
}
func (s *stubBulkOpsRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.BulkOperationID) (*bulkops.BulkOperation, error) {
	return &bulkops.BulkOperation{
		ID:        "bop-1",
		TenantID:  testTenant,
		Type:      "price_update",
		Status:    "pending",
		CreatedAt: time.Now(),
	}, nil
}
func (s *stubBulkOpsRepo) UpdateStatus(_ context.Context, _ kernel.TenantID, _ kernel.BulkOperationID, _ bulkops.OperationStatus) error {
	return nil
}
func (s *stubBulkOpsRepo) UpdateOperation(context.Context, *bulkops.BulkOperation) error { return nil }
func (s *stubBulkOpsRepo) UpdateItem(context.Context, *bulkops.BulkOperationItem) error  { return nil }
func (s *stubBulkOpsRepo) List(_ context.Context, _ kernel.TenantID, _, _ int) (kernel.Paginated[bulkops.BulkOperation], error) {
	return kernel.Paginated[bulkops.BulkOperation]{
		Items: []bulkops.BulkOperation{{ID: "bop-1", Type: "price_update", Status: "pending"}},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (s *stubBulkOpsRepo) ListItems(_ context.Context, _ kernel.TenantID, _ kernel.BulkOperationID, _, _ int) (kernel.Paginated[bulkops.BulkOperationItem], error) {
	return kernel.Paginated[bulkops.BulkOperationItem]{Items: []bulkops.BulkOperationItem{}, Total: 0, Page: 1, TotalPages: 0}, nil
}

// ── Blog stubs ──

type stubBlogRepo struct{}

func (s *stubBlogRepo) CreatePost(_ context.Context, p *blog.BlogPost) error {
	p.ID = "bp-1"
	return nil
}
func (s *stubBlogRepo) GetPostByID(_ context.Context, _ kernel.TenantID, _ kernel.BlogPostID) (*blog.BlogPost, error) {
	return &blog.BlogPost{
		ID:       "bp-1",
		TenantID: testTenant,
		Title:    "My Post",
		Slug:     "my-post",
		Status:   blog.StatusDraft,
	}, nil
}
func (s *stubBlogRepo) GetPostBySlug(_ context.Context, _ kernel.TenantID, _ string) (*blog.BlogPost, error) {
	return nil, blog.ErrPostNotFound
}
func (s *stubBlogRepo) ListPosts(_ context.Context, _ kernel.TenantID, _ blog.ListPostsFilter) (kernel.Paginated[blog.BlogPost], error) {
	return kernel.Paginated[blog.BlogPost]{
		Items: []blog.BlogPost{{ID: "bp-1", Title: "My Post", Slug: "my-post", Status: blog.StatusDraft}},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (s *stubBlogRepo) UpdatePost(context.Context, *blog.BlogPost) error    { return nil }
func (s *stubBlogRepo) DeletePost(context.Context, kernel.TenantID, kernel.BlogPostID) error { panic("unused") }
func (s *stubBlogRepo) CreateCategory(_ context.Context, c *blog.BlogCategory) error {
	c.ID = "bc-1"
	return nil
}
func (s *stubBlogRepo) GetCategoryByID(context.Context, kernel.TenantID, kernel.BlogCategoryID) (*blog.BlogCategory, error) {
	panic("unused")
}
func (s *stubBlogRepo) ListCategories(_ context.Context, _ kernel.TenantID) ([]blog.BlogCategory, error) {
	return []blog.BlogCategory{{ID: "bc-1", Name: "Tech", Slug: "tech"}}, nil
}
func (s *stubBlogRepo) UpdateCategory(context.Context, *blog.BlogCategory) error { panic("unused") }
func (s *stubBlogRepo) DeleteCategory(context.Context, kernel.TenantID, kernel.BlogCategoryID) error {
	panic("unused")
}
func (s *stubBlogRepo) AddPostCategory(context.Context, kernel.TenantID, kernel.BlogPostID, kernel.BlogCategoryID) error {
	panic("unused")
}
func (s *stubBlogRepo) RemovePostCategory(context.Context, kernel.TenantID, kernel.BlogPostID, kernel.BlogCategoryID) error {
	panic("unused")
}
func (s *stubBlogRepo) GetPostCategories(context.Context, kernel.BlogPostID) ([]blog.BlogCategory, error) {
	panic("unused")
}

// ── Collection stubs ──

type stubCollectionDomainRepo struct{}

func (s *stubCollectionDomainRepo) Create(_ context.Context, c *collection.Collection) error {
	c.ID = "col-1"
	return nil
}
func (s *stubCollectionDomainRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.CollectionID) (*collection.Collection, error) {
	panic("unused")
}
func (s *stubCollectionDomainRepo) GetBySlug(_ context.Context, _ kernel.TenantID, _ string) (*collection.Collection, error) {
	return nil, collection.ErrNotFound
}
func (s *stubCollectionDomainRepo) List(_ context.Context, _ kernel.TenantID, _ bool, _ kernel.PaginationOptions) (kernel.Paginated[collection.Collection], error) {
	return kernel.Paginated[collection.Collection]{
		Items: []collection.Collection{{ID: "col-1", Name: "Summer", Slug: "summer", Type: "manual"}},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (s *stubCollectionDomainRepo) Update(_ context.Context, _ *collection.Collection) error {
	panic("unused")
}
func (s *stubCollectionDomainRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.CollectionID) error {
	panic("unused")
}
func (s *stubCollectionDomainRepo) CountProducts(_ context.Context, _ kernel.TenantID, _ kernel.CollectionID) (int, error) {
	return 0, nil
}
func (s *stubCollectionDomainRepo) AddProduct(_ context.Context, _ *collection.CollectionProduct) error {
	return nil
}
func (s *stubCollectionDomainRepo) RemoveProduct(_ context.Context, _ kernel.TenantID, _ kernel.CollectionID, _ string) error {
	panic("unused")
}
func (s *stubCollectionDomainRepo) ListProducts(_ context.Context, _ kernel.TenantID, _ kernel.CollectionID, _ kernel.PaginationOptions) (kernel.Paginated[collection.CollectionProduct], error) {
	panic("unused")
}
func (s *stubCollectionDomainRepo) ReorderProducts(_ context.Context, _ kernel.TenantID, _ kernel.CollectionID, _ []string) error {
	panic("unused")
}

// ── A/B test stubs ──

type stubABTestRepo struct{}

func (s *stubABTestRepo) CreateExperiment(_ context.Context, e *abtest.Experiment) error {
	e.ID = "exp-1"
	return nil
}
func (s *stubABTestRepo) GetExperimentByID(_ context.Context, _ kernel.TenantID, _ kernel.ExperimentID) (*abtest.Experiment, error) {
	return &abtest.Experiment{
		ID:       "exp-1",
		TenantID: testTenant,
		Name:     "Button Color",
		Status:   "running",
		Variants: []abtest.ExperimentVariant{{ID: "v1", Name: "Red", Weight: 50}, {ID: "v2", Name: "Blue", Weight: 50}},
	}, nil
}
func (s *stubABTestRepo) ListExperiments(_ context.Context, _ kernel.TenantID, _ string, _ kernel.PaginationOptions) (kernel.Paginated[abtest.Experiment], error) {
	return kernel.Paginated[abtest.Experiment]{
		Items: []abtest.Experiment{{ID: "exp-1", Name: "Button Color", Status: "running"}},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (s *stubABTestRepo) UpdateExperiment(context.Context, *abtest.Experiment) error {
	panic("unused")
}
func (s *stubABTestRepo) DeleteExperiment(context.Context, kernel.TenantID, kernel.ExperimentID) error {
	panic("unused")
}
func (s *stubABTestRepo) CreateVariant(context.Context, *abtest.ExperimentVariant) error {
	panic("unused")
}
func (s *stubABTestRepo) GetVariantByID(context.Context, kernel.TenantID, kernel.ExperimentVariantID) (*abtest.ExperimentVariant, error) {
	panic("unused")
}
func (s *stubABTestRepo) ListVariants(_ context.Context, _ kernel.TenantID, _ kernel.ExperimentID) ([]abtest.ExperimentVariant, error) {
	return []abtest.ExperimentVariant{
		{ID: "v1", Name: "Red", Visitors: 100, Conversions: 10},
		{ID: "v2", Name: "Blue", Visitors: 100, Conversions: 15},
	}, nil
}
func (s *stubABTestRepo) DeleteVariant(context.Context, kernel.TenantID, kernel.ExperimentVariantID) error {
	panic("unused")
}
func (s *stubABTestRepo) IncrementVariantVisitors(context.Context, kernel.ExperimentVariantID) error {
	panic("unused")
}
func (s *stubABTestRepo) IncrementVariantConversions(context.Context, kernel.ExperimentVariantID, int64) error {
	panic("unused")
}
func (s *stubABTestRepo) CreateAssignment(context.Context, *abtest.ExperimentAssignment) error {
	panic("unused")
}
func (s *stubABTestRepo) GetAssignment(context.Context, kernel.TenantID, kernel.ExperimentID, string) (*abtest.ExperimentAssignment, error) {
	panic("unused")
}
func (s *stubABTestRepo) RecordConversion(context.Context, kernel.TenantID, kernel.ExperimentID, string, int64) error {
	panic("unused")
}


// ── Recommendation stubs ──

type stubRecommendationRepo struct{}

func (s *stubRecommendationRepo) TrackView(context.Context, recommendation.ProductView) error {
	panic("unused")
}
func (s *stubRecommendationRepo) TrackInteraction(context.Context, recommendation.ProductInteraction) error {
	panic("unused")
}
func (s *stubRecommendationRepo) GetFrequentlyBoughtTogether(context.Context, kernel.TenantID, string, int) ([]recommendation.RecommendedProduct, error) {
	panic("unused")
}
func (s *stubRecommendationRepo) GetTrending(_ context.Context, _ kernel.TenantID, _ int, _ time.Duration) ([]recommendation.RecommendedProduct, error) {
	return []recommendation.RecommendedProduct{
		{ProductID: "p-1", Score: 42.5, Reason: "trending"},
	}, nil
}
func (s *stubRecommendationRepo) GetRecentlyViewed(context.Context, kernel.TenantID, string, int) ([]recommendation.RecommendedProduct, error) {
	panic("unused")
}
func (s *stubRecommendationRepo) GetPersonalized(context.Context, kernel.TenantID, string, int) ([]recommendation.RecommendedProduct, error) {
	panic("unused")
}
func (s *stubRecommendationRepo) CreateRule(_ context.Context, r recommendation.RecommendationRule) (recommendation.RecommendationRule, error) {
	r.ID = "rule-1"
	return r, nil
}
func (s *stubRecommendationRepo) GetRuleByID(context.Context, kernel.TenantID, kernel.RecommendationRuleID) (recommendation.RecommendationRule, error) {
	panic("unused")
}
func (s *stubRecommendationRepo) ListRules(_ context.Context, _ kernel.TenantID) ([]recommendation.RecommendationRule, error) {
	return []recommendation.RecommendationRule{{ID: "rule-1", Name: "Similar Items", Type: "similar_products", IsActive: true}}, nil
}
func (s *stubRecommendationRepo) UpdateRule(context.Context, recommendation.RecommendationRule) (recommendation.RecommendationRule, error) {
	panic("unused")
}
func (s *stubRecommendationRepo) DeleteRule(context.Context, kernel.TenantID, kernel.RecommendationRuleID) error {
	panic("unused")
}

// ── Multi-store tests ──

func TestListStorefrontsTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := multistoresrv.New(&stubMultistoreRepo{}, bus)
	tool := &ListStorefrontsTool{multistore: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "default") {
		t.Errorf("expected 'default' in result, got: %s", result)
	}
}

func TestCreateStorefrontTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := multistoresrv.New(&stubMultistoreRepo{}, bus)
	tool := &CreateStorefrontTool{multistore: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"US Store","slug":"us-store","domain":"us.example.com","default_locale":"en","default_currency":"USD"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(result), "storefront") {
		t.Errorf("expected 'storefront' in result, got: %s", result)
	}
}

func TestCreateStorefrontTool_Execute_InvalidJSON(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := multistoresrv.New(&stubMultistoreRepo{}, bus)
	tool := &CreateStorefrontTool{multistore: svc, tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ── Bulk ops tests ──

func TestListBulkOperationsTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := bulkopssrv.New(&stubBulkOpsRepo{}, bus)
	tool := &ListBulkOperationsTool{bulkops: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "price_update") {
		t.Errorf("expected 'price_update' in result, got: %s", result)
	}
}

func TestCreateBulkOperationTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := bulkopssrv.New(&stubBulkOpsRepo{}, bus)
	tool := &CreateBulkOperationTool{bulkops: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"Type":"price_update","ResourceType":"product","ResourceIDs":["p-1","p-2"],"CreatedBy":"admin"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(result), "bulk") || !strings.Contains(strings.ToLower(result), "created") {
		t.Errorf("expected bulk operation created in result, got: %s", result)
	}
}

func TestProcessBulkOperationTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := bulkopssrv.New(&stubBulkOpsRepo{}, bus)
	tool := &ProcessBulkOperationTool{bulkops: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"operation_id":"bop-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(result), "process") {
		t.Errorf("expected 'process' in result, got: %s", result)
	}
}

// ── Blog tests ──

func TestListBlogPostsTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := blogsrv.New(&stubBlogRepo{}, bus)
	tool := &ListBlogPostsTool{blog: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "My Post") {
		t.Errorf("expected 'My Post' in result, got: %s", result)
	}
}

func TestCreateBlogPostTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := blogsrv.New(&stubBlogRepo{}, bus)
	tool := &CreateBlogPostTool{blog: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"title":"New Post","slug":"new-post","content":"Hello world","author":"admin"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "New Post") {
		t.Errorf("expected 'New Post' in result, got: %s", result)
	}
}

func TestPublishBlogPostTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := blogsrv.New(&stubBlogRepo{}, bus)
	tool := &PublishBlogPostTool{blog: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"post_id":"bp-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(result), "publish") {
		t.Errorf("expected 'publish' in result, got: %s", result)
	}
}

func TestListBlogCategoriesTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := blogsrv.New(&stubBlogRepo{}, bus)
	tool := &ListBlogCategoriesTool{blog: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Tech") {
		t.Errorf("expected 'Tech' in result, got: %s", result)
	}
}

// ── Collection tests ──

func TestListCollectionsTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := collectionsrv.New(&stubCollectionDomainRepo{}, bus)
	tool := &ListCollectionsTool{collections: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Summer") {
		t.Errorf("expected 'Summer' in result, got: %s", result)
	}
}

func TestCreateCollectionTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := collectionsrv.New(&stubCollectionDomainRepo{}, bus)
	tool := &CreateCollectionTool{collections: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"Winter Sale","slug":"winter-sale","type":"manual","description":"Winter items"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Winter Sale") {
		t.Errorf("expected 'Winter Sale' in result, got: %s", result)
	}
}

func TestAddCollectionProductTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := collectionsrv.New(&stubCollectionDomainRepo{}, bus)
	tool := &AddCollectionProductTool{collections: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"collection_id":"col-1","product_id":"p-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(result), "added") || !strings.Contains(strings.ToLower(result), "product") {
		t.Errorf("expected 'added product' in result, got: %s", result)
	}
}

// ── A/B testing tests ──

func TestListExperimentsTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := abtestsrv.New(&stubABTestRepo{}, bus)
	tool := &ListExperimentsTool{abtest: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Button Color") {
		t.Errorf("expected 'Button Color' in result, got: %s", result)
	}
}

func TestCreateExperimentTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := abtestsrv.New(&stubABTestRepo{}, bus)
	tool := &CreateExperimentTool{abtest: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"CTA Test","description":"Test CTA buttons","variants":[{"name":"Green","weight":50},{"name":"Orange","weight":50}]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "CTA Test") {
		t.Errorf("expected 'CTA Test' in result, got: %s", result)
	}
}

func TestGetExperimentResultsTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := abtestsrv.New(&stubABTestRepo{}, bus)
	tool := &GetExperimentResultsTool{abtest: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"experiment_id":"exp-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Red") && !strings.Contains(result, "Blue") {
		t.Errorf("expected variant names in result, got: %s", result)
	}
}

// ── Recommendation tests ──

func TestListRecommendationRulesTool_Execute(t *testing.T) {
	svc := recommendationsrv.New(&stubRecommendationRepo{})
	tool := &ListRecommendationRulesTool{recs: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Similar Items") {
		t.Errorf("expected 'Similar Items' in result, got: %s", result)
	}
}

func TestCreateRecommendationRuleTool_Execute(t *testing.T) {
	svc := recommendationsrv.New(&stubRecommendationRepo{})
	tool := &CreateRecommendationRuleTool{recs: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"Cross Sell","type":"cross_sell","max_recommendations":5}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Cross Sell") {
		t.Errorf("expected 'Cross Sell' in result, got: %s", result)
	}
}

func TestGetTrendingProductsTool_Execute(t *testing.T) {
	svc := recommendationsrv.New(&stubRecommendationRepo{})
	tool := &GetTrendingProductsTool{recs: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"limit":10}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "p-1") {
		t.Errorf("expected 'p-1' in result, got: %s", result)
	}
}
