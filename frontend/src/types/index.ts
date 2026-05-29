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
  shipping_address: Address
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

export interface PageMeta {
  description: string
  og_title: string
  og_image: string
  keywords: string[]
}

export interface Page {
  id: string
  tenant_id: string
  slug: string
  title: string
  html: string
  css: string
  meta: PageMeta
  status: PageStatus
  version: number
  created_by: string
  published_at?: string
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
