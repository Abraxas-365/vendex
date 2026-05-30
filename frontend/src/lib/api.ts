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

export function getAccessToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY)
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
  const res = await fetch(`${BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return handleResponse<LoginResponse>(res)
}

export async function handleAuthCallback(provider: string, code: string, state: string): Promise<TokenResponse> {
  const url = new URL(`${BASE_URL}/auth/callback/${provider}`, window.location.origin)
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
  const res = await fetch(`${BASE_URL}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken }),
  })
  const data = await handleResponse<RefreshResponse>(res)
  return data.access_token
}

// Exported refresh function for explicit use
export async function refreshToken(token: string): Promise<RefreshResponse> {
  const res = await fetch(`${BASE_URL}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: token }),
  })
  return handleResponse<RefreshResponse>(res)
}

export async function logout(): Promise<{ message: string }> {
  const res = await fetch(`${BASE_URL}/auth/logout`, {
    method: 'POST',
    headers: buildHeaders(),
  })
  return handleResponse<{ message: string }>(res)
}

export async function getMe(): Promise<MeResponse> {
  const res = await fetch(`${BASE_URL}/auth/me`, {
    method: 'GET',
    headers: buildHeaders(),
  })
  return handleResponse<MeResponse>(res)
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
  return get<PaginatedResult<Category>>('/catalog/categories', params as Record<string, string | number | undefined>)
}

export function getCategory(id: string): Promise<Category> {
  return get<Category>(`/catalog/categories/${id}`)
}

export function createCategory(data: Partial<Category>): Promise<Category> {
  return post<Category>('/catalog/categories', data)
}

export function deleteCategory(id: string): Promise<void> {
  return del(`/catalog/categories/${id}`)
}

// ---------------------------------------------------------------------------
// Catalog — Collections
// ---------------------------------------------------------------------------

export function listCollections(params?: PaginationParams): Promise<PaginatedResult<Collection>> {
  return get<PaginatedResult<Collection>>('/catalog/collections', params as Record<string, string | number | undefined>)
}

export function getCollection(id: string): Promise<Collection> {
  return get<Collection>(`/catalog/collections/${id}`)
}

export function createCollection(data: Partial<Collection>): Promise<Collection> {
  return post<Collection>('/catalog/collections', data)
}

export function deleteCollection(id: string): Promise<void> {
  return del(`/catalog/collections/${id}`)
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
  return get<Page>(`/storefront/pages/by-slug/${slug}`)
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
  return get<PaginatedResult<Promo>>('/promos', params as Record<string, string | number | undefined>)
}

export function createPromo(data: Partial<Promo>): Promise<Promo> {
  return post<Promo>('/promos', data)
}

export function validatePromo(code: string, orderAmount: number): Promise<{ valid: boolean; discount: number }> {
  return post<{ valid: boolean; discount: number }>('/promos/validate', { code, order_amount: orderAmount })
}

export function updatePromo(id: string, data: Partial<Promo>): Promise<Promo> {
  return put<Promo>(`/promos/${id}`, data)
}

export function deactivatePromo(id: string): Promise<Promo> {
  return post<Promo>(`/promos/${id}/deactivate`)
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

  const res = await fetch(`${BASE_URL}/media`, {
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
  return get<BlockType[]>('/block-types', category ? { category } : undefined)
}

export function getBlockType(id: string): Promise<BlockType> {
  return get<BlockType>(`/block-types/${id}`)
}

export function createBlockType(data: Partial<BlockType>): Promise<BlockType> {
  return post<BlockType>('/block-types', data)
}

export function updateBlockType(id: string, data: Partial<BlockType>): Promise<BlockType> {
  return put<BlockType>(`/block-types/${id}`, data)
}

export function deleteBlockType(id: string): Promise<void> {
  return del(`/block-types/${id}`)
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
