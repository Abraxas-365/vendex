/**
 * Store hooks — TanStack Query hooks for public storefront data.
 * These hit /store/* endpoints (tenant resolved from domain).
 */
import { useQuery, type UseQueryResult } from '@tanstack/react-query'
import type { Product, PaginatedResult } from '../types'
import * as storeApi from './store-api'
import type { StoreInfo } from './store-api'

export function useStoreInfo(): UseQueryResult<StoreInfo> {
  return useQuery({
    queryKey: ['store', 'info'],
    queryFn: () => storeApi.getStoreInfo(),
    staleTime: 5 * 60 * 1000, // cache for 5 min
  })
}

export function useStoreProducts(params?: { page?: number; page_size?: number; category_id?: string }): UseQueryResult<PaginatedResult<Product>> {
  return useQuery({
    queryKey: ['store', 'products', params],
    queryFn: () => storeApi.listProducts(params),
  })
}

export function useStoreProduct(id: string): UseQueryResult<Product> {
  return useQuery({
    queryKey: ['store', 'products', id],
    queryFn: () => storeApi.getProduct(id),
    enabled: Boolean(id),
  })
}

export function useStorePageBySlug(slug: string): UseQueryResult<storeApi.StorePage> {
  return useQuery({
    queryKey: ['store', 'pages', slug],
    queryFn: () => storeApi.getPageBySlug(slug),
    enabled: Boolean(slug),
  })
}
