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
// Marketplace
// ---------------------------------------------------------------------------

export function listMarketplacePlugins(params?: PaginationParams): Promise<PaginatedResult<Plugin>> {
  return get<PaginatedResult<Plugin>>('/marketplace/plugins', params as Record<string, string | number | undefined>)
}

export function getMarketplacePlugin(id: string): Promise<Plugin & { latest_version?: PluginVersion }> {
  return get<Plugin & { latest_version?: PluginVersion }>(`/marketplace/plugins/${id}`)
}

export function installPlugin(pluginId: string): Promise<PluginInstallation> {
  return post<PluginInstallation>('/marketplace/install', { plugin_id: pluginId })
}

export function uninstallPlugin(pluginId: string): Promise<void> {
  return post<void>('/marketplace/uninstall', { plugin_id: pluginId })
}

export function listInstalledPlugins(): Promise<PluginInstallation[]> {
  return get<PluginInstallation[]>('/marketplace/installed')
}

export function updatePluginSettings(pluginId: string, settings: Record<string, unknown>): Promise<PluginInstallation> {
  return put<PluginInstallation>(`/marketplace/plugins/${pluginId}/settings`, { settings })
}

// ---------------------------------------------------------------------------
// Plugin Runtime
// ---------------------------------------------------------------------------

export function listPluginManifests(): Promise<PluginManifest[]> {
  return get<PluginManifest[]>('/plugins/manifests')
}

export function getPluginManifest(name: string): Promise<PluginManifest> {
  return get<PluginManifest>(`/plugins/${name}/manifest`)
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
