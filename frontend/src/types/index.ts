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
