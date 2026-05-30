import { useState } from 'react'
import { RotateCcw, ChevronDown, ChevronUp, X } from 'lucide-react'
import type { ReturnStatus, ReturnRequest } from '../../types'
import {
  useReturns,
  useReturn,
  useApproveReturn,
  useRejectReturn,
  useMarkReturnReceived,
  useMarkReturnRefunded,
  useCloseReturn,
} from '../../lib/hooks'

const statusColors: Record<ReturnStatus, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  approved: 'bg-blue-100 text-blue-800',
  received: 'bg-purple-100 text-purple-800',
  refunded: 'bg-green-100 text-green-800',
  closed: 'bg-gray-100 text-gray-800',
  rejected: 'bg-red-100 text-red-800',
}

function ReturnDetailPanel({ id, onClose }: { id: string; onClose: () => void }) {
  const { data: returnReq, isLoading } = useReturn(id)
  const approveReturn = useApproveReturn()
  const rejectReturn = useRejectReturn()
  const markReceived = useMarkReturnReceived()
  const markRefunded = useMarkReturnRefunded()
  const closeReturn = useCloseReturn()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
      </div>
    )
  }
  if (!returnReq) return null

  return (
    <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-gray-900">Return #{returnReq.id.slice(0, 8)}</h2>
        <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
          <X className="h-5 w-5" />
        </button>
      </div>
      <div className="grid grid-cols-2 gap-4 mb-4 text-sm">
        <div>
          <span className="text-gray-500">Order ID:</span>
          <span className="ml-2 font-medium text-gray-900">{returnReq.order_id}</span>
        </div>
        <div>
          <span className="text-gray-500">Status:</span>
          <span className={`ml-2 inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColors[returnReq.status]}`}>
            {returnReq.status}
          </span>
        </div>
        <div>
          <span className="text-gray-500">Customer ID:</span>
          <span className="ml-2 font-medium text-gray-900">{returnReq.customer_id}</span>
        </div>
        {returnReq.refund_amount && (
          <div>
            <span className="text-gray-500">Refund Amount:</span>
            <span className="ml-2 font-medium text-green-700">
              {(returnReq.refund_amount.amount / 100).toFixed(2)} {returnReq.refund_amount.currency}
            </span>
          </div>
        )}
      </div>
      {returnReq.notes && (
        <p className="mb-4 text-sm text-gray-600 italic">{returnReq.notes}</p>
      )}
      {/* Items */}
      <div className="mb-4">
        <h3 className="text-sm font-semibold text-gray-700 mb-2">Items</h3>
        <div className="space-y-2">
          {(returnReq.items ?? []).map((item) => (
            <div key={item.id} className="flex items-center justify-between rounded-lg border border-gray-100 px-4 py-2 text-sm">
              <span className="font-medium text-gray-900">{item.product_name}</span>
              <span className="text-gray-500">Qty: {item.quantity}</span>
              <span className="text-gray-400 italic">{item.reason}</span>
            </div>
          ))}
        </div>
      </div>
      {/* Action Buttons */}
      <div className="flex flex-wrap gap-2">
        {returnReq.status === 'pending' && (
          <>
            <button
              onClick={() => approveReturn.mutate(returnReq.id)}
              className="rounded-lg border border-green-200 px-3 py-1.5 text-sm font-medium text-green-700 hover:bg-green-50 transition-colors"
            >
              Approve
            </button>
            <button
              onClick={() => rejectReturn.mutate(returnReq.id)}
              className="rounded-lg border border-red-200 px-3 py-1.5 text-sm font-medium text-red-700 hover:bg-red-50 transition-colors"
            >
              Reject
            </button>
          </>
        )}
        {returnReq.status === 'approved' && (
          <button
            onClick={() => markReceived.mutate(returnReq.id)}
            className="rounded-lg border border-purple-200 px-3 py-1.5 text-sm font-medium text-purple-700 hover:bg-purple-50 transition-colors"
          >
            Mark Received
          </button>
        )}
        {returnReq.status === 'received' && (
          <button
            onClick={() => markRefunded.mutate(returnReq.id)}
            className="rounded-lg border border-green-200 px-3 py-1.5 text-sm font-medium text-green-700 hover:bg-green-50 transition-colors"
          >
            Mark Refunded
          </button>
        )}
        {(returnReq.status === 'refunded' || returnReq.status === 'rejected') && (
          <button
            onClick={() => closeReturn.mutate(returnReq.id)}
            className="rounded-lg border border-gray-200 px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
          >
            Close
          </button>
        )}
      </div>
    </div>
  )
}

export default function Returns() {
  const [statusFilter, setStatusFilter] = useState<ReturnStatus | ''>('')
  const [expandedId, setExpandedId] = useState<string | null>(null)

  const { data, isLoading } = useReturns(statusFilter ? { status: statusFilter } : undefined)
  const returns: ReturnRequest[] = data?.items ?? []

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Returns & Exchanges</h1>
        <p className="mt-1 text-sm text-gray-500">Manage return requests and RMA workflow</p>
      </div>

      {/* Filter */}
      <div className="flex items-center gap-3">
        <label className="text-sm font-medium text-gray-700">Status:</label>
        <select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value as ReturnStatus | '')}
          className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
        >
          <option value="">All</option>
          <option value="pending">Pending</option>
          <option value="approved">Approved</option>
          <option value="received">Received</option>
          <option value="refunded">Refunded</option>
          <option value="closed">Closed</option>
          <option value="rejected">Rejected</option>
        </select>
      </div>

      {/* Expanded Detail */}
      {expandedId && (
        <ReturnDetailPanel id={expandedId} onClose={() => setExpandedId(null)} />
      )}

      {/* Returns Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : returns.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <RotateCcw className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No return requests found</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">ID</th>
                  <th className="px-6 py-3 font-medium">Order</th>
                  <th className="px-6 py-3 font-medium">Customer</th>
                  <th className="px-6 py-3 font-medium">Items</th>
                  <th className="px-6 py-3 font-medium">Status</th>
                  <th className="px-6 py-3 font-medium">Date</th>
                  <th className="px-6 py-3 font-medium text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {returns.map((r) => (
                  <tr key={r.id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-4 text-sm font-mono text-gray-600">{r.id.slice(0, 8)}</td>
                    <td className="px-6 py-4 text-sm text-gray-600">{r.order_id.slice(0, 8)}</td>
                    <td className="px-6 py-4 text-sm text-gray-600">{r.customer_id.slice(0, 8)}</td>
                    <td className="px-6 py-4 text-sm text-gray-600">{(r.items ?? []).length} items</td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColors[r.status]}`}>
                        {r.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {new Date(r.created_at).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 text-right">
                      <button
                        onClick={() => setExpandedId(expandedId === r.id ? null : r.id)}
                        className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
                        title="View details"
                      >
                        {expandedId === r.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                      </button>
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
