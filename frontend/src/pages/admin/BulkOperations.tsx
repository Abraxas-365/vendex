import { useState } from 'react'
import { Zap, Plus, Play, XCircle, ChevronDown, ChevronUp } from 'lucide-react'
import type { BulkOperation } from '../../types'
import {
  useBulkOperations,
  useCreateBulkOperation,
  useProcessBulkOperation,
  useCancelBulkOperation,
  useBulkOperationItems,
} from '../../lib/hooks'

const emptyForm: Partial<BulkOperation> = {
  type: 'price_update',
}

const statusColors: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-700',
  processing: 'bg-blue-100 text-blue-700',
  completed: 'bg-green-100 text-green-700',
  failed: 'bg-red-100 text-red-700',
  cancelled: 'bg-gray-100 text-gray-600',
}

const opTypes = [
  'price_update',
  'status_update',
  'inventory_update',
  'category_assign',
  'tag_assign',
  'delete',
]

function ItemsPanel({ operationId }: { operationId: string }) {
  const { data, isLoading } = useBulkOperationItems(operationId)
  const items = data ?? []

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-6">
        <div className="h-5 w-5 animate-spin rounded-full border-2 border-gray-300 border-t-indigo-600" />
      </div>
    )
  }

  if (items.length === 0) {
    return <p className="text-sm text-gray-400 py-4">No items in this operation.</p>
  }

  return (
    <table className="w-full text-sm">
      <thead>
        <tr className="border-b border-gray-100">
          <th className="py-2 text-left text-xs font-medium text-gray-400">Resource ID</th>
          <th className="py-2 text-left text-xs font-medium text-gray-400">Status</th>
          <th className="py-2 text-left text-xs font-medium text-gray-400">Error</th>
          <th className="py-2 text-left text-xs font-medium text-gray-400">Processed</th>
        </tr>
      </thead>
      <tbody className="divide-y divide-gray-50">
        {items.map((item) => (
          <tr key={item.id}>
            <td className="py-2 font-mono text-xs text-gray-600">{item.resource_id}</td>
            <td className="py-2">
              <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[item.status] ?? 'bg-gray-100 text-gray-600'}`}>
                {item.status}
              </span>
            </td>
            <td className="py-2 text-xs text-red-500">{item.error || '—'}</td>
            <td className="py-2 text-xs text-gray-500">
              {item.processed_at ? new Date(item.processed_at).toLocaleString() : '—'}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

export default function BulkOperations() {
  const [page] = useState(1)
  const [showDialog, setShowDialog] = useState(false)
  const [form, setForm] = useState<Partial<BulkOperation>>(emptyForm)
  const [expandedId, setExpandedId] = useState<string | null>(null)

  const { data, isLoading } = useBulkOperations({ page, page_size: 20 })
  const operations = data?.items ?? []

  const create = useCreateBulkOperation()
  const process = useProcessBulkOperation()
  const cancel = useCancelBulkOperation()

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    create.mutate(form, { onSuccess: () => setShowDialog(false) })
  }

  function toggleExpand(id: string) {
    setExpandedId((prev) => (prev === id ? null : id))
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Zap className="h-6 w-6 text-indigo-600" />
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Bulk Operations</h1>
            <p className="mt-0.5 text-sm text-gray-500">Run batch updates on products and inventory</p>
          </div>
        </div>
        <button
          onClick={() => {
            setForm(emptyForm)
            setShowDialog(true)
          }}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors"
        >
          <Plus className="h-4 w-4" />
          New Operation
        </button>
      </div>

      {/* Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-indigo-600" />
          </div>
        ) : operations.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Zap className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No bulk operations yet</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead className="border-b border-gray-100 bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Type</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Status</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Progress</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Created By</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Created</th>
                <th className="px-6 py-3 text-right font-medium text-gray-500">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {operations.map((op) => (
                <>
                  <tr key={op.id} className="hover:bg-gray-50/50 transition-colors">
                    <td className="px-6 py-4">
                      <p className="font-medium text-gray-900 capitalize">{op.type.replace('_', ' ')}</p>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[op.status] ?? 'bg-gray-100 text-gray-600'}`}>
                        {op.status}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-2">
                        <div className="flex-1 min-w-16 h-1.5 rounded-full bg-gray-100">
                          <div
                            className="h-full rounded-full bg-indigo-500"
                            style={{
                              width: op.total_items > 0
                                ? `${Math.round((op.processed_items / op.total_items) * 100)}%`
                                : '0%',
                            }}
                          />
                        </div>
                        <span className="text-xs text-gray-500 whitespace-nowrap">
                          {op.processed_items}/{op.total_items}
                        </span>
                      </div>
                      {op.failed_items > 0 && (
                        <p className="text-xs text-red-500 mt-0.5">{op.failed_items} failed</p>
                      )}
                    </td>
                    <td className="px-6 py-4 text-gray-500">{op.created_by || '—'}</td>
                    <td className="px-6 py-4 text-gray-500 text-xs">
                      {new Date(op.created_at).toLocaleString()}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center justify-end gap-1">
                        {op.status === 'pending' && (
                          <button
                            onClick={() => process.mutate(op.id)}
                            className="rounded-md p-1.5 text-gray-400 hover:bg-green-50 hover:text-green-600 transition-colors"
                            title="Process"
                          >
                            <Play className="h-4 w-4" />
                          </button>
                        )}
                        {['pending', 'processing'].includes(op.status) && (
                          <button
                            onClick={() => cancel.mutate(op.id)}
                            className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                            title="Cancel"
                          >
                            <XCircle className="h-4 w-4" />
                          </button>
                        )}
                        <button
                          onClick={() => toggleExpand(op.id)}
                          className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 transition-colors"
                          title="Show items"
                        >
                          {expandedId === op.id ? (
                            <ChevronUp className="h-4 w-4" />
                          ) : (
                            <ChevronDown className="h-4 w-4" />
                          )}
                        </button>
                      </div>
                    </td>
                  </tr>
                  {expandedId === op.id && (
                    <tr key={`${op.id}-items`}>
                      <td colSpan={6} className="px-6 py-4 bg-gray-50/50">
                        <p className="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-3">Items</p>
                        <ItemsPanel operationId={op.id} />
                      </td>
                    </tr>
                  )}
                </>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Dialog */}
      {showDialog && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
          <div className="w-full max-w-md rounded-2xl bg-white shadow-2xl">
            <div className="border-b border-gray-100 px-6 py-4">
              <h2 className="text-lg font-semibold text-gray-900">New Bulk Operation</h2>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 px-6 py-5">
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Operation Type *</label>
                <select
                  value={form.type ?? 'price_update'}
                  onChange={(e) => setForm((f) => ({ ...f, type: e.target.value as BulkOperation['type'] }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                >
                  {opTypes.map((t) => (
                    <option key={t} value={t}>{t.replace(/_/g, ' ')}</option>
                  ))}
                </select>
              </div>
              <div className="rounded-lg bg-amber-50 border border-amber-200 p-3">
                <p className="text-xs text-amber-700 font-medium">Note</p>
                <p className="text-xs text-amber-600 mt-0.5">
                  After creating, add items via the API or upload a CSV, then click Process.
                </p>
              </div>
              {create.error && (
                <p className="text-sm text-red-600">{create.error.message}</p>
              )}
              <div className="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowDialog(false)}
                  className="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={create.isPending}
                  className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 transition-colors"
                >
                  {create.isPending ? 'Creating…' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
