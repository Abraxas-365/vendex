import { useState } from 'react'
import { Plus, Tag, X, Search, ToggleLeft, ToggleRight } from 'lucide-react'
import type { Promo, PromoType } from '../../types'
import { usePromos, useCreatePromo, useUpdatePromo } from '../../lib/hooks'

function formatDate(dateStr?: string): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

function formatValue(type: PromoType, value: number): string {
  switch (type) {
    case 'percentage':
      return `${value}%`
    case 'fixed_amount':
      return `$${value.toFixed(2)}`
    case 'free_shipping':
      return 'Free Shipping'
  }
}

const typeLabels: Record<PromoType, string> = {
  percentage: 'Percentage',
  fixed_amount: 'Fixed Amount',
  free_shipping: 'Free Shipping',
}

interface PromoFormData {
  code: string
  type: PromoType
  value: number
  min_order_amount: string
  max_uses: string
  starts_at: string
  ends_at: string
}

const emptyForm: PromoFormData = {
  code: '',
  type: 'percentage',
  value: 0,
  min_order_amount: '',
  max_uses: '',
  starts_at: '',
  ends_at: '',
}

export default function Promos() {
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<PromoFormData>(emptyForm)
  const [searchQuery, setSearchQuery] = useState('')

  const { data, isLoading } = usePromos()
  const createPromo = useCreatePromo()
  const updatePromo = useUpdatePromo()

  const promos: Promo[] = data?.items ?? []

  const filtered = promos.filter(
    (p) =>
      searchQuery === '' ||
      p.code.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    createPromo.mutate({
      code: formData.code.toUpperCase(),
      type: formData.type,
      value: formData.value,
      min_order_amount: formData.min_order_amount
        ? parseFloat(formData.min_order_amount)
        : undefined,
      max_uses: formData.max_uses ? parseInt(formData.max_uses) : undefined,
      starts_at: formData.starts_at || undefined,
      ends_at: formData.ends_at || undefined,
      active: true,
    })
    setShowForm(false)
    setFormData(emptyForm)
  }

  function handleToggleActive(promo: Promo) {
    updatePromo.mutate({ id: promo.id, data: { active: !promo.active } })
  }

  function handleCancel() {
    setShowForm(false)
    setFormData(emptyForm)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Promos</h1>
          <p className="mt-1 text-sm text-gray-500">
            Create and manage promotional codes
          </p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-gray-800 transition-colors"
        >
          <Plus className="h-4 w-4" />
          Create Promo
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          placeholder="Search by promo code..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full rounded-lg border border-gray-200 bg-white py-2.5 pl-10 pr-4 text-sm text-gray-900 placeholder:text-gray-400 focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
        />
      </div>

      {/* Create Promo Form */}
      {showForm && (
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">New Promo Code</h2>
            <button onClick={handleCancel} className="text-gray-400 hover:text-gray-600">
              <X className="h-5 w-5" />
            </button>
          </div>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <div>
              <label className="block text-sm font-medium text-gray-700">Code</label>
              <input
                type="text"
                required
                value={formData.code}
                onChange={(e) => setFormData({ ...formData, code: e.target.value })}
                placeholder="SUMMER2025"
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm uppercase tracking-wider focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Type</label>
              <select
                value={formData.type}
                onChange={(e) => setFormData({ ...formData, type: e.target.value as PromoType })}
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              >
                <option value="percentage">Percentage</option>
                <option value="fixed_amount">Fixed Amount</option>
                <option value="free_shipping">Free Shipping</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Value</label>
              <input
                type="number"
                step="0.01"
                min="0"
                required={formData.type !== 'free_shipping'}
                disabled={formData.type === 'free_shipping'}
                value={formData.type === 'free_shipping' ? 0 : formData.value}
                onChange={(e) => setFormData({ ...formData, value: parseFloat(e.target.value) })}
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm disabled:bg-gray-50 disabled:text-gray-400 focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Min Order Amount</label>
              <input
                type="number"
                step="0.01"
                min="0"
                value={formData.min_order_amount}
                onChange={(e) => setFormData({ ...formData, min_order_amount: e.target.value })}
                placeholder="Optional"
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Max Uses</label>
              <input
                type="number"
                min="1"
                value={formData.max_uses}
                onChange={(e) => setFormData({ ...formData, max_uses: e.target.value })}
                placeholder="Unlimited"
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
            </div>
            <div>{/* spacer for grid alignment */}</div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Start Date</label>
              <input
                type="date"
                value={formData.starts_at}
                onChange={(e) => setFormData({ ...formData, starts_at: e.target.value })}
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">End Date</label>
              <input
                type="date"
                value={formData.ends_at}
                onChange={(e) => setFormData({ ...formData, ends_at: e.target.value })}
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
            </div>
            <div className="flex items-end sm:col-span-2 lg:col-span-3">
              <div className="flex items-center gap-3">
                <button
                  type="submit"
                  className="rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-gray-800 transition-colors"
                >
                  Create Promo
                </button>
                <button
                  type="button"
                  onClick={handleCancel}
                  className="rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
              </div>
            </div>
          </form>
        </div>
      )}

      {/* Promos Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Tag className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No promo codes found</p>
            <p className="mt-1 text-xs">Create your first promo code to get started.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">Code</th>
                  <th className="px-6 py-3 font-medium">Type</th>
                  <th className="px-6 py-3 font-medium">Value</th>
                  <th className="px-6 py-3 font-medium">Uses</th>
                  <th className="px-6 py-3 font-medium">Active</th>
                  <th className="px-6 py-3 font-medium">Dates</th>
                  <th className="px-6 py-3 font-medium text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((promo) => (
                  <tr key={promo.id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-4">
                      <span className="inline-flex rounded-md bg-gray-100 px-2.5 py-1 font-mono text-sm font-semibold text-gray-900 tracking-wider">
                        {promo.code}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">{typeLabels[promo.type]}</td>
                    <td className="px-6 py-4 text-sm font-medium text-gray-900">
                      {formatValue(promo.type, promo.value)}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {promo.used_count}
                      {promo.max_uses != null ? ` / ${promo.max_uses}` : ' / unlimited'}
                    </td>
                    <td className="px-6 py-4">
                      <span
                        className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${
                          promo.active
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-500'
                        }`}
                      >
                        {promo.active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {formatDate(promo.starts_at)} — {formatDate(promo.ends_at)}
                    </td>
                    <td className="px-6 py-4 text-right">
                      <button
                        onClick={() => handleToggleActive(promo)}
                        className={`rounded-md p-1.5 transition-colors ${
                          promo.active
                            ? 'text-green-600 hover:bg-red-50 hover:text-red-600'
                            : 'text-gray-400 hover:bg-green-50 hover:text-green-600'
                        }`}
                        title={promo.active ? 'Deactivate' : 'Activate'}
                      >
                        {promo.active ? (
                          <ToggleRight className="h-5 w-5" />
                        ) : (
                          <ToggleLeft className="h-5 w-5" />
                        )}
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
