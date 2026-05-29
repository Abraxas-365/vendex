import {
  Package,
  ShoppingCart,
  FileText,
  Tag,
  TrendingUp,
  Users,
  DollarSign,
  ArrowUpRight,
} from 'lucide-react'

interface StatCard {
  title: string
  value: string | number
  icon: React.ReactNode
  change?: string
  changeType?: 'positive' | 'negative' | 'neutral'
}

function StatCardComponent({ title, value, icon, change, changeType = 'neutral' }: StatCard) {
  const changeColor =
    changeType === 'positive'
      ? 'text-green-600'
      : changeType === 'negative'
        ? 'text-red-600'
        : 'text-gray-500'

  return (
    <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <div className="flex items-center justify-between">
        <div className="rounded-lg bg-gray-100 p-3 text-gray-600">{icon}</div>
        {change && (
          <span className={`flex items-center gap-1 text-sm font-medium ${changeColor}`}>
            {change}
            <ArrowUpRight className="h-3 w-3" />
          </span>
        )}
      </div>
      <div className="mt-4">
        <p className="text-sm font-medium text-gray-500">{title}</p>
        <p className="mt-1 text-3xl font-bold text-gray-900">{value}</p>
      </div>
    </div>
  )
}

const stats: StatCard[] = [
  {
    title: 'Total Products',
    value: 128,
    icon: <Package className="h-5 w-5" />,
    change: '+12%',
    changeType: 'positive',
  },
  {
    title: 'Total Orders',
    value: 342,
    icon: <ShoppingCart className="h-5 w-5" />,
    change: '+8%',
    changeType: 'positive',
  },
  {
    title: 'Pending Pages',
    value: 7,
    icon: <FileText className="h-5 w-5" />,
    change: '3 new',
    changeType: 'neutral',
  },
  {
    title: 'Active Promos',
    value: 4,
    icon: <Tag className="h-5 w-5" />,
    change: '-1',
    changeType: 'negative',
  },
]

const secondaryStats: StatCard[] = [
  {
    title: 'Revenue (MTD)',
    value: '$12,480',
    icon: <DollarSign className="h-5 w-5" />,
    change: '+23%',
    changeType: 'positive',
  },
  {
    title: 'Customers',
    value: 1204,
    icon: <Users className="h-5 w-5" />,
    change: '+5%',
    changeType: 'positive',
  },
  {
    title: 'Conversion Rate',
    value: '3.2%',
    icon: <TrendingUp className="h-5 w-5" />,
    change: '+0.4%',
    changeType: 'positive',
  },
]

interface RecentOrder {
  id: string
  customer: string
  total: string
  status: string
  date: string
}

const recentOrders: RecentOrder[] = [
  { id: 'ORD-001', customer: 'Alice Johnson', total: '$89.00', status: 'delivered', date: '2 hours ago' },
  { id: 'ORD-002', customer: 'Bob Smith', total: '$245.50', status: 'shipped', date: '4 hours ago' },
  { id: 'ORD-003', customer: 'Carol White', total: '$32.00', status: 'pending', date: '6 hours ago' },
  { id: 'ORD-004', customer: 'Dan Brown', total: '$178.00', status: 'confirmed', date: '8 hours ago' },
  { id: 'ORD-005', customer: 'Eve Davis', total: '$56.99', status: 'processing', date: '12 hours ago' },
]

const statusColors: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  confirmed: 'bg-blue-100 text-blue-800',
  processing: 'bg-indigo-100 text-indigo-800',
  shipped: 'bg-purple-100 text-purple-800',
  delivered: 'bg-green-100 text-green-800',
  cancelled: 'bg-red-100 text-red-800',
}

export default function Dashboard() {
  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-1 text-sm text-gray-500">
          Welcome back. Here's what's happening with your store today.
        </p>
      </div>

      {/* Primary Stats */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <StatCardComponent key={stat.title} {...stat} />
        ))}
      </div>

      {/* Secondary Stats */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-3">
        {secondaryStats.map((stat) => (
          <StatCardComponent key={stat.title} {...stat} />
        ))}
      </div>

      {/* Recent Orders */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        <div className="border-b border-gray-200 px-6 py-4">
          <h2 className="text-lg font-semibold text-gray-900">Recent Orders</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                <th className="px-6 py-3 font-medium">Order ID</th>
                <th className="px-6 py-3 font-medium">Customer</th>
                <th className="px-6 py-3 font-medium">Total</th>
                <th className="px-6 py-3 font-medium">Status</th>
                <th className="px-6 py-3 font-medium">Date</th>
              </tr>
            </thead>
            <tbody>
              {recentOrders.map((order) => (
                <tr key={order.id} className="border-b border-gray-50 last:border-0">
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{order.id}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{order.customer}</td>
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{order.total}</td>
                  <td className="px-6 py-4">
                    <span
                      className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColors[order.status] ?? 'bg-gray-100 text-gray-800'}`}
                    >
                      {order.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-500">{order.date}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
