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
} from 'lucide-react'

// Store pages
import Navbar from './components/store/Navbar'
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
      <footer className="bg-white border-t border-gray-100 py-8 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center text-sm text-gray-400">
          © {new Date().getFullYear()} Hada Store. All rights reserved.
        </div>
      </footer>
    </div>
  ),
})

// ─── Admin layout route ──────────────────────────────────────────────────────

interface NavItem {
  to: string
  label: string
  icon: React.ComponentType<{ size?: number; className?: string }>
}

const adminNavItems: NavItem[] = [
  { to: '/admin', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/admin/products', label: 'Products', icon: Package },
  { to: '/admin/catalog', label: 'Catalog', icon: Tags },
  { to: '/admin/orders', label: 'Orders', icon: ShoppingCart },
  { to: '/admin/customers', label: 'Customers', icon: Users },
  { to: '/admin/pages', label: 'Pages', icon: FileText },
  { to: '/admin/promos', label: 'Promos', icon: Tag },
  { to: '/admin/shipping', label: 'Shipping', icon: Truck },
  { to: '/admin/tax', label: 'Tax', icon: Receipt },
  { to: '/admin/payments', label: 'Payments', icon: CreditCard },
  { to: '/admin/customer-groups', label: 'Customer Groups', icon: UsersRound },
  { to: '/admin/import-export', label: 'Import / Export', icon: ArrowUpDown },
  { to: '/admin/media', label: 'Media', icon: Image },
  { to: '/admin/agent', label: 'Agent Chat', icon: Bot },
  { to: '/admin/marketplace', label: 'Marketplace', icon: Puzzle },
  { to: '/admin/theme', label: 'Theme', icon: Palette },
  { to: '/admin/settings', label: 'Settings', icon: Settings2 },
]

function AdminLayout() {
  const { user, logout, isLoading } = useAuth()

  return (
    <div className="flex h-screen overflow-hidden bg-slate-50 font-sans antialiased">
      {/* Sidebar */}
      <aside className="flex w-60 shrink-0 flex-col border-r border-slate-200 bg-white">
        <div className="flex h-16 items-center border-b border-slate-200 px-4">
          <Store size={22} className="text-indigo-600" />
          <span className="ml-2 text-base font-semibold text-slate-800 tracking-tight">
            Hada Commerce
          </span>
        </div>
        <nav className="flex-1 overflow-y-auto px-3 py-4">
          <p className="mb-1 px-3 text-xs font-semibold uppercase tracking-widest text-slate-400">
            Admin
          </p>
          <ul className="space-y-0.5">
            {adminNavItems.map((item) => {
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
          <h1 className="text-sm font-medium text-slate-500">Hada Commerce Admin</h1>
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
  adminMediaRoute,
  adminAgentRoute,
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
