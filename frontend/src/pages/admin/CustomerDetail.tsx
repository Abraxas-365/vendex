import { useState } from 'react'
import { useParams, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, Users, MapPin, ShoppingCart, Trash2 } from 'lucide-react'
import type { OrderStatus } from '../../types'
import { useCustomer, useDeleteCustomer, useCustomerOrders } from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Route: /admin/customers/:id
// ---------------------------------------------------------------------------

const orderStatusColors: Record<OrderStatus, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  confirmed: 'bg-blue-100 text-blue-800',
  processing: 'bg-indigo-100 text-indigo-800',
  shipped: 'bg-purple-100 text-purple-800',
  delivered: 'bg-green-100 text-green-800',
  cancelled: 'bg-red-100 text-red-800',
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

export default function CustomerDetail() {
  const { id } = useParams({ from: '/_admin/admin/customers/$id' })
  const navigate = useNavigate()
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)

  const { data: customer, isLoading, error } = useCustomer(id)
  const deleteCustomer = useDeleteCustomer()
  const { data: ordersData } = useCustomerOrders(id)

  const orders = ordersData?.items ?? []

  function handleBack() {
    void navigate({ to: '/admin/customers' })
  }

  function handleDelete() {
    deleteCustomer.mutate(id, {
      onSuccess: () => {
        void navigate({ to: '/admin/customers' })
      },
    })
  }

  // ── Loading state ──
  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-24">
        <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
      </div>
    )
  }

  // ── Error / not found state ──
  if (error || !customer) {
    return (
      <div className="space-y-4">
        <button
          onClick={handleBack}
          className="inline-flex items-center gap-2 text-sm text-gray-500 hover:text-gray-700 transition-colors"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to Customers
        </button>
        <div className="flex flex-col items-center gap-3 rounded-xl border border-red-200 bg-red-50 p-12 text-center">
          <Users className="h-8 w-8 text-red-400" />
          <p className="text-sm font-medium text-red-700">Customer not found</p>
          <p className="text-xs text-red-500">
            {error instanceof Error ? error.message : 'This customer could not be loaded.'}
          </p>
        </div>
      </div>
    )
  }

  const defaultAddress = customer.addresses.find((a) => a.is_default) ?? customer.addresses[0]

  return (
    <div className="space-y-6">
      {/* Back button */}
      <button
        onClick={handleBack}
        className="inline-flex items-center gap-2 text-sm text-gray-500 hover:text-gray-700 transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Customers
      </button>

      {/* Customer Info Card */}
      <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-indigo-100 text-indigo-600">
              <Users className="h-6 w-6" />
            </div>
            <div>
              <h1 className="text-xl font-bold text-gray-900">{customer.name}</h1>
              <p className="text-sm text-gray-500">Member since {formatDate(customer.created_at)}</p>
            </div>
          </div>
          <button
            onClick={() => setShowDeleteConfirm(true)}
            className="inline-flex items-center gap-2 rounded-lg border border-red-200 px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50 transition-colors"
          >
            <Trash2 className="h-4 w-4" />
            Delete Customer
          </button>
        </div>

        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div>
            <p className="text-xs font-semibold uppercase tracking-wider text-gray-400">Email</p>
            <p className="mt-1 text-sm text-gray-900">{customer.email}</p>
          </div>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wider text-gray-400">Phone</p>
            <p className="mt-1 text-sm text-gray-900">{customer.phone || <span className="text-gray-300">—</span>}</p>
          </div>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wider text-gray-400">Default Address</p>
            <p className="mt-1 text-sm text-gray-900">
              {defaultAddress
                ? `${defaultAddress.city}, ${defaultAddress.country}`
                : <span className="text-gray-300">None</span>}
            </p>
          </div>
        </div>
      </div>

      {/* Addresses Section */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        <div className="flex items-center gap-2 border-b border-gray-100 px-6 py-4">
          <MapPin className="h-4 w-4 text-gray-400" />
          <h2 className="text-sm font-semibold text-gray-900">Addresses</h2>
          <span className="ml-auto text-xs text-gray-400">{customer.addresses.length}</span>
        </div>

        {customer.addresses.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-10 text-gray-400">
            <MapPin className="mb-2 h-8 w-8" />
            <p className="text-sm">No addresses on file</p>
          </div>
        ) : (
          <ul className="divide-y divide-gray-50">
            {customer.addresses.map((address, idx) => (
              <li key={idx} className="px-6 py-4">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="text-sm text-gray-900">
                      {address.street}
                    </p>
                    <p className="mt-0.5 text-sm text-gray-500">
                      {address.city}, {address.state} {address.postal_code}
                    </p>
                    <p className="mt-0.5 text-sm text-gray-500">{address.country}</p>
                  </div>
                  {address.is_default && (
                    <span className="inline-flex rounded-full bg-indigo-50 px-2.5 py-0.5 text-xs font-medium text-indigo-700">
                      Default
                    </span>
                  )}
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>

      {/* Orders Section */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        <div className="flex items-center gap-2 border-b border-gray-100 px-6 py-4">
          <ShoppingCart className="h-4 w-4 text-gray-400" />
          <h2 className="text-sm font-semibold text-gray-900">Orders</h2>
          <span className="ml-auto text-xs text-gray-400">{orders.length}</span>
        </div>

        {orders.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-10 text-gray-400">
            <ShoppingCart className="mb-2 h-8 w-8" />
            <p className="text-sm">No orders yet</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-xs text-gray-500">
                  <th className="px-6 py-3 font-medium">Order ID</th>
                  <th className="px-6 py-3 font-medium">Items</th>
                  <th className="px-6 py-3 font-medium">Total</th>
                  <th className="px-6 py-3 font-medium">Status</th>
                  <th className="px-6 py-3 font-medium">Date</th>
                </tr>
              </thead>
              <tbody>
                {orders.map((order) => (
                  <tr key={order.id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-3 text-sm font-mono font-medium text-gray-900">
                      {order.id.slice(0, 8)}...
                    </td>
                    <td className="px-6 py-3 text-sm text-gray-600">
                      {order.items.length} item{order.items.length !== 1 && 's'}
                    </td>
                    <td className="px-6 py-3 text-sm font-medium text-gray-900">
                      ${order.total_amount.amount.toFixed(2)}
                    </td>
                    <td className="px-6 py-3">
                      <span
                        className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${orderStatusColors[order.status]}`}
                      >
                        {order.status}
                      </span>
                    </td>
                    <td className="px-6 py-3 text-sm text-gray-500">
                      {formatDate(order.created_at)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Delete Confirmation Dialog */}
      {showDeleteConfirm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
          <div className="w-full max-w-sm rounded-xl border border-gray-200 bg-white p-6 shadow-xl">
            <h3 className="text-base font-semibold text-gray-900">Delete Customer</h3>
            <p className="mt-2 text-sm text-gray-500">
              Are you sure you want to delete <strong>{customer.name}</strong>? This action cannot be undone.
            </p>
            <div className="mt-5 flex justify-end gap-3">
              <button
                onClick={() => setShowDeleteConfirm(false)}
                className="rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleDelete}
                disabled={deleteCustomer.isPending}
                className="inline-flex items-center gap-2 rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 disabled:opacity-50 transition-colors"
              >
                {deleteCustomer.isPending ? (
                  <div className="h-3.5 w-3.5 animate-spin rounded-full border border-white border-t-transparent" />
                ) : (
                  <Trash2 className="h-3.5 w-3.5" />
                )}
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
