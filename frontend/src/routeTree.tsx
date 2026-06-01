import { useState } from 'react'
import {
  createRootRoute,
  createRoute,
  Outlet,
  Link,
} from '@tanstack/react-router'
import {
  LayoutDashboard,
  Package,
  ShoppingCart,
  Users,
  FileText,
  Tag,
  Tags,
  Image,
  Bot,
  ExternalLink,
  Store,
  Puzzle,
  Settings2,
  LogOut,
  Palette,
  Truck,
  Receipt,
  CreditCard,
  ArrowUpDown,
  UsersRound,
  Gift,
  MailQuestion,
  DollarSign,
  Languages,
  RefreshCw,
  Warehouse,
  Star,
  RotateCcw,
  Webhook,
  ClipboardList,
  Award,
  Layers,
  BarChart2,
  Share2,
  Bell,
  Globe2,
  Newspaper,
  FolderTree,
  FlaskConical,
  Sparkles,
  Zap,
  MonitorPlay,
  ShieldCheck,
  Brain,
  ChevronDown,
} from 'lucide-react'

// Store pages
import Navbar from './components/store/Navbar'
import Footer from './components/store/Footer'
import Home from './pages/store/Home'
import ProductList from './pages/store/ProductList'
import ProductDetail from './pages/store/ProductDetail'
import Cart from './pages/store/Cart'
import Checkout from './pages/store/Checkout'
import DynamicPage from './pages/store/DynamicPage'

// Admin pages
import Dashboard from './pages/admin/Dashboard'
import Products from './pages/admin/Products'
import Orders from './pages/admin/Orders'
import OrderDetail from './pages/admin/OrderDetail'
import Customers from './pages/admin/Customers'
import CustomerDetail from './pages/admin/CustomerDetail'
import Pages from './pages/admin/Pages'
import Promos from './pages/admin/Promos'
import Media from './pages/admin/Media'
import AgentChat from './pages/admin/AgentChat'
import Marketplace from './pages/admin/Marketplace'
import PluginView from './pages/admin/PluginView'
import Catalog from './pages/admin/Catalog'
import Settings from './pages/admin/Settings'
import PageEditor from './pages/admin/PageEditor'
import ThemeEditor from './pages/admin/ThemeEditor'
import Shipping from './pages/admin/Shipping'
import Tax from './pages/admin/Tax'
import Payments from './pages/admin/Payments'
import ImportExport from './pages/admin/ImportExport'
import CustomerGroups from './pages/admin/CustomerGroups'
import GiftCards from './pages/admin/GiftCards'
import CartRecovery from './pages/admin/CartRecovery'
import CurrencyRates from './pages/admin/CurrencyRates'
import Translations from './pages/admin/Translations'
import Subscriptions from './pages/admin/Subscriptions'
import Inventory from './pages/admin/Inventory'
import Reviews from './pages/admin/Reviews'
import Returns from './pages/admin/Returns'
import Webhooks from './pages/admin/Webhooks'
import AuditLogs from './pages/admin/AuditLogs'
import Loyalty from './pages/admin/Loyalty'
import Bundles from './pages/admin/Bundles'
import Reporting from './pages/admin/Reporting'
import SocialAccounts from './pages/admin/SocialAccounts'
import Notifications from './pages/admin/Notifications'
import Multistores from './pages/admin/Multistores'
import Blog from './pages/admin/Blog'
import Collections from './pages/admin/Collections'
import ABTesting from './pages/admin/ABTesting'
import Recommendations from './pages/admin/Recommendations'
import BulkOperations from './pages/admin/BulkOperations'
import Presets from './pages/admin/Presets'
import Workspaces from './pages/admin/Workspaces'
import WorkspaceViewPage from './pages/admin/WorkspaceView'
import Approvals from './pages/admin/Approvals'
import AgentMemory from './pages/admin/AgentMemory'
import AgentTriggers from './pages/admin/AgentTriggers'

// Auth pages
import Login from './pages/auth/Login'
import Callback from './pages/auth/Callback'

// Auth context
import { useAuth } from './lib/auth'

// ─── Root route (bare) ───────────────────────────────────────────────────────

const rootRoute = createRootRoute({
  component: () => <Outlet />,
})

// ─── Store layout route ──────────────────────────────────────────────────────

const storeLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: '_store',
  component: () => (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <Navbar />
      <main className="flex-1">
        <Outlet />
      </main>
      <Footer />
    </div>
  ),
})

// ─── Admin layout route ──────────────────────────────────────────────────────

interface NavItem {
  to: string
  label: string
  icon: React.ComponentType<{ size?: number; className?: string }>
}

interface NavSection {
  id: string
  label: string
  items: NavItem[]
}

const adminNavSections: NavSection[] = [
  {
    id: 'commerce',
    label: 'Commerce',
    items: [
      { to: '/admin/products', label: 'Products', icon: Package },
      { to: '/admin/catalog', label: 'Catalog', icon: Tags },
      { to: '/admin/collections', label: 'Collections', icon: FolderTree },
      { to: '/admin/orders', label: 'Orders', icon: ShoppingCart },
      { to: '/admin/inventory', label: 'Inventory', icon: Warehouse },
      { to: '/admin/bundles', label: 'Bundles', icon: Layers },
      { to: '/admin/subscriptions', label: 'Subscriptions', icon: RefreshCw },
    ],
  },
  {
    id: 'customers',
    label: 'Customers',
    items: [
      { to: '/admin/customers', label: 'Customers', icon: Users },
      { to: '/admin/customer-groups', label: 'Customer Groups', icon: UsersRound },
      { to: '/admin/reviews', label: 'Reviews', icon: Star },
      { to: '/admin/loyalty', label: 'Loyalty', icon: Award },
      { to: '/admin/cart-recovery', label: 'Cart Recovery', icon: MailQuestion },
    ],
  },
  {
    id: 'content',
    label: 'Content',
    items: [
      { to: '/admin/pages', label: 'Pages', icon: FileText },
      { to: '/admin/blog', label: 'Blog', icon: Newspaper },
      { to: '/admin/media', label: 'Media', icon: Image },
      { to: '/admin/translations', label: 'Translations', icon: Languages },
    ],
  },
  {
    id: 'marketing',
    label: 'Marketing',
    items: [
      { to: '/admin/promos', label: 'Promos', icon: Tag },
      { to: '/admin/gift-cards', label: 'Gift Cards', icon: Gift },
      { to: '/admin/recommendations', label: 'Recommendations', icon: Sparkles },
      { to: '/admin/ab-testing', label: 'A/B Testing', icon: FlaskConical },
      { to: '/admin/social-accounts', label: 'Social Accounts', icon: Share2 },
    ],
  },
  {
    id: 'operations',
    label: 'Operations',
    items: [
      { to: '/admin/shipping', label: 'Shipping', icon: Truck },
      { to: '/admin/tax', label: 'Tax', icon: Receipt },
      { to: '/admin/payments', label: 'Payments', icon: CreditCard },
      { to: '/admin/currency-rates', label: 'Currency Rates', icon: DollarSign },
      { to: '/admin/import-export', label: 'Import / Export', icon: ArrowUpDown },
      { to: '/admin/returns', label: 'Returns', icon: RotateCcw },
      { to: '/admin/bulk-operations', label: 'Bulk Ops', icon: Zap },
    ],
  },
  {
    id: 'ai-automation',
    label: 'AI & Automation',
    items: [
      { to: '/admin/agent', label: 'Agent Chat', icon: Bot },
      { to: '/admin/agent-memory', label: 'Agent Memory', icon: Brain },
      { to: '/admin/agent-triggers', label: 'Triggers', icon: Zap },
      { to: '/admin/workspaces', label: 'Workspaces', icon: MonitorPlay },
      { to: '/admin/approvals', label: 'Approvals', icon: ShieldCheck },
      { to: '/admin/notifications', label: 'Notifications', icon: Bell },
      { to: '/admin/webhooks', label: 'Webhooks', icon: Webhook },
    ],
  },
  {
    id: 'settings',
    label: 'Settings',
    items: [
      { to: '/admin/theme', label: 'Theme', icon: Palette },
      { to: '/admin/settings', label: 'Settings', icon: Settings2 },
      { to: '/admin/storefronts', label: 'Storefronts', icon: Globe2 },
      { to: '/admin/presets', label: 'Preset Marketplace', icon: Sparkles },
      { to: '/admin/audit-logs', label: 'Audit Logs', icon: ClipboardList },
      { to: '/admin/marketplace', label: 'Marketplace', icon: Puzzle },
      { to: '/admin/reporting', label: 'Reporting', icon: BarChart2 },
    ],
  },
]

const COLLAPSED_STORAGE_KEY = 'admin-nav-collapsed-sections'

function loadCollapsedSections(): Set<string> {
  try {
    const stored = localStorage.getItem(COLLAPSED_STORAGE_KEY)
    if (stored) return new Set(JSON.parse(stored) as string[])
  } catch {
    // ignore parse errors
  }
  return new Set()
}

function saveCollapsedSections(collapsed: Set<string>): void {
  try {
    localStorage.setItem(COLLAPSED_STORAGE_KEY, JSON.stringify([...collapsed]))
  } catch {
    // ignore storage errors
  }
}

function NavSectionGroup({ section }: { section: NavSection }) {
  const [collapsed, setCollapsed] = useState<boolean>(() => {
    const stored = loadCollapsedSections()
    return stored.has(section.id)
  })

  function toggleCollapsed() {
    setCollapsed((prev) => {
      const next = !prev
      const stored = loadCollapsedSections()
      if (next) {
        stored.add(section.id)
      } else {
        stored.delete(section.id)
      }
      saveCollapsedSections(stored)
      return next
    })
  }

  return (
    <div className="mb-1">
      <button
        onClick={toggleCollapsed}
        className="flex w-full items-center justify-between px-3 py-1.5 text-xs font-semibold uppercase tracking-widest text-slate-400 hover:text-slate-600 transition-colors rounded-md hover:bg-slate-50"
      >
        <span>{section.label}</span>
        <ChevronDown
          size={13}
          className={`shrink-0 transition-transform duration-200 ${collapsed ? '-rotate-90' : ''}`}
        />
      </button>
      {!collapsed && (
        <ul className="space-y-0.5 mt-0.5">
          {section.items.map((item) => {
            const Icon = item.icon
            return (
              <li key={item.to}>
                <Link
                  to={item.to}
                  className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-slate-900 [&.active]:bg-indigo-50 [&.active]:text-indigo-700"
                >
                  <Icon size={18} className="shrink-0" />
                  {item.label}
                </Link>
              </li>
            )
          })}
        </ul>
      )}
    </div>
  )
}

function AdminLayout() {
  const { user, logout, isLoading } = useAuth()

  return (
    <div className="flex h-screen overflow-hidden bg-slate-50 font-sans antialiased">
      {/* Sidebar */}
      <aside className="flex w-60 shrink-0 flex-col border-r border-slate-200 bg-white">
        <div className="flex h-16 items-center border-b border-slate-200 px-4">
          <Store size={22} className="text-indigo-600" />
          <span className="ml-2 text-base font-semibold text-slate-800 tracking-tight">
            Vendex
          </span>
        </div>
        <nav className="flex-1 overflow-y-auto px-3 py-4">
          {/* Dashboard — always visible, ungrouped */}
          <ul className="mb-3 space-y-0.5">
            <li>
              <Link
                to="/admin"
                className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-slate-900 [&.active]:bg-indigo-50 [&.active]:text-indigo-700"
              >
                <LayoutDashboard size={18} className="shrink-0" />
                Dashboard
              </Link>
            </li>
          </ul>
          {/* Grouped sections */}
          <div className="space-y-2">
            {adminNavSections.map((section) => (
              <NavSectionGroup key={section.id} section={section} />
            ))}
          </div>
        </nav>
        <div className="border-t border-slate-200 px-3 py-4 space-y-1">
          <a
            href="/"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-500 transition-colors hover:bg-slate-100 hover:text-slate-900"
          >
            <ExternalLink size={18} className="shrink-0" />
            View Store
          </a>
          {!isLoading && (
            <button
              onClick={() => void logout()}
              className="w-full flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-500 transition-colors hover:bg-slate-100 hover:text-slate-900 text-left"
            >
              <LogOut size={18} className="shrink-0" />
              Sign Out
            </button>
          )}
        </div>
      </aside>
      <div className="flex flex-1 flex-col overflow-hidden">
        <header className="flex h-16 shrink-0 items-center justify-between border-b border-slate-200 bg-white px-6">
          <h1 className="text-sm font-medium text-slate-500">Vendex Admin</h1>
          {user && (
            <div className="flex items-center gap-3">
              {user.picture ? (
                <img
                  src={user.picture}
                  alt={user.name}
                  className="w-8 h-8 rounded-full object-cover"
                />
              ) : (
                <div className="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-700 text-sm font-semibold">
                  {user.name.charAt(0).toUpperCase()}
                </div>
              )}
              <div className="text-right">
                <p className="text-sm font-medium text-slate-700 leading-none">{user.name}</p>
                <p className="text-xs text-slate-400 mt-0.5">{user.email}</p>
              </div>
            </div>
          )}
        </header>
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

// Guard component that redirects to /login if not authenticated
function AdminGuard() {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-slate-50">
        <div className="flex flex-col items-center gap-3">
          <div className="w-8 h-8 border-2 border-indigo-600 border-t-transparent rounded-full animate-spin" />
          <p className="text-sm text-slate-500">Loading…</p>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    window.location.replace('/login')
    return null
  }

  return <AdminLayout />
}

const adminLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: '_admin',
  component: AdminGuard,
})

// ─── Auth routes ──────────────────────────────────────────────────────────────

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: Login,
})

// /auth/callback/:provider
const authCallbackRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/auth/callback/$provider',
  component: function AuthCallbackPage() {
    const { provider } = authCallbackRoute.useParams()
    return <Callback provider={provider} />
  },
})

// ─── Store routes ────────────────────────────────────────────────────────────

const homeRoute = createRoute({
  getParentRoute: () => storeLayoutRoute,
  path: '/',
  component: Home,
})

const productsRoute = createRoute({
  getParentRoute: () => storeLayoutRoute,
  path: '/products',
  component: ProductList,
})

const productDetailRoute = createRoute({
  getParentRoute: () => storeLayoutRoute,
  path: '/products/$id',
  component: ProductDetail,
})

const cartRoute = createRoute({
  getParentRoute: () => storeLayoutRoute,
  path: '/cart',
  component: Cart,
})

const checkoutRoute = createRoute({
  getParentRoute: () => storeLayoutRoute,
  path: '/checkout',
  component: Checkout,
})

const dynamicPageRoute = createRoute({
  getParentRoute: () => storeLayoutRoute,
  path: '/pages/$slug',
  component: DynamicPage,
})

// ─── Admin routes ────────────────────────────────────────────────────────────

const adminDashRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin',
  component: Dashboard,
})

const adminProductsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/products',
  component: Products,
})

const adminCatalogRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/catalog',
  component: Catalog,
})

const adminOrdersRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/orders',
  component: Orders,
})

const adminOrderDetailRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/orders/$id',
  component: OrderDetail,
})

const adminCustomersRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/customers',
  component: Customers,
})

const adminCustomerDetailRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/customers/$id',
  component: CustomerDetail,
})

const adminPagesRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/pages',
  component: Pages,
})

const adminPromosRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/promos',
  component: Promos,
})

const adminMediaRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/media',
  component: Media,
})

const adminAgentRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/agent',
  component: AgentChat,
})

const adminMarketplaceRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/marketplace',
  component: Marketplace,
})

const adminPluginViewRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/plugins/$name',
  component: PluginView,
})

const adminSettingsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/settings',
  component: Settings,
})

const adminPageEditorRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/pages/$pageId/edit',
  component: PageEditor,
})

const adminNewBlockPageRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/pages/new-block',
  component: PageEditor,
})

const adminThemeRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/theme',
  component: ThemeEditor,
})

const adminShippingRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/shipping',
  component: Shipping,
})

const adminTaxRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/tax',
  component: Tax,
})

const adminPaymentsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/payments',
  component: Payments,
})

const adminImportExportRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/import-export',
  component: ImportExport,
})

const adminCustomerGroupsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/customer-groups',
  component: CustomerGroups,
})

const adminGiftCardsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/gift-cards',
  component: GiftCards,
})

const adminCartRecoveryRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/cart-recovery',
  component: CartRecovery,
})

const adminCurrencyRatesRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/currency-rates',
  component: CurrencyRates,
})

const adminTranslationsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/translations',
  component: Translations,
})

const adminSubscriptionsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/subscriptions',
  component: Subscriptions,
})

const adminInventoryRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/inventory',
  component: Inventory,
})

const adminReviewsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/reviews',
  component: Reviews,
})

const adminReturnsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/returns',
  component: Returns,
})

const adminWebhooksRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/webhooks',
  component: Webhooks,
})

const adminAuditLogsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/audit-logs',
  component: AuditLogs,
})

const adminLoyaltyRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/loyalty',
  component: Loyalty,
})

const adminBundlesRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/bundles',
  component: Bundles,
})

const adminReportingRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/reporting',
  component: Reporting,
})

const adminSocialAccountsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/social-accounts',
  component: SocialAccounts,
})

const adminNotificationsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/notifications',
  component: Notifications,
})

const adminStorefrontsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/storefronts',
  component: Multistores,
})

const adminBlogRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/blog',
  component: Blog,
})

const adminCollectionsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/collections',
  component: Collections,
})

const adminABTestingRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/ab-testing',
  component: ABTesting,
})

const adminRecommendationsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/recommendations',
  component: Recommendations,
})

const adminBulkOperationsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/bulk-operations',
  component: BulkOperations,
})

const adminPresetsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/presets',
  component: Presets,
})

const adminWorkspacesRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/workspaces',
  component: Workspaces,
})

const adminWorkspaceViewRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/workspaces/$id',
  component: function WorkspaceViewWrapper() {
    const { id } = adminWorkspaceViewRoute.useParams()
    return <WorkspaceViewPage sessionId={id} />
  },
})

const adminApprovalsRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/approvals',
  component: Approvals,
})

const adminAgentMemoryRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/agent-memory',
  component: AgentMemory,
})

const adminAgentTriggersRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: '/admin/agent-triggers',
  component: AgentTriggers,
})

// ─── Route tree ───────────────────────────────────────────────────────────────

const storeTree = storeLayoutRoute.addChildren([
  homeRoute,
  productsRoute,
  productDetailRoute,
  cartRoute,
  checkoutRoute,
  dynamicPageRoute,
])

const adminTree = adminLayoutRoute.addChildren([
  adminDashRoute,
  adminProductsRoute,
  adminCatalogRoute,
  adminOrdersRoute,
  adminOrderDetailRoute,
  adminCustomersRoute,
  adminCustomerDetailRoute,
  adminPagesRoute,
  adminNewBlockPageRoute,
  adminPageEditorRoute,
  adminPromosRoute,
  adminShippingRoute,
  adminTaxRoute,
  adminPaymentsRoute,
  adminImportExportRoute,
  adminCustomerGroupsRoute,
  adminGiftCardsRoute,
  adminCartRecoveryRoute,
  adminCurrencyRatesRoute,
  adminTranslationsRoute,
  adminSubscriptionsRoute,
  adminInventoryRoute,
  adminReviewsRoute,
  adminReturnsRoute,
  adminWebhooksRoute,
  adminAuditLogsRoute,
  adminLoyaltyRoute,
  adminBundlesRoute,
  adminReportingRoute,
  adminSocialAccountsRoute,
  adminNotificationsRoute,
  adminStorefrontsRoute,
  adminBlogRoute,
  adminCollectionsRoute,
  adminABTestingRoute,
  adminRecommendationsRoute,
  adminBulkOperationsRoute,
  adminMediaRoute,
  adminAgentRoute,
  adminPresetsRoute,
  adminWorkspacesRoute,
  adminWorkspaceViewRoute,
  adminApprovalsRoute,
  adminAgentMemoryRoute,
  adminAgentTriggersRoute,
  adminMarketplaceRoute,
  adminPluginViewRoute,
  adminThemeRoute,
  adminSettingsRoute,
])

export const routeTree = rootRoute.addChildren([
  loginRoute,
  authCallbackRoute,
  storeTree,
  adminTree,
])
