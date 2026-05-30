import { useState } from 'react'
import { CreditCard, X, ChevronDown, ChevronUp } from 'lucide-react'
import type { Payment, PaymentStatus, Refund } from '../../types'
import { useOrderPayments, useCreateRefund, useRefunds } from '../../lib/hooks'
import { useQuery } from '@tanstack/react-query'
import { listOrders } from '../../lib/api'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function formatMoney(amount: number, currency: string): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount / 100)
}

function formatDate(dateStr?: string): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function truncateId(id: string): string {
  return id.slice(0, 8) + '…'
}

const statusColors: Record<PaymentStatus, string> = {
  pending: 'bg-yellow-50 text-yellow-700',
  processing: 'bg-blue-50 text-blue-700',
  completed: 'bg-green-50 text-green-700',
  failed: 'bg-red-50 text-red-700',
  refunded: 'bg-slate-100 text-slate-500',
}

// ---------------------------------------------------------------------------
// Refund Dialog
// ---------------------------------------------------------------------------

interface RefundDialogProps {
  payment: Payment
  onClose: () => void
}

function RefundDialog({ payment, onClose }: RefundDialogProps) {
  const [amount, setAmount] = useState('')
  const [reason, setReason] = useState('')
  const createRefund = useCreateRefund()

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const amountCents = Math.round(parseFloat(amount) * 100)
    createRefund.mutate(
      { paymentId: payment.id, amount: amountCents, reason: reason || undefined },
      { onSuccess: () => onClose() },
    )
  }

  const maxAmount = (payment.amount.amount / 100).toFixed(2)

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div className="w-full max-w-md rounded-xl bg-white shadow-xl">
        <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <h2 className="text-base font-semibold text-slate-800">Issue Refund</h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
            <X size={18} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 px-6 py-4">
          <div className="rounded-lg bg-slate-50 p-3 text-sm">
            <p className="text-slate-500">Payment</p>
            <p className="font-medium text-slate-800">{payment.id}</p>
            <p className="text-slate-500 mt-1">Original amount</p>
            <p className="font-semibold text-slate-800">
              {formatMoney(payment.amount.amount, payment.amount.currency)}
            </p>
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">
              Refund Amount ($) <span className="font-normal text-slate-400">(max {maxAmount})</span>
            </label>
            <input
              type="number"
              required
              step="0.01"
              min="0.01"
              max={maxAmount}
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
              placeholder="0.00"
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">
              Reason <span className="font-normal text-slate-400">(optional)</span>
            </label>
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              rows={2}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none resize-none"
              placeholder="Customer requested refund…"
            />
          </div>

          {createRefund.error && (
            <p className="text-sm text-red-600">{createRefund.error.message}</p>
          )}

          <div className="flex gap-2 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={createRefund.isPending}
              className="flex-1 rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 disabled:opacity-60"
            >
              {createRefund.isPending ? 'Processing…' : 'Issue Refund'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Payment Row (expandable)
// ---------------------------------------------------------------------------

interface PaymentRowProps {
  payment: Payment
}

function PaymentRow({ payment }: PaymentRowProps) {
  const [expanded, setExpanded] = useState(false)
  const [showRefundDialog, setShowRefundDialog] = useState(false)

  const { data: refunds = [], isLoading: refundsLoading } = useRefunds(
    expanded ? payment.id : '',
  )

  return (
    <>
      <tr
        className="hover:bg-slate-50 cursor-pointer"
        onClick={() => setExpanded((v) => !v)}
      >
        <td className="px-5 py-3 font-mono text-xs text-slate-500">{truncateId(payment.id)}</td>
        <td className="px-5 py-3 font-mono text-xs text-slate-500">{truncateId(payment.order_id)}</td>
        <td className="px-5 py-3 font-medium text-slate-800">
          {formatMoney(payment.amount.amount, payment.amount.currency)}
        </td>
        <td className="px-5 py-3">
          <span
            className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColors[payment.status]}`}
          >
            {payment.status}
          </span>
        </td>
        <td className="px-5 py-3 text-slate-600">{payment.provider}</td>
        <td className="px-5 py-3 text-slate-500">{payment.method ?? '—'}</td>
        <td className="px-5 py-3 text-slate-400 text-xs">{formatDate(payment.created_at)}</td>
        <td className="px-5 py-3 text-right">
          {expanded ? (
            <ChevronUp size={16} className="text-slate-400 inline" />
          ) : (
            <ChevronDown size={16} className="text-slate-400 inline" />
          )}
        </td>
      </tr>

      {expanded && (
        <tr className="bg-slate-50">
          <td colSpan={8} className="px-5 py-4">
            <div className="grid grid-cols-2 gap-6">
              {/* Payment details */}
              <div className="space-y-2">
                <h4 className="text-xs font-semibold uppercase tracking-widest text-slate-400 mb-3">
                  Payment Details
                </h4>
                <dl className="space-y-1 text-sm">
                  <div className="flex justify-between">
                    <dt className="text-slate-500">Payment ID</dt>
                    <dd className="font-mono text-xs text-slate-700">{payment.id}</dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-slate-500">Order ID</dt>
                    <dd className="font-mono text-xs text-slate-700">{payment.order_id}</dd>
                  </div>
                  {payment.provider_payment_id && (
                    <div className="flex justify-between">
                      <dt className="text-slate-500">Provider Ref</dt>
                      <dd className="font-mono text-xs text-slate-700">{payment.provider_payment_id}</dd>
                    </div>
                  )}
                  {payment.paid_at && (
                    <div className="flex justify-between">
                      <dt className="text-slate-500">Paid At</dt>
                      <dd className="text-slate-700">{formatDate(payment.paid_at)}</dd>
                    </div>
                  )}
                  {payment.error_message && (
                    <div className="mt-2 rounded-lg bg-red-50 p-2 text-xs text-red-600">
                      {payment.error_message}
                    </div>
                  )}
                </dl>

                {payment.status === 'completed' && (
                  <button
                    onClick={(e) => {
                      e.stopPropagation()
                      setShowRefundDialog(true)
                    }}
                    className="mt-3 rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50"
                  >
                    Issue Refund
                  </button>
                )}
              </div>

              {/* Refunds */}
              <div>
                <h4 className="text-xs font-semibold uppercase tracking-widest text-slate-400 mb-3">
                  Refunds
                </h4>
                {refundsLoading && (
                  <p className="text-xs text-slate-400">Loading refunds…</p>
                )}
                {!refundsLoading && refunds.length === 0 && (
                  <p className="text-xs text-slate-400">No refunds issued.</p>
                )}
                {refunds.length > 0 && (
                  <div className="space-y-2">
                    {refunds.map((refund: Refund) => (
                      <div
                        key={refund.id}
                        className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm"
                      >
                        <div className="flex items-center justify-between">
                          <span className="font-medium text-slate-700">
                            {formatMoney(refund.amount.amount, refund.amount.currency)}
                          </span>
                          <span
                            className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                              refund.status === 'completed'
                                ? 'bg-green-50 text-green-700'
                                : refund.status === 'failed'
                                  ? 'bg-red-50 text-red-700'
                                  : 'bg-yellow-50 text-yellow-700'
                            }`}
                          >
                            {refund.status}
                          </span>
                        </div>
                        {refund.reason && (
                          <p className="mt-0.5 text-xs text-slate-500">{refund.reason}</p>
                        )}
                        <p className="mt-0.5 text-xs text-slate-400">
                          {formatDate(refund.created_at)}
                        </p>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </td>
        </tr>
      )}

      {showRefundDialog && (
        <tr>
          <td>
            <RefundDialog payment={payment} onClose={() => setShowRefundDialog(false)} />
          </td>
        </tr>
      )}
    </>
  )
}

// ---------------------------------------------------------------------------
// Main Payments page
// ---------------------------------------------------------------------------

export default function Payments() {
  // Fetch all orders and collect payments per order
  const { data: ordersData, isLoading: ordersLoading } = useQuery({
    queryKey: ['orders', 'list', { page: 1, page_size: 100 }],
    queryFn: () => listOrders({ page: 1, page_size: 100 }),
  })

  const orders = ordersData?.items ?? []

  // We'll display payments aggregated across recent orders by fetching each order's payments
  // Since there's no "list all payments" endpoint, we'll show a per-order view
  const [selectedOrderId, setSelectedOrderId] = useState<string | null>(null)

  const { data: payments = [], isLoading: paymentsLoading } = useOrderPayments(
    selectedOrderId ?? '',
  )

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-xl font-semibold text-slate-800">Payments</h1>
        <p className="mt-0.5 text-sm text-slate-500">
          View payment details and issue refunds for orders.
        </p>
      </div>

      <div className="flex gap-6">
        {/* Left: Order selector */}
        <div className="w-72 shrink-0">
          <div className="rounded-xl border border-slate-200 bg-white overflow-hidden">
            <div className="border-b border-slate-100 px-4 py-3">
              <p className="text-xs font-semibold uppercase tracking-widest text-slate-400">
                Orders
              </p>
            </div>

            {ordersLoading && (
              <div className="py-8 text-center text-sm text-slate-400">Loading orders…</div>
            )}

            {!ordersLoading && orders.length === 0 && (
              <div className="py-8 text-center text-sm text-slate-400">No orders found.</div>
            )}

            <ul className="divide-y divide-slate-100 max-h-[70vh] overflow-y-auto">
              {orders.map((order) => (
                <li key={order.id}>
                  <button
                    onClick={() => setSelectedOrderId(order.id)}
                    className={`flex w-full items-center gap-3 px-4 py-3 text-left transition-colors hover:bg-slate-50 ${
                      selectedOrderId === order.id
                        ? 'bg-indigo-50 text-indigo-700'
                        : 'text-slate-700'
                    }`}
                  >
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium font-mono">{truncateId(order.id)}</p>
                      <p className="text-xs text-slate-400">
                        {formatMoney(order.total_amount.amount, order.total_amount.currency)} ·{' '}
                        <span className="capitalize">{order.status}</span>
                      </p>
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Right: Payments for selected order */}
        <div className="flex-1">
          {selectedOrderId ? (
            <div className="rounded-xl border border-slate-200 bg-white overflow-hidden">
              <div className="border-b border-slate-100 px-5 py-4">
                <p className="text-xs font-semibold uppercase tracking-widest text-slate-400">
                  Payments for Order
                </p>
                <p className="mt-0.5 font-mono text-sm text-slate-700">{selectedOrderId}</p>
              </div>

              {paymentsLoading && (
                <div className="py-8 text-center text-sm text-slate-400">Loading payments…</div>
              )}

              {!paymentsLoading && payments.length === 0 && (
                <div className="py-16 text-center">
                  <CreditCard size={40} className="mx-auto mb-3 text-slate-300" />
                  <p className="text-sm text-slate-500">No payments found for this order</p>
                </div>
              )}

              {payments.length > 0 && (
                <table className="w-full text-sm">
                  <thead className="bg-slate-50 text-xs font-semibold uppercase tracking-wide text-slate-500">
                    <tr>
                      <th className="px-5 py-3 text-left">Payment ID</th>
                      <th className="px-5 py-3 text-left">Order ID</th>
                      <th className="px-5 py-3 text-left">Amount</th>
                      <th className="px-5 py-3 text-left">Status</th>
                      <th className="px-5 py-3 text-left">Provider</th>
                      <th className="px-5 py-3 text-left">Method</th>
                      <th className="px-5 py-3 text-left">Date</th>
                      <th className="px-5 py-3" />
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {payments.map((payment) => (
                      <PaymentRow key={payment.id} payment={payment} />
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          ) : (
            <div className="flex h-full min-h-48 items-center justify-center rounded-xl border border-dashed border-slate-200 bg-white text-center">
              <div>
                <CreditCard size={32} className="mx-auto mb-2 text-slate-300" />
                <p className="text-sm text-slate-400">Select an order to view its payments</p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
