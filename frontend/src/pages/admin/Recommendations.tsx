import { useState } from 'react'
import { Sparkles, Plus, Trash2, Edit2, Search } from 'lucide-react'
import type { RecommendationRule, RecommendedProduct } from '../../types'
import {
  useRecommendationRules,
  useCreateRecommendationRule,
  useUpdateRecommendationRule,
  useDeleteRecommendationRule,
  useRecommendations,
} from '../../lib/hooks'

const emptyForm: Partial<RecommendationRule> = {
  name: '',
  type: 'similar',
  source_product_id: '',
  weight: 1,
  is_active: true,
}

const typeColors: Record<string, string> = {
  similar: 'bg-blue-100 text-blue-700',
  frequently_bought: 'bg-green-100 text-green-700',
  trending: 'bg-orange-100 text-orange-700',
  personalized: 'bg-purple-100 text-purple-700',
}

function TestPanel() {
  const [productId, setProductId] = useState('')
  const [query, setQuery] = useState('')

  const { data, isLoading, error } = useRecommendations(query)

  function handleSearch(e: React.FormEvent) {
    e.preventDefault()
    setQuery(productId)
  }

  return (
    <div className="rounded-xl border border-gray-200 bg-white shadow-sm p-6 space-y-4">
      <h2 className="text-base font-semibold text-gray-900">Test Recommendations</h2>
      <form onSubmit={handleSearch} className="flex gap-3">
        <input
          value={productId}
          onChange={(e) => setProductId(e.target.value)}
          placeholder="Enter a product ID…"
          className="flex-1 rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <button
          type="submit"
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors"
        >
          <Search className="h-4 w-4" />
          Get Recommendations
        </button>
      </form>

      {isLoading && (
        <div className="flex items-center justify-center py-6">
          <div className="h-5 w-5 animate-spin rounded-full border-2 border-gray-300 border-t-indigo-600" />
        </div>
      )}

      {error && <p className="text-sm text-red-600">{error.message}</p>}

      {!isLoading && query && (data ?? []).length === 0 && (
        <p className="text-sm text-gray-400 py-4">No recommendations found for this product.</p>
      )}

      {(data ?? []).length > 0 && (
        <div className="space-y-2">
          {(data as RecommendedProduct[]).map((rec) => (
            <div
              key={rec.product_id}
              className="flex items-center justify-between rounded-lg border border-gray-100 p-3"
            >
              <div>
                <p className="text-sm font-medium text-gray-900">{rec.name}</p>
                <p className="text-xs text-gray-400 mt-0.5">{rec.reason}</p>
              </div>
              <div className="text-right">
                <p className="text-xs font-medium text-indigo-600">Score</p>
                <p className="text-sm font-bold text-gray-900">{rec.score.toFixed(3)}</p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default function Recommendations() {
  const [page] = useState(1)
  const [showDialog, setShowDialog] = useState(false)
  const [editing, setEditing] = useState<RecommendationRule | null>(null)
  const [form, setForm] = useState<Partial<RecommendationRule>>(emptyForm)

  const { data, isLoading } = useRecommendationRules({ page, page_size: 20 })
  const rules = data?.items ?? []

  const create = useCreateRecommendationRule()
  const update = useUpdateRecommendationRule()
  const remove = useDeleteRecommendationRule()

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setShowDialog(true)
  }

  function openEdit(rule: RecommendationRule) {
    setEditing(rule)
    setForm({
      name: rule.name,
      type: rule.type,
      source_product_id: rule.source_product_id,
      weight: rule.weight,
      is_active: rule.is_active,
    })
    setShowDialog(true)
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (editing) {
      update.mutate(
        { id: editing.id, data: form },
        { onSuccess: () => setShowDialog(false) },
      )
    } else {
      create.mutate(form, { onSuccess: () => setShowDialog(false) })
    }
  }

  const isPending = create.isPending || update.isPending

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Sparkles className="h-6 w-6 text-indigo-600" />
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Recommendations</h1>
            <p className="mt-0.5 text-sm text-gray-500">Configure product recommendation rules</p>
          </div>
        </div>
        <button
          onClick={openCreate}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors"
        >
          <Plus className="h-4 w-4" />
          New Rule
        </button>
      </div>

      {/* Rules Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-indigo-600" />
          </div>
        ) : rules.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Sparkles className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No recommendation rules yet</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead className="border-b border-gray-100 bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Name</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Type</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Source Product</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Weight</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Status</th>
                <th className="px-6 py-3 text-right font-medium text-gray-500">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {rules.map((rule) => (
                <tr key={rule.id} className="hover:bg-gray-50/50 transition-colors">
                  <td className="px-6 py-4 font-medium text-gray-900">{rule.name}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${typeColors[rule.type] ?? 'bg-gray-100 text-gray-600'}`}>
                      {rule.type.replace('_', ' ')}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-gray-500 text-xs font-mono">{rule.source_product_id || '—'}</td>
                  <td className="px-6 py-4 text-gray-600">{rule.weight}</td>
                  <td className="px-6 py-4">
                    <span
                      className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${
                        rule.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'
                      }`}
                    >
                      {rule.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        onClick={() => openEdit(rule)}
                        className="rounded-md p-1.5 text-gray-400 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                        title="Edit"
                      >
                        <Edit2 className="h-4 w-4" />
                      </button>
                      <button
                        onClick={() => remove.mutate(rule.id)}
                        className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                        title="Delete"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Test Panel */}
      <TestPanel />

      {/* Dialog */}
      {showDialog && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
          <div className="w-full max-w-md rounded-2xl bg-white shadow-2xl">
            <div className="border-b border-gray-100 px-6 py-4">
              <h2 className="text-lg font-semibold text-gray-900">
                {editing ? 'Edit Rule' : 'New Rule'}
              </h2>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 px-6 py-5">
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Name *</label>
                <input
                  required
                  value={form.name ?? ''}
                  onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Type</label>
                <select
                  value={form.type ?? 'similar'}
                  onChange={(e) =>
                    setForm((f) => ({
                      ...f,
                      type: e.target.value as RecommendationRule['type'],
                    }))
                  }
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                >
                  <option value="similar">Similar</option>
                  <option value="frequently_bought">Frequently Bought Together</option>
                  <option value="trending">Trending</option>
                  <option value="personalized">Personalized</option>
                </select>
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Source Product ID</label>
                <input
                  value={form.source_product_id ?? ''}
                  onChange={(e) => setForm((f) => ({ ...f, source_product_id: e.target.value }))}
                  placeholder="Leave empty for global rule"
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Weight</label>
                <input
                  type="number"
                  min={0}
                  step={0.1}
                  value={form.weight ?? 1}
                  onChange={(e) => setForm((f) => ({ ...f, weight: parseFloat(e.target.value) || 1 }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="rec_is_active"
                  checked={form.is_active ?? true}
                  onChange={(e) => setForm((f) => ({ ...f, is_active: e.target.checked }))}
                  className="rounded border-gray-300"
                />
                <label htmlFor="rec_is_active" className="text-sm text-gray-700">Active</label>
              </div>
              {(create.error ?? update.error) && (
                <p className="text-sm text-red-600">{(create.error ?? update.error)?.message}</p>
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
                  disabled={isPending}
                  className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 transition-colors"
                >
                  {isPending ? 'Saving…' : editing ? 'Save Changes' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
