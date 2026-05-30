import { useState } from 'react'
import { Plus, Pencil, Trash2, X, Webhook as WebhookIcon, RefreshCw, ChevronDown, ChevronUp } from 'lucide-react'
import type { Webhook } from '../../types'
import {
  useWebhooks,
  useCreateWebhook,
  useUpdateWebhook,
  useDeleteWebhook,
  useToggleWebhook,
  useWebhookDeliveries,
  useRetryWebhookDelivery,
} from '../../lib/hooks'

interface WebhookFormData {
  url: string
  events: string
  active: boolean
}

const emptyForm: WebhookFormData = { url: '', events: '', active: true }

function DeliveryList({ webhookId }: { webhookId: string }) {
  const { data: deliveries = [], isLoading } = useWebhookDeliveries(webhookId)
  const retry = useRetryWebhookDelivery()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-4">
        <div className="h-5 w-5 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
      </div>
    )
  }

  if (deliveries.length === 0) {
    return <p className="py-4 text-sm text-gray-400 text-center">No deliveries yet</p>
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-gray-100 text-left text-xs text-gray-500">
            <th className="px-4 py-2 font-medium">Event</th>
            <th className="px-4 py-2 font-medium">Status</th>
            <th className="px-4 py-2 font-medium">Response</th>
            <th className="px-4 py-2 font-medium">Attempts</th>
            <th className="px-4 py-2 font-medium">Date</th>
            <th className="px-4 py-2 font-medium text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {deliveries.map((d) => (
            <tr key={d.id} className="border-b border-gray-50 last:border-0">
              <td className="px-4 py-2 font-mono text-xs text-gray-700">{d.event}</td>
              <td className="px-4 py-2">
                <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${
                  d.status === 'success' ? 'bg-green-100 text-green-800' :
                  d.status === 'failed' ? 'bg-red-100 text-red-800' :
                  'bg-yellow-100 text-yellow-800'
                }`}>
                  {d.status}
                </span>
              </td>
              <td className="px-4 py-2 text-gray-600">{d.response_code ?? '—'}</td>
              <td className="px-4 py-2 text-gray-600">{d.attempts}</td>
              <td className="px-4 py-2 text-gray-500">{new Date(d.created_at).toLocaleDateString()}</td>
              <td className="px-4 py-2 text-right">
                {d.status === 'failed' && (
                  <button
                    onClick={() => retry.mutate({ deliveryId: d.id, webhookId })}
                    className="inline-flex items-center gap-1 rounded-md border border-gray-200 px-2 py-1 text-xs font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                  >
                    <RefreshCw className="h-3 w-3" />
                    Retry
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function Webhooks() {
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<WebhookFormData>(emptyForm)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [expandedId, setExpandedId] = useState<string | null>(null)

  const { data: webhooks = [], isLoading } = useWebhooks()
  const createWebhook = useCreateWebhook()
  const updateWebhook = useUpdateWebhook()
  const deleteWebhook = useDeleteWebhook()
  const toggleWebhook = useToggleWebhook()

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const payload = {
      url: formData.url,
      events: formData.events.split(',').map((s) => s.trim()).filter(Boolean),
      active: formData.active,
    }
    if (editingId) {
      updateWebhook.mutate({ id: editingId, data: payload })
    } else {
      createWebhook.mutate(payload)
    }
    setShowForm(false)
    setFormData(emptyForm)
    setEditingId(null)
  }

  function handleEdit(w: Webhook) {
    setFormData({ url: w.url, events: w.events.join(', '), active: w.active })
    setEditingId(w.id)
    setShowForm(true)
  }

  function handleDelete(id: string) {
    if (confirm('Delete this webhook?')) deleteWebhook.mutate(id)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Webhooks</h1>
          <p className="mt-1 text-sm text-gray-500">Register HTTP callbacks for events</p>
        </div>
        <button
          onClick={() => { setShowForm(true); setEditingId(null); setFormData(emptyForm) }}
          className="inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-gray-800 transition-colors"
        >
          <Plus className="h-4 w-4" />
          Add Webhook
        </button>
      </div>

      {/* Form */}
      {showForm && (
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">
              {editingId ? 'Edit Webhook' : 'New Webhook'}
            </h2>
            <button onClick={() => setShowForm(false)} className="text-gray-400 hover:text-gray-600">
              <X className="h-5 w-5" />
            </button>
          </div>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="sm:col-span-2">
              <label className="block text-sm font-medium text-gray-700">URL</label>
              <input
                type="url"
                required
                value={formData.url}
                onChange={(e) => setFormData({ ...formData, url: e.target.value })}
                placeholder="https://your-server.com/webhook"
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
            </div>
            <div className="sm:col-span-2">
              <label className="block text-sm font-medium text-gray-700">Events (comma-separated)</label>
              <input
                type="text"
                required
                value={formData.events}
                onChange={(e) => setFormData({ ...formData, events: e.target.value })}
                placeholder="order.created, order.updated, payment.success"
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
            </div>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="active"
                checked={formData.active}
                onChange={(e) => setFormData({ ...formData, active: e.target.checked })}
                className="rounded border-gray-300"
              />
              <label htmlFor="active" className="text-sm font-medium text-gray-700">Active</label>
            </div>
            <div className="flex items-end gap-3 sm:col-span-2">
              <button
                type="submit"
                className="rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-gray-800 transition-colors"
              >
                {editingId ? 'Update Webhook' : 'Create Webhook'}
              </button>
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Webhooks Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : webhooks.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <WebhookIcon className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No webhooks registered</p>
          </div>
        ) : (
          <div className="divide-y divide-gray-50">
            {webhooks.map((w) => (
              <div key={w.id}>
                <div className="flex items-center gap-4 px-6 py-4">
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium text-gray-900 truncate">{w.url}</div>
                    <div className="mt-0.5 flex flex-wrap gap-1">
                      {w.events.map((ev) => (
                        <span key={ev} className="inline-flex rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600">{ev}</span>
                      ))}
                    </div>
                  </div>
                  <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${w.active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'}`}>
                    {w.active ? 'Active' : 'Inactive'}
                  </span>
                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => toggleWebhook.mutate(w.id)}
                      className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
                      title={w.active ? 'Deactivate' : 'Activate'}
                    >
                      <RefreshCw className="h-4 w-4" />
                    </button>
                    <button
                      onClick={() => handleEdit(w)}
                      className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
                    >
                      <Pencil className="h-4 w-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(w.id)}
                      className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                    <button
                      onClick={() => setExpandedId(expandedId === w.id ? null : w.id)}
                      className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
                      title="Delivery history"
                    >
                      {expandedId === w.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                    </button>
                  </div>
                </div>
                {expandedId === w.id && (
                  <div className="border-t border-gray-50 bg-gray-50 px-6 py-3">
                    <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-gray-400">Delivery History</p>
                    <DeliveryList webhookId={w.id} />
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
