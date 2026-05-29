import { useState } from 'react'
import { Search, Users } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import type { Customer } from '../../types'
import { useCustomers } from '../../lib/hooks'

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

export default function Customers() {
  const [searchQuery, setSearchQuery] = useState('')

  const { data, isLoading } = useCustomers()
  const customers: Customer[] = data?.items ?? []

  const filtered = customers.filter(
    (c) =>
      c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      c.email.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Customers</h1>
        <p className="mt-1 text-sm text-gray-500">
          View and manage your customer accounts
        </p>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          placeholder="Search by name or email..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full rounded-lg border border-gray-200 bg-white py-2.5 pl-10 pr-4 text-sm text-gray-900 placeholder:text-gray-400 focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
        />
      </div>

      {/* Customers Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Users className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No customers found</p>
            <p className="mt-1 text-xs">Customers will appear here when they register.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">Name</th>
                  <th className="px-6 py-3 font-medium">Email</th>
                  <th className="px-6 py-3 font-medium">Phone</th>
                  <th className="px-6 py-3 font-medium">Member Since</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((customer) => (
                  <tr
                    key={customer.id}
                    className="border-b border-gray-50 last:border-0 cursor-pointer hover:bg-gray-50/50 transition-colors"
                  >
                    <td className="px-6 py-4">
                      <Link
                        to="/admin/customers/$id"
                        params={{ id: customer.id }}
                        className="block"
                      >
                        <div className="text-sm font-medium text-gray-900">{customer.name}</div>
                        <div className="mt-0.5 text-xs text-gray-400">
                          {customer.addresses.length} address{customer.addresses.length !== 1 && 'es'}
                        </div>
                      </Link>
                    </td>
                    <td className="px-6 py-4">
                      <Link
                        to="/admin/customers/$id"
                        params={{ id: customer.id }}
                        className="block text-sm text-gray-600"
                      >
                        {customer.email}
                      </Link>
                    </td>
                    <td className="px-6 py-4">
                      <Link
                        to="/admin/customers/$id"
                        params={{ id: customer.id }}
                        className="block text-sm text-gray-600"
                      >
                        {customer.phone || <span className="text-gray-300">—</span>}
                      </Link>
                    </td>
                    <td className="px-6 py-4">
                      <Link
                        to="/admin/customers/$id"
                        params={{ id: customer.id }}
                        className="block text-sm text-gray-500"
                      >
                        {formatDate(customer.created_at)}
                      </Link>
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
