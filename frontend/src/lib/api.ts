import type {
  Product,
  Order,
  OrderStatus,
  Customer,
  Category,
  Collection,
  Page,
  PageVersion,
  Promo,
  Media,
  PaginatedResult,
  Plugin,
  PluginVersion,
  PluginInstallation,
  PluginManifest,
  DashboardStats,
  RevenuePoint,
  TopProduct,
  OrderStatusBreakdown,
  RecentOrder,
  StoreSettings,
  LoginResponse,
  RefreshResponse,
  TokenResponse,
  MeResponse,
  BlockType,
  Theme,
  ShippingZone,
  ShippingRate,
  TaxRate,
  Payment,
  Refund,
  ProductOption,
  ProductVariant,
  ImportResult,
  CustomerGroup,
  GroupRules,
  GroupMembership,
  GiftCard,
  GiftCardTransaction,
  RecoveryEmail,
  RecoveryStats,
  CurrencyRate,
  ConvertResult,
  TranslationBundle,
  Subscription,
  BillingRecord,
  Warehouse,
  StockLevel,
  StockMovement,
  LowStockAlert,
  Review,
  ReviewStatus,
  ReturnRequest,
  ReturnStatus,
  Webhook,
  WebhookDelivery,
  AuditLog,
  AuditStats,
  LoyaltyReward,
  LoyaltyAccount,
  LoyaltyTransaction,
  Bundle,
  BundleItem,
  BundlePrice,
  SalesOverview,
  RevenueData,
  TopProductReport,
  CustomerStats,
  FunnelStep,
  SocialAccount,
  SocialProvider,
  AdminNotification,
  UnreadCount,
  Storefront,
  StorefrontCatalog,
  BulkOperation,
  BulkOperationItem,
  BlogPost,
  BlogCategory,
  AdminCollection,
  CollectionProduct,
  Experiment,
  ExperimentVariant,
  ExperimentResults,
  RecommendationRule,
  RecommendedProduct,
  Preset,
  PresetInstall,
  AgentSession,
  ChatMessage,
  ApprovalRequest,
  AgentMemory,
  CreateMemoryRequest,
  UpdateMemoryRequest,
  AgentTrigger,
  CreateTriggerRequest,
  UpdateTriggerRequest,
  TriggerLog,
  PasswordlessTenantsResponse,
  PasswordlessInitiateResponse,
  PasswordlessVerifyResponse,
  PasswordlessResendResponse,
} from '../types'

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api/v1'

export interface PaginationParams {
  page?: number
  page_size?: number
}

// ---------------------------------------------------------------------------
// Token management (localStorage)
// ---------------------------------------------------------------------------

const TOKEN_KEY = 'hada_access_token'
const REFRESH_TOKEN_KEY = 'hada_refresh_token'
const TENANT_ID_KEY = 'hada_tenant_id'

export function getAccessToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY)
}

export function getTenantId(): string | null {
  return localStorage.getItem(TENANT_ID_KEY)
}

export function setTenantId(tenantId: string): void {
  localStorage.setItem(TENANT_ID_KEY, tenantId)
}

export function setTokens(accessToken: string, refreshToken?: string): void {
  localStorage.setItem(TOKEN_KEY, accessToken)
  if (refreshToken) {
    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
  }
}

export function clearTokens(): void {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
  localStorage.removeItem(TENANT_ID_KEY)
}

// ---------------------------------------------------------------------------
// Core fetch helpers
// ---------------------------------------------------------------------------

function buildHeaders(extra?: Record<string, string>): HeadersInit {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...extra,
  }
  const token = getAccessToken()
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  const tenantId = getTenantId()
  if (tenantId) {
    headers['X-Tenant-ID'] = tenantId
  }
  return headers
}

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    let message = `HTTP ${res.status} ${res.statusText}`
    try {
      const body = await res.json()
      if (body.error) message = body.error
      else if (body.message) message = body.message
    } catch {
      // ignore parse errors — keep the status message
    }
    throw new Error(message)
  }
  // 204 No Content — nothing to parse
  if (res.status === 204) return undefined as unknown as T
  return res.json() as Promise<T>
}

// Attempt a token refresh and retry the original request once
let _refreshPromise: Promise<string> | null = null

async function withRefreshRetry<T>(fn: () => Promise<Response>): Promise<T> {
  const res = await fn()
  if (res.status !== 401) {
    return handleResponse<T>(res)
  }

  // 401 — try refreshing the token (deduplicate concurrent refreshes)
  if (!_refreshPromise) {
    const refreshToken = getRefreshToken()
    if (!refreshToken) {
      // No refresh token — clear state and throw
      clearTokens()
      throw new Error('Session expired. Please log in again.')
    }
    _refreshPromise = refreshTokenApi(refreshToken)
      .then((newToken) => {
        setTokens(newToken)
        return newToken
      })
      .finally(() => {
        _refreshPromise = null
      })
  }

  try {
    await _refreshPromise
  } catch {
    clearTokens()
    throw new Error('Session expired. Please log in again.')
  }

  // Retry with new token
  const retried = await fn()
  return handleResponse<T>(retried)
}

export async function get<T>(path: string, params?: Record<string, string | number | undefined>): Promise<T> {
  const url = new URL(`${BASE_URL}${path}`, window.location.origin)
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (v !== undefined) url.searchParams.set(k, String(v))
    }
  }
  return withRefreshRetry<T>(() =>
    fetch(url.toString(), {
      method: 'GET',
      headers: buildHeaders(),
    }),
  )
}

export async function post<T>(path: string, body?: unknown): Promise<T> {
  return withRefreshRetry<T>(() =>
    fetch(`${BASE_URL}${path}`, {
      method: 'POST',
      headers: buildHeaders(),
      body: body !== undefined ? JSON.stringify(body) : undefined,
    }),
  )
}

export async function put<T>(path: string, body?: unknown): Promise<T> {
  return withRefreshRetry<T>(() =>
    fetch(`${BASE_URL}${path}`, {
      method: 'PUT',
      headers: buildHeaders(),
      body: body !== undefined ? JSON.stringify(body) : undefined,
    }),
  )
}

export async function del(path: string): Promise<void> {
  return withRefreshRetry<void>(() =>
    fetch(`${BASE_URL}${path}`, {
      method: 'DELETE',
      headers: buildHeaders(),
    }),
  )
}

// ---------------------------------------------------------------------------
// Auth API — these do NOT use withRefreshRetry to avoid infinite loops
// ---------------------------------------------------------------------------

export async function initiateLogin(provider: 'GOOGLE' | 'MICROSOFT', invitationToken?: string): Promise<LoginResponse> {
  const body: { provider: string; invitation_token?: string } = { provider }
  if (invitationToken) body.invitation_token = invitationToken
  const res = await fetch(`/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return handleResponse<LoginResponse>(res)
}

export async function handleAuthCallback(provider: string, code: string, state: string): Promise<TokenResponse> {
  const url = new URL(`/auth/callback/${provider}`, window.location.origin)
  url.searchParams.set('code', code)
  url.searchParams.set('state', state)
  const res = await fetch(url.toString(), {
    method: 'GET',
    credentials: 'include', // include cookies
  })
  return handleResponse<TokenResponse>(res)
}

// Internal helper used only by withRefreshRetry — returns new access token
async function refreshTokenApi(refreshToken: string): Promise<string> {
  const res = await fetch(`/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken }),
  })
  const data = await handleResponse<RefreshResponse>(res)
  return data.access_token
}

// Exported refresh function for explicit use
export async function refreshToken(token: string): Promise<RefreshResponse> {
  const res = await fetch(`/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: token }),
  })
  return handleResponse<RefreshResponse>(res)
}

export async function logout(): Promise<{ message: string }> {
  const res = await fetch(`/auth/logout`, {
    method: 'POST',
    headers: buildHeaders(),
  })
  return handleResponse<{ message: string }>(res)
}

export async function getMe(): Promise<MeResponse> {
  const res = await fetch(`/auth/me`, {
    method: 'GET',
    headers: buildHeaders(),
  })
  return handleResponse<MeResponse>(res)
}

// ---------------------------------------------------------------------------
// Passwordless / OTP Auth — these do NOT use withRefreshRetry (pre-auth flows)
// ---------------------------------------------------------------------------

export async function getPasswordlessTenants(email: string): Promise<PasswordlessTenantsResponse> {
  const res = await fetch(`/auth/passwordless/tenants`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email }),
  })
  return handleResponse<PasswordlessTenantsResponse>(res)
}

export async function initiatePasswordlessLogin(
  email: string,
  tenantId: string,
): Promise<PasswordlessInitiateResponse> {
  const res = await fetch(`/auth/passwordless/login/initiate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, tenant_id: tenantId }),
  })
  return handleResponse<PasswordlessInitiateResponse>(res)
}

export async function verifyPasswordlessLogin(
  email: string,
  code: string,
  tenantId: string,
): Promise<PasswordlessVerifyResponse> {
  const res = await fetch(`/auth/passwordless/login/verify`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, code, tenant_id: tenantId }),
  })
  return handleResponse<PasswordlessVerifyResponse>(res)
}

export async function resendOTP(email: string, tenantId: string): Promise<PasswordlessResendResponse> {
  const res = await fetch(`/auth/passwordless/resend-otp`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, tenant_id: tenantId, purpose: 'login' }),
  })
  return handleResponse<PasswordlessResendResponse>(res)
}

// ---------------------------------------------------------------------------
// Products
// ---------------------------------------------------------------------------

export function listProducts(params?: PaginationParams): Promise<PaginatedResult<Product>> {
  return get<PaginatedResult<Product>>('/products', params as Record<string, string | number | undefined>)
}

export function getProduct(id: string): Promise<Product> {
  return get<Product>(`/products/${id}`)
}

export function createProduct(data: Partial<Product>): Promise<Product> {
  return post<Product>('/products', data)
}

export function updateProduct(id: string, data: Partial<Product>): Promise<Product> {
  return put<Product>(`/products/${id}`, data)
}

export function deleteProduct(id: string): Promise<void> {
  return del(`/products/${id}`)
}

// ---------------------------------------------------------------------------
// Orders
// ---------------------------------------------------------------------------

export function listOrders(params?: PaginationParams): Promise<PaginatedResult<Order>> {
  return get<PaginatedResult<Order>>('/orders', params as Record<string, string | number | undefined>)
}

export function getOrder(id: string): Promise<Order> {
  return get<Order>(`/orders/${id}`)
}

export function updateOrderStatus(id: string, status: OrderStatus): Promise<Order> {
  return put<Order>(`/orders/${id}/status`, { status })
}

export function cancelOrder(id: string): Promise<Order> {
  return post<Order>(`/orders/${id}/cancel`)
}

// ---------------------------------------------------------------------------
// Customers
// ---------------------------------------------------------------------------

export function listCustomers(params?: PaginationParams): Promise<PaginatedResult<Customer>> {
  return get<PaginatedResult<Customer>>('/customers', params as Record<string, string | number | undefined>)
}

export function getCustomer(id: string): Promise<Customer> {
  return get<Customer>(`/customers/${id}`)
}

export function createCustomer(data: Partial<Customer>): Promise<Customer> {
  return post<Customer>('/customers', data)
}

export function deleteCustomer(id: string): Promise<void> {
  return del(`/customers/${id}`)
}

export function listOrdersByCustomer(customerId: string, params?: PaginationParams): Promise<PaginatedResult<Order>> {
  return get<PaginatedResult<Order>>('/orders', {
    ...(params as Record<string, string | number | undefined>),
    customer_id: customerId,
  })
}

// ---------------------------------------------------------------------------
// Catalog — Categories
// ---------------------------------------------------------------------------

export function listCategories(params?: PaginationParams): Promise<PaginatedResult<Category>> {
  return get<PaginatedResult<Category>>('/categories', params as Record<string, string | number | undefined>)
}

export function getCategory(id: string): Promise<Category> {
  return get<Category>(`/categories/${id}`)
}

export function createCategory(data: Partial<Category>): Promise<Category> {
  return post<Category>('/categories', data)
}

export function deleteCategory(id: string): Promise<void> {
  return del(`/categories/${id}`)
}

// ---------------------------------------------------------------------------
// Catalog — Collections
// ---------------------------------------------------------------------------

export function listCollections(params?: PaginationParams): Promise<PaginatedResult<Collection>> {
  return get<PaginatedResult<Collection>>('/collections', params as Record<string, string | number | undefined>)
}

export function getCollection(id: string): Promise<Collection> {
  return get<Collection>(`/collections/${id}`)
}

export function createCollection(data: Partial<Collection>): Promise<Collection> {
  return post<Collection>('/collections', data)
}

export function deleteCollection(id: string): Promise<void> {
  return del(`/collections/${id}`)
}

// ---------------------------------------------------------------------------
// Storefront — Pages
// ---------------------------------------------------------------------------

export function listPages(params?: PaginationParams): Promise<PaginatedResult<Page>> {
  return get<PaginatedResult<Page>>('/storefront/pages', params as Record<string, string | number | undefined>)
}

export function getPage(id: string): Promise<Page> {
  return get<Page>(`/storefront/pages/${id}`)
}

export function getPageBySlug(slug: string): Promise<Page> {
  return get<Page>(`/storefront/pages/slug/${slug}`)
}

export function createPage(data: Partial<Page>): Promise<Page> {
  return post<Page>('/storefront/pages', data)
}

export function updatePage(id: string, data: Partial<Page>): Promise<Page> {
  return put<Page>(`/storefront/pages/${id}`, data)
}

export function publishPage(id: string): Promise<Page> {
  return post<Page>(`/storefront/pages/${id}/publish`)
}

export function unpublishPage(id: string): Promise<Page> {
  return post<Page>(`/storefront/pages/${id}/unpublish`)
}

export function archivePage(id: string): Promise<Page> {
  return post<Page>(`/storefront/pages/${id}/archive`)
}

export function deletePage(id: string): Promise<void> {
  return del(`/storefront/pages/${id}`)
}

export function getPageVersions(id: string): Promise<PageVersion[]> {
  return get<PageVersion[]>(`/storefront/pages/${id}/versions`)
}

// ---------------------------------------------------------------------------
// Promos
// ---------------------------------------------------------------------------

export function listPromos(params?: PaginationParams): Promise<PaginatedResult<Promo>> {
  return get<PaginatedResult<Promo>>('/admin/promos', params as Record<string, string | number | undefined>)
}

export function createPromo(data: Partial<Promo>): Promise<Promo> {
  return post<Promo>('/admin/promos', data)
}

export function validatePromo(code: string, orderAmount: number): Promise<{ valid: boolean; discount: number }> {
  return post<{ valid: boolean; discount: number }>('/promos/validate', { code, order_amount: orderAmount })
}

export function updatePromo(id: string, data: Partial<Promo>): Promise<Promo> {
  return put<Promo>(`/admin/promos/${id}`, data)
}

export function deactivatePromo(id: string): Promise<Promo> {
  return post<Promo>(`/admin/promos/${id}/deactivate`)
}

// ---------------------------------------------------------------------------
// Media
// ---------------------------------------------------------------------------

export function listMedia(params?: PaginationParams): Promise<PaginatedResult<Media>> {
  return get<PaginatedResult<Media>>('/media', params as Record<string, string | number | undefined>)
}

export async function uploadMedia(file: File, alt?: string): Promise<Media> {
  const formData = new FormData()
  formData.append('file', file)
  if (alt) formData.append('alt', alt)

  const headers: Record<string, string> = {
    // Do NOT set Content-Type — browser sets multipart/form-data with boundary automatically
  }
  const token = getAccessToken()
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${BASE_URL}/media/upload`, {
    method: 'POST',
    headers,
    body: formData,
  })
  return handleResponse<Media>(res)
}

export function deleteMedia(id: string): Promise<void> {
  return del(`/media/${id}`)
}

// ---------------------------------------------------------------------------
// Plugins (marketplace + runtime)
// ---------------------------------------------------------------------------

export function listMarketplacePlugins(params?: PaginationParams): Promise<PaginatedResult<Plugin>> {
  return get<PaginatedResult<Plugin>>('/plugins', params as Record<string, string | number | undefined>)
}

export function getMarketplacePlugin(id: string): Promise<Plugin & { latest_version?: PluginVersion }> {
  return get<Plugin & { latest_version?: PluginVersion }>(`/plugins/${id}`)
}

export function installPlugin(pluginId: string): Promise<PluginInstallation> {
  return post<PluginInstallation>('/plugins/install', { plugin_id: pluginId })
}

export function uninstallPlugin(pluginId: string): Promise<void> {
  return post<void>(`/plugins/${pluginId}/uninstall`)
}

export function listInstalledPlugins(params?: PaginationParams): Promise<PaginatedResult<PluginInstallation>> {
  return get<PaginatedResult<PluginInstallation>>('/plugins/installed', params as Record<string, string | number | undefined>)
}

export function updatePluginSettings(pluginId: string, settings: Record<string, unknown>): Promise<PluginInstallation> {
  return put<PluginInstallation>(`/plugins/${pluginId}/settings`, settings)
}

export function enablePlugin(pluginId: string): Promise<PluginInstallation> {
  return post<PluginInstallation>(`/plugins/${pluginId}/enable`)
}

export function disablePlugin(pluginId: string): Promise<PluginInstallation> {
  return post<PluginInstallation>(`/plugins/${pluginId}/disable`)
}

export function getPluginManifest(id: string): Promise<{ manifest: PluginManifest }> {
  return get<{ manifest: PluginManifest }>(`/plugins/${id}/manifest`)
}

// ---------------------------------------------------------------------------
// Analytics
// ---------------------------------------------------------------------------

export function getDashboardStats(): Promise<DashboardStats> {
  return get<DashboardStats>('/analytics/dashboard')
}

export function getRevenueTimeline(days: number = 30): Promise<RevenuePoint[]> {
  return get<RevenuePoint[]>('/analytics/revenue', { days })
}

export function getTopProducts(limit: number = 5): Promise<TopProduct[]> {
  return get<TopProduct[]>('/analytics/top-products', { limit })
}

export function getOrderStatusBreakdown(): Promise<OrderStatusBreakdown[]> {
  return get<OrderStatusBreakdown[]>('/analytics/order-status')
}

export function getRecentOrders(limit: number = 5): Promise<RecentOrder[]> {
  return get<RecentOrder[]>('/analytics/recent-orders', { limit })
}

// ---------------------------------------------------------------------------
// Settings
// ---------------------------------------------------------------------------

export function getSettings(): Promise<StoreSettings> {
  return get<StoreSettings>('/settings')
}

export function updateSettings(data: Partial<StoreSettings>): Promise<StoreSettings> {
  return put<StoreSettings>('/settings', data)
}

// ---------------------------------------------------------------------------
// Block Types
// ---------------------------------------------------------------------------

export function listBlockTypes(category?: string): Promise<BlockType[]> {
  return get<BlockType[]>('/storefront/block-types', category ? { category } : undefined)
}

export function getBlockType(id: string): Promise<BlockType> {
  return get<BlockType>(`/storefront/block-types/${id}`)
}

export function createBlockType(data: Partial<BlockType>): Promise<BlockType> {
  return post<BlockType>('/storefront/block-types', data)
}

export function updateBlockType(id: string, data: Partial<BlockType>): Promise<BlockType> {
  return put<BlockType>(`/storefront/block-types/${id}`, data)
}

export function deleteBlockType(id: string): Promise<void> {
  return del(`/storefront/block-types/${id}`)
}

// ---------------------------------------------------------------------------
// Themes
// ---------------------------------------------------------------------------

export function listThemes(): Promise<Theme[]> {
  return get<Theme[]>('/themes')
}

export function getActiveTheme(): Promise<Theme> {
  return get<Theme>('/themes/active')
}

export function getTheme(id: string): Promise<Theme> {
  return get<Theme>(`/themes/${id}`)
}

export function createTheme(data: Partial<Theme>): Promise<Theme> {
  return post<Theme>('/themes', data)
}

export function updateTheme(id: string, data: Partial<Theme>): Promise<Theme> {
  return put<Theme>(`/themes/${id}`, data)
}

export function activateTheme(id: string): Promise<Theme> {
  return post<Theme>(`/themes/${id}/activate`)
}

export function duplicateTheme(id: string, name: string): Promise<Theme> {
  return post<Theme>(`/themes/${id}/duplicate`, { name })
}

export function deleteTheme(id: string): Promise<void> {
  return del(`/themes/${id}`)
}

// ---------------------------------------------------------------------------
// Shipping — Zones
// ---------------------------------------------------------------------------

export function listShippingZones(): Promise<ShippingZone[]> {
  return get<ShippingZone[]>('/shipping/zones')
}

export function getShippingZone(id: string): Promise<ShippingZone> {
  return get<ShippingZone>(`/shipping/zones/${id}`)
}

export function createShippingZone(data: Partial<ShippingZone>): Promise<ShippingZone> {
  return post<ShippingZone>('/shipping/zones', data)
}

export function updateShippingZone(id: string, data: Partial<ShippingZone>): Promise<ShippingZone> {
  return put<ShippingZone>(`/shipping/zones/${id}`, data)
}

export function deleteShippingZone(id: string): Promise<void> {
  return del(`/shipping/zones/${id}`)
}

// ---------------------------------------------------------------------------
// Shipping — Rates
// ---------------------------------------------------------------------------

export function listShippingRates(zoneId: string): Promise<ShippingRate[]> {
  return get<ShippingRate[]>(`/shipping/zones/${zoneId}/rates`)
}

export function createShippingRate(zoneId: string, data: Partial<ShippingRate>): Promise<ShippingRate> {
  return post<ShippingRate>(`/shipping/zones/${zoneId}/rates`, data)
}

export function updateShippingRate(id: string, data: Partial<ShippingRate>): Promise<ShippingRate> {
  return put<ShippingRate>(`/shipping/rates/${id}`, data)
}

export function deleteShippingRate(id: string): Promise<void> {
  return del(`/shipping/rates/${id}`)
}

// ---------------------------------------------------------------------------
// Tax — Rates
// ---------------------------------------------------------------------------

export function listTaxRates(): Promise<TaxRate[]> {
  return get<TaxRate[]>('/tax/rates')
}

export function createTaxRate(data: Partial<TaxRate>): Promise<TaxRate> {
  return post<TaxRate>('/tax/rates', data)
}

export function updateTaxRate(id: string, data: Partial<TaxRate>): Promise<TaxRate> {
  return put<TaxRate>(`/tax/rates/${id}`, data)
}

export function deleteTaxRate(id: string): Promise<void> {
  return del(`/tax/rates/${id}`)
}

// ---------------------------------------------------------------------------
// Payments
// ---------------------------------------------------------------------------

export function getPayment(id: string): Promise<Payment> {
  return get<Payment>(`/payments/${id}`)
}

export function listOrderPayments(orderId: string): Promise<Payment[]> {
  return get<Payment[]>(`/payments/order/${orderId}`)
}

export function createPayment(data: {
  order_id: string
  amount: number
  currency: string
  provider: string
  method: string
}): Promise<Payment> {
  return post<Payment>('/payments', data)
}

export function processPayment(id: string, data?: { token?: string }): Promise<Payment> {
  return post<Payment>(`/payments/${id}/process`, data ?? {})
}

export function createRefund(paymentId: string, data: { amount: number; reason?: string }): Promise<Refund> {
  return post<Refund>(`/payments/${paymentId}/refund`, data)
}

export function listRefunds(paymentId: string): Promise<Refund[]> {
  return get<Refund[]>(`/payments/${paymentId}/refunds`)
}

// ---------------------------------------------------------------------------
// Product — Options
// ---------------------------------------------------------------------------

export function listProductOptions(productId: string): Promise<ProductOption[]> {
  return get<ProductOption[]>(`/products/${productId}/options`)
}

export function createProductOption(
  productId: string,
  data: { name: string; position: number; values: string[] },
): Promise<ProductOption> {
  return post<ProductOption>(`/products/${productId}/options`, data)
}

export function updateProductOption(optionId: string, data: Partial<ProductOption>): Promise<ProductOption> {
  return put<ProductOption>(`/products/options/${optionId}`, data)
}

export function deleteProductOption(optionId: string): Promise<void> {
  return del(`/products/options/${optionId}`)
}

// ---------------------------------------------------------------------------
// Product — Variants
// ---------------------------------------------------------------------------

export function listProductVariants(productId: string): Promise<ProductVariant[]> {
  return get<ProductVariant[]>(`/products/${productId}/variants`)
}

export function createProductVariant(
  productId: string,
  data: {
    sku: string
    price_amount: number
    price_currency: string
    stock: number
    options: Record<string, string>
  },
): Promise<ProductVariant> {
  return post<ProductVariant>(`/products/${productId}/variants`, data)
}

export function updateProductVariant(variantId: string, data: Partial<ProductVariant>): Promise<ProductVariant> {
  return put<ProductVariant>(`/products/variants/${variantId}`, data)
}

export function deleteProductVariant(variantId: string): Promise<void> {
  return del(`/products/variants/${variantId}`)
}

// ---------------------------------------------------------------------------
// Import / Export
// ---------------------------------------------------------------------------

export async function exportProducts(): Promise<Blob> {
  const res = await fetch(`${BASE_URL}/import-export/products/export`, {
    headers: buildHeaders(),
  })
  if (!res.ok) throw new Error('Export failed')
  return res.blob()
}

export async function exportOrders(): Promise<Blob> {
  const res = await fetch(`${BASE_URL}/import-export/orders/export`, {
    headers: buildHeaders(),
  })
  if (!res.ok) throw new Error('Export failed')
  return res.blob()
}

export async function exportCustomers(): Promise<Blob> {
  const res = await fetch(`${BASE_URL}/import-export/customers/export`, {
    headers: buildHeaders(),
  })
  if (!res.ok) throw new Error('Export failed')
  return res.blob()
}

export async function importProducts(file: File): Promise<ImportResult> {
  const formData = new FormData()
  formData.append('file', file)
  // Note: do NOT set Content-Type here — browser sets it with the boundary for multipart
  const headers: Record<string, string> = {}
  const token = getAccessToken()
  if (token) headers['Authorization'] = `Bearer ${token}`
  const res = await fetch(`${BASE_URL}/import-export/products/import`, {
    method: 'POST',
    headers,
    body: formData,
  })
  if (!res.ok) {
    let message = `HTTP ${res.status} ${res.statusText}`
    try {
      const body = await res.json()
      if (body.error) message = body.error
      else if (body.message) message = body.message
    } catch { /* ignore */ }
    throw new Error(message)
  }
  return res.json() as Promise<ImportResult>
}

// ---------------------------------------------------------------------------
// Customer Groups
// ---------------------------------------------------------------------------

export interface CustomerGroupInput {
  name: string
  description: string
  rules: GroupRules
}

export function listCustomerGroups(): Promise<CustomerGroup[]> {
  return get<CustomerGroup[]>('/customer-groups/')
}

export function getCustomerGroup(id: string): Promise<CustomerGroup> {
  return get<CustomerGroup>(`/customer-groups/${id}`)
}

export function createCustomerGroup(data: CustomerGroupInput): Promise<CustomerGroup> {
  return post<CustomerGroup>('/customer-groups/', data)
}

export function updateCustomerGroup(id: string, data: CustomerGroupInput): Promise<CustomerGroup> {
  return put<CustomerGroup>(`/customer-groups/${id}`, data)
}

export function deleteCustomerGroup(id: string): Promise<void> {
  return del(`/customer-groups/${id}`)
}

export function listGroupMembers(groupId: string): Promise<GroupMembership[]> {
  return get<GroupMembership[]>(`/customer-groups/${groupId}/members`)
}

export function addGroupMember(groupId: string, customerId: string): Promise<GroupMembership> {
  return post<GroupMembership>(`/customer-groups/${groupId}/members`, { customer_id: customerId })
}

export function removeGroupMember(groupId: string, customerId: string): Promise<void> {
  return del(`/customer-groups/${groupId}/members/${customerId}`)
}

// ---------------------------------------------------------------------------
// Gift Cards
// ---------------------------------------------------------------------------

export interface CreateGiftCardInput {
  code?: string
  initial_amount: number
  currency: string
  expires_at?: string
}

export function listGiftCards(params?: PaginationParams): Promise<PaginatedResult<GiftCard>> {
  return get<PaginatedResult<GiftCard>>('/admin/gift-cards', params as Record<string, string | number | undefined>)
}

export function getGiftCard(id: string): Promise<GiftCard> {
  return get<GiftCard>(`/admin/gift-cards/${id}`)
}

export function createGiftCard(data: CreateGiftCardInput): Promise<GiftCard> {
  return post<GiftCard>('/admin/gift-cards', data)
}

export function updateGiftCard(id: string, data: Partial<GiftCard>): Promise<GiftCard> {
  return put<GiftCard>(`/admin/gift-cards/${id}`, data)
}

export function deleteGiftCard(id: string): Promise<void> {
  return del(`/admin/gift-cards/${id}`)
}

export function disableGiftCard(id: string): Promise<GiftCard> {
  return post<GiftCard>(`/admin/gift-cards/${id}/disable`)
}

export function listGiftCardTransactions(id: string): Promise<GiftCardTransaction[]> {
  return get<GiftCardTransaction[]>(`/admin/gift-cards/${id}/transactions`)
}

// ---------------------------------------------------------------------------
// Cart Recovery
// ---------------------------------------------------------------------------

export function listRecoveryEmails(params?: PaginationParams): Promise<PaginatedResult<RecoveryEmail>> {
  return get<PaginatedResult<RecoveryEmail>>('/cart-recovery', params as Record<string, string | number | undefined>)
}

export function getRecoveryStats(): Promise<RecoveryStats> {
  return get<RecoveryStats>('/cart-recovery/stats')
}

export function updateRecoveryStatus(id: string, status: string): Promise<RecoveryEmail> {
  return put<RecoveryEmail>(`/cart-recovery/${id}/status`, { status })
}

// ---------------------------------------------------------------------------
// Currency Rates
// ---------------------------------------------------------------------------

export interface SetCurrencyRateInput {
  base_currency: string
  target_currency: string
  rate: number
  auto_update?: boolean
}

export function listCurrencyRates(): Promise<CurrencyRate[]> {
  return get<CurrencyRate[]>('/currency-rates')
}

export function setCurrencyRate(data: SetCurrencyRateInput): Promise<CurrencyRate> {
  return post<CurrencyRate>('/currency-rates', data)
}

export function deleteCurrencyRate(id: string): Promise<void> {
  return del(`/currency-rates/${id}`)
}

export function convertCurrency(data: { amount: number; currency: string; target_currency: string }): Promise<ConvertResult> {
  return post<ConvertResult>('/currency/convert', data)
}

export function listSupportedCurrencies(): Promise<string[]> {
  return get<string[]>('/currencies')
}

// ---------------------------------------------------------------------------
// Translations / I18n
// ---------------------------------------------------------------------------

export function getTranslationBundle(entityType: string, entityId: string, locale: string): Promise<TranslationBundle> {
  return get<TranslationBundle>(`/i18n/${entityType}/${entityId}/${locale}`)
}

export function setTranslations(
  entityType: string,
  entityId: string,
  locale: string,
  fields: Record<string, string>,
): Promise<TranslationBundle> {
  return put<TranslationBundle>(`/i18n/${entityType}/${entityId}/${locale}`, { fields })
}

export function listEntityLocales(entityType: string, entityId: string): Promise<string[]> {
  return get<string[]>(`/i18n/${entityType}/${entityId}/locales`)
}

export function deleteTranslationField(entityType: string, entityId: string, locale: string, field: string): Promise<void> {
  return del(`/i18n/${entityType}/${entityId}/${locale}/${field}`)
}

export function deleteAllTranslations(entityType: string, entityId: string): Promise<void> {
  return del(`/i18n/${entityType}/${entityId}`)
}

export function listSupportedLocales(): Promise<string[]> {
  return get<string[]>('/i18n/supported-locales')
}

// ---------------------------------------------------------------------------
// Subscriptions
// ---------------------------------------------------------------------------

export interface CreateSubscriptionInput {
  customer_id: string
  product_id: string
  variant_id?: string
  price: { amount: number; currency: string }
  interval: 'weekly' | 'monthly' | 'quarterly' | 'yearly'
  trial_ends_at?: string
}

export function listSubscriptions(params?: PaginationParams): Promise<PaginatedResult<Subscription>> {
  return get<PaginatedResult<Subscription>>('/subscriptions/', params as Record<string, string | number | undefined>)
}

export function listDueSubscriptions(): Promise<Subscription[]> {
  return get<Subscription[]>('/subscriptions/due')
}

export function getSubscription(id: string): Promise<Subscription> {
  return get<Subscription>(`/subscriptions/${id}`)
}

export function createSubscription(data: CreateSubscriptionInput): Promise<Subscription> {
  return post<Subscription>('/subscriptions/', data)
}

export function cancelSubscription(id: string): Promise<Subscription> {
  return post<Subscription>(`/subscriptions/${id}/cancel`)
}

export function pauseSubscription(id: string): Promise<Subscription> {
  return post<Subscription>(`/subscriptions/${id}/pause`)
}

export function resumeSubscription(id: string): Promise<Subscription> {
  return post<Subscription>(`/subscriptions/${id}/resume`)
}

export function listBillingRecords(subscriptionId: string, params?: PaginationParams): Promise<PaginatedResult<BillingRecord>> {
  return get<PaginatedResult<BillingRecord>>(`/subscriptions/${subscriptionId}/billing`, params as Record<string, string | number | undefined>)
}

// ---------------------------------------------------------------------------
// Inventory — Warehouses (027)
// ---------------------------------------------------------------------------

export function listWarehouses(): Promise<Warehouse[]> {
  return get<Warehouse[]>('/inventory/warehouses')
}

export function createWarehouse(data: Partial<Warehouse>): Promise<Warehouse> {
  return post<Warehouse>('/inventory/warehouses', data)
}

export function updateWarehouse(id: string, data: Partial<Warehouse>): Promise<Warehouse> {
  return put<Warehouse>(`/inventory/warehouses/${id}`, data)
}

export function deleteWarehouse(id: string): Promise<void> {
  return del(`/inventory/warehouses/${id}`)
}

export function getStockLevels(productId: string): Promise<StockLevel[]> {
  return get<StockLevel[]>(`/inventory/stock/${productId}`)
}

export function adjustStock(data: { product_id: string; warehouse_id: string; quantity: number; note?: string }): Promise<StockMovement> {
  return post<StockMovement>('/inventory/stock/adjust', data)
}

export function getLowStockAlerts(): Promise<LowStockAlert[]> {
  return get<LowStockAlert[]>('/inventory/stock/low')
}

export function getStockMovements(productId: string): Promise<StockMovement[]> {
  return get<StockMovement[]>(`/inventory/movements/${productId}`)
}

// ---------------------------------------------------------------------------
// Reviews (028)
// ---------------------------------------------------------------------------

export function listReviews(params?: PaginationParams & { status?: ReviewStatus }): Promise<PaginatedResult<Review>> {
  return get<PaginatedResult<Review>>('/reviews', params as Record<string, string | number | undefined>)
}

export function getReview(id: string): Promise<Review> {
  return get<Review>(`/reviews/${id}`)
}

export function approveReview(id: string): Promise<Review> {
  return put<Review>(`/reviews/${id}/approve`)
}

export function rejectReview(id: string): Promise<Review> {
  return put<Review>(`/reviews/${id}/reject`)
}

export function respondToReview(id: string, response: string): Promise<Review> {
  return post<Review>(`/reviews/${id}/respond`, { response })
}

// ---------------------------------------------------------------------------
// Returns (029)
// ---------------------------------------------------------------------------

export function listReturns(params?: PaginationParams & { status?: ReturnStatus }): Promise<PaginatedResult<ReturnRequest>> {
  return get<PaginatedResult<ReturnRequest>>('/returns', params as Record<string, string | number | undefined>)
}

export function getReturn(id: string): Promise<ReturnRequest> {
  return get<ReturnRequest>(`/returns/${id}`)
}

export function approveReturn(id: string): Promise<ReturnRequest> {
  return put<ReturnRequest>(`/returns/${id}/approve`)
}

export function rejectReturn(id: string): Promise<ReturnRequest> {
  return put<ReturnRequest>(`/returns/${id}/reject`)
}

export function markReturnReceived(id: string): Promise<ReturnRequest> {
  return put<ReturnRequest>(`/returns/${id}/received`)
}

export function markReturnRefunded(id: string): Promise<ReturnRequest> {
  return put<ReturnRequest>(`/returns/${id}/refunded`)
}

export function closeReturn(id: string): Promise<ReturnRequest> {
  return put<ReturnRequest>(`/returns/${id}/close`)
}

// ---------------------------------------------------------------------------
// Webhooks (030)
// ---------------------------------------------------------------------------

export function listWebhooks(): Promise<Webhook[]> {
  return get<Webhook[]>('/webhooks')
}

export function createWebhook(data: Partial<Webhook>): Promise<Webhook> {
  return post<Webhook>('/webhooks', data)
}

export function updateWebhook(id: string, data: Partial<Webhook>): Promise<Webhook> {
  return put<Webhook>(`/webhooks/${id}`, data)
}

export function deleteWebhook(id: string): Promise<void> {
  return del(`/webhooks/${id}`)
}

export function toggleWebhook(id: string): Promise<Webhook> {
  return put<Webhook>(`/webhooks/${id}/toggle`)
}

export function listWebhookDeliveries(webhookId: string): Promise<WebhookDelivery[]> {
  return get<WebhookDelivery[]>(`/webhooks/${webhookId}/deliveries`)
}

export function retryWebhookDelivery(deliveryId: string): Promise<WebhookDelivery> {
  return post<WebhookDelivery>(`/webhooks/deliveries/${deliveryId}/retry`)
}

// ---------------------------------------------------------------------------
// Audit Logs (031)
// ---------------------------------------------------------------------------

export function listAuditLogs(params?: PaginationParams & {
  user_id?: string
  action?: string
  resource_type?: string
  from?: string
  to?: string
}): Promise<PaginatedResult<AuditLog>> {
  return get<PaginatedResult<AuditLog>>('/audit', params as Record<string, string | number | undefined>)
}

export function getAuditLog(id: string): Promise<AuditLog> {
  return get<AuditLog>(`/audit/${id}`)
}

export function getAuditStats(): Promise<AuditStats> {
  return get<AuditStats>('/audit/stats')
}

// ---------------------------------------------------------------------------
// Loyalty (032)
// ---------------------------------------------------------------------------

export function listLoyaltyRewards(): Promise<LoyaltyReward[]> {
  return get<LoyaltyReward[]>('/admin/loyalty/rewards')
}

export function createLoyaltyReward(data: Partial<LoyaltyReward>): Promise<LoyaltyReward> {
  return post<LoyaltyReward>('/admin/loyalty/rewards', data)
}

export function updateLoyaltyReward(id: string, data: Partial<LoyaltyReward>): Promise<LoyaltyReward> {
  return put<LoyaltyReward>(`/admin/loyalty/rewards/${id}`, data)
}

export function listLoyaltyAccounts(params?: PaginationParams): Promise<PaginatedResult<LoyaltyAccount>> {
  return get<PaginatedResult<LoyaltyAccount>>('/admin/loyalty/accounts', params as Record<string, string | number | undefined>)
}

export function getLoyaltyAccount(id: string): Promise<LoyaltyAccount> {
  return get<LoyaltyAccount>(`/admin/loyalty/accounts/${id}`)
}

export function adjustLoyaltyPoints(id: string, data: { points: number; note: string }): Promise<LoyaltyAccount> {
  return post<LoyaltyAccount>(`/admin/loyalty/accounts/${id}/adjust`, data)
}

export function getLoyaltyTransactions(id: string): Promise<LoyaltyTransaction[]> {
  return get<LoyaltyTransaction[]>(`/admin/loyalty/accounts/${id}/transactions`)
}

// ---------------------------------------------------------------------------
// Bundles (033)
// ---------------------------------------------------------------------------

export function listBundles(params?: PaginationParams): Promise<PaginatedResult<Bundle>> {
  return get<PaginatedResult<Bundle>>('/bundles', params as Record<string, string | number | undefined>)
}

export function createBundle(data: Partial<Bundle>): Promise<Bundle> {
  return post<Bundle>('/bundles', data)
}

export function updateBundle(id: string, data: Partial<Bundle>): Promise<Bundle> {
  return put<Bundle>(`/bundles/${id}`, data)
}

export function deleteBundle(id: string): Promise<void> {
  return del(`/bundles/${id}`)
}

export function addBundleItem(bundleId: string, data: { product_id: string; quantity: number; discount_percent: number }): Promise<BundleItem> {
  return post<BundleItem>(`/bundles/${bundleId}/items`, data)
}

export function removeBundleItem(bundleId: string, itemId: string): Promise<void> {
  return del(`/bundles/${bundleId}/items/${itemId}`)
}

export function getBundlePrice(bundleId: string): Promise<BundlePrice> {
  return get<BundlePrice>(`/bundles/${bundleId}/price`)
}

// ---------------------------------------------------------------------------
// Dashboard Reporting (034)
// ---------------------------------------------------------------------------

export function getDashboardSales(params?: { from?: string; to?: string }): Promise<SalesOverview> {
  return get<SalesOverview>('/dashboard/sales', params as Record<string, string | undefined>)
}

export function getDashboardTopProducts(params?: { from?: string; to?: string }): Promise<TopProductReport[]> {
  return get<TopProductReport[]>('/dashboard/top-products', params as Record<string, string | undefined>)
}

export function getDashboardRevenue(params?: { from?: string; to?: string }): Promise<RevenueData[]> {
  return get<RevenueData[]>('/dashboard/revenue', params as Record<string, string | undefined>)
}

export function getDashboardCustomers(params?: { from?: string; to?: string }): Promise<CustomerStats> {
  return get<CustomerStats>('/dashboard/customers', params as Record<string, string | undefined>)
}

export function getDashboardFunnel(params?: { from?: string; to?: string }): Promise<FunnelStep[]> {
  return get<FunnelStep[]>('/dashboard/funnel', params as Record<string, string | undefined>)
}

// ---------------------------------------------------------------------------
// Social Accounts (035)
// ---------------------------------------------------------------------------

export function listSocialAccounts(params?: PaginationParams & { provider?: SocialProvider }): Promise<PaginatedResult<SocialAccount>> {
  return get<PaginatedResult<SocialAccount>>('/social-accounts', params as Record<string, string | number | undefined>)
}

export function listCustomerSocialAccounts(customerId: string): Promise<SocialAccount[]> {
  return get<SocialAccount[]>(`/social-accounts/customer/${customerId}`)
}

// ---------------------------------------------------------------------------
// Notifications (036)
// ---------------------------------------------------------------------------

export function listNotifications(params?: PaginationParams & { read?: boolean }): Promise<PaginatedResult<AdminNotification>> {
  const { read, ...rest } = params ?? {}
  const queryParams: Record<string, string | number | undefined> = { ...rest }
  if (read !== undefined) queryParams.read = read ? '1' : '0'
  return get<PaginatedResult<AdminNotification>>('/notifications', queryParams)
}

export function getUnreadCount(): Promise<UnreadCount> {
  return get<UnreadCount>('/notifications/unread-count')
}

export function markAllNotificationsRead(): Promise<void> {
  return put<void>('/notifications/read-all')
}

export function markNotificationRead(id: string): Promise<AdminNotification> {
  return put<AdminNotification>(`/notifications/${id}/read`)
}

export function deleteNotification(id: string): Promise<void> {
  return del(`/notifications/${id}`)
}

// ---------------------------------------------------------------------------
// Storefronts / Multistore (037)
// ---------------------------------------------------------------------------

export function listStorefronts(params?: PaginationParams): Promise<PaginatedResult<Storefront>> {
  return get<PaginatedResult<Storefront>>('/storefronts', params as Record<string, string | number | undefined>)
}

export function getStorefront(id: string): Promise<Storefront> {
  return get<Storefront>(`/storefronts/${id}`)
}

export function createStorefront(data: Partial<Storefront>): Promise<Storefront> {
  return post<Storefront>('/storefronts', data)
}

export function updateStorefront(id: string, data: Partial<Storefront>): Promise<Storefront> {
  return put<Storefront>(`/storefronts/${id}`, data)
}

export function deleteStorefront(id: string): Promise<void> {
  return del(`/storefronts/${id}`)
}

export function setDefaultStorefront(id: string): Promise<Storefront> {
  return put<Storefront>(`/storefronts/${id}/default`)
}

export function getStorefrontCatalogs(id: string): Promise<StorefrontCatalog[]> {
  return get<StorefrontCatalog[]>(`/storefronts/${id}/catalogs`)
}

export function addStorefrontCatalog(id: string, data: { catalog_id: string; display_order?: number }): Promise<StorefrontCatalog> {
  return post<StorefrontCatalog>(`/storefronts/${id}/catalogs`, data)
}

export function removeStorefrontCatalog(id: string, catalogId: string): Promise<void> {
  return del(`/storefronts/${id}/catalogs/${catalogId}`)
}

// ---------------------------------------------------------------------------
// Bulk Operations (038)
// ---------------------------------------------------------------------------

export function listBulkOperations(params?: PaginationParams): Promise<PaginatedResult<BulkOperation>> {
  return get<PaginatedResult<BulkOperation>>('/bulk-operations', params as Record<string, string | number | undefined>)
}

export function getBulkOperation(id: string): Promise<BulkOperation> {
  return get<BulkOperation>(`/bulk-operations/${id}`)
}

export function createBulkOperation(data: Partial<BulkOperation>): Promise<BulkOperation> {
  return post<BulkOperation>('/bulk-operations', data)
}

export function getBulkOperationItems(id: string): Promise<BulkOperationItem[]> {
  return get<BulkOperationItem[]>(`/bulk-operations/${id}/items`)
}

export function processBulkOperation(id: string): Promise<BulkOperation> {
  return post<BulkOperation>(`/bulk-operations/${id}/process`)
}

export function cancelBulkOperation(id: string): Promise<BulkOperation> {
  return post<BulkOperation>(`/bulk-operations/${id}/cancel`)
}

// ---------------------------------------------------------------------------
// Blog (039)
// ---------------------------------------------------------------------------

export function listBlogPosts(params?: PaginationParams): Promise<PaginatedResult<BlogPost>> {
  return get<PaginatedResult<BlogPost>>('/blog/posts', params as Record<string, string | number | undefined>)
}

export function getBlogPost(id: string): Promise<BlogPost> {
  return get<BlogPost>(`/blog/posts/${id}`)
}

export function createBlogPost(data: Partial<BlogPost>): Promise<BlogPost> {
  return post<BlogPost>('/blog/posts', data)
}

export function updateBlogPost(id: string, data: Partial<BlogPost>): Promise<BlogPost> {
  return put<BlogPost>(`/blog/posts/${id}`, data)
}

export function deleteBlogPost(id: string): Promise<void> {
  return del(`/blog/posts/${id}`)
}

export function publishBlogPost(id: string): Promise<BlogPost> {
  return put<BlogPost>(`/blog/posts/${id}/publish`)
}

export function archiveBlogPost(id: string): Promise<BlogPost> {
  return put<BlogPost>(`/blog/posts/${id}/archive`)
}

export function listBlogCategories(): Promise<BlogCategory[]> {
  return get<BlogCategory[]>('/blog/categories')
}

export function createBlogCategory(data: Partial<BlogCategory>): Promise<BlogCategory> {
  return post<BlogCategory>('/blog/categories', data)
}

export function updateBlogCategory(id: string, data: Partial<BlogCategory>): Promise<BlogCategory> {
  return put<BlogCategory>(`/blog/categories/${id}`, data)
}

export function deleteBlogCategory(id: string): Promise<void> {
  return del(`/blog/categories/${id}`)
}

// ---------------------------------------------------------------------------
// Admin Collections (040)
// ---------------------------------------------------------------------------

export function listAdminCollections(params?: PaginationParams): Promise<PaginatedResult<AdminCollection>> {
  return get<PaginatedResult<AdminCollection>>('/collections', params as Record<string, string | number | undefined>)
}

export function getAdminCollection(id: string): Promise<AdminCollection> {
  return get<AdminCollection>(`/collections/${id}`)
}

export function createAdminCollection(data: Partial<AdminCollection>): Promise<AdminCollection> {
  return post<AdminCollection>('/collections', data)
}

export function updateAdminCollection(id: string, data: Partial<AdminCollection>): Promise<AdminCollection> {
  return put<AdminCollection>(`/collections/${id}`, data)
}

export function deleteAdminCollection(id: string): Promise<void> {
  return del(`/collections/${id}`)
}

export function getCollectionProducts(id: string): Promise<CollectionProduct[]> {
  return get<CollectionProduct[]>(`/collections/${id}/products`)
}

export function addCollectionProduct(id: string, data: { product_id: string; sort_order?: number }): Promise<CollectionProduct> {
  return post<CollectionProduct>(`/collections/${id}/products`, data)
}

export function removeCollectionProduct(id: string, productId: string): Promise<void> {
  return del(`/collections/${id}/products/${productId}`)
}

// ---------------------------------------------------------------------------
// A/B Testing — Experiments (041)
// ---------------------------------------------------------------------------

export function listExperiments(params?: PaginationParams): Promise<PaginatedResult<Experiment>> {
  return get<PaginatedResult<Experiment>>('/experiments', params as Record<string, string | number | undefined>)
}

export function getExperiment(id: string): Promise<Experiment> {
  return get<Experiment>(`/experiments/${id}`)
}

export function createExperiment(data: Partial<Experiment>): Promise<Experiment> {
  return post<Experiment>('/experiments', data)
}

export function updateExperiment(id: string, data: Partial<Experiment>): Promise<Experiment> {
  return put<Experiment>(`/experiments/${id}`, data)
}

export function deleteExperiment(id: string): Promise<void> {
  return del(`/experiments/${id}`)
}

export function startExperiment(id: string): Promise<Experiment> {
  return put<Experiment>(`/experiments/${id}/start`)
}

export function pauseExperiment(id: string): Promise<Experiment> {
  return put<Experiment>(`/experiments/${id}/pause`)
}

export function completeExperiment(id: string): Promise<Experiment> {
  return put<Experiment>(`/experiments/${id}/complete`)
}

export function getExperimentResults(id: string): Promise<ExperimentResults> {
  return get<ExperimentResults>(`/experiments/${id}/results`)
}

export function addExperimentVariant(id: string, data: Partial<ExperimentVariant>): Promise<ExperimentVariant> {
  return post<ExperimentVariant>(`/experiments/${id}/variants`, data)
}

export function deleteExperimentVariant(id: string, variantId: string): Promise<void> {
  return del(`/experiments/${id}/variants/${variantId}`)
}

// ---------------------------------------------------------------------------
// Recommendations (042)
// ---------------------------------------------------------------------------

export function listRecommendationRules(params?: PaginationParams): Promise<PaginatedResult<RecommendationRule>> {
  return get<PaginatedResult<RecommendationRule>>('/recommendations/rules', params as Record<string, string | number | undefined>)
}

export function createRecommendationRule(data: Partial<RecommendationRule>): Promise<RecommendationRule> {
  return post<RecommendationRule>('/recommendations/rules', data)
}

export function updateRecommendationRule(id: string, data: Partial<RecommendationRule>): Promise<RecommendationRule> {
  return put<RecommendationRule>(`/recommendations/rules/${id}`, data)
}

export function deleteRecommendationRule(id: string): Promise<void> {
  return del(`/recommendations/rules/${id}`)
}

export function getRecommendations(productId: string): Promise<RecommendedProduct[]> {
  return get<RecommendedProduct[]>(`/recommendations/product/${productId}`)
}

// ---------------------------------------------------------------------------
// Preset Marketplace (050)
// ---------------------------------------------------------------------------

export function listMarketplacePresets(params?: { category?: string; search?: string }): Promise<PaginatedResult<Preset>> {
  return get<PaginatedResult<Preset>>('/marketplace/presets', params as Record<string, string | undefined>)
}

export function getMarketplacePreset(id: string): Promise<Preset> {
  return get<Preset>(`/marketplace/presets/${id}`)
}

export function fetchInstalledPresets(): Promise<PaginatedResult<PresetInstall>> {
  return get<PaginatedResult<PresetInstall>>('/marketplace/installations')
}

export function installPreset(id: string): Promise<PresetInstall> {
  return post<PresetInstall>('/marketplace/installations', { preset_id: id })
}

export function uninstallPreset(id: string): Promise<void> {
  return del(`/marketplace/installations/${id}`)
}

// ---------------------------------------------------------------------------
// Agent Sessions / Workspaces (050)
// ---------------------------------------------------------------------------

export function listAgentSessions(): Promise<PaginatedResult<AgentSession>> {
  return get<PaginatedResult<AgentSession>>('/agent/sessions')
}

export function getAgentSession(id: string): Promise<AgentSession> {
  return get<AgentSession>(`/agent/sessions/${id}`)
}

export function createAgentSession(data: { preset_id: string; name?: string }): Promise<AgentSession> {
  return post<AgentSession>('/agent/sessions', data)
}

export function stopAgentSession(id: string): Promise<void> {
  return del(`/agent/sessions/${id}`)
}

export function fetchSessionHistory(id: string): Promise<ChatMessage[]> {
  return get<ChatMessage[]>(`/agent/sessions/${id}/messages`)
}

export function sendSessionMessage(id: string, message: string): Promise<ChatMessage> {
  return post<ChatMessage>(`/agent/sessions/${id}/messages`, { message })
}

// ---------------------------------------------------------------------------
// Approvals
// ---------------------------------------------------------------------------

export function listApprovals(params?: { status?: string; page?: number; page_size?: number }): Promise<PaginatedResult<ApprovalRequest>> {
  return get<PaginatedResult<ApprovalRequest>>('/approvals', params as Record<string, string | number | undefined>)
}

export function getApproval(id: string): Promise<ApprovalRequest> {
  return get<ApprovalRequest>(`/approvals/${id}`)
}

export function getApprovalCount(): Promise<{ count: number }> {
  return get<{ count: number }>('/approvals/count')
}

export function approveRequest(id: string, reason?: string): Promise<ApprovalRequest> {
  return post<ApprovalRequest>(`/approvals/${id}/approve`, { reason: reason ?? '' })
}

export function rejectRequest(id: string, reason?: string): Promise<ApprovalRequest> {
  return post<ApprovalRequest>(`/approvals/${id}/reject`, { reason: reason ?? '' })
}

// ---------------------------------------------------------------------------
// Agent Memory
// ---------------------------------------------------------------------------

export function listAgentMemories(params?: { page?: number; page_size?: number }): Promise<PaginatedResult<AgentMemory>> {
  return get<PaginatedResult<AgentMemory>>('/agent/memories', params as Record<string, string | number | undefined>)
}

export function searchAgentMemories(params: { q?: string; category?: string; tags?: string }): Promise<PaginatedResult<AgentMemory>> {
  return get<PaginatedResult<AgentMemory>>('/agent/memories/search', params as Record<string, string | undefined>)
}

export function createAgentMemory(data: CreateMemoryRequest): Promise<AgentMemory> {
  return post<AgentMemory>('/agent/memories', data)
}

export function getAgentMemory(id: string): Promise<AgentMemory> {
  return get<AgentMemory>(`/agent/memories/${id}`)
}

export function updateAgentMemory(id: string, data: UpdateMemoryRequest): Promise<AgentMemory> {
  return put<AgentMemory>(`/agent/memories/${id}`, data)
}

export function deleteAgentMemory(id: string): Promise<void> {
  return del(`/agent/memories/${id}`)
}

// ---------------------------------------------------------------------------
// Agent Triggers
// ---------------------------------------------------------------------------

export function listAgentTriggers(params?: { page?: number; page_size?: number }): Promise<PaginatedResult<AgentTrigger>> {
  return get<PaginatedResult<AgentTrigger>>('/agent/triggers', params as Record<string, string | number | undefined>)
}

export function createAgentTrigger(data: CreateTriggerRequest): Promise<AgentTrigger> {
  return post<AgentTrigger>('/agent/triggers', data)
}

export function getAgentTrigger(id: string): Promise<AgentTrigger> {
  return get<AgentTrigger>(`/agent/triggers/${id}`)
}

export function updateAgentTrigger(id: string, data: UpdateTriggerRequest): Promise<AgentTrigger> {
  return put<AgentTrigger>(`/agent/triggers/${id}`, data)
}

export function deleteAgentTrigger(id: string): Promise<void> {
  return del(`/agent/triggers/${id}`)
}

export function enableAgentTrigger(id: string): Promise<AgentTrigger> {
  return post<AgentTrigger>(`/agent/triggers/${id}/enable`)
}

export function disableAgentTrigger(id: string): Promise<AgentTrigger> {
  return post<AgentTrigger>(`/agent/triggers/${id}/disable`)
}

export function listTriggerLogs(id: string, params?: { page?: number; page_size?: number }): Promise<PaginatedResult<TriggerLog>> {
  return get<PaginatedResult<TriggerLog>>(`/agent/triggers/${id}/logs`, params as Record<string, string | number | undefined>)
}

export function listTriggerEventTypes(): Promise<{ event_types: string[] }> {
  return get<{ event_types: string[] }>('/agent/triggers/event-types')
}
