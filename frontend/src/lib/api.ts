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
} from '../types'

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api/v1'

// TODO: replace with auth context / tenant resolver once multi-tenancy UI is built
const TENANT_ID = import.meta.env.VITE_TENANT_ID ?? 'default'

export interface PaginationParams {
  page?: number
  page_size?: number
}

// ---------------------------------------------------------------------------
// Core fetch helpers
// ---------------------------------------------------------------------------

function buildHeaders(extra?: Record<string, string>): HeadersInit {
  return {
    'Content-Type': 'application/json',
    'X-Tenant-ID': TENANT_ID,
    ...extra,
  }
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

export async function get<T>(path: string, params?: Record<string, string | number | undefined>): Promise<T> {
  const url = new URL(`${BASE_URL}${path}`, window.location.origin)
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (v !== undefined) url.searchParams.set(k, String(v))
    }
  }
  const res = await fetch(url.toString(), {
    method: 'GET',
    headers: buildHeaders(),
  })
  return handleResponse<T>(res)
}

export async function post<T>(path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'POST',
    headers: buildHeaders(),
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })
  return handleResponse<T>(res)
}

export async function put<T>(path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'PUT',
    headers: buildHeaders(),
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })
  return handleResponse<T>(res)
}

export async function del(path: string): Promise<void> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'DELETE',
    headers: buildHeaders(),
  })
  await handleResponse<void>(res)
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

  const res = await fetch(`${BASE_URL}/media`, {
    method: 'POST',
    headers: {
      // Do NOT set Content-Type — browser sets multipart/form-data with boundary automatically
      'X-Tenant-ID': TENANT_ID,
    },
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
