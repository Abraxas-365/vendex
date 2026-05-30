import { useState } from 'react'
import { RefreshCw, X, ChevronRight, PauseCircle, PlayCircle, XCircle } from 'lucide-react'
import type { Subscription, BillingRecord } from '../../types'
import {
  useSubscriptions,
  useDueSubscriptions,
  useBillingRecords,
  useCancelSubscription,
  usePauseSubscription,
  useResumeSubscription,
} from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function formatMoney(amount: number, currency: string): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount / 100)
}

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

const statusColors: Record<Subscription['status'], string> = {
  active: 'bg-green-50 text-green-700',
  paused: 'bg-amber-50 text-amber-700',
  cancelled: 'bg-red-50 text-red-700',
  expired: 'bg-slate-100 text-slate-500',
}

const billingStatusColors: Record<BillingRecord['status'], string> = {
  success: 'bg-green-50 text-green-700',
  failed: 'bg-red-50 text-red-700',
  pending: 'bg-amber-50 text-amber-700',
}

// ---------------------------------------------------------------------------
// Billing Records Panel
// ---------------------------------------------------------------------------

interface BillingPanelProps {
  subscription: Subscription
  onClose: () => void
}

function BillingPanel({ subscription, onClose }: BillingPanelProps) {
  const { data, isPending } = useBillingRecords(subscription.id)
  const records = data?.items ?? []

  return (
    <div className="flex flex-col rounded-xl border border-slate-200 bg-white">
      <div className="flex items-center justify-between border-b border-slate-100 px-5 py-3">
        <div>
          <p className="text-sm font-semibold text-slate-800">Billing History</p>
          <p className="text-xs text-slate-400">Subscription {subscription.id.slice(0, 8)}…</p>
        </div>
        <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
          <X size={18} />
        </button>
      </div>
      <div className="flex-1 overflow-auto">
        {isPending ? (
          <p className="p-4 text-sm text-slate-500">Loading…</p>
        ) : records.length === 0 ? (
          <p className="p-4 text-sm text-slate-400">No billing records yet.</p>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-100">
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-slate-400">Amount</th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-slate-400">Status</th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-slate-400">Billed At</th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-slate-400">Note</th>
              </tr>
            </thead>
            <tbody>
              {records.map((record) => (
                <tr key={record.id} className="border-b border-slate-50 hover:bg-slate-50">
                  <td className="px-4 py-2 font-medium">{formatMoney(record.amount, record.currency)}</td>
                  <td className="px-4 py-2">
                    <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${billingStatusColors[record.status]}`}>
                      {record.status}
                    </span>
                  </td>
                  <td className="px-4 py-2 text-slate-400 text-xs">{formatDate(record.billed_at)}</td>
                  <td className="px-4 py-2 text-slate-500 text-xs">{record.failure_reason ?? '—'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Subscription Row Actions
// ---------------------------------------------------------------------------

interface RowActionsProps {
  subscription: Subscription
}

function RowActions({ subscription }: RowActionsProps) {
  const cancel = useCancelSubscription()
  const pause = usePauseSubscription()
  const resume = useResumeSubscription()

  function handleCancel(e: React.MouseEvent) {
    e.stopPropagation()
    if (!confirm('Cancel this subscription?')) return
    cancel.mutate(subscription.id)
  }

  function handlePause(e: React.MouseEvent) {
    e.stopPropagation()
    pause.mutate(subscription.id)
  }

  function handleResume(e: React.MouseEvent) {
    e.stopPropagation()
    resume.mutate(subscription.id)
  }

  return (
    <div className="flex items-center gap-1">
      {subscription.status === 'active' && (
        <button
          onClick={handlePause}
          disabled={pause.isPending}
          title="Pause"
          className="rounded p-1 text-slate-400 hover:bg-amber-50 hover:text-amber-600 disabled:opacity-40"
        >
          <PauseCircle size={14} />
        </button>
      )}
      {subscription.status === 'paused' && (
        <button
          onClick={handleResume}
          disabled={resume.isPending}
          title="Resume"
          className="rounded p-1 text-slate-400 hover:bg-green-50 hover:text-green-600 disabled:opacity-40"
        >
          <PlayCircle size={14} />
        </button>
      )}
      {(subscription.status === 'active' || subscription.status === 'paused') && (
        <button
          onClick={handleCancel}
          disabled={cancel.isPending}
          title="Cancel"
          className="rounded p-1 text-slate-400 hover:bg-red-50 hover:text-red-600 disabled:opacity-40"
        >
          <XCircle size={14} />
        </button>
      )}
      <ChevronRight size={14} className="text-slate-300 ml-0.5" />
    </div>
  )
}

// ---------------------------------------------------------------------------
// Subscription Table
// ---------------------------------------------------------------------------

interface SubscriptionTableProps {
  subscriptions: Subscription[]
  selectedId: string | null
  onSelect: (sub: Subscription) => void
}

function SubscriptionTable({ subscriptions, selectedId, onSelect }: SubscriptionTableProps) {
  if (subscriptions.length === 0) {
    return <p className="p-6 text-sm text-slate-400">No subscriptions found.</p>
  }

  return (
    <table className="w-full text-sm">
      <thead>
        <tr className="border-b border-slate-100">
          <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Customer</th>
          <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Product</th>
          <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Interval</th>
          <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Price</th>
          <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Status</th>
          <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Next Billing</th>
          <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Created</th>
          <th className="px-4 py-3" />
        </tr>
      </thead>
      <tbody>
        {subscriptions.map((sub) => (
          <tr
            key={sub.id}
            onClick={() => onSelect(sub)}
            className={`cursor-pointer border-b border-slate-50 hover:bg-slate-50 ${
              selectedId === sub.id ? 'bg-indigo-50' : ''
            }`}
          >
            <td className="px-4 py-3 font-mono text-xs text-slate-600">{sub.customer_id.slice(0, 8)}…</td>
            <td className="px-4 py-3 font-mono text-xs text-slate-600">{sub.product_id.slice(0, 8)}…</td>
            <td className="px-4 py-3 text-slate-600 capitalize">{sub.interval}</td>
            <td className="px-4 py-3 text-slate-700">{formatMoney(sub.price.amount, sub.price.currency)}</td>
            <td className="px-4 py-3">
              <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[sub.status]}`}>
                {sub.status}
              </span>
            </td>
            <td className="px-4 py-3 text-slate-500 text-xs">{formatDate(sub.next_billing_date)}</td>
            <td className="px-4 py-3 text-slate-400 text-xs">{formatDate(sub.created_at)}</td>
            <td className="px-4 py-3" onClick={(e) => e.stopPropagation()}>
              <RowActions subscription={sub} />
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

// ---------------------------------------------------------------------------
// Main Page
// ---------------------------------------------------------------------------

type Tab = 'all' | 'due'

export default function Subscriptions() {
  const [tab, setTab] = useState<Tab>('all')
  const [page, setPage] = useState(1)
  const [selectedSub, setSelectedSub] = useState<Subscription | null>(null)

  const { data: allData, isPending: allPending, error: allError } = useSubscriptions({ page, page_size: 20 })
  const { data: dueData, isPending: duePending, error: dueError } = useDueSubscriptions()

  const isPending = tab === 'all' ? allPending : duePending
  const error = tab === 'all' ? allError : dueError

  const subscriptions: Subscription[] =
    tab === 'all' ? (allData?.items ?? []) : (dueData ?? [])

  function handleSelect(sub: Subscription) {
    setSelectedSub(selectedSub?.id === sub.id ? null : sub)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-2">
        <RefreshCw size={20} className="text-indigo-600" />
        <h1 className="text-lg font-semibold text-slate-800">Subscriptions</h1>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 rounded-lg border border-slate-200 bg-slate-50 p-1 w-fit">
        {(['all', 'due'] as Tab[]).map((t) => (
          <button
            key={t}
            onClick={() => { setTab(t); setSelectedSub(null) }}
            className={`rounded-md px-4 py-1.5 text-sm font-medium transition-colors ${
              tab === t
                ? 'bg-white text-slate-800 shadow-sm'
                : 'text-slate-500 hover:text-slate-700'
            }`}
          >
            {t === 'all' ? 'All Subscriptions' : 'Due for Billing'}
          </button>
        ))}
      </div>

      <div className={`grid gap-6 ${selectedSub ? 'grid-cols-[1fr_420px]' : 'grid-cols-1'}`}>
        {/* Table */}
        <div className="rounded-xl border border-slate-200 bg-white">
          {isPending ? (
            <p className="p-6 text-sm text-slate-500">Loading…</p>
          ) : error ? (
            <p className="p-6 text-sm text-red-600">{error.message}</p>
          ) : (
            <SubscriptionTable
              subscriptions={subscriptions}
              selectedId={selectedSub?.id ?? null}
              onSelect={handleSelect}
            />
          )}

          {/* Pagination (all tab only) */}
          {tab === 'all' && allData && allData.total_pages > 1 && (
            <div className="flex items-center justify-between border-t border-slate-100 px-4 py-3">
              <p className="text-xs text-slate-400">
                Page {allData.page} of {allData.total_pages} · {allData.total} subscriptions
              </p>
              <div className="flex gap-2">
                <button
                  disabled={page <= 1}
                  onClick={() => setPage((p) => p - 1)}
                  className="rounded border border-slate-200 px-3 py-1 text-xs disabled:opacity-40"
                >
                  Prev
                </button>
                <button
                  disabled={page >= allData.total_pages}
                  onClick={() => setPage((p) => p + 1)}
                  className="rounded border border-slate-200 px-3 py-1 text-xs disabled:opacity-40"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </div>

        {/* Billing Panel */}
        {selectedSub && (
          <BillingPanel subscription={selectedSub} onClose={() => setSelectedSub(null)} />
        )}
      </div>
    </div>
  )
}
