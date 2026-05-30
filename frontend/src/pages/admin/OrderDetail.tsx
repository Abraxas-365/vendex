import { useParams, Link } from '@tanstack/react-router'
import {
  ArrowLeft,
  Loader2,
  AlertCircle,
  ShoppingCart,
  MapPin,
  User,
  CheckCircle2,
  XCircle,
  CreditCard,
  Truck,
} from 'lucide-react'
import type { OrderStatus, PaymentStatus } from '../../types'
import { useOrder, useUpdateOrderStatus, useCancelOrder, useOrderPayments } from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Route: /admin/orders/$id
// ---------------------------------------------------------------------------

const statusColors: Record<OrderStatus, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  confirmed: 'bg-blue-100 text-blue-800',
  processing: 'bg-indigo-100 text-indigo-800',
  shipped: 'bg-purple-100 text-purple-800',
  delivered: 'bg-green-100 text-green-800',
  cancelled: 'bg-red-100 text-red-800',
}

const paymentStatusColors: Record<PaymentStatus, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  processing: 'bg-blue-100 text-blue-800',
  completed: 'bg-green-100 text-green-800',
  failed: 'bg-red-100 text-red-800',
  refunded: 'bg-gray-100 text-gray-600',
}

/** Status transitions: key → allowed next statuses */
const STATUS_TRANSITIONS: Partial<Record<OrderStatus, OrderStatus[]>> = {
  pending: ['confirmed', 'cancelled'],
  confirmed: ['processing', 'cancelled'],
  processing: ['shipped', 'cancelled'],
  shipped: ['delivered'],
}

/** Human-readable labels for transition buttons */
const TRANSITION_LABELS: Record<OrderStatus, string> = {
  confirmed: 'Confirm Order',
  processing: 'Start Processing',
  shipped: 'Mark Shipped',
  delivered: 'Mark Delivered',
  cancelled: 'Cancel Order',
  pending: 'Reset to Pending',
}

function formatMoney(amount: number, currency: string): string {
  // Backend sends cents as int64 — divide by 100 to get dollars
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency || 'USD',
  }).format(amount / 100)
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export default function OrderDetail() {
  const { id } = useParams({ from: '/_admin/admin/orders/$id' })

  const { data: order, isLoading, error } = useOrder(id)
  const updateStatus = useUpdateOrderStatus()
  const cancelOrder = useCancelOrder()
  const { data: payments } = useOrderPayments(id)

  const isPending = updateStatus.isPending || cancelOrder.isPending
  const latestPayment = payments?.[0]

  function handleTransition(nextStatus: OrderStatus) {
    if (!order) return
    if (nextStatus === 'cancelled') {
      cancelOrder.mutate(order.id)
    } else {
      updateStatus.mutate({ id: order.id, status: nextStatus })
    }
  }

  // ── Loading ──
  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center py-24">
        <div className="flex flex-col items-center gap-3 text-slate-500">
          <Loader2 size={32} className="animate-spin" />
          <p className="text-sm">Loading order…</p>
        </div>
      </div>
    )
  }

  // ── Error ──
  if (error || !order) {
    return (
      <div className="space-y-4">
        <Link
          to="/admin/orders"
          className="inline-flex items-center gap-2 text-sm text-slate-500 hover:text-slate-700 transition-colors"
        >
          <ArrowLeft size={16} />
          Back to Orders
        </Link>
        <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-center text-red-600">
          <AlertCircle size={32} className="mx-auto mb-3" />
          <p className="font-medium">Order not found</p>
          <p className="mt-1 text-sm text-red-500">{error?.message ?? 'The order could not be loaded.'}</p>
        </div>
      </div>
    )
  }

  const allowedTransitions = STATUS_TRANSITIONS[order.status] ?? []
  const isTerminal = order.status === 'delivered' || order.status === 'cancelled'

  return (
    <div className="space-y-6">
      {/* ── Header ── */}
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-center gap-3">
          <Link
            to="/admin/orders"
            className="inline-flex items-center gap-1.5 rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-sm font-medium text-slate-600 shadow-sm transition-colors hover:bg-slate-50 hover:text-slate-900"
          >
            <ArrowLeft size={15} />
            Orders
          </Link>
          <div>
            <h1 className="text-xl font-bold text-gray-900 font-mono">
              Order #{order.id.slice(0, 8)}
            </h1>
            <p className="mt-0.5 text-xs text-gray-400">
              {formatDate(order.created_at)}
            </p>
          </div>
        </div>

        {/* Status badge */}
        <span
          className={`inline-flex items-center rounded-full px-3 py-1 text-sm font-medium capitalize ${statusColors[order.status]}`}
        >
          {order.status}
        </span>
      </div>

      {/* ── Status actions bar ── */}
      {!isTerminal && allowedTransitions.length > 0 && (
        <div className="flex flex-wrap items-center gap-3 rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
          <p className="text-sm font-medium text-slate-600 mr-2">Actions:</p>
          {allowedTransitions.map((nextStatus) => (
            <button
              key={nextStatus}
              disabled={isPending}
              onClick={() => handleTransition(nextStatus)}
              className={`inline-flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-colors disabled:opacity-50 ${
                nextStatus === 'cancelled'
                  ? 'border border-red-200 bg-white text-red-600 hover:bg-red-50'
                  : 'bg-indigo-600 text-white hover:bg-indigo-700 shadow-sm'
              }`}
            >
              {isPending ? (
                <Loader2 size={14} className="animate-spin" />
              ) : nextStatus === 'cancelled' ? (
                <XCircle size={14} />
              ) : (
                <CheckCircle2 size={14} />
              )}
              {TRANSITION_LABELS[nextStatus]}
            </button>
          ))}
        </div>
      )}

      {/* Terminal state notice */}
      {isTerminal && (
        <div
          className={`flex items-center gap-3 rounded-xl border px-4 py-3 text-sm font-medium ${
            order.status === 'delivered'
              ? 'border-green-200 bg-green-50 text-green-700'
              : 'border-red-200 bg-red-50 text-red-700'
          }`}
        >
          {order.status === 'delivered' ? (
            <CheckCircle2 size={16} />
          ) : (
            <XCircle size={16} />
          )}
          {order.status === 'delivered'
            ? 'This order has been delivered and is complete.'
            : 'This order has been cancelled.'}
        </div>
      )}

      {/* ── Mutation error ── */}
      {(updateStatus.error || cancelOrder.error) && (
        <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-600">
          <AlertCircle size={14} className="mr-1.5 inline-block" />
          {updateStatus.error?.message ?? cancelOrder.error?.message}
        </div>
      )}

      <div className="grid gap-6 lg:grid-cols-3">
        {/* ── Left column: items + summary ── */}
        <div className="space-y-6 lg:col-span-2">
          {/* Order items table */}
          <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
            <div className="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
              <ShoppingCart size={16} className="text-slate-500" />
              <h2 className="text-sm font-semibold text-gray-900">Order Items</h2>
              <span className="ml-auto text-xs text-gray-400">
                {order.items.length} item{order.items.length !== 1 && 's'}
              </span>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-100 text-left text-xs font-medium text-gray-500">
                    <th className="px-5 py-3">Product</th>
                    <th className="px-5 py-3 text-right">Qty</th>
                    <th className="px-5 py-3 text-right">Unit Price</th>
                    <th className="px-5 py-3 text-right">Total</th>
                  </tr>
                </thead>
                <tbody>
                  {order.items.map((item) => (
                    <tr key={item.id} className="border-b border-gray-50 last:border-0">
                      <td className="px-5 py-3.5 text-sm text-gray-900">{item.product_name}</td>
                      <td className="px-5 py-3.5 text-right text-sm text-gray-600">{item.quantity}</td>
                      <td className="px-5 py-3.5 text-right text-sm text-gray-600">
                        {formatMoney(item.unit_price.amount, item.unit_price.currency)}
                      </td>
                      <td className="px-5 py-3.5 text-right text-sm font-medium text-gray-900">
                        {formatMoney(item.total.amount, item.total.currency)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            {/* Summary */}
            <div className="border-t border-gray-100 px-5 py-4 space-y-2">
              {order.subtotal && (
                <div className="flex items-center justify-between text-sm text-gray-600">
                  <span>Subtotal</span>
                  <span>{formatMoney(order.subtotal.amount, order.subtotal.currency)}</span>
                </div>
              )}
              {order.shipping_amount && (
                <div className="flex items-center justify-between text-sm text-gray-600">
                  <span>Shipping</span>
                  <span>{formatMoney(order.shipping_amount.amount, order.shipping_amount.currency)}</span>
                </div>
              )}
              {order.tax_amount && (
                <div className="flex items-center justify-between text-sm text-gray-600">
                  <span>Tax</span>
                  <span>{formatMoney(order.tax_amount.amount, order.tax_amount.currency)}</span>
                </div>
              )}
              {order.discount_amount && order.discount_amount.amount > 0 && (
                <div className="flex items-center justify-between text-sm text-green-600">
                  <span>Discount</span>
                  <span>−{formatMoney(order.discount_amount.amount, order.discount_amount.currency)}</span>
                </div>
              )}
              <div className="flex items-center justify-between border-t border-gray-100 pt-2 text-sm font-bold text-gray-900">
                <span>Total</span>
                <span>{formatMoney(order.total_amount.amount, order.total_amount.currency)}</span>
              </div>
            </div>
          </div>
        </div>

        {/* ── Right column: customer + shipping ── */}
        <div className="space-y-4">
          {/* Customer info */}
          <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
            <div className="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
              <User size={16} className="text-slate-500" />
              <h2 className="text-sm font-semibold text-gray-900">Customer</h2>
            </div>
            <div className="px-5 py-4">
              <p className="text-xs font-medium uppercase tracking-wider text-gray-400">Customer ID</p>
              <p className="mt-1 font-mono text-sm text-gray-700 break-all">{order.customer_id}</p>
              <p className="mt-3 text-xs font-medium uppercase tracking-wider text-gray-400">Order ID</p>
              <p className="mt-1 font-mono text-sm text-gray-700 break-all">{order.id}</p>
              <p className="mt-3 text-xs font-medium uppercase tracking-wider text-gray-400">Last Updated</p>
              <p className="mt-1 text-sm text-gray-600">{formatDate(order.updated_at)}</p>
            </div>
          </div>

          {/* Shipping address */}
          {order.shipping_address && (
            <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
              <div className="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
                <MapPin size={16} className="text-slate-500" />
                <h2 className="text-sm font-semibold text-gray-900">Shipping Address</h2>
              </div>
              <div className="px-5 py-4 text-sm text-gray-600 leading-relaxed">
                <p>{order.shipping_address.street}</p>
                <p>
                  {order.shipping_address.city}, {order.shipping_address.state}{' '}
                  {order.shipping_address.postal_code}
                </p>
                <p className="font-medium text-gray-800">{order.shipping_address.country}</p>
              </div>
            </div>
          )}

          {/* Shipping method & tracking */}
          {(order.shipping_method || order.tracking_number) && (
            <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
              <div className="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
                <Truck size={16} className="text-slate-500" />
                <h2 className="text-sm font-semibold text-gray-900">Shipping Info</h2>
              </div>
              <div className="px-5 py-4 space-y-2">
                {order.shipping_method && (
                  <div>
                    <p className="text-xs font-medium uppercase tracking-wider text-gray-400">Method</p>
                    <p className="mt-1 text-sm text-gray-700">{order.shipping_method}</p>
                  </div>
                )}
                {order.carrier && (
                  <div>
                    <p className="text-xs font-medium uppercase tracking-wider text-gray-400">Carrier</p>
                    <p className="mt-1 text-sm text-gray-700">{order.carrier}</p>
                  </div>
                )}
                {order.tracking_number && (
                  <div>
                    <p className="text-xs font-medium uppercase tracking-wider text-gray-400">Tracking</p>
                    <p className="mt-1 font-mono text-sm text-gray-700 break-all">{order.tracking_number}</p>
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Payment */}
          {latestPayment && (
            <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
              <div className="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
                <CreditCard size={16} className="text-slate-500" />
                <h2 className="text-sm font-semibold text-gray-900">Payment</h2>
              </div>
              <div className="px-5 py-4 space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium uppercase tracking-wider text-gray-400">Status</span>
                  <span
                    className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${paymentStatusColors[latestPayment.status]}`}
                  >
                    {latestPayment.status}
                  </span>
                </div>
                {latestPayment.provider && (
                  <div className="flex items-center justify-between">
                    <span className="text-xs font-medium uppercase tracking-wider text-gray-400">Provider</span>
                    <span className="text-sm text-gray-700 capitalize">{latestPayment.provider}</span>
                  </div>
                )}
                {latestPayment.method && (
                  <div className="flex items-center justify-between">
                    <span className="text-xs font-medium uppercase tracking-wider text-gray-400">Method</span>
                    <span className="text-sm text-gray-700 capitalize">{latestPayment.method}</span>
                  </div>
                )}
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium uppercase tracking-wider text-gray-400">Amount</span>
                  <span className="text-sm font-medium text-gray-900">
                    {formatMoney(latestPayment.amount.amount, latestPayment.amount.currency)}
                  </span>
                </div>
                <div className="pt-2">
                  <Link
                    to="/admin/payments"
                    className="text-xs text-indigo-600 hover:text-indigo-800 hover:underline"
                  >
                    View all payments →
                  </Link>
                </div>
              </div>
            </div>
          )}

          {/* Payment status from order (if no payment fetched) */}
          {!latestPayment && order.payment_status && (
            <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
              <div className="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
                <CreditCard size={16} className="text-slate-500" />
                <h2 className="text-sm font-semibold text-gray-900">Payment</h2>
              </div>
              <div className="px-5 py-4">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium uppercase tracking-wider text-gray-400">Status</span>
                  <span
                    className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${paymentStatusColors[order.payment_status]}`}
                  >
                    {order.payment_status}
                  </span>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
