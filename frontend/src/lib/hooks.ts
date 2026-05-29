import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryResult,
  type UseMutationResult,
} from '@tanstack/react-query'
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
  PluginInstallation,
  PluginManifest,
} from '../types'
import * as api from './api'
import type { PaginationParams } from './api'

// ---------------------------------------------------------------------------
// Query key factory — centralised to avoid typos and simplify invalidation
// ---------------------------------------------------------------------------

export const queryKeys = {
  products: {
    all: ['products'] as const,
    list: (params?: PaginationParams) => ['products', 'list', params] as const,
    detail: (id: string) => ['products', 'detail', id] as const,
  },
  orders: {
    all: ['orders'] as const,
    list: (params?: PaginationParams) => ['orders', 'list', params] as const,
    detail: (id: string) => ['orders', 'detail', id] as const,
  },
  customers: {
    all: ['customers'] as const,
    list: (params?: PaginationParams) => ['customers', 'list', params] as const,
    detail: (id: string) => ['customers', 'detail', id] as const,
  },
  categories: {
    all: ['categories'] as const,
    list: (params?: PaginationParams) => ['categories', 'list', params] as const,
    detail: (id: string) => ['categories', 'detail', id] as const,
  },
  collections: {
    all: ['collections'] as const,
    list: (params?: PaginationParams) => ['collections', 'list', params] as const,
    detail: (id: string) => ['collections', 'detail', id] as const,
  },
  pages: {
    all: ['pages'] as const,
    list: (params?: PaginationParams) => ['pages', 'list', params] as const,
    detail: (id: string) => ['pages', 'detail', id] as const,
    bySlug: (slug: string) => ['pages', 'slug', slug] as const,
    versions: (id: string) => ['pages', 'versions', id] as const,
  },
  promos: {
    all: ['promos'] as const,
    list: (params?: PaginationParams) => ['promos', 'list', params] as const,
  },
  media: {
    all: ['media'] as const,
    list: (params?: PaginationParams) => ['media', 'list', params] as const,
  },
} as const

// ---------------------------------------------------------------------------
// Products
// ---------------------------------------------------------------------------

export function useProducts(params?: PaginationParams): UseQueryResult<PaginatedResult<Product>> {
  return useQuery({
    queryKey: queryKeys.products.list(params),
    queryFn: () => api.listProducts(params),
  })
}

export function useProduct(id: string): UseQueryResult<Product> {
  return useQuery({
    queryKey: queryKeys.products.detail(id),
    queryFn: () => api.getProduct(id),
    enabled: Boolean(id),
  })
}

export function useCreateProduct(): UseMutationResult<Product, Error, Partial<Product>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createProduct(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.products.all })
    },
  })
}

export function useUpdateProduct(): UseMutationResult<Product, Error, { id: string; data: Partial<Product> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateProduct(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.products.all })
      qc.setQueryData(queryKeys.products.detail(updated.id), updated)
    },
  })
}

export function useDeleteProduct(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteProduct(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.products.all })
    },
  })
}

// ---------------------------------------------------------------------------
// Orders
// ---------------------------------------------------------------------------

export function useOrders(params?: PaginationParams): UseQueryResult<PaginatedResult<Order>> {
  return useQuery({
    queryKey: queryKeys.orders.list(params),
    queryFn: () => api.listOrders(params),
  })
}

export function useOrder(id: string): UseQueryResult<Order> {
  return useQuery({
    queryKey: queryKeys.orders.detail(id),
    queryFn: () => api.getOrder(id),
    enabled: Boolean(id),
  })
}

export function useUpdateOrderStatus(): UseMutationResult<Order, Error, { id: string; status: OrderStatus }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }) => api.updateOrderStatus(id, status),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.orders.all })
      qc.setQueryData(queryKeys.orders.detail(updated.id), updated)
    },
  })
}

// ---------------------------------------------------------------------------
// Customers
// ---------------------------------------------------------------------------

export function useCustomers(params?: PaginationParams): UseQueryResult<PaginatedResult<Customer>> {
  return useQuery({
    queryKey: queryKeys.customers.list(params),
    queryFn: () => api.listCustomers(params),
  })
}

export function useCustomer(id: string): UseQueryResult<Customer> {
  return useQuery({
    queryKey: queryKeys.customers.detail(id),
    queryFn: () => api.getCustomer(id),
    enabled: Boolean(id),
  })
}

export function useCreateCustomer(): UseMutationResult<Customer, Error, Partial<Customer>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createCustomer(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.customers.all })
    },
  })
}

// ---------------------------------------------------------------------------
// Categories
// ---------------------------------------------------------------------------

export function useCategories(params?: PaginationParams): UseQueryResult<PaginatedResult<Category>> {
  return useQuery({
    queryKey: queryKeys.categories.list(params),
    queryFn: () => api.listCategories(params),
  })
}

export function useCategory(id: string): UseQueryResult<Category> {
  return useQuery({
    queryKey: queryKeys.categories.detail(id),
    queryFn: () => api.getCategory(id),
    enabled: Boolean(id),
  })
}

export function useCreateCategory(): UseMutationResult<Category, Error, Partial<Category>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createCategory(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.categories.all })
    },
  })
}

// ---------------------------------------------------------------------------
// Collections
// ---------------------------------------------------------------------------

export function useCollections(params?: PaginationParams): UseQueryResult<PaginatedResult<Collection>> {
  return useQuery({
    queryKey: queryKeys.collections.list(params),
    queryFn: () => api.listCollections(params),
  })
}

export function useCollection(id: string): UseQueryResult<Collection> {
  return useQuery({
    queryKey: queryKeys.collections.detail(id),
    queryFn: () => api.getCollection(id),
    enabled: Boolean(id),
  })
}

export function useCreateCollection(): UseMutationResult<Collection, Error, Partial<Collection>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createCollection(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.collections.all })
    },
  })
}

// ---------------------------------------------------------------------------
// Pages
// ---------------------------------------------------------------------------

export function usePages(params?: PaginationParams): UseQueryResult<PaginatedResult<Page>> {
  return useQuery({
    queryKey: queryKeys.pages.list(params),
    queryFn: () => api.listPages(params),
  })
}

export function usePage(id: string): UseQueryResult<Page> {
  return useQuery({
    queryKey: queryKeys.pages.detail(id),
    queryFn: () => api.getPage(id),
    enabled: Boolean(id),
  })
}

export function usePageBySlug(slug: string): UseQueryResult<Page> {
  return useQuery({
    queryKey: queryKeys.pages.bySlug(slug),
    queryFn: () => api.getPageBySlug(slug),
    enabled: Boolean(slug),
  })
}

export function usePageVersions(id: string): UseQueryResult<PageVersion[]> {
  return useQuery({
    queryKey: queryKeys.pages.versions(id),
    queryFn: () => api.getPageVersions(id),
    enabled: Boolean(id),
  })
}

export function useCreatePage(): UseMutationResult<Page, Error, Partial<Page>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createPage(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.pages.all })
    },
  })
}

export function useUpdatePage(): UseMutationResult<Page, Error, { id: string; data: Partial<Page> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updatePage(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.pages.all })
      qc.setQueryData(queryKeys.pages.detail(updated.id), updated)
    },
  })
}

export function usePublishPage(): UseMutationResult<Page, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.publishPage(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.pages.all })
      qc.setQueryData(queryKeys.pages.detail(updated.id), updated)
    },
  })
}

export function useUnpublishPage(): UseMutationResult<Page, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.unpublishPage(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.pages.all })
      qc.setQueryData(queryKeys.pages.detail(updated.id), updated)
    },
  })
}

export function useArchivePage(): UseMutationResult<Page, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.archivePage(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.pages.all })
      qc.setQueryData(queryKeys.pages.detail(updated.id), updated)
    },
  })
}

export function useDeletePage(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deletePage(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.pages.all })
    },
  })
}

// ---------------------------------------------------------------------------
// Promos
// ---------------------------------------------------------------------------

export function usePromos(params?: PaginationParams): UseQueryResult<PaginatedResult<Promo>> {
  return useQuery({
    queryKey: queryKeys.promos.list(params),
    queryFn: () => api.listPromos(params),
  })
}

export function useCreatePromo(): UseMutationResult<Promo, Error, Partial<Promo>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createPromo(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.promos.all })
    },
  })
}

export function useUpdatePromo(): UseMutationResult<Promo, Error, { id: string; data: Partial<Promo> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updatePromo(id, data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.promos.all })
    },
  })
}

// ---------------------------------------------------------------------------
// Media
// ---------------------------------------------------------------------------

export function useMedia(params?: PaginationParams): UseQueryResult<PaginatedResult<Media>> {
  return useQuery({
    queryKey: queryKeys.media.list(params),
    queryFn: () => api.listMedia(params),
  })
}

export function useUploadMedia(): UseMutationResult<Media, Error, { file: File; alt?: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ file, alt }) => api.uploadMedia(file, alt),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.media.all })
    },
  })
}

export function useDeleteMedia(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteMedia(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.media.all })
    },
  })
}

// ---------------------------------------------------------------------------
// Marketplace
// ---------------------------------------------------------------------------

export function useMarketplacePlugins(params?: PaginationParams): UseQueryResult<PaginatedResult<Plugin>> {
  return useQuery({
    queryKey: ['marketplace', 'plugins', params],
    queryFn: () => api.listMarketplacePlugins(params),
  })
}

export function useInstallPlugin(): UseMutationResult<PluginInstallation, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (pluginId: string) => api.installPlugin(pluginId),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['marketplace'] })
      void qc.invalidateQueries({ queryKey: ['plugins'] })
    },
  })
}

export function useUninstallPlugin(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (pluginId: string) => api.uninstallPlugin(pluginId),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['marketplace'] })
      void qc.invalidateQueries({ queryKey: ['plugins'] })
    },
  })
}

export function useInstalledPlugins(): UseQueryResult<PluginInstallation[]> {
  return useQuery({
    queryKey: ['marketplace', 'installed'],
    queryFn: () => api.listInstalledPlugins(),
  })
}

export function usePluginManifests(): UseQueryResult<PluginManifest[]> {
  return useQuery({
    queryKey: ['plugins', 'manifests'],
    queryFn: () => api.listPluginManifests(),
  })
}
