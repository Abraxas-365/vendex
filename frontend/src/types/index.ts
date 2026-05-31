// ---------------------------------------------------------------------------
// Shared value types
// ---------------------------------------------------------------------------

export interface Money {
  amount: number
  currency: string
}

// ---------------------------------------------------------------------------
// Address
// ---------------------------------------------------------------------------

export interface Address {
  id?: string
  street: string
  city: string
  state: string
  country: string
  postal_code: string
  is_default: boolean
}

// ---------------------------------------------------------------------------
// Product Variants (defined before Product so Product can reference them)
// ---------------------------------------------------------------------------

export interface ProductOption {
  id: string
  product_id: string
  tenant_id: string
  name: string
  position: number
  values: string[]
  created_at: string
  updated_at: string
}

export interface ProductVariant {
  id: string
  product_id: string
  tenant_id: string
  sku: string
  price: Money
  stock: number
  options: Record<string, string>
  active: boolean
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Product
// ---------------------------------------------------------------------------

export type ProductStatus = 'draft' | 'active' | 'archived'

export interface Product {
  id: string
  tenant_id: string
  name: string
  description: string
  sku: string
  price: Money
  images: string[]
  category_id: string
  tags: string[]
  status: ProductStatus
  stock: number
  has_variants?: boolean
  options?: ProductOption[]
  variants?: ProductVariant[]
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Order
// ---------------------------------------------------------------------------

export type OrderStatus =
  | 'pending'
  | 'confirmed'
  | 'processing'
  | 'shipped'
  | 'delivered'
  | 'cancelled'

export interface OrderItem {
  id: string
  product_id: string
  product_name: string
  quantity: number
  unit_price: Money
  total: Money
}

export interface Order {
  id: string
  tenant_id: string
  customer_id: string
  items: OrderItem[]
  status: OrderStatus
  total_amount: Money
  subtotal?: Money
  shipping_amount?: Money
  tax_amount?: Money
  discount_amount?: Money
  shipping_address: Address
  shipping_method?: string
  tracking_number?: string
  carrier?: string
  payment_status?: PaymentStatus
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Customer
// ---------------------------------------------------------------------------

export interface Customer {
  id: string
  tenant_id: string
  email: string
  name: string
  phone: string
  addresses: Address[]
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Catalog — Category & Collection
// ---------------------------------------------------------------------------

export interface Category {
  id: string
  tenant_id: string
  name: string
  slug: string
  parent_id: string | null
  description: string
  created_at: string
  updated_at: string
}

export interface Collection {
  id: string
  tenant_id: string
  name: string
  slug: string
  description: string
  product_ids: string[]
  is_automatic: boolean
  rules: Record<string, unknown>
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Storefront — Pages
// ---------------------------------------------------------------------------

export type PageStatus = 'draft' | 'pending_review' | 'published' | 'archived'
export type ContentType = 'html' | 'blocks'

export interface PageMeta {
  description: string
  og_title: string
  og_image: string
  keywords: string[]
}

export interface Block {
  id: string
  type: string
  settings: Record<string, unknown>
}

export interface Section {
  id: string
  type: string
  settings: Record<string, unknown>
  blocks: Block[]
}

export interface Page {
  id: string
  tenant_id: string
  slug: string
  title: string
  html: string
  css: string
  meta: PageMeta
  content_type: ContentType
  sections: Section[]
  status: PageStatus
  version: number
  created_by: string
  published_at?: string
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Block Types
// ---------------------------------------------------------------------------

export type BlockCategory = 'content' | 'commerce' | 'media' | 'layout'

export interface BlockType {
  id: string
  name: string
  display_name: string
  category: BlockCategory
  schema: Record<string, unknown>
  default_settings: Record<string, unknown>
  icon: string
  plugin_id?: string
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Themes
// ---------------------------------------------------------------------------

export interface ThemeColors {
  primary: string
  secondary: string
  background: string
  surface: string
  text: string
  text_muted: string
  border: string
  error: string
  success: string
  warning: string
  info: string
}

export interface ThemeTypography {
  font_heading: string
  font_body: string
  base_size: string
  scale_ratio: number
}

export interface ThemeSpacing {
  unit: string
  section_padding: string
}

export interface ThemeBorders {
  radius_sm: string
  radius_md: string
  radius_lg: string
  radius_full: string
}

export interface ThemeShadows {
  sm: string
  md: string
  lg: string
}

export interface ThemeTokens {
  colors: ThemeColors
  typography: ThemeTypography
  spacing: ThemeSpacing
  borders: ThemeBorders
  shadows: ThemeShadows
}

export interface Theme {
  id: string
  tenant_id: string
  name: string
  is_active: boolean
  tokens: ThemeTokens
  created_at: string
  updated_at: string
}

export interface PageVersion {
  id: string
  page_id: string
  version: number
  html: string
  css: string
  edited_by: string
  comment: string
  created_at: string
}

// ---------------------------------------------------------------------------
// Promos
// ---------------------------------------------------------------------------

export type PromoType = 'percentage' | 'fixed_amount' | 'free_shipping'

export interface Promo {
  id: string
  tenant_id: string
  code: string
  type: PromoType
  value: number
  min_order_amount?: number
  max_uses?: number
  used_count: number
  starts_at?: string
  ends_at?: string
  active: boolean
  created_at: string
}

// ---------------------------------------------------------------------------
// Media
// ---------------------------------------------------------------------------

export interface Media {
  id: string
  tenant_id: string
  filename: string
  content_type: string
  size: number
  url: string
  alt: string
  uploaded_by: string
  created_at: string
}

// ---------------------------------------------------------------------------
// Pagination
// ---------------------------------------------------------------------------

export interface PaginatedResult<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// ---------------------------------------------------------------------------
// Agent chat
// ---------------------------------------------------------------------------

export interface AgentMessage {
  role: 'user' | 'assistant'
  content: string
  timestamp: string
}

// ---------------------------------------------------------------------------
// Marketplace types
// ---------------------------------------------------------------------------

export type PluginCategory = 'official' | 'community' | 'custom'
export type InstallationStatus = 'active' | 'inactive' | 'failed'

export interface Plugin {
  id: string
  name: string
  display_name: string
  description: string
  author: string
  icon: string
  category: PluginCategory
  tags: string[]
  created_at: string
  updated_at: string
}

export interface PluginVersion {
  id: string
  plugin_id: string
  version: string
  changelog: string
  permissions: string[]
  manifest_json: string
  frontend_url: string
  min_platform_ver: string
  created_at: string
}

export interface PluginInstallation {
  id: string
  tenant_id: string
  plugin_id: string
  version_id: string
  status: InstallationStatus
  settings: Record<string, unknown>
  installed_at: string
  updated_at: string
}

export interface PluginManifest {
  name: string
  display_name: string
  version: string
  description: string
  author: string
  permissions: string[]
  ui: {
    tabs: Array<{ label: string; icon: string; entry: string }>
    widgets: Array<{ slot: string; component: string; entry: string }>
  }
  tools: Array<{ name: string; description: string }>
  migrations: string[]
}

// ---------------------------------------------------------------------------
// Analytics
// ---------------------------------------------------------------------------

export interface DashboardStats {
  total_products: number
  total_orders: number
  total_customers: number
  total_revenue: number
  currency: string
  pending_orders: number
  active_promos: number
  pending_pages: number
}

export interface RevenuePoint {
  date: string
  amount: number
}

export interface TopProduct {
  product_id: string
  product_name: string
  total_sold: number
  revenue: number
}

export interface OrderStatusBreakdown {
  status: string
  count: number
}

export interface RecentOrder {
  id: string
  customer_id: string
  status: string
  total_amount: number
  currency: string
  created_at: string
}

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

export interface AuthUser {
  id: string
  tenant_id: string
  email: string
  name: string
  picture?: string
  status: string
  scopes: string[]
}

export interface AuthTenant {
  id: string
  name: string
  slug: string
  plan: string
  is_active: boolean
}

export interface TokenResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: AuthUser
  tenant: AuthTenant
}

export interface LoginResponse {
  auth_url: string
  state: string
}

export interface RefreshResponse {
  access_token: string
  token_type: string
  expires_in: number
}

export interface MeResponse {
  user: AuthUser
  tenant: AuthTenant
}

// ---------------------------------------------------------------------------
// Shipping
// ---------------------------------------------------------------------------

export interface ShippingZone {
  id: string
  tenant_id: string
  name: string
  countries: string[]
  states: string[]
  created_at: string
  updated_at: string
}

export interface ShippingRate {
  id: string
  zone_id: string
  tenant_id: string
  name: string
  type: 'flat' | 'weight_based' | 'price_based' | 'free'
  price: Money
  min_weight?: number
  max_weight?: number
  min_order_amount?: number
  max_order_amount?: number
  est_days_min?: number
  est_days_max?: number
  active: boolean
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Tax
// ---------------------------------------------------------------------------

export interface TaxRate {
  id: string
  tenant_id: string
  name: string
  rate: number
  country: string
  state?: string
  city?: string
  zip_code?: string
  priority: number
  compound: boolean
  includes_shipping: boolean
  active: boolean
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Payment
// ---------------------------------------------------------------------------

export type PaymentStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'refunded'
export type RefundStatus = 'pending' | 'completed' | 'failed'

export interface Payment {
  id: string
  tenant_id: string
  order_id: string
  amount: Money
  status: PaymentStatus
  provider: string
  provider_payment_id?: string
  method?: string
  error_message?: string
  paid_at?: string
  created_at: string
  updated_at: string
}

export interface Refund {
  id: string
  tenant_id: string
  payment_id: string
  order_id: string
  amount: Money
  reason?: string
  status: RefundStatus
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Settings
// ---------------------------------------------------------------------------

export interface StoreAddress {
  street: string
  city: string
  state: string
  country: string
  zip: string
}

export interface SocialLinks {
  instagram: string
  twitter: string
  facebook: string
}

export interface CheckoutConfig {
  guest_checkout: boolean
  require_phone: boolean
}

export interface StoreSettings {
  tenant_id: string
  store_name: string
  store_email: string
  store_phone: string
  currency: string
  timezone: string
  address: StoreAddress
  logo_url: string
  favicon_url: string
  social_links: SocialLinks
  checkout_config: CheckoutConfig
  updated_at: string
}

// ---------------------------------------------------------------------------
// Import / Export
// ---------------------------------------------------------------------------

export interface ImportResult {
  total: number
  imported: number
  errors: { row: number; error: string }[]
}

// ---------------------------------------------------------------------------
// Customer Groups
// ---------------------------------------------------------------------------

export interface GroupRules {
  min_order_count?: number
  min_total_spent?: number
  tags?: string[]
}

export interface CustomerGroup {
  id: string
  tenant_id: string
  name: string
  description: string
  rules: GroupRules
  member_count: number
  created_at: string
  updated_at: string
}

export interface GroupMembership {
  id: string
  group_id: string
  customer_id: string
  joined_at: string
}

// ---------------------------------------------------------------------------
// Gift Cards
// ---------------------------------------------------------------------------

export interface GiftCard {
  id: string
  tenant_id: string
  code: string
  initial_amount: number
  balance: number
  currency: string
  expires_at: string | null
  active: boolean
  created_by: string
  created_at: string
  updated_at: string
}

export interface GiftCardTransaction {
  id: string
  gift_card_id: string
  tenant_id: string
  type: 'credit' | 'debit'
  amount: number
  currency: string
  order_id: string | null
  note: string
  created_at: string
}

// ---------------------------------------------------------------------------
// Cart Recovery
// ---------------------------------------------------------------------------

export interface RecoveryEmail {
  id: string
  tenant_id: string
  cart_id: string
  customer_id: string
  email: string
  step: number
  status: string
  discount_code: string | null
  sent_at: string | null
  clicked_at: string | null
  converted_at: string | null
  created_at: string
}

export interface RecoveryStats {
  total: number
  sent: number
  clicked: number
  converted: number
  conversion_rate: number
}

// ---------------------------------------------------------------------------
// Currency
// ---------------------------------------------------------------------------

export interface CurrencyRate {
  id: string
  tenant_id: string
  base_currency: string
  target_currency: string
  rate: number
  auto_update: boolean
  updated_at: string
  created_at: string
}

export interface ConvertResult {
  from_amount: number
  from_currency: string
  to_amount: number
  to_currency: string
  rate: number
}

// ---------------------------------------------------------------------------
// I18n / Translations
// ---------------------------------------------------------------------------

export interface TranslationBundle {
  entity_type: string
  entity_id: string
  locale: string
  fields: Record<string, string>
}

// ---------------------------------------------------------------------------
// Subscriptions
// ---------------------------------------------------------------------------

export interface Subscription {
  id: string
  tenant_id: string
  customer_id: string
  product_id: string
  variant_id: string | null
  price: { amount: number; currency: string }
  interval: 'weekly' | 'monthly' | 'quarterly' | 'yearly'
  status: 'active' | 'paused' | 'cancelled' | 'expired'
  next_billing_date: string
  last_billed_at: string | null
  cancelled_at: string | null
  paused_at: string | null
  trial_ends_at: string | null
  metadata: Record<string, string>
  created_at: string
  updated_at: string
}

export interface BillingRecord {
  id: string
  subscription_id: string
  tenant_id: string
  amount: number
  currency: string
  status: 'success' | 'failed' | 'pending'
  order_id: string | null
  failure_reason: string | null
  billed_at: string
  created_at: string
}

// ---------------------------------------------------------------------------
// Inventory (027)
// ---------------------------------------------------------------------------

export interface Warehouse {
  id: string
  tenant_id: string
  name: string
  address: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface StockLevel {
  product_id: string
  warehouse_id: string
  warehouse_name: string
  quantity: number
  reserved: number
  available: number
}

export interface StockMovement {
  id: string
  tenant_id: string
  product_id: string
  warehouse_id: string
  type: 'adjustment' | 'sale' | 'return' | 'transfer'
  quantity: number
  note: string
  created_at: string
}

export interface LowStockAlert {
  product_id: string
  product_name: string
  sku: string
  warehouse_id: string
  warehouse_name: string
  quantity: number
  threshold: number
}

// ---------------------------------------------------------------------------
// Reviews (028)
// ---------------------------------------------------------------------------

export type ReviewStatus = 'pending' | 'approved' | 'rejected'

export interface Review {
  id: string
  tenant_id: string
  product_id: string
  product_name: string
  customer_id: string
  customer_name: string
  rating: number
  title: string
  body: string
  status: ReviewStatus
  admin_response: string | null
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Returns (029)
// ---------------------------------------------------------------------------

export type ReturnStatus = 'pending' | 'approved' | 'received' | 'refunded' | 'closed' | 'rejected'

export interface ReturnItem {
  id: string
  return_id: string
  product_id: string
  product_name: string
  quantity: number
  reason: string
}

export interface ReturnRequest {
  id: string
  tenant_id: string
  order_id: string
  customer_id: string
  status: ReturnStatus
  items: ReturnItem[]
  refund_amount: Money | null
  notes: string
  created_at: string
  updated_at: string
}

// ---------------------------------------------------------------------------
// Webhooks (030)
// ---------------------------------------------------------------------------

export interface Webhook {
  id: string
  tenant_id: string
  url: string
  events: string[]
  active: boolean
  secret: string
  created_at: string
  updated_at: string
}

export type WebhookDeliveryStatus = 'success' | 'failed' | 'pending'

export interface WebhookDelivery {
  id: string
  webhook_id: string
  event: string
  payload: string
  status: WebhookDeliveryStatus
  response_code: number | null
  response_body: string | null
  attempts: number
  created_at: string
  delivered_at: string | null
}

// ---------------------------------------------------------------------------
// Audit Logs (031)
// ---------------------------------------------------------------------------

export interface AuditLog {
  id: string
  tenant_id: string
  user_id: string
  user_email: string
  action: string
  resource_type: string
  resource_id: string
  changes: Record<string, unknown>
  ip_address: string
  created_at: string
}

export interface AuditStats {
  total_actions: number
  actions_by_type: Record<string, number>
  actions_by_user: Array<{ user_email: string; count: number }>
  recent_activity: number
}

// ---------------------------------------------------------------------------
// Loyalty (032)
// ---------------------------------------------------------------------------

export type RewardType = 'points_multiplier' | 'fixed_discount' | 'free_shipping' | 'free_product'

export interface LoyaltyReward {
  id: string
  tenant_id: string
  name: string
  description: string
  type: RewardType
  points_cost: number
  value: number
  active: boolean
  created_at: string
  updated_at: string
}

export type LoyaltyTier = 'bronze' | 'silver' | 'gold' | 'platinum'

export interface LoyaltyAccount {
  id: string
  tenant_id: string
  customer_id: string
  customer_name: string
  customer_email: string
  points: number
  lifetime_points: number
  tier: LoyaltyTier
  created_at: string
  updated_at: string
}

export interface LoyaltyTransaction {
  id: string
  account_id: string
  tenant_id: string
  type: 'earn' | 'redeem' | 'adjust' | 'expire'
  points: number
  note: string
  order_id: string | null
  created_at: string
}

// ---------------------------------------------------------------------------
// Bundles (033)
// ---------------------------------------------------------------------------

export interface BundleItem {
  id: string
  bundle_id: string
  product_id: string
  product_name: string
  quantity: number
  discount_percent: number
}

export interface Bundle {
  id: string
  tenant_id: string
  name: string
  description: string
  price: Money | null
  discount_percent: number
  active: boolean
  items: BundleItem[]
  created_at: string
  updated_at: string
}

export interface BundlePrice {
  original_total: Money
  discount_amount: Money
  final_price: Money
  savings_percent: number
}

// ---------------------------------------------------------------------------
// Dashboard Reporting (034)
// ---------------------------------------------------------------------------

export interface SalesOverview {
  total_revenue: Money
  total_orders: number
  average_order_value: Money
  period: string
}

export interface RevenueData {
  date: string
  revenue: number
  orders: number
}

export interface TopProductReport {
  product_id: string
  product_name: string
  units_sold: number
  revenue: number
  currency: string
}

export interface CustomerStats {
  total_customers: number
  new_customers: number
  returning_customers: number
  average_lifetime_value: Money
}

export interface FunnelStep {
  step: string
  count: number
  conversion_rate: number
}

// ---------------------------------------------------------------------------
// Social Accounts (035)
// ---------------------------------------------------------------------------

export type SocialProvider = 'google' | 'facebook'

export interface SocialAccount {
  id: string
  tenant_id: string
  customer_id: string
  customer_name: string
  customer_email: string
  provider: SocialProvider
  provider_user_id: string
  created_at: string
}

// ---------------------------------------------------------------------------
// Notifications (036)
// ---------------------------------------------------------------------------

export type NotificationSeverity = 'info' | 'warning' | 'error' | 'success'

export interface AdminNotification {
  id: string
  tenant_id: string
  title: string
  body: string
  severity: NotificationSeverity
  read: boolean
  resource_type: string | null
  resource_id: string | null
  created_at: string
  read_at: string | null
}

export interface UnreadCount {
  count: number
}

// ---------------------------------------------------------------------------
// Storefronts / Multistore (037)
// ---------------------------------------------------------------------------

export interface Storefront {
  id: string
  tenant_id: string
  name: string
  slug: string
  domain: string
  description: string
  logo_url: string
  favicon_url: string
  theme: string
  default_locale: string
  default_currency: string
  is_default: boolean
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface StorefrontCatalog {
  storefront_id: string
  catalog_id: string
  display_order: number
  created_at: string
}

// ---------------------------------------------------------------------------
// Bulk Operations (038)
// ---------------------------------------------------------------------------

export type BulkOperationStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled'
export type BulkOperationType =
  | 'price_update'
  | 'status_update'
  | 'inventory_update'
  | 'category_assign'
  | 'tag_assign'
  | 'delete'

export interface BulkOperation {
  id: string
  tenant_id: string
  type: BulkOperationType
  status: BulkOperationStatus
  total_items: number
  processed_items: number
  failed_items: number
  errors: string[]
  created_by: string
  started_at: string | null
  completed_at: string | null
  created_at: string
}

export interface BulkOperationItem {
  id: string
  operation_id: string
  resource_id: string
  status: BulkOperationStatus
  error: string | null
  processed_at: string | null
}

// ---------------------------------------------------------------------------
// Blog (039)
// ---------------------------------------------------------------------------

export type BlogPostStatus = 'draft' | 'published' | 'archived'

export interface BlogPost {
  id: string
  tenant_id: string
  title: string
  slug: string
  content: string
  excerpt: string
  featured_image_url: string
  author: string
  status: BlogPostStatus
  published_at: string | null
  tags: string[]
  seo_title: string
  seo_description: string
  created_at: string
  updated_at: string
}

export interface BlogCategory {
  id: string
  tenant_id: string
  name: string
  slug: string
  description: string
  parent_id: string | null
  sort_order: number
  created_at: string
}

// ---------------------------------------------------------------------------
// Collections — extended (040)
// ---------------------------------------------------------------------------

export type CollectionType = 'manual' | 'automated'

export interface CollectionRule {
  field: string
  operator: string
  value: string
}

export interface AdminCollection {
  id: string
  tenant_id: string
  name: string
  slug: string
  description: string
  image_url: string
  type: CollectionType
  rules: CollectionRule[]
  sort_order: number
  is_active: boolean
  product_count: number
  created_at: string
  updated_at: string
}

export interface CollectionProduct {
  collection_id: string
  product_id: string
  sort_order: number
  added_at: string
}

// ---------------------------------------------------------------------------
// A/B Testing — Experiments (041)
// ---------------------------------------------------------------------------

export type ExperimentStatus = 'draft' | 'running' | 'paused' | 'completed'

export interface ExperimentVariant {
  id: string
  experiment_id: string
  name: string
  description: string
  is_control: boolean
  traffic_weight: number
  conversions: number
  impressions: number
  revenue: number
  created_at: string
}

export interface Experiment {
  id: string
  tenant_id: string
  name: string
  description: string
  type: string
  status: ExperimentStatus
  traffic_percentage: number
  variants: ExperimentVariant[]
  started_at: string | null
  completed_at: string | null
  created_at: string
  updated_at: string
}

export interface VariantResult {
  variant_id: string
  name: string
  impressions: number
  conversions: number
  revenue: number
  conversion_rate: number
  revenue_per_visitor: number
  is_winner: boolean
}

export interface ExperimentResults {
  experiment_id: string
  variants: VariantResult[]
  total_impressions: number
}

// ---------------------------------------------------------------------------
// Recommendations (042)
// ---------------------------------------------------------------------------

export type RecommendationType = 'similar' | 'frequently_bought' | 'trending' | 'personalized'

export interface RecommendationRule {
  id: string
  tenant_id: string
  name: string
  type: RecommendationType
  source_product_id: string | null
  weight: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface RecommendedProduct {
  product_id: string
  name: string
  score: number
  reason: string
}

// ---------------------------------------------------------------------------
// Agent Presets (050)
// ---------------------------------------------------------------------------

export type PresetStatus = 'draft' | 'published' | 'archived'
export type PresetVisibility = 'public' | 'private' | 'tenant_only'
export type PresetCategory = 'webdev' | 'research' | 'content' | 'analytics' | 'store-manager'

export interface Preset {
  id: string
  name: string
  slug: string
  description: string
  version: string
  image: string
  frontend_port: number
  system_prompt: string
  tools_manifest: unknown[]
  status: PresetStatus
  visibility: PresetVisibility
  category: PresetCategory
  author: string
  icon_url?: string
  created_at: string
  updated_at: string
}

export interface PresetInstall {
  id: string
  tenant_id: string
  preset_id: string
  installed_at: string
}

// Agent Sessions (Workspaces)
export type SessionStatus = 'creating' | 'running' | 'stopped' | 'failed'

export interface AgentSession {
  id: string
  tenant_id: string
  preset_id: string
  name: string
  container_id: string
  status: SessionStatus
  frontend_url: string
  metadata: Record<string, unknown>
  created_at: string
  updated_at: string
  stopped_at?: string
}

export interface ChatMessage {
  id: string
  session_id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  tool_calls?: unknown[]
  created_at: string
}

// Agent SSE Events
export interface AgentEvent {
  kind: 'text_delta' | 'tool_start' | 'tool_end' | 'turn_end' | 'error'
  text?: string
  tool_name?: string
  tool_input?: string
  result?: string
  error?: string
  timestamp: string
}
