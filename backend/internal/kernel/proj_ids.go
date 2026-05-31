package kernel

// Commerce domain IDs

type ProductID string

func NewProductID(id string) ProductID { return ProductID(id) }
func (p ProductID) String() string     { return string(p) }
func (p ProductID) IsEmpty() bool      { return string(p) == "" }

type OrderID string

func NewOrderID(id string) OrderID { return OrderID(id) }
func (o OrderID) String() string   { return string(o) }
func (o OrderID) IsEmpty() bool    { return string(o) == "" }

type OrderItemID string

func NewOrderItemID(id string) OrderItemID { return OrderItemID(id) }
func (o OrderItemID) String() string       { return string(o) }
func (o OrderItemID) IsEmpty() bool        { return string(o) == "" }

type CustomerID string

func NewCustomerID(id string) CustomerID { return CustomerID(id) }
func (c CustomerID) String() string      { return string(c) }
func (c CustomerID) IsEmpty() bool       { return string(c) == "" }

type CategoryID string

func NewCategoryID(id string) CategoryID { return CategoryID(id) }
func (c CategoryID) String() string      { return string(c) }
func (c CategoryID) IsEmpty() bool       { return string(c) == "" }

type CollectionID string

func NewCollectionID(id string) CollectionID { return CollectionID(id) }
func (c CollectionID) String() string        { return string(c) }
func (c CollectionID) IsEmpty() bool         { return string(c) == "" }

type PageID string

func NewPageID(id string) PageID { return PageID(id) }
func (p PageID) String() string  { return string(p) }
func (p PageID) IsEmpty() bool   { return string(p) == "" }

type PageVersionID string

func NewPageVersionID(id string) PageVersionID { return PageVersionID(id) }
func (p PageVersionID) String() string         { return string(p) }
func (p PageVersionID) IsEmpty() bool          { return string(p) == "" }

type PromoID string

func NewPromoID(id string) PromoID { return PromoID(id) }
func (p PromoID) String() string   { return string(p) }
func (p PromoID) IsEmpty() bool    { return string(p) == "" }

type MediaID string

func NewMediaID(id string) MediaID { return MediaID(id) }
func (m MediaID) String() string   { return string(m) }
func (m MediaID) IsEmpty() bool    { return string(m) == "" }

type PluginID string

func NewPluginID(id string) PluginID { return PluginID(id) }
func (p PluginID) String() string    { return string(p) }
func (p PluginID) IsEmpty() bool     { return string(p) == "" }

type PluginVersionID string

func NewPluginVersionID(id string) PluginVersionID { return PluginVersionID(id) }
func (p PluginVersionID) String() string           { return string(p) }
func (p PluginVersionID) IsEmpty() bool            { return string(p) == "" }

type InstallationID string

func NewInstallationID(id string) InstallationID { return InstallationID(id) }
func (i InstallationID) String() string          { return string(i) }
func (i InstallationID) IsEmpty() bool           { return string(i) == "" }

type SettingID string

func NewSettingID(id string) SettingID { return SettingID(id) }
func (s SettingID) String() string     { return string(s) }
func (s SettingID) IsEmpty() bool      { return string(s) == "" }

type VendorID string

func NewVendorID(id string) VendorID { return VendorID(id) }
func (v VendorID) String() string    { return string(v) }
func (v VendorID) IsEmpty() bool     { return string(v) == "" }

type BlockTypeID string

func NewBlockTypeID(id string) BlockTypeID { return BlockTypeID(id) }
func (b BlockTypeID) String() string       { return string(b) }
func (b BlockTypeID) IsEmpty() bool        { return string(b) == "" }

type ThemeID string

func NewThemeID(id string) ThemeID { return ThemeID(id) }
func (t ThemeID) String() string   { return string(t) }
func (t ThemeID) IsEmpty() bool    { return string(t) == "" }

type BlockID string

func NewBlockID(id string) BlockID { return BlockID(id) }
func (b BlockID) String() string   { return string(b) }
func (b BlockID) IsEmpty() bool    { return string(b) == "" }

type CartID string

func NewCartID(id string) CartID { return CartID(id) }
func (c CartID) String() string  { return string(c) }
func (c CartID) IsEmpty() bool   { return string(c) == "" }

type CartItemID string

func NewCartItemID(id string) CartItemID { return CartItemID(id) }
func (c CartItemID) String() string      { return string(c) }
func (c CartItemID) IsEmpty() bool       { return string(c) == "" }

type ShippingZoneID string

func NewShippingZoneID(id string) ShippingZoneID { return ShippingZoneID(id) }
func (s ShippingZoneID) String() string          { return string(s) }
func (s ShippingZoneID) IsEmpty() bool           { return string(s) == "" }

type ShippingRateID string

func NewShippingRateID(id string) ShippingRateID { return ShippingRateID(id) }
func (s ShippingRateID) String() string          { return string(s) }
func (s ShippingRateID) IsEmpty() bool           { return string(s) == "" }

type TaxRateID string

func NewTaxRateID(id string) TaxRateID { return TaxRateID(id) }
func (t TaxRateID) String() string     { return string(t) }
func (t TaxRateID) IsEmpty() bool      { return string(t) == "" }

type PaymentID string

func NewPaymentID(id string) PaymentID { return PaymentID(id) }
func (p PaymentID) String() string     { return string(p) }
func (p PaymentID) IsEmpty() bool      { return string(p) == "" }

type RefundID string

func NewRefundID(id string) RefundID { return RefundID(id) }
func (r RefundID) String() string    { return string(r) }
func (r RefundID) IsEmpty() bool     { return string(r) == "" }

type VariantID string

func NewVariantID(id string) VariantID { return VariantID(id) }
func (v VariantID) String() string     { return string(v) }
func (v VariantID) IsEmpty() bool      { return string(v) == "" }

type OptionID string

func NewOptionID(id string) OptionID { return OptionID(id) }
func (o OptionID) String() string    { return string(o) }
func (o OptionID) IsEmpty() bool     { return string(o) == "" }

type WishlistID string

func NewWishlistID(id string) WishlistID { return WishlistID(id) }
func (w WishlistID) String() string      { return string(w) }
func (w WishlistID) IsEmpty() bool       { return string(w) == "" }

type WishlistItemID string

func NewWishlistItemID(id string) WishlistItemID { return WishlistItemID(id) }
func (w WishlistItemID) String() string          { return string(w) }
func (w WishlistItemID) IsEmpty() bool           { return string(w) == "" }

type CustomerGroupID string

func NewCustomerGroupID(id string) CustomerGroupID { return CustomerGroupID(id) }
func (c CustomerGroupID) String() string           { return string(c) }
func (c CustomerGroupID) IsEmpty() bool            { return string(c) == "" }

type CustomerGroupMembershipID string

func NewCustomerGroupMembershipID(id string) CustomerGroupMembershipID {
	return CustomerGroupMembershipID(id)
}
func (c CustomerGroupMembershipID) String() string { return string(c) }
func (c CustomerGroupMembershipID) IsEmpty() bool  { return string(c) == "" }

type GiftCardID string

func NewGiftCardID(id string) GiftCardID { return GiftCardID(id) }
func (g GiftCardID) String() string      { return string(g) }
func (g GiftCardID) IsEmpty() bool       { return string(g) == "" }

type GiftCardTransactionID string

func NewGiftCardTransactionID(id string) GiftCardTransactionID { return GiftCardTransactionID(id) }
func (g GiftCardTransactionID) String() string                 { return string(g) }
func (g GiftCardTransactionID) IsEmpty() bool                  { return string(g) == "" }

type RecoveryID string

func NewRecoveryID(id string) RecoveryID { return RecoveryID(id) }
func (r RecoveryID) String() string      { return string(r) }
func (r RecoveryID) IsEmpty() bool       { return string(r) == "" }

type CurrencyRateID string

func NewCurrencyRateID(id string) CurrencyRateID { return CurrencyRateID(id) }
func (c CurrencyRateID) String() string          { return string(c) }
func (c CurrencyRateID) IsEmpty() bool           { return string(c) == "" }

type TranslationID string

func NewTranslationID(id string) TranslationID { return TranslationID(id) }
func (t TranslationID) String() string         { return string(t) }
func (t TranslationID) IsEmpty() bool          { return string(t) == "" }

type SubscriptionID string

func NewSubscriptionID(id string) SubscriptionID { return SubscriptionID(id) }
func (s SubscriptionID) String() string          { return string(s) }
func (s SubscriptionID) IsEmpty() bool           { return string(s) == "" }

type BillingRecordID string

func NewBillingRecordID(id string) BillingRecordID { return BillingRecordID(id) }
func (b BillingRecordID) String() string           { return string(b) }
func (b BillingRecordID) IsEmpty() bool            { return string(b) == "" }

type WarehouseID string

func NewWarehouseID(id string) WarehouseID { return WarehouseID(id) }
func (w WarehouseID) String() string       { return string(w) }
func (w WarehouseID) IsEmpty() bool        { return string(w) == "" }

type StockLevelID string

func NewStockLevelID(id string) StockLevelID { return StockLevelID(id) }
func (s StockLevelID) String() string        { return string(s) }
func (s StockLevelID) IsEmpty() bool         { return string(s) == "" }

type StockMovementID string

func NewStockMovementID(id string) StockMovementID { return StockMovementID(id) }
func (s StockMovementID) String() string           { return string(s) }
func (s StockMovementID) IsEmpty() bool            { return string(s) == "" }

type ReviewID string

func NewReviewID(id string) ReviewID { return ReviewID(id) }
func (r ReviewID) String() string    { return string(r) }
func (r ReviewID) IsEmpty() bool     { return string(r) == "" }
type ReturnID string

func NewReturnID(id string) ReturnID { return ReturnID(id) }
func (r ReturnID) String() string    { return string(r) }
func (r ReturnID) IsEmpty() bool     { return string(r) == "" }

type ReturnItemID string

func NewReturnItemID(id string) ReturnItemID { return ReturnItemID(id) }
func (r ReturnItemID) String() string        { return string(r) }
func (r ReturnItemID) IsEmpty() bool         { return string(r) == "" }
type WebhookID string

func NewWebhookID(id string) WebhookID { return WebhookID(id) }
func (w WebhookID) String() string     { return string(w) }
func (w WebhookID) IsEmpty() bool      { return string(w) == "" }

type WebhookDeliveryID string

func NewWebhookDeliveryID(id string) WebhookDeliveryID { return WebhookDeliveryID(id) }
func (w WebhookDeliveryID) String() string             { return string(w) }
func (w WebhookDeliveryID) IsEmpty() bool              { return string(w) == "" }
type AuditEntryID string

func NewAuditEntryID(id string) AuditEntryID { return AuditEntryID(id) }
func (a AuditEntryID) String() string        { return string(a) }
func (a AuditEntryID) IsEmpty() bool         { return string(a) == "" }

type LoyaltyAccountID string

func NewLoyaltyAccountID(id string) LoyaltyAccountID { return LoyaltyAccountID(id) }
func (l LoyaltyAccountID) String() string            { return string(l) }
func (l LoyaltyAccountID) IsEmpty() bool             { return string(l) == "" }

type LoyaltyTransactionID string

func NewLoyaltyTransactionID(id string) LoyaltyTransactionID { return LoyaltyTransactionID(id) }
func (l LoyaltyTransactionID) String() string                { return string(l) }
func (l LoyaltyTransactionID) IsEmpty() bool                 { return string(l) == "" }

type RewardID string

func NewRewardID(id string) RewardID { return RewardID(id) }
func (r RewardID) String() string    { return string(r) }
func (r RewardID) IsEmpty() bool     { return string(r) == "" }
type BundleID string

func NewBundleID(id string) BundleID { return BundleID(id) }
func (b BundleID) String() string    { return string(b) }
func (b BundleID) IsEmpty() bool     { return string(b) == "" }

type BundleItemID string

func NewBundleItemID(id string) BundleItemID { return BundleItemID(id) }
func (b BundleItemID) String() string        { return string(b) }
func (b BundleItemID) IsEmpty() bool         { return string(b) == "" }
type SocialAccountID string

func NewSocialAccountID(id string) SocialAccountID { return SocialAccountID(id) }
func (s SocialAccountID) String() string           { return string(s) }
func (s SocialAccountID) IsEmpty() bool            { return string(s) == "" }
type NotificationID string

func NewNotificationID(id string) NotificationID { return NotificationID(id) }
func (n NotificationID) String() string          { return string(n) }
func (n NotificationID) IsEmpty() bool           { return string(n) == "" }

type StorefrontEntryID string

func NewStorefrontEntryID(id string) StorefrontEntryID { return StorefrontEntryID(id) }
func (s StorefrontEntryID) String() string             { return string(s) }
func (s StorefrontEntryID) IsEmpty() bool              { return string(s) == "" }

type StorefrontCatalogID string

func NewStorefrontCatalogID(id string) StorefrontCatalogID { return StorefrontCatalogID(id) }
func (s StorefrontCatalogID) String() string               { return string(s) }
func (s StorefrontCatalogID) IsEmpty() bool                { return string(s) == "" }
type BulkOperationID string

func NewBulkOperationID(id string) BulkOperationID { return BulkOperationID(id) }
func (b BulkOperationID) String() string           { return string(b) }
func (b BulkOperationID) IsEmpty() bool            { return string(b) == "" }

type BulkOperationItemID string

func NewBulkOperationItemID(id string) BulkOperationItemID { return BulkOperationItemID(id) }
func (b BulkOperationItemID) String() string               { return string(b) }
func (b BulkOperationItemID) IsEmpty() bool                { return string(b) == "" }
type BlogPostID string

func NewBlogPostID(id string) BlogPostID { return BlogPostID(id) }
func (b BlogPostID) String() string      { return string(b) }
func (b BlogPostID) IsEmpty() bool       { return string(b) == "" }

type BlogCategoryID string

func NewBlogCategoryID(id string) BlogCategoryID { return BlogCategoryID(id) }
func (b BlogCategoryID) String() string          { return string(b) }
func (b BlogCategoryID) IsEmpty() bool           { return string(b) == "" }
type CollectionProductID string

func NewCollectionProductID(id string) CollectionProductID { return CollectionProductID(id) }
func (c CollectionProductID) String() string               { return string(c) }
func (c CollectionProductID) IsEmpty() bool                { return string(c) == "" }
type ExperimentID string

func NewExperimentID(id string) ExperimentID { return ExperimentID(id) }
func (e ExperimentID) String() string        { return string(e) }
func (e ExperimentID) IsEmpty() bool         { return string(e) == "" }

type ExperimentVariantID string

func NewExperimentVariantID(id string) ExperimentVariantID { return ExperimentVariantID(id) }
func (e ExperimentVariantID) String() string               { return string(e) }
func (e ExperimentVariantID) IsEmpty() bool                { return string(e) == "" }

type ExperimentAssignmentID string

func NewExperimentAssignmentID(id string) ExperimentAssignmentID { return ExperimentAssignmentID(id) }
func (e ExperimentAssignmentID) String() string                  { return string(e) }
func (e ExperimentAssignmentID) IsEmpty() bool                   { return string(e) == "" }
type ProductViewID string

func NewProductViewID(id string) ProductViewID { return ProductViewID(id) }
func (p ProductViewID) String() string         { return string(p) }
func (p ProductViewID) IsEmpty() bool          { return string(p) == "" }

type ProductInteractionID string

func NewProductInteractionID(id string) ProductInteractionID { return ProductInteractionID(id) }
func (p ProductInteractionID) String() string                { return string(p) }
func (p ProductInteractionID) IsEmpty() bool                 { return string(p) == "" }

type RecommendationRuleID string

func NewRecommendationRuleID(id string) RecommendationRuleID { return RecommendationRuleID(id) }
func (r RecommendationRuleID) String() string                { return string(r) }
func (r RecommendationRuleID) IsEmpty() bool                 { return string(r) == "" }

// Agent preset and session IDs

type PresetID string

func NewPresetID(id string) PresetID { return PresetID(id) }
func (p PresetID) String() string    { return string(p) }
func (p PresetID) IsEmpty() bool     { return string(p) == "" }

type AgentSessionID string

func NewAgentSessionID(id string) AgentSessionID { return AgentSessionID(id) }
func (a AgentSessionID) String() string          { return string(a) }
func (a AgentSessionID) IsEmpty() bool           { return string(a) == "" }
