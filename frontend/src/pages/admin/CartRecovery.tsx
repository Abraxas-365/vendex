import { useState } from 'react'
import { MailQuestion } from 'lucide-react'
import type { RecoveryEmail } from '../../types'
import { useRecoveryEmails, useRecoveryStats, useUpdateRecoveryStatus } from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const statusColors: Record<string, string> = {
  sent: 'bg-blue-50 text-blue-700',
  clicked: 'bg-amber-50 text-amber-700',
  converted: 'bg-green-50 text-green-700',
  pending: 'bg-slate-100 text-slate-500',
}

// ---------------------------------------------------------------------------
// Stats Card
// ---------------------------------------------------------------------------

interface StatCardProps {
  label: string
  value: string | number
}

function StatCard({ label, value }: StatCardProps) {
  return (
    <div className="rounded-xl border border-slate-200 bg-white p-4">
      <p className="text-xs font-medium uppercase tracking-wide text-slate-400">{label}</p>
      <p className="mt-1 text-2xl font-semibold text-slate-800">{value}</p>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Status Dropdown
// ---------------------------------------------------------------------------

interface StatusSelectProps {
  email: RecoveryEmail
}

function StatusSelect({ email }: StatusSelectProps) {
  const updateStatus = useUpdateRecoveryStatus()

  function handleChange(e: React.ChangeEvent<HTMLSelectElement>) {
    e.stopPropagation()
    updateStatus.mutate({ id: email.id, status: e.target.value })
  }

  return (
    <select
      value={email.status}
      onChange={handleChange}
      onClick={(e) => e.stopPropagation()}
      className="rounded border border-slate-200 px-2 py-1 text-xs text-slate-700 outline-none focus:border-indigo-400"
    >
      <option value="sent">sent</option>
      <option value="clicked">clicked</option>
      <option value="converted">converted</option>
    </select>
  )
}

// ---------------------------------------------------------------------------
// Main Page
// ---------------------------------------------------------------------------

export default function CartRecovery() {
  const [page, setPage] = useState(1)

  const { data: stats, isPending: statsLoading } = useRecoveryStats()
  const { data, isPending, error } = useRecoveryEmails({ page, page_size: 20 })

  const emails = data?.items ?? []

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-2">
        <MailQuestion size={20} className="text-indigo-600" />
        <h1 className="text-lg font-semibold text-slate-800">Cart Recovery</h1>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-5">
        {statsLoading ? (
          <p className="col-span-5 text-sm text-slate-500">Loading stats…</p>
        ) : stats ? (
          <>
            <StatCard label="Total" value={stats.total} />
            <StatCard label="Sent" value={stats.sent} />
            <StatCard label="Clicked" value={stats.clicked} />
            <StatCard label="Converted" value={stats.converted} />
            <StatCard label="Conversion Rate" value={`${(stats.conversion_rate * 100).toFixed(1)}%`} />
          </>
        ) : null}
      </div>

      {/* Table */}
      <div className="rounded-xl border border-slate-200 bg-white">
        <div className="border-b border-slate-100 px-5 py-3">
          <h2 className="text-sm font-semibold text-slate-700">Recovery Emails</h2>
        </div>

        {isPending ? (
          <p className="p-6 text-sm text-slate-500">Loading…</p>
        ) : error ? (
          <p className="p-6 text-sm text-red-600">{error.message}</p>
        ) : emails.length === 0 ? (
          <p className="p-6 text-sm text-slate-400">No recovery emails yet.</p>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-100">
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Email</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Cart ID</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Step</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Status</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Sent At</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Clicked At</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Converted At</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Update</th>
              </tr>
            </thead>
            <tbody>
              {emails.map((email) => (
                <tr key={email.id} className="border-b border-slate-50 hover:bg-slate-50">
                  <td className="px-4 py-3 font-medium text-slate-700">{email.email}</td>
                  <td className="px-4 py-3 font-mono text-xs text-slate-500">{email.cart_id}</td>
                  <td className="px-4 py-3 text-slate-600">{email.step}</td>
                  <td className="px-4 py-3">
                    <span
                      className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                        statusColors[email.status] ?? 'bg-slate-100 text-slate-500'
                      }`}
                    >
                      {email.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-slate-400 text-xs">{formatDate(email.sent_at)}</td>
                  <td className="px-4 py-3 text-slate-400 text-xs">{formatDate(email.clicked_at)}</td>
                  <td className="px-4 py-3 text-slate-400 text-xs">{formatDate(email.converted_at)}</td>
                  <td className="px-4 py-3">
                    <StatusSelect email={email} />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}

        {/* Pagination */}
        {data && data.total_pages > 1 && (
          <div className="flex items-center justify-between border-t border-slate-100 px-4 py-3">
            <p className="text-xs text-slate-400">
              Page {data.page} of {data.total_pages} · {data.total} emails
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
                disabled={page >= data.total_pages}
                onClick={() => setPage((p) => p + 1)}
                className="rounded border border-slate-200 px-3 py-1 text-xs disabled:opacity-40"
              >
                Next
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
