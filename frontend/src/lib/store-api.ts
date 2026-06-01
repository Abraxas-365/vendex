/**
 * Store API — public storefront endpoints under /store/*.
 * Tenant is resolved server-side from the Host header (subdomain or custom domain).
 * In development, we pass X-Tenant-ID explicitly.
 */
import type { Product, PaginatedResult } from '../types'

const STORE_URL = import.meta.env.VITE_STORE_BASE_URL ?? '/store'
const DEV_TENANT_ID = import.meta.env.VITE_TENANT_ID ?? ''

function storeHeaders(): HeadersInit {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  // In development, send tenant ID explicitly since Host is localhost
  if (DEV_TENANT_ID) {
    headers['X-Tenant-ID'] = DEV_TENANT_ID
  }
  return headers
}

async function storeGet<T>(path: string, params?: Record<string, string | number | undefined>): Promise<T> {
  const url = new URL(`${STORE_URL}${path}`, window.location.origin)
  if (params) {
    Object.entries(params).forEach(([k, v]) => {
      if (v !== undefined && v !== null) url.searchParams.set(k, String(v))
    })
  }
  const res = await fetch(url.toString(), { headers: storeHeaders() })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

// ---------------------------------------------------------------------------
// Products
// ---------------------------------------------------------------------------

export function listProducts(params?: { page?: number; page_size?: number; category_id?: string }): Promise<PaginatedResult<Product>> {
  return storeGet<PaginatedResult<Product>>('/products', params as Record<string, string | number | undefined>)
}

export function getProduct(id: string): Promise<Product> {
  return storeGet<Product>(`/products/${id}`)
}

export function getProductBySlug(slug: string): Promise<Product> {
  return storeGet<Product>(`/products/slug/${slug}`)
}

// ---------------------------------------------------------------------------
// Pages (CMS)
// ---------------------------------------------------------------------------

export interface StorePage {
  id: string
  title: string
  slug: string
  content: string
  html: string
  css: string
  status: string
  is_published: boolean
  meta?: { description?: string }
}

export function getPageBySlug(slug: string): Promise<StorePage> {
  return storeGet<StorePage>(`/pages/slug/${slug}`)
}

// ---------------------------------------------------------------------------
// Store Info (branding)
// ---------------------------------------------------------------------------

export interface TrustBadge {
  icon: string
  title: string
  desc: string
}

export interface StoreInfo {
  store_name: string
  store_email: string
  logo_url: string
  currency: string
  tagline: string
  hero_title: string
  hero_subtitle: string
  accent_color: string
  bg_style: string
  announcement: string
  trust_badges: TrustBadge[]
}

export function getStoreInfo(): Promise<StoreInfo> {
  return storeGet<StoreInfo>('/info')
}

// ---------------------------------------------------------------------------
// Categories
// ---------------------------------------------------------------------------

export interface Category {
  id: string
  name: string
  slug?: string
  description?: string
}

export async function listCategories(): Promise<Category[]> {
  const result = await storeGet<{ items: Category[] }>('/categories')
  return result.items ?? []
}

// ---------------------------------------------------------------------------
// Pages (nav/footer links)
// ---------------------------------------------------------------------------

export interface PageLink {
  slug: string
  title: string
}

export function listPages(): Promise<PageLink[]> {
  return storeGet<PageLink[]>('/pages')
}
