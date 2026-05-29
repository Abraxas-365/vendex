import { useState } from 'react'
import { Link } from '@tanstack/react-router'
import {
  Package,
  ShoppingCart,
  FileText,
  Tag,
  Users,
  DollarSign,
  Loader2,
  AlertCircle,
} from 'lucide-react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts'
import {
  useDashboardStats,
  useRevenueTimeline,
  useTopProducts,
  useOrderStatusBreakdown,
  useRecentOrders,
} from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function formatMoney(cents: number, currency = 'USD'): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(cents / 100)
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// ---------------------------------------------------------------------------
// Sub-components
// ---------------------------------------------------------------------------

function StatCard({
  title,
  value,
  icon,
  loading,
}: {
  title: string
  value: string | number
  icon: React.ReactNode
  loading?: boolean
}) {
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <div className="flex items-center justify-between">
        <div className="rounded-lg bg-gray-100 p-3 text-gray-600">{icon}</div>
      </div>
      <div className="mt-4">
        <p className="text-sm font-medium text-gray-500">{title}</p>
        {loading ? (
          <div className="mt-2 h-8 w-20 animate-pulse rounded bg-gray-200" />
        ) : (
          <p className="mt-1 text-3xl font-bold text-gray-900">{value}</p>
        )}
      </div>
    </div>
  )
}

const STATUS_COLORS: Record<string, string> = {
  pending: '#eab308',
  confirmed: '#3b82f6',
  processing: '#6366f1',
  shipped: '#a855f7',
  delivered: '#22c55e',
  cancelled: '#ef4444',
}

const STATUS_BADGE: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  confirmed: 'bg-blue-100 text-blue-800',
  processing: 'bg-indigo-100 text-indigo-800',
  shipped: 'bg-purple-100 text-purple-800',
  delivered: 'bg-green-100 text-green-800',
  cancelled: 'bg-red-100 text-red-800',
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

export default function Dashboard() {
  const [revenueDays, setRevenueDays] = useState(30)

  const stats = useDashboardStats()
  const revenue = useRevenueTimeline(revenueDays)
  const topProducts = useTopProducts(5)
  const statusBreakdown = useOrderStatusBreakdown()
  const recentOrders = useRecentOrders(5)

  const s = stats.data

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-1 text-sm text-gray-500">
          Welcome back. Here's what's happening with your store.
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Total Products"
          value={s?.total_products ?? 0}
          icon={<Package className="h-5 w-5" />}
          loading={stats.isLoading}
        />
        <StatCard
          title="Total Orders"
          value={s?.total_orders ?? 0}
          icon={<ShoppingCart className="h-5 w-5" />}
          loading={stats.isLoading}
        />
        <StatCard
          title="Total Customers"
          value={s?.total_customers ?? 0}
          icon={<Users className="h-5 w-5" />}
          loading={stats.isLoading}
        />
        <StatCard
          title="Total Revenue"
          value={s ? formatMoney(s.total_revenue, s.currency || 'USD') : '$0'}
          icon={<DollarSign className="h-5 w-5" />}
          loading={stats.isLoading}
        />
      </div>

      <div className="grid grid-cols-1 gap-5 sm:grid-cols-3">
        <StatCard
          title="Pending Orders"
          value={s?.pending_orders ?? 0}
          icon={<ShoppingCart className="h-5 w-5" />}
          loading={stats.isLoading}
        />
        <StatCard
          title="Active Promos"
          value={s?.active_promos ?? 0}
          icon={<Tag className="h-5 w-5" />}
          loading={stats.isLoading}
        />
        <StatCard
          title="Pending Pages"
          value={s?.pending_pages ?? 0}
          icon={<FileText className="h-5 w-5" />}
          loading={stats.isLoading}
        />
      </div>

      {/* Revenue Chart */}
      <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900">Revenue</h2>
          <div className="flex gap-1 rounded-lg bg-gray-100 p-1">
            {[7, 30, 90].map((d) => (
              <button
                key={d}
                onClick={() => setRevenueDays(d)}
                className={`rounded-md px-3 py-1 text-xs font-medium transition-colors ${
                  revenueDays === d
                    ? 'bg-white text-gray-900 shadow-sm'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                {d}d
              </button>
            ))}
          </div>
        </div>
        {revenue.isLoading ? (
          <div className="flex h-64 items-center justify-center">
            <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
          </div>
        ) : revenue.isError ? (
          <div className="flex h-64 items-center justify-center gap-2 text-sm text-gray-500">
            <AlertCircle className="h-4 w-4" />
            Failed to load revenue data
          </div>
        ) : (
          <ResponsiveContainer width="100%" height={280}>
            <LineChart data={revenue.data ?? []}>
              <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
              <XAxis
                dataKey="date"
                tickFormatter={(v: string) =>
                  new Date(v).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
                }
                tick={{ fontSize: 12 }}
                stroke="#94a3b8"
              />
              <YAxis
                tickFormatter={(v: number) => `$${(v / 100).toLocaleString()}`}
                tick={{ fontSize: 12 }}
                stroke="#94a3b8"
                width={70}
              />
              <Tooltip
                formatter={(value) => [formatMoney(Number(value)), 'Revenue']}
                labelFormatter={(label) =>
                  new Date(String(label)).toLocaleDateString('en-US', {
                    weekday: 'short',
                    month: 'short',
                    day: 'numeric',
                  })
                }
                contentStyle={{
                  borderRadius: '8px',
                  border: '1px solid #e2e8f0',
                  fontSize: '13px',
                }}
              />
              <Line
                type="monotone"
                dataKey="amount"
                stroke="#4f46e5"
                strokeWidth={2}
                dot={false}
                activeDot={{ r: 4, fill: '#4f46e5' }}
              />
            </LineChart>
          </ResponsiveContainer>
        )}
      </div>

      {/* Two-column: Top Products + Order Status Breakdown */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Top Products */}
        <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="border-b border-gray-200 px-6 py-4">
            <h2 className="text-lg font-semibold text-gray-900">Top Products</h2>
          </div>
          {topProducts.isLoading ? (
            <div className="flex h-48 items-center justify-center">
              <Loader2 className="h-5 w-5 animate-spin text-gray-400" />
            </div>
          ) : !topProducts.data?.length ? (
            <div className="p-6 text-center text-sm text-gray-500">No sales data yet</div>
          ) : (
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">Product</th>
                  <th className="px-6 py-3 font-medium text-right">Sold</th>
                  <th className="px-6 py-3 font-medium text-right">Revenue</th>
                </tr>
              </thead>
              <tbody>
                {topProducts.data.map((p) => (
                  <tr key={p.product_id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-3 text-sm font-medium text-gray-900">
                      {p.product_name}
                    </td>
                    <td className="px-6 py-3 text-right text-sm text-gray-600">{p.total_sold}</td>
                    <td className="px-6 py-3 text-right text-sm font-medium text-gray-900">
                      {formatMoney(p.revenue)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>

        {/* Order Status Breakdown */}
        <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="border-b border-gray-200 px-6 py-4">
            <h2 className="text-lg font-semibold text-gray-900">Order Status</h2>
          </div>
          {statusBreakdown.isLoading ? (
            <div className="flex h-48 items-center justify-center">
              <Loader2 className="h-5 w-5 animate-spin text-gray-400" />
            </div>
          ) : !statusBreakdown.data?.length ? (
            <div className="p-6 text-center text-sm text-gray-500">No orders yet</div>
          ) : (
            <div className="flex items-center gap-6 p-6">
              <div className="h-48 w-48 shrink-0">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={statusBreakdown.data}
                      dataKey="count"
                      nameKey="status"
                      cx="50%"
                      cy="50%"
                      innerRadius={40}
                      outerRadius={70}
                      strokeWidth={2}
                    >
                      {statusBreakdown.data.map((entry) => (
                        <Cell
                          key={entry.status}
                          fill={STATUS_COLORS[entry.status] ?? '#94a3b8'}
                        />
                      ))}
                    </Pie>
                    <Tooltip
                      formatter={(value) => [Number(value), '']}
                      contentStyle={{
                        borderRadius: '8px',
                        border: '1px solid #e2e8f0',
                        fontSize: '13px',
                      }}
                    />
                  </PieChart>
                </ResponsiveContainer>
              </div>
              <div className="flex flex-col gap-2">
                {statusBreakdown.data.map((s) => (
                  <div key={s.status} className="flex items-center gap-2 text-sm">
                    <span
                      className="h-3 w-3 rounded-full"
                      style={{ backgroundColor: STATUS_COLORS[s.status] ?? '#94a3b8' }}
                    />
                    <span className="capitalize text-gray-700">{s.status}</span>
                    <span className="font-medium text-gray-900">{s.count}</span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Recent Orders */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        <div className="border-b border-gray-200 px-6 py-4">
          <h2 className="text-lg font-semibold text-gray-900">Recent Orders</h2>
        </div>
        {recentOrders.isLoading ? (
          <div className="flex h-32 items-center justify-center">
            <Loader2 className="h-5 w-5 animate-spin text-gray-400" />
          </div>
        ) : !recentOrders.data?.length ? (
          <div className="p-6 text-center text-sm text-gray-500">No orders yet</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">Order ID</th>
                  <th className="px-6 py-3 font-medium">Total</th>
                  <th className="px-6 py-3 font-medium">Status</th>
                  <th className="px-6 py-3 font-medium">Date</th>
                </tr>
              </thead>
              <tbody>
                {recentOrders.data.map((order) => (
                  <tr key={order.id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-4 text-sm font-medium text-indigo-600">
                      <Link to="/admin/orders/$id" params={{ id: order.id }}>
                        {order.id.slice(0, 8)}...
                      </Link>
                    </td>
                    <td className="px-6 py-4 text-sm font-medium text-gray-900">
                      {formatMoney(order.total_amount, order.currency || 'USD')}
                    </td>
                    <td className="px-6 py-4">
                      <span
                        className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${STATUS_BADGE[order.status] ?? 'bg-gray-100 text-gray-800'}`}
                      >
                        {order.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {formatDate(order.created_at)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
