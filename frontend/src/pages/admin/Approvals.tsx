import { useState } from 'react'
import {
  ShieldCheck,
  CheckCircle2,
  XCircle,
  Clock,
  ChevronDown,
  ChevronRight,
} from 'lucide-react'
import {
  useApprovals,
  useApprovalCount,
  useApproveRequest,
  useRejectRequest,
} from '../../lib/hooks'
import type { ApprovalRequest } from '../../types'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

type StatusFilter = 'pending' | 'all'

function formatDate(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleString()
  } catch {
    return dateStr
  }
}

function StatusBadge({ status }: { status: ApprovalRequest['status'] }) {
  const classes = {
    pending: 'bg-amber-100 text-amber-700',
    approved: 'bg-green-100 text-green-700',
    rejected: 'bg-red-100 text-red-700',
  }
  const icons = {
    pending: <Clock size={11} className="shrink-0" />,
    approved: <CheckCircle2 size={11} className="shrink-0" />,
    rejected: <XCircle size={11} className="shrink-0" />,
  }
  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium capitalize ${classes[status]}`}
    >
      {icons[status]}
      {status}
    </span>
  )
}

// ---------------------------------------------------------------------------
// JSON Viewer (collapsible)
// ---------------------------------------------------------------------------

function JsonViewer({ data }: { data: Record<string, unknown> }) {
  const [expanded, setExpanded] = useState(false)
  const json = JSON.stringify(data, null, 2)

  return (
    <div className="mt-1">
      <button
        onClick={() => setExpanded(!expanded)}
        className="inline-flex items-center gap-1 text-xs text-indigo-600 hover:text-indigo-800 transition-colors"
      >
        {expanded ? <ChevronDown size={12} /> : <ChevronRight size={12} />}
        {expanded ? 'Hide' : 'View'} tool input
      </button>
      {expanded && (
        <pre className="mt-2 rounded-lg bg-slate-900 p-3 text-xs text-slate-100 overflow-x-auto max-h-48">
          {json}
        </pre>
      )}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Approve/Reject dialog
// ---------------------------------------------------------------------------

interface ActionDialogProps {
  approval: ApprovalRequest
  action: 'approve' | 'reject'
  onClose: () => void
}

function ActionDialog({ approval, action, onClose }: ActionDialogProps) {
  const [reason, setReason] = useState('')
  const approveMutation = useApproveRequest()
  const rejectMutation = useRejectRequest()

  const isPending = approveMutation.isPending || rejectMutation.isPending
  const error = approveMutation.error ?? rejectMutation.error

  function handleSubmit() {
    const mutation = action === 'approve' ? approveMutation : rejectMutation
    mutation.mutate(
      { id: approval.id, reason },
      {
        onSuccess: onClose,
      },
    )
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="w-full max-w-md rounded-xl bg-white shadow-xl p-6 mx-4">
        <h2 className="text-base font-semibold text-slate-800 mb-1">
          {action === 'approve' ? 'Approve' : 'Reject'} Request
        </h2>
        <p className="text-sm text-slate-500 mb-4">
          Tool: <span className="font-mono font-medium text-slate-700">{approval.tool_name}</span>
        </p>

        <label className="block text-xs font-medium text-slate-600 mb-1">
          Reason (optional)
        </label>
        <textarea
          value={reason}
          onChange={(e) => setReason(e.target.value)}
          placeholder={`Reason for ${action}ing this request…`}
          rows={3}
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100 resize-none"
        />

        {error && (
          <p className="mt-2 text-xs text-red-600">{error.message}</p>
        )}

        <div className="mt-4 flex items-center justify-end gap-2">
          <button
            onClick={onClose}
            disabled={isPending}
            className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50 disabled:opacity-50 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={isPending}
            className={`rounded-lg px-4 py-2 text-sm font-medium text-white disabled:opacity-50 transition-colors ${
              action === 'approve'
                ? 'bg-green-600 hover:bg-green-700'
                : 'bg-red-600 hover:bg-red-700'
            }`}
          >
            {isPending
              ? action === 'approve'
                ? 'Approving…'
                : 'Rejecting…'
              : action === 'approve'
              ? 'Approve'
              : 'Reject'}
          </button>
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Approval row
// ---------------------------------------------------------------------------

interface ApprovalRowProps {
  approval: ApprovalRequest
  onAction: (approval: ApprovalRequest, action: 'approve' | 'reject') => void
}

function ApprovalRow({ approval, onAction }: ApprovalRowProps) {
  return (
    <tr className="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors">
      <td className="py-3 px-4">
        <p className="text-sm font-medium text-slate-800 font-mono">{approval.tool_name}</p>
        <JsonViewer data={approval.tool_input} />
      </td>
      <td className="py-3 px-4 text-sm text-slate-500">
        {approval.requested_by || <span className="text-slate-300 italic">—</span>}
      </td>
      <td className="py-3 px-4">
        <span className="font-mono text-xs text-slate-400 truncate max-w-[120px] block" title={approval.session_id}>
          {approval.session_id.slice(0, 12)}…
        </span>
      </td>
      <td className="py-3 px-4 text-sm text-slate-500 whitespace-nowrap">
        {formatDate(approval.created_at)}
      </td>
      <td className="py-3 px-4">
        <StatusBadge status={approval.status} />
        {approval.reason && (
          <p className="mt-1 text-xs text-slate-400 italic max-w-xs truncate" title={approval.reason}>
            {approval.reason}
          </p>
        )}
      </td>
      <td className="py-3 px-4">
        {approval.status === 'pending' ? (
          <div className="flex items-center justify-end gap-2">
            <button
              onClick={() => onAction(approval, 'approve')}
              className="inline-flex items-center gap-1 rounded-lg bg-green-50 border border-green-200 px-2.5 py-1.5 text-xs font-medium text-green-700 hover:bg-green-100 transition-colors"
            >
              <CheckCircle2 size={12} />
              Approve
            </button>
            <button
              onClick={() => onAction(approval, 'reject')}
              className="inline-flex items-center gap-1 rounded-lg bg-red-50 border border-red-200 px-2.5 py-1.5 text-xs font-medium text-red-700 hover:bg-red-100 transition-colors"
            >
              <XCircle size={12} />
              Reject
            </button>
          </div>
        ) : (
          <div className="flex justify-end">
            <span className="text-xs text-slate-400 italic">
              by {approval.reviewed_by || '—'}
            </span>
          </div>
        )}
      </td>
    </tr>
  )
}

// ---------------------------------------------------------------------------
// Main page
// ---------------------------------------------------------------------------

export default function Approvals() {
  const [tab, setTab] = useState<StatusFilter>('pending')
  const [dialogState, setDialogState] = useState<{
    approval: ApprovalRequest
    action: 'approve' | 'reject'
  } | null>(null)

  const { data: countData } = useApprovalCount()
  const pendingCount = countData?.count ?? 0

  const { data: approvalsPage, isLoading, error } = useApprovals(
    tab === 'pending' ? { status: 'pending' } : undefined,
  )

  const approvals = approvalsPage?.items ?? []

  function handleAction(approval: ApprovalRequest, action: 'approve' | 'reject') {
    setDialogState({ approval, action })
  }

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <ShieldCheck size={24} className="text-indigo-600" />
          <div>
            <h1 className="text-xl font-semibold text-slate-800">Approvals</h1>
            <p className="text-sm text-slate-500">Review and manage agent tool approval requests</p>
          </div>
        </div>
        {pendingCount > 0 && (
          <span className="rounded-full bg-amber-100 px-3 py-1 text-sm font-medium text-amber-700">
            {pendingCount} pending
          </span>
        )}
      </div>

      {/* Tabs */}
      <div className="flex gap-1 rounded-lg border border-slate-200 bg-white p-1 w-fit">
        {(['pending', 'all'] as const).map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`rounded-md px-4 py-1.5 text-sm font-medium capitalize transition-colors ${
              tab === t
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'text-slate-600 hover:text-slate-800 hover:bg-slate-100'
            }`}
          >
            {t === 'pending' ? (
              <>
                Pending
                {pendingCount > 0 && (
                  <span className="ml-1.5 rounded-full bg-amber-500/20 px-1.5 text-xs">
                    {pendingCount}
                  </span>
                )}
              </>
            ) : (
              'All'
            )}
          </button>
        ))}
      </div>

      {/* Table */}
      <div className="rounded-xl border border-slate-200 bg-white overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-sm text-slate-500">Loading…</div>
        ) : error ? (
          <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-600 m-4">
            Failed to load approvals. Please try again.
          </div>
        ) : approvals.length === 0 ? (
          <div className="p-12 text-center">
            <ShieldCheck size={32} className="mx-auto mb-3 text-slate-300" />
            <p className="text-sm text-slate-500">
              {tab === 'pending' ? 'No pending approvals' : 'No approvals found'}
            </p>
          </div>
        ) : (
          <table className="w-full">
            <thead>
              <tr className="border-b border-slate-200 bg-slate-50">
                <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Tool Name
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Requested By
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Session ID
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Created At
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Status
                </th>
                <th className="py-3 px-4 text-right text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              {approvals.map((approval) => (
                <ApprovalRow
                  key={approval.id}
                  approval={approval}
                  onAction={handleAction}
                />
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Action dialog */}
      {dialogState && (
        <ActionDialog
          approval={dialogState.approval}
          action={dialogState.action}
          onClose={() => setDialogState(null)}
        />
      )}
    </div>
  )
}
