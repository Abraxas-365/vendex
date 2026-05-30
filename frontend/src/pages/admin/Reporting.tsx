import { useState } from 'react'
import { TrendingUp, Users, ShoppingCart, DollarSign } from 'lucide-react'
import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  BarChart,
  Bar,
} from 'recharts'
import {
  useDashboardSales,
  useDashboardRevenue,
  useDashboardTopProducts,
  useDashboardCustomers,
  useDashboardFunnel,
} from '../../lib/hooks'

const COLORS = ['#6366f1', '#8b5cf6', '#a78bfa', '#c4b5fd', '#ddd6fe']

function formatMoney(cents: number, currency = 'USD') {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency, maximumFractionDigits: 0 }).format(cents / 100)
}

function StatCard({ label, value, icon: Icon }: { label: string; value: string; icon: React.ElementType }) {
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
      <div className="flex items-center justify-between">
        <p className="text-sm font-medium text-gray-500">{label}</p>
        <div className="rounded-lg bg-gray-100 p-2">
          <Icon className="h-4 w-4 text-gray-600" />
        </div>
      </div>
      <p className="mt-3 text-2xl font-bold text-gray-900">{value}</p>
    </div>
  )
}

export default function Reporting() {
  const today = new Date()
  const thirtyDaysAgo = new Date(today)
  thirtyDaysAgo.setDate(today.getDate() - 30)

  const [dateFrom, setDateFrom] = useState(thirtyDaysAgo.toISOString().split('T')[0])
  const [dateTo, setDateTo] = useState(today.toISOString().split('T')[0])

  const dateParams = { from: dateFrom, to: dateTo }

  const { data: sales, isLoading: salesLoading } = useDashboardSales(dateParams)
  const { data: revenue = [], isLoading: revenueLoading } = useDashboardRevenue(dateParams)
  const { data: topProducts = [], isLoading: topLoading } = useDashboardTopProducts(dateParams)
  const { data: customers, isLoading: customersLoading } = useDashboardCustomers(dateParams)
  const { data: funnel = [], isLoading: funnelLoading } = useDashboardFunnel(dateParams)

  const isLoading = salesLoading || revenueLoading || topLoading || customersLoading || funnelLoading

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Sales Reporting</h1>
          <p className="mt-1 text-sm text-gray-500">KPIs, conversion funnels and revenue charts</p>
        </div>
        {/* Date Range */}
        <div className="flex items-center gap-3">
          <input
            type="date"
            value={dateFrom}
            onChange={(e) => setDateFrom(e.target.value)}
            className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
          />
          <span className="text-sm text-gray-400">to</span>
          <input
            type="date"
            value={dateTo}
            onChange={(e) => setDateTo(e.target.value)}
            className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
          />
        </div>
      </div>

      {isLoading && (
        <div className="flex items-center justify-center py-12">
          <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
        </div>
      )}

      {!isLoading && (
        <>
          {/* KPI Cards */}
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
            <StatCard
              label="Total Revenue"
              value={sales ? formatMoney(sales.total_revenue.amount, sales.total_revenue.currency) : '—'}
              icon={DollarSign}
            />
            <StatCard
              label="Orders"
              value={sales ? sales.total_orders.toString() : '—'}
              icon={ShoppingCart}
            />
            <StatCard
              label="Avg Order Value"
              value={sales ? formatMoney(sales.average_order_value.amount, sales.average_order_value.currency) : '—'}
              icon={TrendingUp}
            />
            <StatCard
              label="New Customers"
              value={customers ? customers.new_customers.toString() : '—'}
              icon={Users}
            />
          </div>

          {/* Revenue Chart */}
          <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
            <h2 className="mb-4 text-base font-semibold text-gray-900">Revenue Over Time</h2>
            <ResponsiveContainer width="100%" height={260}>
              <LineChart data={revenue}>
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis dataKey="date" tick={{ fontSize: 12 }} tickLine={false} />
                <YAxis
                  tick={{ fontSize: 12 }}
                  tickLine={false}
                  axisLine={false}
                  tickFormatter={(v: number) => `$${(v / 100).toFixed(0)}`}
                />
                <Tooltip formatter={(value: unknown) => [`$${(Number(value) / 100).toFixed(2)}`, 'Revenue']} />
                <Line type="monotone" dataKey="revenue" stroke="#6366f1" strokeWidth={2} dot={false} />
              </LineChart>
            </ResponsiveContainer>
          </div>

          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            {/* Top Products */}
            <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
              <h2 className="mb-4 text-base font-semibold text-gray-900">Top Products</h2>
              {topProducts.length > 0 && (
                <ResponsiveContainer width="100%" height={200}>
                  <BarChart data={topProducts.slice(0, 7)} layout="vertical">
                    <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                    <XAxis type="number" tick={{ fontSize: 11 }} tickLine={false} />
                    <YAxis dataKey="product_name" type="category" width={100} tick={{ fontSize: 11 }} tickLine={false} axisLine={false} />
                    <Tooltip />
                    <Bar dataKey="revenue" fill="#6366f1" radius={[0, 4, 4, 0]} />
                  </BarChart>
                </ResponsiveContainer>
              )}
              <div className="mt-3 overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-gray-100 text-left text-xs text-gray-500">
                      <th className="py-2 pr-4 font-medium">Product</th>
                      <th className="py-2 pr-4 font-medium">Units</th>
                      <th className="py-2 font-medium">Revenue</th>
                    </tr>
                  </thead>
                  <tbody>
                    {topProducts.slice(0, 5).map((p, i) => (
                      <tr key={i} className="border-b border-gray-50 last:border-0">
                        <td className="py-2 pr-4 font-medium text-gray-900">{p.product_name}</td>
                        <td className="py-2 pr-4 text-gray-600">{p.units_sold}</td>
                        <td className="py-2 font-medium text-gray-900">{formatMoney(p.revenue)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            {/* Customer Stats + Funnel */}
            <div className="space-y-4">
              {customers && (
                <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
                  <h2 className="mb-3 text-base font-semibold text-gray-900">Customer Stats</h2>
                  <div className="grid grid-cols-2 gap-3">
                    <div className="rounded-lg bg-gray-50 p-3 text-center">
                      <p className="text-xs text-gray-500">Total</p>
                      <p className="text-xl font-bold text-gray-900">{customers.total_customers}</p>
                    </div>
                    <div className="rounded-lg bg-gray-50 p-3 text-center">
                      <p className="text-xs text-gray-500">New</p>
                      <p className="text-xl font-bold text-indigo-600">{customers.new_customers}</p>
                    </div>
                    <div className="rounded-lg bg-gray-50 p-3 text-center">
                      <p className="text-xs text-gray-500">Returning</p>
                      <p className="text-xl font-bold text-purple-600">{customers.returning_customers}</p>
                    </div>
                    <div className="rounded-lg bg-gray-50 p-3 text-center">
                      <p className="text-xs text-gray-500">Avg LTV</p>
                      <p className="text-xl font-bold text-gray-900">
                        {formatMoney(customers.average_lifetime_value.amount, customers.average_lifetime_value.currency)}
                      </p>
                    </div>
                  </div>
                </div>
              )}

              {funnel.length > 0 && (
                <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
                  <h2 className="mb-3 text-base font-semibold text-gray-900">Conversion Funnel</h2>
                  <div className="space-y-2">
                    {funnel.map((step, i) => (
                      <div key={i} className="flex items-center gap-3">
                        <span className="w-28 truncate text-xs text-gray-600">{step.step}</span>
                        <div className="flex-1 rounded-full bg-gray-100 h-5 overflow-hidden">
                          <div
                            className="h-full rounded-full transition-all"
                            style={{
                              width: `${step.conversion_rate}%`,
                              backgroundColor: COLORS[i % COLORS.length],
                            }}
                          />
                        </div>
                        <span className="w-10 text-right text-xs font-medium text-gray-700">
                          {step.conversion_rate.toFixed(1)}%
                        </span>
                        <span className="w-14 text-right text-xs text-gray-400">
                          {step.count.toLocaleString()}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </>
      )}
    </div>
  )
}
