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
  giftCards: {
    all: ['gift-cards'] as const,
    list: (params?: PaginationParams) => ['gift-cards', 'list', params] as const,
    detail: (id: string) => ['gift-cards', 'detail', id] as const,
    transactions: (id: string) => ['gift-cards', id, 'transactions'] as const,
  },
  cartRecovery: {
    all: ['cart-recovery'] as const,
    list: (params?: PaginationParams) => ['cart-recovery', 'list', params] as const,
    stats: ['cart-recovery', 'stats'] as const,
  },
  currencyRates: {
    all: ['currency-rates'] as const,
    list: () => ['currency-rates', 'list'] as const,
    currencies: ['currencies'] as const,
  },
  translations: {
    bundle: (entityType: string, entityId: string, locale: string) =>
      ['i18n', entityType, entityId, locale] as const,
    locales: (entityType: string, entityId: string) =>
      ['i18n', entityType, entityId, 'locales'] as const,
    supportedLocales: ['i18n', 'supported-locales'] as const,
  },
  subscriptions: {
    all: ['subscriptions'] as const,
    list: (params?: PaginationParams) => ['subscriptions', 'list', params] as const,
    due: ['subscriptions', 'due'] as const,
    detail: (id: string) => ['subscriptions', 'detail', id] as const,
    billing: (id: string, params?: PaginationParams) => ['subscriptions', id, 'billing', params] as const,
  },
  inventory: {
    warehouses: ['inventory', 'warehouses'] as const,
    stock: (productId: string) => ['inventory', 'stock', productId] as const,
    lowStock: ['inventory', 'low-stock'] as const,
    movements: (productId: string) => ['inventory', 'movements', productId] as const,
  },
  reviews: {
    all: ['reviews'] as const,
    list: (params?: PaginationParams & { status?: ReviewStatus }) => ['reviews', 'list', params] as const,
    detail: (id: string) => ['reviews', 'detail', id] as const,
  },
  returns: {
    all: ['returns'] as const,
    list: (params?: PaginationParams & { status?: ReturnStatus }) => ['returns', 'list', params] as const,
    detail: (id: string) => ['returns', 'detail', id] as const,
  },
  webhooks: {
    all: ['webhooks'] as const,
    list: ['webhooks', 'list'] as const,
    deliveries: (webhookId: string) => ['webhooks', webhookId, 'deliveries'] as const,
  },
  auditLogs: {
    all: ['audit'] as const,
    list: (params?: Record<string, string | number | undefined>) => ['audit', 'list', params] as const,
    detail: (id: string) => ['audit', 'detail', id] as const,
    stats: ['audit', 'stats'] as const,
  },
  loyalty: {
    rewards: ['loyalty', 'rewards'] as const,
    accounts: {
      all: ['loyalty', 'accounts'] as const,
      list: (params?: PaginationParams) => ['loyalty', 'accounts', 'list', params] as const,
      detail: (id: string) => ['loyalty', 'accounts', id] as const,
      transactions: (id: string) => ['loyalty', 'accounts', id, 'transactions'] as const,
    },
  },
  bundles: {
    all: ['bundles'] as const,
    list: (params?: PaginationParams) => ['bundles', 'list', params] as const,
    detail: (id: string) => ['bundles', 'detail', id] as const,
    price: (id: string) => ['bundles', id, 'price'] as const,
  },
  dashboardReporting: {
    sales: (params?: { from?: string; to?: string }) => ['dashboard-reporting', 'sales', params] as const,
    topProducts: (params?: { from?: string; to?: string }) => ['dashboard-reporting', 'top-products', params] as const,
    revenue: (params?: { from?: string; to?: string }) => ['dashboard-reporting', 'revenue', params] as const,
    customers: (params?: { from?: string; to?: string }) => ['dashboard-reporting', 'customers', params] as const,
    funnel: (params?: { from?: string; to?: string }) => ['dashboard-reporting', 'funnel', params] as const,
  },
  socialAccounts: {
    all: ['social-accounts'] as const,
    list: (params?: PaginationParams & { provider?: SocialProvider }) => ['social-accounts', 'list', params] as const,
    byCustomer: (customerId: string) => ['social-accounts', 'customer', customerId] as const,
  },
  notifications: {
    all: ['notifications'] as const,
    list: (params?: PaginationParams & { read?: boolean }) => ['notifications', 'list', params] as const,
    unreadCount: ['notifications', 'unread-count'] as const,
  },
  storefronts: {
    all: ['storefronts'] as const,
    list: (params?: PaginationParams) => ['storefronts', 'list', params] as const,
    detail: (id: string) => ['storefronts', 'detail', id] as const,
    catalogs: (id: string) => ['storefronts', id, 'catalogs'] as const,
  },
  bulkOperations: {
    all: ['bulk-operations'] as const,
    list: (params?: PaginationParams) => ['bulk-operations', 'list', params] as const,
    detail: (id: string) => ['bulk-operations', 'detail', id] as const,
    items: (id: string) => ['bulk-operations', id, 'items'] as const,
  },
  blog: {
    posts: {
      all: ['blog-posts'] as const,
      list: (params?: PaginationParams) => ['blog-posts', 'list', params] as const,
      detail: (id: string) => ['blog-posts', 'detail', id] as const,
    },
    categories: ['blog-categories'] as const,
  },
  adminCollections: {
    all: ['admin-collections'] as const,
    list: (params?: PaginationParams) => ['admin-collections', 'list', params] as const,
    detail: (id: string) => ['admin-collections', 'detail', id] as const,
    products: (id: string) => ['admin-collections', id, 'products'] as const,
  },
  experiments: {
    all: ['experiments'] as const,
    list: (params?: PaginationParams) => ['experiments', 'list', params] as const,
    detail: (id: string) => ['experiments', 'detail', id] as const,
    results: (id: string) => ['experiments', id, 'results'] as const,
  },
  recommendations: {
    all: ['recommendations'] as const,
    list: (params?: PaginationParams) => ['recommendations', 'list', params] as const,
    forProduct: (productId: string) => ['recommendations', 'product', productId] as const,
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

// ---------------------------------------------------------------------------
// Gift Cards
// ---------------------------------------------------------------------------

export function useGiftCards(params?: PaginationParams): UseQueryResult<PaginatedResult<GiftCard>> {
  return useQuery({
    queryKey: queryKeys.giftCards.list(params),
    queryFn: () => api.listGiftCards(params),
  })
}

export function useGiftCard(id: string): UseQueryResult<GiftCard> {
  return useQuery({
    queryKey: queryKeys.giftCards.detail(id),
    queryFn: () => api.getGiftCard(id),
    enabled: Boolean(id),
  })
}

export function useGiftCardTransactions(id: string): UseQueryResult<GiftCardTransaction[]> {
  return useQuery({
    queryKey: queryKeys.giftCards.transactions(id),
    queryFn: () => api.listGiftCardTransactions(id),
    enabled: Boolean(id),
  })
}

export function useCreateGiftCard(): UseMutationResult<GiftCard, Error, api.CreateGiftCardInput> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createGiftCard(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.giftCards.all })
    },
  })
}

export function useUpdateGiftCard(): UseMutationResult<GiftCard, Error, { id: string; data: Partial<GiftCard> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateGiftCard(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.giftCards.all })
      qc.setQueryData(queryKeys.giftCards.detail(updated.id), updated)
    },
  })
}

export function useDeleteGiftCard(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteGiftCard(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.giftCards.all })
    },
  })
}

export function useDisableGiftCard(): UseMutationResult<GiftCard, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.disableGiftCard(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.giftCards.all })
      qc.setQueryData(queryKeys.giftCards.detail(updated.id), updated)
    },
  })
}

// ---------------------------------------------------------------------------
// Cart Recovery
// ---------------------------------------------------------------------------

export function useRecoveryEmails(params?: PaginationParams): UseQueryResult<PaginatedResult<RecoveryEmail>> {
  return useQuery({
    queryKey: queryKeys.cartRecovery.list(params),
    queryFn: () => api.listRecoveryEmails(params),
  })
}

export function useRecoveryStats(): UseQueryResult<RecoveryStats> {
  return useQuery({
    queryKey: queryKeys.cartRecovery.stats,
    queryFn: () => api.getRecoveryStats(),
  })
}

export function useUpdateRecoveryStatus(): UseMutationResult<RecoveryEmail, Error, { id: string; status: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }) => api.updateRecoveryStatus(id, status),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.cartRecovery.all })
    },
  })
}

// ---------------------------------------------------------------------------
// Currency Rates
// ---------------------------------------------------------------------------

export function useCurrencyRates(): UseQueryResult<CurrencyRate[]> {
  return useQuery({
    queryKey: queryKeys.currencyRates.list(),
    queryFn: () => api.listCurrencyRates(),
  })
}

export function useSupportedCurrencies(): UseQueryResult<string[]> {
  return useQuery({
    queryKey: queryKeys.currencyRates.currencies,
    queryFn: () => api.listSupportedCurrencies(),
  })
}

export function useSetCurrencyRate(): UseMutationResult<CurrencyRate, Error, api.SetCurrencyRateInput> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.setCurrencyRate(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.currencyRates.all })
    },
  })
}

export function useDeleteCurrencyRate(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteCurrencyRate(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.currencyRates.all })
    },
  })
}

export function useConvertCurrency(): UseMutationResult<
  ConvertResult,
  Error,
  { amount: number; currency: string; target_currency: string }
> {
  return useMutation({
    mutationFn: (data) => api.convertCurrency(data),
  })
}

// ---------------------------------------------------------------------------
// Translations / I18n
// ---------------------------------------------------------------------------

export function useTranslationBundle(
  entityType: string,
  entityId: string,
  locale: string,
): UseQueryResult<TranslationBundle> {
  return useQuery({
    queryKey: queryKeys.translations.bundle(entityType, entityId, locale),
    queryFn: () => api.getTranslationBundle(entityType, entityId, locale),
    enabled: Boolean(entityType) && Boolean(entityId) && Boolean(locale),
  })
}

export function useEntityLocales(entityType: string, entityId: string): UseQueryResult<string[]> {
  return useQuery({
    queryKey: queryKeys.translations.locales(entityType, entityId),
    queryFn: () => api.listEntityLocales(entityType, entityId),
    enabled: Boolean(entityType) && Boolean(entityId),
  })
}

export function useSupportedLocales(): UseQueryResult<string[]> {
  return useQuery({
    queryKey: queryKeys.translations.supportedLocales,
    queryFn: () => api.listSupportedLocales(),
  })
}

export function useSetTranslations(): UseMutationResult<
  TranslationBundle,
  Error,
  { entityType: string; entityId: string; locale: string; fields: Record<string, string> }
> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ entityType, entityId, locale, fields }) =>
      api.setTranslations(entityType, entityId, locale, fields),
    onSuccess: (_data, { entityType, entityId, locale }) => {
      void qc.invalidateQueries({
        queryKey: queryKeys.translations.bundle(entityType, entityId, locale),
      })
      void qc.invalidateQueries({
        queryKey: queryKeys.translations.locales(entityType, entityId),
      })
    },
  })
}

export function useDeleteTranslationField(): UseMutationResult<
  void,
  Error,
  { entityType: string; entityId: string; locale: string; field: string }
> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ entityType, entityId, locale, field }) =>
      api.deleteTranslationField(entityType, entityId, locale, field),
    onSuccess: (_data, { entityType, entityId, locale }) => {
      void qc.invalidateQueries({
        queryKey: queryKeys.translations.bundle(entityType, entityId, locale),
      })
    },
  })
}

// ---------------------------------------------------------------------------
// Subscriptions
// ---------------------------------------------------------------------------

export function useSubscriptions(params?: PaginationParams): UseQueryResult<PaginatedResult<Subscription>> {
  return useQuery({
    queryKey: queryKeys.subscriptions.list(params),
    queryFn: () => api.listSubscriptions(params),
  })
}

export function useDueSubscriptions(): UseQueryResult<Subscription[]> {
  return useQuery({
    queryKey: queryKeys.subscriptions.due,
    queryFn: () => api.listDueSubscriptions(),
  })
}

export function useSubscription(id: string): UseQueryResult<Subscription> {
  return useQuery({
    queryKey: queryKeys.subscriptions.detail(id),
    queryFn: () => api.getSubscription(id),
    enabled: Boolean(id),
  })
}

export function useCreateSubscription(): UseMutationResult<Subscription, Error, api.CreateSubscriptionInput> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createSubscription(data),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.subscriptions.all })
    },
  })
}

export function useCancelSubscription(): UseMutationResult<Subscription, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.cancelSubscription(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.subscriptions.all })
      qc.setQueryData(queryKeys.subscriptions.detail(updated.id), updated)
    },
  })
}

export function usePauseSubscription(): UseMutationResult<Subscription, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.pauseSubscription(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.subscriptions.all })
      qc.setQueryData(queryKeys.subscriptions.detail(updated.id), updated)
    },
  })
}

export function useResumeSubscription(): UseMutationResult<Subscription, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.resumeSubscription(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.subscriptions.all })
      qc.setQueryData(queryKeys.subscriptions.detail(updated.id), updated)
    },
  })
}

export function useBillingRecords(subscriptionId: string, params?: PaginationParams): UseQueryResult<PaginatedResult<BillingRecord>> {
  return useQuery({
    queryKey: queryKeys.subscriptions.billing(subscriptionId, params),
    queryFn: () => api.listBillingRecords(subscriptionId, params),
    enabled: Boolean(subscriptionId),
  })
}

// ---------------------------------------------------------------------------
// Inventory (027)
// ---------------------------------------------------------------------------

export function useWarehouses(): UseQueryResult<Warehouse[]> {
  return useQuery({
    queryKey: queryKeys.inventory.warehouses,
    queryFn: () => api.listWarehouses(),
  })
}

export function useCreateWarehouse(): UseMutationResult<Warehouse, Error, Partial<Warehouse>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createWarehouse(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.inventory.warehouses }) },
  })
}

export function useUpdateWarehouse(): UseMutationResult<Warehouse, Error, { id: string; data: Partial<Warehouse> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateWarehouse(id, data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.inventory.warehouses }) },
  })
}

export function useDeleteWarehouse(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteWarehouse(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.inventory.warehouses }) },
  })
}

export function useStockLevels(productId: string): UseQueryResult<StockLevel[]> {
  return useQuery({
    queryKey: queryKeys.inventory.stock(productId),
    queryFn: () => api.getStockLevels(productId),
    enabled: Boolean(productId),
  })
}

export function useAdjustStock(): UseMutationResult<StockMovement, Error, { product_id: string; warehouse_id: string; quantity: number; note?: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.adjustStock(data),
    onSuccess: (_data, vars) => {
      void qc.invalidateQueries({ queryKey: queryKeys.inventory.stock(vars.product_id) })
      void qc.invalidateQueries({ queryKey: queryKeys.inventory.lowStock })
    },
  })
}

export function useLowStockAlerts(): UseQueryResult<LowStockAlert[]> {
  return useQuery({
    queryKey: queryKeys.inventory.lowStock,
    queryFn: () => api.getLowStockAlerts(),
  })
}

export function useStockMovements(productId: string): UseQueryResult<StockMovement[]> {
  return useQuery({
    queryKey: queryKeys.inventory.movements(productId),
    queryFn: () => api.getStockMovements(productId),
    enabled: Boolean(productId),
  })
}

// ---------------------------------------------------------------------------
// Reviews (028)
// ---------------------------------------------------------------------------

export function useReviews(params?: PaginationParams & { status?: ReviewStatus }): UseQueryResult<PaginatedResult<Review>> {
  return useQuery({
    queryKey: queryKeys.reviews.list(params),
    queryFn: () => api.listReviews(params),
  })
}

export function useReview(id: string): UseQueryResult<Review> {
  return useQuery({
    queryKey: queryKeys.reviews.detail(id),
    queryFn: () => api.getReview(id),
    enabled: Boolean(id),
  })
}

export function useApproveReview(): UseMutationResult<Review, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.approveReview(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.reviews.all }) },
  })
}

export function useRejectReview(): UseMutationResult<Review, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.rejectReview(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.reviews.all }) },
  })
}

export function useRespondToReview(): UseMutationResult<Review, Error, { id: string; response: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, response }) => api.respondToReview(id, response),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.reviews.all }) },
  })
}

// ---------------------------------------------------------------------------
// Returns (029)
// ---------------------------------------------------------------------------

export function useReturns(params?: PaginationParams & { status?: ReturnStatus }): UseQueryResult<PaginatedResult<ReturnRequest>> {
  return useQuery({
    queryKey: queryKeys.returns.list(params),
    queryFn: () => api.listReturns(params),
  })
}

export function useReturn(id: string): UseQueryResult<ReturnRequest> {
  return useQuery({
    queryKey: queryKeys.returns.detail(id),
    queryFn: () => api.getReturn(id),
    enabled: Boolean(id),
  })
}

export function useApproveReturn(): UseMutationResult<ReturnRequest, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.approveReturn(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.returns.all }) },
  })
}

export function useRejectReturn(): UseMutationResult<ReturnRequest, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.rejectReturn(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.returns.all }) },
  })
}

export function useMarkReturnReceived(): UseMutationResult<ReturnRequest, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.markReturnReceived(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.returns.all }) },
  })
}

export function useMarkReturnRefunded(): UseMutationResult<ReturnRequest, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.markReturnRefunded(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.returns.all }) },
  })
}

export function useCloseReturn(): UseMutationResult<ReturnRequest, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.closeReturn(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.returns.all }) },
  })
}

// ---------------------------------------------------------------------------
// Webhooks (030)
// ---------------------------------------------------------------------------

export function useWebhooks(): UseQueryResult<Webhook[]> {
  return useQuery({
    queryKey: queryKeys.webhooks.list,
    queryFn: () => api.listWebhooks(),
  })
}

export function useCreateWebhook(): UseMutationResult<Webhook, Error, Partial<Webhook>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createWebhook(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.webhooks.all }) },
  })
}

export function useUpdateWebhook(): UseMutationResult<Webhook, Error, { id: string; data: Partial<Webhook> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateWebhook(id, data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.webhooks.all }) },
  })
}

export function useDeleteWebhook(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteWebhook(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.webhooks.all }) },
  })
}

export function useToggleWebhook(): UseMutationResult<Webhook, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.toggleWebhook(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.webhooks.all }) },
  })
}

export function useWebhookDeliveries(webhookId: string): UseQueryResult<WebhookDelivery[]> {
  return useQuery({
    queryKey: queryKeys.webhooks.deliveries(webhookId),
    queryFn: () => api.listWebhookDeliveries(webhookId),
    enabled: Boolean(webhookId),
  })
}

export function useRetryWebhookDelivery(): UseMutationResult<WebhookDelivery, Error, { deliveryId: string; webhookId: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ deliveryId }) => api.retryWebhookDelivery(deliveryId),
    onSuccess: (_data, { webhookId }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.webhooks.deliveries(webhookId) })
    },
  })
}

// ---------------------------------------------------------------------------
// Audit Logs (031)
// ---------------------------------------------------------------------------

export function useAuditLogs(params?: PaginationParams & {
  user_id?: string; action?: string; resource_type?: string; from?: string; to?: string
}): UseQueryResult<PaginatedResult<AuditLog>> {
  return useQuery({
    queryKey: queryKeys.auditLogs.list(params as Record<string, string | number | undefined>),
    queryFn: () => api.listAuditLogs(params),
  })
}

export function useAuditLog(id: string): UseQueryResult<AuditLog> {
  return useQuery({
    queryKey: queryKeys.auditLogs.detail(id),
    queryFn: () => api.getAuditLog(id),
    enabled: Boolean(id),
  })
}

export function useAuditStats(): UseQueryResult<AuditStats> {
  return useQuery({
    queryKey: queryKeys.auditLogs.stats,
    queryFn: () => api.getAuditStats(),
  })
}

// ---------------------------------------------------------------------------
// Loyalty (032)
// ---------------------------------------------------------------------------

export function useLoyaltyRewards(): UseQueryResult<LoyaltyReward[]> {
  return useQuery({
    queryKey: queryKeys.loyalty.rewards,
    queryFn: () => api.listLoyaltyRewards(),
  })
}

export function useCreateLoyaltyReward(): UseMutationResult<LoyaltyReward, Error, Partial<LoyaltyReward>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createLoyaltyReward(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.loyalty.rewards }) },
  })
}

export function useUpdateLoyaltyReward(): UseMutationResult<LoyaltyReward, Error, { id: string; data: Partial<LoyaltyReward> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateLoyaltyReward(id, data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.loyalty.rewards }) },
  })
}

export function useLoyaltyAccounts(params?: PaginationParams): UseQueryResult<PaginatedResult<LoyaltyAccount>> {
  return useQuery({
    queryKey: queryKeys.loyalty.accounts.list(params),
    queryFn: () => api.listLoyaltyAccounts(params),
  })
}

export function useLoyaltyAccount(id: string): UseQueryResult<LoyaltyAccount> {
  return useQuery({
    queryKey: queryKeys.loyalty.accounts.detail(id),
    queryFn: () => api.getLoyaltyAccount(id),
    enabled: Boolean(id),
  })
}

export function useAdjustLoyaltyPoints(): UseMutationResult<LoyaltyAccount, Error, { id: string; points: number; note: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }) => api.adjustLoyaltyPoints(id, data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.loyalty.accounts.all }) },
  })
}

export function useLoyaltyTransactions(id: string): UseQueryResult<LoyaltyTransaction[]> {
  return useQuery({
    queryKey: queryKeys.loyalty.accounts.transactions(id),
    queryFn: () => api.getLoyaltyTransactions(id),
    enabled: Boolean(id),
  })
}

// ---------------------------------------------------------------------------
// Bundles (033)
// ---------------------------------------------------------------------------

export function useBundles(params?: PaginationParams): UseQueryResult<PaginatedResult<Bundle>> {
  return useQuery({
    queryKey: queryKeys.bundles.list(params),
    queryFn: () => api.listBundles(params),
  })
}

export function useCreateBundle(): UseMutationResult<Bundle, Error, Partial<Bundle>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createBundle(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.bundles.all }) },
  })
}

export function useUpdateBundle(): UseMutationResult<Bundle, Error, { id: string; data: Partial<Bundle> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateBundle(id, data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.bundles.all }) },
  })
}

export function useDeleteBundle(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteBundle(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.bundles.all }) },
  })
}

export function useAddBundleItem(): UseMutationResult<BundleItem, Error, { bundleId: string; product_id: string; quantity: number; discount_percent: number }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ bundleId, ...data }) => api.addBundleItem(bundleId, data),
    onSuccess: (_data, { bundleId }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.bundles.detail(bundleId) })
      void qc.invalidateQueries({ queryKey: queryKeys.bundles.price(bundleId) })
    },
  })
}

export function useRemoveBundleItem(): UseMutationResult<void, Error, { bundleId: string; itemId: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ bundleId, itemId }) => api.removeBundleItem(bundleId, itemId),
    onSuccess: (_data, { bundleId }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.bundles.detail(bundleId) })
      void qc.invalidateQueries({ queryKey: queryKeys.bundles.price(bundleId) })
    },
  })
}

export function useBundlePrice(bundleId: string): UseQueryResult<BundlePrice> {
  return useQuery({
    queryKey: queryKeys.bundles.price(bundleId),
    queryFn: () => api.getBundlePrice(bundleId),
    enabled: Boolean(bundleId),
  })
}

// ---------------------------------------------------------------------------
// Dashboard Reporting (034)
// ---------------------------------------------------------------------------

export function useDashboardSales(params?: { from?: string; to?: string }): UseQueryResult<SalesOverview> {
  return useQuery({
    queryKey: queryKeys.dashboardReporting.sales(params),
    queryFn: () => api.getDashboardSales(params),
  })
}

export function useDashboardTopProducts(params?: { from?: string; to?: string }): UseQueryResult<TopProductReport[]> {
  return useQuery({
    queryKey: queryKeys.dashboardReporting.topProducts(params),
    queryFn: () => api.getDashboardTopProducts(params),
  })
}

export function useDashboardRevenue(params?: { from?: string; to?: string }): UseQueryResult<RevenueData[]> {
  return useQuery({
    queryKey: queryKeys.dashboardReporting.revenue(params),
    queryFn: () => api.getDashboardRevenue(params),
  })
}

export function useDashboardCustomers(params?: { from?: string; to?: string }): UseQueryResult<CustomerStats> {
  return useQuery({
    queryKey: queryKeys.dashboardReporting.customers(params),
    queryFn: () => api.getDashboardCustomers(params),
  })
}

export function useDashboardFunnel(params?: { from?: string; to?: string }): UseQueryResult<FunnelStep[]> {
  return useQuery({
    queryKey: queryKeys.dashboardReporting.funnel(params),
    queryFn: () => api.getDashboardFunnel(params),
  })
}

// ---------------------------------------------------------------------------
// Social Accounts (035)
// ---------------------------------------------------------------------------

export function useSocialAccounts(params?: PaginationParams & { provider?: SocialProvider }): UseQueryResult<PaginatedResult<SocialAccount>> {
  return useQuery({
    queryKey: queryKeys.socialAccounts.list(params),
    queryFn: () => api.listSocialAccounts(params),
  })
}

export function useCustomerSocialAccounts(customerId: string): UseQueryResult<SocialAccount[]> {
  return useQuery({
    queryKey: queryKeys.socialAccounts.byCustomer(customerId),
    queryFn: () => api.listCustomerSocialAccounts(customerId),
    enabled: Boolean(customerId),
  })
}

// ---------------------------------------------------------------------------
// Notifications (036)
// ---------------------------------------------------------------------------

export function useNotifications(params?: PaginationParams & { read?: boolean }): UseQueryResult<PaginatedResult<AdminNotification>> {
  return useQuery({
    queryKey: queryKeys.notifications.list(params),
    queryFn: () => api.listNotifications(params),
  })
}

export function useUnreadCount(): UseQueryResult<UnreadCount> {
  return useQuery({
    queryKey: queryKeys.notifications.unreadCount,
    queryFn: () => api.getUnreadCount(),
  })
}

export function useMarkAllNotificationsRead(): UseMutationResult<void, Error, void> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.markAllNotificationsRead(),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.notifications.all }) },
  })
}

export function useMarkNotificationRead(): UseMutationResult<AdminNotification, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.markNotificationRead(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.notifications.all }) },
  })
}

export function useDeleteNotification(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteNotification(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.notifications.all }) },
  })
}

// ---------------------------------------------------------------------------
// Storefronts / Multistore (037)
// ---------------------------------------------------------------------------

export function useStorefronts(params?: PaginationParams): UseQueryResult<PaginatedResult<Storefront>> {
  return useQuery({
    queryKey: queryKeys.storefronts.list(params),
    queryFn: () => api.listStorefronts(params),
  })
}

export function useStorefront(id: string): UseQueryResult<Storefront> {
  return useQuery({
    queryKey: queryKeys.storefronts.detail(id),
    queryFn: () => api.getStorefront(id),
    enabled: Boolean(id),
  })
}

export function useCreateStorefront(): UseMutationResult<Storefront, Error, Partial<Storefront>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createStorefront(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.storefronts.all }) },
  })
}

export function useUpdateStorefront(): UseMutationResult<Storefront, Error, { id: string; data: Partial<Storefront> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateStorefront(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.storefronts.all })
      qc.setQueryData(queryKeys.storefronts.detail(updated.id), updated)
    },
  })
}

export function useDeleteStorefront(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteStorefront(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.storefronts.all }) },
  })
}

export function useSetDefaultStorefront(): UseMutationResult<Storefront, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.setDefaultStorefront(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.storefronts.all }) },
  })
}

export function useStorefrontCatalogs(id: string): UseQueryResult<StorefrontCatalog[]> {
  return useQuery({
    queryKey: queryKeys.storefronts.catalogs(id),
    queryFn: () => api.getStorefrontCatalogs(id),
    enabled: Boolean(id),
  })
}

// ---------------------------------------------------------------------------
// Bulk Operations (038)
// ---------------------------------------------------------------------------

export function useBulkOperations(params?: PaginationParams): UseQueryResult<PaginatedResult<BulkOperation>> {
  return useQuery({
    queryKey: queryKeys.bulkOperations.list(params),
    queryFn: () => api.listBulkOperations(params),
  })
}

export function useBulkOperation(id: string): UseQueryResult<BulkOperation> {
  return useQuery({
    queryKey: queryKeys.bulkOperations.detail(id),
    queryFn: () => api.getBulkOperation(id),
    enabled: Boolean(id),
  })
}

export function useCreateBulkOperation(): UseMutationResult<BulkOperation, Error, Partial<BulkOperation>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createBulkOperation(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.bulkOperations.all }) },
  })
}

export function useBulkOperationItems(id: string): UseQueryResult<BulkOperationItem[]> {
  return useQuery({
    queryKey: queryKeys.bulkOperations.items(id),
    queryFn: () => api.getBulkOperationItems(id),
    enabled: Boolean(id),
  })
}

export function useProcessBulkOperation(): UseMutationResult<BulkOperation, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.processBulkOperation(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.bulkOperations.all })
      qc.setQueryData(queryKeys.bulkOperations.detail(updated.id), updated)
    },
  })
}

export function useCancelBulkOperation(): UseMutationResult<BulkOperation, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.cancelBulkOperation(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.bulkOperations.all })
      qc.setQueryData(queryKeys.bulkOperations.detail(updated.id), updated)
    },
  })
}

// ---------------------------------------------------------------------------
// Blog (039)
// ---------------------------------------------------------------------------

export function useBlogPosts(params?: PaginationParams): UseQueryResult<PaginatedResult<BlogPost>> {
  return useQuery({
    queryKey: queryKeys.blog.posts.list(params),
    queryFn: () => api.listBlogPosts(params),
  })
}

export function useBlogPost(id: string): UseQueryResult<BlogPost> {
  return useQuery({
    queryKey: queryKeys.blog.posts.detail(id),
    queryFn: () => api.getBlogPost(id),
    enabled: Boolean(id),
  })
}

export function useCreateBlogPost(): UseMutationResult<BlogPost, Error, Partial<BlogPost>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createBlogPost(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.blog.posts.all }) },
  })
}

export function useUpdateBlogPost(): UseMutationResult<BlogPost, Error, { id: string; data: Partial<BlogPost> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateBlogPost(id, data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.blog.posts.all }) },
  })
}

export function useDeleteBlogPost(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteBlogPost(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.blog.posts.all }) },
  })
}

export function usePublishBlogPost(): UseMutationResult<BlogPost, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.publishBlogPost(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.blog.posts.all }) },
  })
}

export function useArchiveBlogPost(): UseMutationResult<BlogPost, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.archiveBlogPost(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.blog.posts.all }) },
  })
}

export function useBlogCategories(): UseQueryResult<BlogCategory[]> {
  return useQuery({
    queryKey: queryKeys.blog.categories,
    queryFn: () => api.listBlogCategories(),
  })
}

export function useCreateBlogCategory(): UseMutationResult<BlogCategory, Error, Partial<BlogCategory>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createBlogCategory(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.blog.categories }) },
  })
}

export function useUpdateBlogCategory(): UseMutationResult<BlogCategory, Error, { id: string; data: Partial<BlogCategory> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateBlogCategory(id, data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.blog.categories }) },
  })
}

export function useDeleteBlogCategory(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteBlogCategory(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.blog.categories }) },
  })
}

// ---------------------------------------------------------------------------
// Admin Collections (040)
// ---------------------------------------------------------------------------

export function useAdminCollections(params?: PaginationParams): UseQueryResult<PaginatedResult<AdminCollection>> {
  return useQuery({
    queryKey: queryKeys.adminCollections.list(params),
    queryFn: () => api.listAdminCollections(params),
  })
}

export function useAdminCollection(id: string): UseQueryResult<AdminCollection> {
  return useQuery({
    queryKey: queryKeys.adminCollections.detail(id),
    queryFn: () => api.getAdminCollection(id),
    enabled: Boolean(id),
  })
}

export function useCreateAdminCollection(): UseMutationResult<AdminCollection, Error, Partial<AdminCollection>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createAdminCollection(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.adminCollections.all }) },
  })
}

export function useUpdateAdminCollection(): UseMutationResult<AdminCollection, Error, { id: string; data: Partial<AdminCollection> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateAdminCollection(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.adminCollections.all })
      qc.setQueryData(queryKeys.adminCollections.detail(updated.id), updated)
    },
  })
}

export function useDeleteAdminCollection(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteAdminCollection(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.adminCollections.all }) },
  })
}

export function useCollectionProducts(id: string): UseQueryResult<CollectionProduct[]> {
  return useQuery({
    queryKey: queryKeys.adminCollections.products(id),
    queryFn: () => api.getCollectionProducts(id),
    enabled: Boolean(id),
  })
}

export function useAddCollectionProduct(): UseMutationResult<CollectionProduct, Error, { id: string; product_id: string; sort_order?: number }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...data }) => api.addCollectionProduct(id, data),
    onSuccess: (_data, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.adminCollections.products(id) })
      void qc.invalidateQueries({ queryKey: queryKeys.adminCollections.all })
    },
  })
}

export function useRemoveCollectionProduct(): UseMutationResult<void, Error, { id: string; productId: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, productId }) => api.removeCollectionProduct(id, productId),
    onSuccess: (_data, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.adminCollections.products(id) })
      void qc.invalidateQueries({ queryKey: queryKeys.adminCollections.all })
    },
  })
}

// ---------------------------------------------------------------------------
// A/B Testing — Experiments (041)
// ---------------------------------------------------------------------------

export function useExperiments(params?: PaginationParams): UseQueryResult<PaginatedResult<Experiment>> {
  return useQuery({
    queryKey: queryKeys.experiments.list(params),
    queryFn: () => api.listExperiments(params),
  })
}

export function useExperiment(id: string): UseQueryResult<Experiment> {
  return useQuery({
    queryKey: queryKeys.experiments.detail(id),
    queryFn: () => api.getExperiment(id),
    enabled: Boolean(id),
  })
}

export function useCreateExperiment(): UseMutationResult<Experiment, Error, Partial<Experiment>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createExperiment(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.experiments.all }) },
  })
}

export function useUpdateExperiment(): UseMutationResult<Experiment, Error, { id: string; data: Partial<Experiment> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateExperiment(id, data),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.experiments.all })
      qc.setQueryData(queryKeys.experiments.detail(updated.id), updated)
    },
  })
}

export function useDeleteExperiment(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteExperiment(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.experiments.all }) },
  })
}

export function useStartExperiment(): UseMutationResult<Experiment, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.startExperiment(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.experiments.all })
      qc.setQueryData(queryKeys.experiments.detail(updated.id), updated)
    },
  })
}

export function usePauseExperiment(): UseMutationResult<Experiment, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.pauseExperiment(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.experiments.all })
      qc.setQueryData(queryKeys.experiments.detail(updated.id), updated)
    },
  })
}

export function useCompleteExperiment(): UseMutationResult<Experiment, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.completeExperiment(id),
    onSuccess: (updated) => {
      void qc.invalidateQueries({ queryKey: queryKeys.experiments.all })
      qc.setQueryData(queryKeys.experiments.detail(updated.id), updated)
    },
  })
}

export function useExperimentResults(id: string): UseQueryResult<ExperimentResults> {
  return useQuery({
    queryKey: queryKeys.experiments.results(id),
    queryFn: () => api.getExperimentResults(id),
    enabled: Boolean(id),
  })
}

export function useAddExperimentVariant(): UseMutationResult<ExperimentVariant, Error, { id: string; data: Partial<ExperimentVariant> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.addExperimentVariant(id, data),
    onSuccess: (_data, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.experiments.detail(id) })
    },
  })
}

export function useDeleteExperimentVariant(): UseMutationResult<void, Error, { id: string; variantId: string }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, variantId }) => api.deleteExperimentVariant(id, variantId),
    onSuccess: (_data, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.experiments.detail(id) })
    },
  })
}

// ---------------------------------------------------------------------------
// Recommendations (042)
// ---------------------------------------------------------------------------

export function useRecommendationRules(params?: PaginationParams): UseQueryResult<PaginatedResult<RecommendationRule>> {
  return useQuery({
    queryKey: queryKeys.recommendations.list(params),
    queryFn: () => api.listRecommendationRules(params),
  })
}

export function useCreateRecommendationRule(): UseMutationResult<RecommendationRule, Error, Partial<RecommendationRule>> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data) => api.createRecommendationRule(data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.recommendations.all }) },
  })
}

export function useUpdateRecommendationRule(): UseMutationResult<RecommendationRule, Error, { id: string; data: Partial<RecommendationRule> }> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }) => api.updateRecommendationRule(id, data),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.recommendations.all }) },
  })
}

export function useDeleteRecommendationRule(): UseMutationResult<void, Error, string> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id) => api.deleteRecommendationRule(id),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: queryKeys.recommendations.all }) },
  })
}

export function useRecommendations(productId: string): UseQueryResult<RecommendedProduct[]> {
  return useQuery({
    queryKey: queryKeys.recommendations.forProduct(productId),
    queryFn: () => api.getRecommendations(productId),
    enabled: Boolean(productId),
  })
}
