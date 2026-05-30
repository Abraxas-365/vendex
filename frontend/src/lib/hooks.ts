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
  DashboardStats,
  RevenuePoint,
  TopProduct,
  OrderStatusBreakdown,
  RecentOrder,
  StoreSettings,
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
  CustomerGroup,
  GroupMembership,
} from '../types'
import * as api from './api'
import type { PaginationParams } from './api'
import { useAuth } from './auth'

// ---------------------------------------------------------------------------
// Query key factory — centralised to avoid typos and simplify invalidation
// ---------------------------------------------------------------------------

export const queryKeys = {
  products: {
    all: ['products'] as const,
    list: (params?: PaginationParams) => ['products', 'list', params] as const,
    detail: (id: string) => ['products', 'detail', id] as const,
    options: (id: string) => ['products', 'options', id] as const,
    variants: (id: string) => ['products', 'variants', id] as const,
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
  blockTypes: {
    all: ['block-types'] as const,
    list: (category?: string) => ['block-types', 'list', category] as const,
    detail: (id: string) => ['block-types', 'detail', id] as const,
  },
  themes: {
    all: ['themes'] as const,
    list: () => ['themes', 'list'] as const,
    detail: (id: string) => ['themes', 'detail', id] as const,
    active: () => ['themes', 'active'] as const,
  },
  shipping: {
    zones: ['shipping', 'zones'] as const,
    zone: (id: string) => ['shipping', 'zones', id] as const,
    rates: (zoneId: string) => ['shipping', 'zones', zoneId, 'rates'] as const,
  },
  tax: {
    rates: ['tax', 'rates'] as const,
  },
  payments: {
    all: ['payments'] as const,
    detail: (id: string) => ['payments', 'detail', id] as const,
    byOrder: (orderId: string) => ['payments', 'order', orderId] as const,
    refunds: (paymentId: string) => ['payments', 'refunds', paymentId] as const,
  },
  customerGroups: {
    all: ['customer-groups'] as const,
    list: () => ['customer-groups', 'list'] as const,
    detail: (id: string) => ['customer-groups', 'detail', id] as const,
    members: (groupId: string) => ['customer-groups', groupId, 'members'] as const,
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

export function useCancelOrder(): UseMutationResult<Order, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.cancelOrder(id),
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

export function useDeleteCustomer(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteCustomer(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.customers.all })
    },
  })
}

export function useCustomerOrders(customerId: string): UseQueryResult<PaginatedResult<Order>> {
  return useQuery({
    queryKey: ['orders', 'by-customer', customerId],
    queryFn: () => api.listOrdersByCustomer(customerId),
    enabled: Boolean(customerId),
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

export function useDeleteCategory(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteCategory(id),
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

export function useDeleteCollection(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteCollection(id),
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

export function useInstalledPlugins(params?: PaginationParams): UseQueryResult<PaginatedResult<PluginInstallation>> {
  return useQuery({
    queryKey: ['marketplace', 'installed', params],
    queryFn: () => api.listInstalledPlugins(params),
  })
}

export function useEnablePlugin(): UseMutationResult<PluginInstallation, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (pluginId: string) => api.enablePlugin(pluginId),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['marketplace'] })
      void qc.invalidateQueries({ queryKey: ['plugins'] })
    },
  })
}

export function useDisablePlugin(): UseMutationResult<PluginInstallation, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (pluginId: string) => api.disablePlugin(pluginId),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['marketplace'] })
      void qc.invalidateQueries({ queryKey: ['plugins'] })
    },
  })
}

// ---------------------------------------------------------------------------
// Analytics
// ---------------------------------------------------------------------------

export function useDashboardStats(): UseQueryResult<DashboardStats> {
  return useQuery({
    queryKey: ['analytics', 'dashboard'],
    queryFn: () => api.getDashboardStats(),
  })
}

export function useRevenueTimeline(days: number = 30): UseQueryResult<RevenuePoint[]> {
  return useQuery({
    queryKey: ['analytics', 'revenue', days],
    queryFn: () => api.getRevenueTimeline(days),
  })
}

export function useTopProducts(limit: number = 5): UseQueryResult<TopProduct[]> {
  return useQuery({
    queryKey: ['analytics', 'top-products', limit],
    queryFn: () => api.getTopProducts(limit),
  })
}

export function useOrderStatusBreakdown(): UseQueryResult<OrderStatusBreakdown[]> {
  return useQuery({
    queryKey: ['analytics', 'order-status'],
    queryFn: () => api.getOrderStatusBreakdown(),
  })
}

export function useRecentOrders(limit: number = 5): UseQueryResult<RecentOrder[]> {
  return useQuery({
    queryKey: ['analytics', 'recent-orders', limit],
    queryFn: () => api.getRecentOrders(limit),
  })
}

// ---------------------------------------------------------------------------
// Settings
// ---------------------------------------------------------------------------

export function useSettings(): UseQueryResult<StoreSettings> {
  return useQuery({
    queryKey: ['settings'],
    queryFn: () => api.getSettings(),
  })
}

export function useUpdateSettings(): UseMutationResult<StoreSettings, Error, Partial<StoreSettings>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<StoreSettings>) => api.updateSettings(data),
    onSuccess: (updated) => {
      qc.setQueryData(['settings'], updated)
    },
  })
}

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

export function useCurrentUser(): UseQueryResult<MeResponse> {
  return useQuery({
    queryKey: ['auth', 'me'],
    queryFn: () => api.getMe(),
    staleTime: 5 * 60_000, // 5 minutes
    retry: false, // Don't retry auth failures
  })
}

export function useLogout(): UseMutationResult<{ message: string }, Error, void> {
  const qc = useQueryClient()
  const { logout } = useAuth()
  return useMutation({
    mutationFn: () => api.logout(),
    onSuccess: () => {
      // Clear all cached queries after logout
      qc.clear()
      void logout()
    },
    onError: () => {
      // Even on error, clear local state
      qc.clear()
      void logout()
    },
  })
}

// ---------------------------------------------------------------------------
// Block Types
// ---------------------------------------------------------------------------

export function useBlockTypes(category?: string): UseQueryResult<BlockType[]> {
  return useQuery({
    queryKey: queryKeys.blockTypes.list(category),
    queryFn: () => api.listBlockTypes(category),
  })
}

export function useBlockType(id: string): UseQueryResult<BlockType> {
  return useQuery({
    queryKey: queryKeys.blockTypes.detail(id),
    queryFn: () => api.getBlockType(id),
    enabled: Boolean(id),
  })
}

export function useCreateBlockType(): UseMutationResult<BlockType, Error, Partial<BlockType>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createBlockType(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.blockTypes.all })
    },
  })
}

export function useUpdateBlockType(): UseMutationResult<BlockType, Error, { id: string; data: Partial<BlockType> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateBlockType(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.blockTypes.all })
      qc.setQueryData(queryKeys.blockTypes.detail(updated.id), updated)
    },
  })
}

export function useDeleteBlockType(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteBlockType(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.blockTypes.all })
    },
  })
}

// ---------------------------------------------------------------------------
// Themes
// ---------------------------------------------------------------------------

export function useThemes(): UseQueryResult<Theme[]> {
  return useQuery({
    queryKey: queryKeys.themes.list(),
    queryFn: () => api.listThemes(),
  })
}

export function useActiveTheme(): UseQueryResult<Theme> {
  return useQuery({
    queryKey: queryKeys.themes.active(),
    queryFn: () => api.getActiveTheme(),
  })
}

export function useTheme(id: string): UseQueryResult<Theme> {
  return useQuery({
    queryKey: queryKeys.themes.detail(id),
    queryFn: () => api.getTheme(id),
    enabled: Boolean(id),
  })
}

export function useCreateTheme(): UseMutationResult<Theme, Error, Partial<Theme>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createTheme(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.themes.all })
    },
  })
}

export function useUpdateTheme(): UseMutationResult<Theme, Error, { id: string; data: Partial<Theme> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateTheme(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.themes.all })
      qc.setQueryData(queryKeys.themes.detail(updated.id), updated)
    },
  })
}

export function useActivateTheme(): UseMutationResult<Theme, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.activateTheme(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.themes.all })
    },
  })
}

export function useDuplicateTheme(): UseMutationResult<Theme, Error, { id: string; name: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, name }) => api.duplicateTheme(id, name),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.themes.all })
    },
  })
}

export function useDeleteTheme(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteTheme(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.themes.all })
    },
  })
}

// ---------------------------------------------------------------------------
// Shipping — Zones
// ---------------------------------------------------------------------------

export function useShippingZones(): UseQueryResult<ShippingZone[]> {
  return useQuery({
    queryKey: queryKeys.shipping.zones,
    queryFn: () => api.listShippingZones(),
  })
}

export function useShippingZone(id: string): UseQueryResult<ShippingZone> {
  return useQuery({
    queryKey: queryKeys.shipping.zone(id),
    queryFn: () => api.getShippingZone(id),
    enabled: Boolean(id),
  })
}

export function useCreateShippingZone(): UseMutationResult<ShippingZone, Error, Partial<ShippingZone>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createShippingZone(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.shipping.zones })
    },
  })
}

export function useUpdateShippingZone(): UseMutationResult<ShippingZone, Error, { id: string; data: Partial<ShippingZone> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateShippingZone(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.shipping.zones })
      qc.setQueryData(queryKeys.shipping.zone(updated.id), updated)
    },
  })
}

export function useDeleteShippingZone(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteShippingZone(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.shipping.zones })
    },
  })
}

// ---------------------------------------------------------------------------
// Shipping — Rates
// ---------------------------------------------------------------------------

export function useShippingRates(zoneId: string): UseQueryResult<ShippingRate[]> {
  return useQuery({
    queryKey: queryKeys.shipping.rates(zoneId),
    queryFn: () => api.listShippingRates(zoneId),
    enabled: Boolean(zoneId),
  })
}

export function useCreateShippingRate(): UseMutationResult<ShippingRate, Error, { zoneId: string; data: Partial<ShippingRate> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ zoneId, data }) => api.createShippingRate(zoneId, data),
    onSuccess: (created) => {
      void qc.invalidateQueries({ queryKey: queryKeys.shipping.rates(created.zone_id) })
    },
  })
}

export function useUpdateShippingRate(): UseMutationResult<ShippingRate, Error, { id: string; data: Partial<ShippingRate> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateShippingRate(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.shipping.rates(updated.zone_id) })
    },
  })
}

export function useDeleteShippingRate(): UseMutationResult<void, Error, { id: string; zoneId: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id }) => api.deleteShippingRate(id),
    onSuccess: (_data, { zoneId }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.shipping.rates(zoneId) })
    },
  })
}

// ---------------------------------------------------------------------------
// Tax — Rates
// ---------------------------------------------------------------------------

export function useTaxRates(): UseQueryResult<TaxRate[]> {
  return useQuery({
    queryKey: queryKeys.tax.rates,
    queryFn: () => api.listTaxRates(),
  })
}

export function useCreateTaxRate(): UseMutationResult<TaxRate, Error, Partial<TaxRate>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createTaxRate(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.tax.rates })
    },
  })
}

export function useUpdateTaxRate(): UseMutationResult<TaxRate, Error, { id: string; data: Partial<TaxRate> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateTaxRate(id, data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.tax.rates })
    },
  })
}

export function useDeleteTaxRate(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteTaxRate(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.tax.rates })
    },
  })
}

// ---------------------------------------------------------------------------
// Payments
// ---------------------------------------------------------------------------

export function usePayment(id: string): UseQueryResult<Payment> {
  return useQuery({
    queryKey: queryKeys.payments.detail(id),
    queryFn: () => api.getPayment(id),
    enabled: Boolean(id),
  })
}

export function useOrderPayments(orderId: string): UseQueryResult<Payment[]> {
  return useQuery({
    queryKey: queryKeys.payments.byOrder(orderId),
    queryFn: () => api.listOrderPayments(orderId),
    enabled: Boolean(orderId),
  })
}

export function useCreateRefund(): UseMutationResult<Refund, Error, { paymentId: string; amount: number; reason?: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ paymentId, amount, reason }) => api.createRefund(paymentId, { amount, reason }),
    onSuccess: (_data, { paymentId }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.payments.refunds(paymentId) })
      void qc.invalidateQueries({ queryKey: queryKeys.payments.detail(paymentId) })
    },
  })
}

export function useRefunds(paymentId: string): UseQueryResult<Refund[]> {
  return useQuery({
    queryKey: queryKeys.payments.refunds(paymentId),
    queryFn: () => api.listRefunds(paymentId),
    enabled: Boolean(paymentId),
  })
}

// ---------------------------------------------------------------------------
// Product — Options
// ---------------------------------------------------------------------------

export function useProductOptions(productId: string): UseQueryResult<ProductOption[]> {
  return useQuery({
    queryKey: queryKeys.products.options(productId),
    queryFn: () => api.listProductOptions(productId),
    enabled: Boolean(productId),
  })
}

export function useCreateProductOption(): UseMutationResult<
  ProductOption,
  Error,
  { productId: string; name: string; position: number; values: string[] }
> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ productId, ...data }) => api.createProductOption(productId, data),
    onSuccess: (created) => {
      void qc.invalidateQueries({ queryKey: queryKeys.products.options(created.product_id) })
    },
  })
}

export function useDeleteProductOption(): UseMutationResult<void, Error, { optionId: string; productId: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ optionId }) => api.deleteProductOption(optionId),
    onSuccess: (_data, { productId }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.products.options(productId) })
    },
  })
}

// ---------------------------------------------------------------------------
// Product — Variants
// ---------------------------------------------------------------------------

export function useProductVariants(productId: string): UseQueryResult<ProductVariant[]> {
  return useQuery({
    queryKey: queryKeys.products.variants(productId),
    queryFn: () => api.listProductVariants(productId),
    enabled: Boolean(productId),
  })
}

export function useCreateProductVariant(): UseMutationResult<
  ProductVariant,
  Error,
  {
    productId: string
    sku: string
    price_amount: number
    price_currency: string
    stock: number
    options: Record<string, string>
  }
> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ productId, ...data }) => api.createProductVariant(productId, data),
    onSuccess: (created) => {
      void qc.invalidateQueries({ queryKey: queryKeys.products.variants(created.product_id) })
    },
  })
}

export function useDeleteProductVariant(): UseMutationResult<void, Error, { variantId: string; productId: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ variantId }) => api.deleteProductVariant(variantId),
    onSuccess: (_data, { productId }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.products.variants(productId) })
    },
  })
}

// ---------------------------------------------------------------------------
// Customer Groups
// ---------------------------------------------------------------------------

export function useCustomerGroups(): UseQueryResult<CustomerGroup[]> {
  return useQuery({
    queryKey: queryKeys.customerGroups.list(),
    queryFn: () => api.listCustomerGroups(),
  })
}

export function useCustomerGroup(id: string): UseQueryResult<CustomerGroup> {
  return useQuery({
    queryKey: queryKeys.customerGroups.detail(id),
    queryFn: () => api.getCustomerGroup(id),
    enabled: Boolean(id),
  })
}

export function useGroupMembers(groupId: string): UseQueryResult<GroupMembership[]> {
  return useQuery({
    queryKey: queryKeys.customerGroups.members(groupId),
    queryFn: () => api.listGroupMembers(groupId),
    enabled: Boolean(groupId),
  })
}

export function useCreateCustomerGroup(): UseMutationResult<CustomerGroup, Error, api.CustomerGroupInput> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createCustomerGroup(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.customerGroups.all })
    },
  })
}

export function useUpdateCustomerGroup(): UseMutationResult<CustomerGroup, Error, { id: string } & api.CustomerGroupInput> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }) => api.updateCustomerGroup(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.customerGroups.all })
      qc.setQueryData(queryKeys.customerGroups.detail(updated.id), updated)
    },
  })
}

export function useDeleteCustomerGroup(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteCustomerGroup(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.customerGroups.all })
    },
  })
}

export function useAddGroupMember(): UseMutationResult<GroupMembership, Error, { groupId: string; customerId: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ groupId, customerId }) => api.addGroupMember(groupId, customerId),
    onSuccess: (_data, { groupId }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.customerGroups.members(groupId) })
      void qc.invalidateQueries({ queryKey: queryKeys.customerGroups.all })
    },
  })
}

export function useRemoveGroupMember(): UseMutationResult<void, Error, { groupId: string; customerId: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ groupId, customerId }) => api.removeGroupMember(groupId, customerId),
    onSuccess: (_data, { groupId }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.customerGroups.members(groupId) })
      void qc.invalidateQueries({ queryKey: queryKeys.customerGroups.all })
    },
  })
}
